package administracion

import (
	g "github.com/sisoputnfrba/tp-golang/memoria/estructuras"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"io"
	"os"
)

func RecolectarEntradasParaSwap(pid int) (entradas map[int]g.EntradaSwap) {
	contador := 0
	entradas = make(map[int]g.EntradaSwap)

	g.MutexProcesosPorPID.Lock()
	proceso := g.ProcesosPorPID[pid]
	g.MutexProcesosPorPID.Unlock()

	for _, subtabla := range proceso.TablaRaiz {
		RecorrerTablaPaginasParaSwap(subtabla, entradas, &contador)
	}

	return
}

func RecorrerTablaPaginasParaSwap(tabla *g.TablaPagina, entradas map[int]g.EntradaSwap, contador *int) {
	if tabla.Subtabla != nil {
		for _, subTabla := range tabla.Subtabla {
			RecorrerTablaPaginasParaSwap(subTabla, entradas, contador)
		}
		return
	}

	for _, entrada := range tabla.EntradasPaginas {
		tamanioPagina := g.MemoryConfig.PagSize

		entrada.EstaPresente = false
		inicio := entrada.NumeroFrame * tamanioPagina
		fin := inicio + tamanioPagina

		vacio := make([]byte, g.MemoryConfig.PagSize)
		datos := make([]byte, g.MemoryConfig.PagSize)

		g.MutexMemoriaPrincipal.Lock()
		copy(datos, g.MemoriaPrincipal[inicio:fin])
		copy(g.MemoriaPrincipal[inicio:fin], vacio)
		g.MutexMemoriaPrincipal.Unlock()

		logger.Debug("Lo copiado es: %q", datos)

		entraditaNueva := g.EntradaSwap{
			NumeroPagina: *contador,
			Datos:        datos,
		}

		MarcarLibreFrame(entrada.NumeroFrame)

		entradas[entraditaNueva.NumeroPagina] = entraditaNueva
		*contador++
	}
}

func CargarEntradasASwap(pid int, entradas map[int]g.EntradaSwap) error {

	archivo, err := os.OpenFile(g.MemoryConfig.SwapfilePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer func(archivo *os.File) {
		err := archivo.Close()
		if err != nil {
			logger.Error("Error al cerrar: %v", err)
		}
	}(archivo)

	_, errSeek := archivo.Seek(0, io.SeekEnd) // siempre se apunta al final del archivo! (pense que no, alto bobo)
	if errSeek != nil {
		logger.Error("Error al setear el puntero para SWAP: %v", errSeek)
		return errSeek
	}

	var info = &g.SwapProcesoInfo{
		Entradas:             make(map[int]*g.EntradaSwapInfo),
		NumerosDePaginas:     make([]int, 0),
		InstruccionesEnBytes: make(map[int][]byte),
	}

	g.MutexProcesosPorPID.Lock()
	info.InstruccionesEnBytes = g.ProcesosPorPID[pid].InstruccionesEnBytes
	g.ProcesosPorPID[pid].InstruccionesEnBytes = nil
	g.MutexProcesosPorPID.Unlock()

	posicionPunteroArchivo := g.PunteroSwap

	for i := 0; i < len(entradas); i++ {
		entrada := entradas[i]

		_, errWrite := archivo.Write(entrada.Datos)
		if errWrite != nil {
			logger.Error("Error al escribir el archivo: %v", errWrite)
			return errWrite
		}

		longitudEscrito := len(entrada.Datos)
		posicionPunteroArchivo += longitudEscrito

		if entrada.NumeroPagina == 0 {
			posicionPunteroArchivo = g.PunteroSwap
		}

		info.Entradas[entrada.NumeroPagina] = &g.EntradaSwapInfo{
			NumeroPagina:   entrada.NumeroPagina,
			Tamanio:        longitudEscrito,
			PosicionInicio: posicionPunteroArchivo,
		}
		info.NumerosDePaginas = append(info.NumerosDePaginas, entrada.NumeroPagina)

		VerificarLecturaDesdeSwap(archivo, posicionPunteroArchivo, longitudEscrito)

		logger.Warn("PunteroSWAP actual: %d", g.PunteroSwap)
		logger.Warn("Posición guardada en SwapIndex: %d", posicionPunteroArchivo)

		logger.Info("## PID: <%d> - <MEMORIA A SWAP> - Entrada: <%d> - Posición en SWAP: <%d> - Tamaño: <%d>",
			pid,
			entrada.NumeroPagina,
			posicionPunteroArchivo,
			longitudEscrito,
		)
	}
	g.MutexSwapIndex.Lock()
	g.SwapIndex[pid] = info
	g.MutexSwapIndex.Unlock()

	g.PunteroSwap = posicionPunteroArchivo

	return nil
}
