package administracion

import (
	g "github.com/sisoputnfrba/tp-golang/memoria/estructuras"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"io"
	"os"
	"sync"
)

func RecolectarEntradasProcesoSwap(proceso g.Proceso) (resultados []int) {
	cantidadEntradas := g.MemoryConfig.EntriesPerPage
	var wg sync.WaitGroup
	canal := make(chan int, cantidadEntradas)

	for _, subtabla := range proceso.TablaRaiz {
		wg.Add(1)
		go func(st *g.TablaPagina) {
			defer wg.Done()
			RecorrerTablaPaginaDeFormaConcurrenteSwap(st, canal, proceso.PID)
		}(subtabla)
	}

	go func() {
		wg.Wait()
		close(canal)
	}()

	for numeroFrame := range canal {
		resultados = append(resultados, numeroFrame)
	}

	return
}

func RecorrerTablaPaginaDeFormaConcurrenteSwap(tabla *g.TablaPagina, canal chan int, pid int) {

	if tabla.Subtabla != nil {
		for _, subTabla := range tabla.Subtabla {
			RecorrerTablaPaginaDeFormaConcurrenteSwap(subTabla, canal, pid)
		}
		return
	}
	for i, entrada := range tabla.EntradasPaginas {
		if tabla.EntradasPaginas[i].EstaPresente {
			canal <- entrada.NumeroFrame
			entrada.EstaPresente = false
			MarcarLibreFrame(entrada.NumeroFrame)
		}
	}
}

func CargarEntradasDeMemoria(pid int) (resultados map[int]g.EntradaSwap, err error) {
	resultados = make(map[int]g.EntradaSwap)
	g.MutexProcesosPorPID.Lock()
	proceso := g.ProcesosPorPID[pid]
	g.MutexProcesosPorPID.Unlock()

	if proceso == nil {
		logger.Error("No existe el proceso solicitado para SWAP")
		return resultados, logger.ErrNoInstance
	}

	var entradas []int

	entradas = RecolectarEntradasProcesoSwap(*proceso)

	tamanioPagina := g.MemoryConfig.PagSize
	for i := 0; i < len(entradas); i++ {
		numeroFrame := entradas[i]
		inicio := numeroFrame * tamanioPagina
		fin := inicio + tamanioPagina

		if fin > len(g.MemoriaPrincipal) {
			logger.Error("Acceso fuera de rango al hacer dump del frame %d con PID: %d", numeroFrame, pid)
			fin = len(g.MemoriaPrincipal) - 1
			continue
		}
		vacio := make([]byte, tamanioPagina)
		g.MutexMemoriaPrincipal.Lock()
		datos := g.MemoriaPrincipal[inicio:fin]
		copy(g.MemoriaPrincipal[inicio:fin], vacio)

		entradita := g.EntradaSwap{
			NumeroPagina: numeroFrame,
			Datos:        datos,
			Tamanio:      len(datos),
		}

		g.MutexDump.Lock()
		resultados[numeroFrame] = entradita
		g.MutexDump.Unlock()
	}

	return
}

func CargarEntradasASwap(pid int, entradas map[int]g.EntradaSwap) (err error) {
	tamanioPagina := g.MemoryConfig.PagSize
	err = nil

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

	pos, err := file.Seek(0, io.SeekEnd) // siempre se apunta al final del archivo! (pense que no, alto bobo)
	if err != nil {
		logger.Error("Error al setear el puntero para SWAP: %v", err)
		return err
	}
	var info = &g.SwapProcesoInfo{
		Entradas:         make(map[int]*g.EntradaSwapInfo),
		NumerosDePaginas: make([]int, 0),
	}
	for _, entrada := range entradas {
		_, err = file.Write(entrada.Datos)
		if err != nil {
			logger.Error("Error al escribir el archivo: %v", err)
			return err
		}
		info.Entradas[entrada.NumeroPagina] = &g.EntradaSwapInfo{
			NumeroPagina:   entrada.NumeroPagina,
			Tamanio:        entrada.Tamanio,
			PosicionInicio: int(pos),
		}
		info.NumerosDePaginas = append(info.NumerosDePaginas, entrada.NumeroPagina)

		g.MutexSwapIndex.Lock()
		g.SwapIndex[pid] = info
		g.MutexSwapIndex.Unlock()

		logger.Info("## PID: <%d> - <MEMORIA A SWAP> - Posición en SWAP: <%d> - Tamaño: <%d>",
			pid,
			entrada.NumeroPagina*tamanioPagina,
			entrada.Tamanio,
		)
	}

	return nil
}
