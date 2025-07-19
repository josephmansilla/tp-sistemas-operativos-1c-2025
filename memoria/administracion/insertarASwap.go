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

		g.MutexMemoriaPrincipal.Lock()
		datos := g.MemoriaPrincipal[inicio:fin]
		copy(g.MemoriaPrincipal[inicio:fin], vacio)
		g.MutexMemoriaPrincipal.Unlock()

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

	file, err := os.OpenFile(g.MemoryConfig.SwapfilePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logger.Error("Error al cerrar: %v", err)
		}
	}(file)

	_, errSeek := file.Seek(0, io.SeekEnd) // siempre se apunta al final del archivo! (pense que no, alto bobo)
	if errSeek != nil {
		logger.Error("Error al setear el puntero para SWAP: %v", errSeek)
		return errSeek
	}

	var info = &g.SwapProcesoInfo{
		Entradas:         make(map[int]*g.EntradaSwapInfo),
		NumerosDePaginas: make([]int, 0),
	}
	var posicionPunteroArchivo = 0
	for i := 0; i < len(entradas); i++ {
		entrada := entradas[i]

		_, errWrite := file.Write(entrada.Datos)
		if errWrite != nil {
			logger.Error("Error al escribir el archivo: %v", errWrite)
			return errWrite
		}

		longitudEscrito := len(entrada.Datos)
		posicionPunteroArchivo += longitudEscrito

		if entrada.NumeroPagina == 0 {
			posicionPunteroArchivo = 0
		}

		info.Entradas[entrada.NumeroPagina] = &g.EntradaSwapInfo{
			NumeroPagina:   entrada.NumeroPagina,
			Tamanio:        longitudEscrito,
			PosicionInicio: posicionPunteroArchivo,
		}
		info.NumerosDePaginas = append(info.NumerosDePaginas, entrada.NumeroPagina)

		g.MutexSwapIndex.Lock()
		g.SwapIndex[pid] = info
		g.MutexSwapIndex.Unlock()

		logger.Info("## PID: <%d> - <MEMORIA A SWAP> - Entrada: <%d> - Posición en SWAP: <%d> - Tamaño: <%d>",
			pid,
			entrada.NumeroPagina,
			posicionPunteroArchivo,
			longitudEscrito,
		)
	}

	return nil
}
