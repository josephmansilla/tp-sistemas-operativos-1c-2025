package administracion

import (
	"encoding/json"
	"fmt"
	g "github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/data"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
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
			LiberarEntradaPagina(entrada.NumeroFrame)

		}
	}
}

func CargarEntradasDeMemoria(pid int) (resultados map[int]g.EntradaSwap, err error) {
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
		copy(g.MemoriaPrincipal[inicio:fin], vacio) // TODO: preguntar o cambiar impl de leer entrada entera
		g.MutexMemoriaPrincipal.Unlock()

		entradita := g.EntradaSwap{
			NumeroFrame: numeroFrame,
			Datos:       datos,
			Tamanio:     len(datos),
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
	defer file.Close()

	pos, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}
	var info = &g.SwapProcesoInfo{
		Entradas:     make(map[int]*g.EntradaSwapInfo),
		NumerosFrame: make([]int, 0),
	}
	for _, entrada := range entradas {
		_, err = file.Write(entrada.Datos)
		if err != nil {
			return err
		}
		info.Entradas[entrada.NumeroFrame] = &g.EntradaSwapInfo{
			NumeroFrame:    entrada.NumeroFrame,
			Tamanio:        entrada.Tamanio,
			PosicionInicio: int(pos),
		}
		info.NumerosFrame = append(info.NumerosFrame, entrada.NumeroFrame)

		g.MutexSwapIndex.Lock()
		g.SwapIndex[pid] = info
		g.MutexSwapIndex.Unlock()

		logger.Info("## PID: <%d> - <Escritura> - Dir. Física: <%d> - Tamaño: <%d>",
			pid,
			entrada.NumeroFrame*tamanioPagina,
			entrada.Tamanio,
		)
	}

	return nil
}

func SuspensionProcesoHandler(w http.ResponseWriter, r *http.Request) {
	inicio := time.Now()
	retrasoSwap := time.Duration(g.MemoryConfig.SwapDelay) * time.Second
	ignore := 0
	var mensaje g.PedidoKernel
	if err := data.LeerJson(w, r, &mensaje); err != nil {
		return
	}
	respuesta := g.RespuestaMemoria{
		Exito:   true,
		Mensaje: "Proceso cargado a SWAP",
	}
	g.MutexSwapBool.Lock()
	estaProcesoEnSwap := g.EstaEnSwap[mensaje.PID]
	g.MutexSwapBool.Unlock()
	if estaProcesoEnSwap {
		respuesta = g.RespuestaMemoria{Exito: false, Mensaje: "Ya esta en SWAP"}
		ignore = 1
	}

	if ignore != 1 {
		entradas, errEntradas := CargarEntradasDeMemoria(mensaje.PID)
		if errEntradas != nil {
			logger.Error("Error: %v", errEntradas)
			http.Error(w, "error: %v", http.StatusNoContent)
			respuesta = g.RespuestaMemoria{Exito: false, Mensaje: fmt.Sprintf("Errror: %s", errEntradas.Error())}
			return
		}
		errSwap := CargarEntradasASwap(mensaje.PID, entradas) // REQUIERE ACTUALIZAR ESTRUCTURAS
		if errSwap != nil {
			logger.Error("Error: %v", errSwap)
			http.Error(w, "error: %v", http.StatusConflict)
			respuesta = g.RespuestaMemoria{Exito: false, Mensaje: fmt.Sprintf("Error: %s", errEntradas.Error())}
			return
		}

		g.MutexProcesosPorPID.Lock()
		proceso := g.ProcesosPorPID[mensaje.PID]
		g.MutexProcesosPorPID.Unlock()
		IncrementarMetrica(proceso, 1, IncrementarBajadasSwap)

		g.EstaEnSwap[mensaje.PID] = true

		tiempoTranscurrido := time.Now().Sub(inicio)
		g.CalcularEjecutarSleep(tiempoTranscurrido, retrasoSwap)

	}

	if err := json.NewEncoder(w).Encode(respuesta); err != nil {
		logger.Error("Error al serializar mock de espacio: %v", err)
	}
	w.WriteHeader(http.StatusOK)
	//w.Write([]byte("Respuesta devuelta"))
}

// ==========================================================================

func CargarEntradasDeSwap(pid int) (entradas map[int]g.EntradaSwap, err error) {

	g.MutexSwapIndex.Lock()
	info, existe := g.SwapIndex[pid]
	if !existe {
		return nil, logger.ErrNoInstance
	}
	g.MutexSwapIndex.Unlock()

	file, err := os.Open(g.MemoryConfig.SwapfilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	entradas = make(map[int]g.EntradaSwap, len(info.NumerosFrame))

	for i, entrada := range info.Entradas {
		_, errArch := file.Seek(int64(entrada.PosicionInicio), io.SeekStart)
		if errArch != nil {
			return nil, err
		}
		datos := make([]byte, 0)
		enviarEntrada := g.EntradaSwap{
			NumeroFrame: info.NumerosFrame[i],
			Datos:       datos,
			Tamanio:     entrada.Tamanio,
		}
		entradas[info.NumerosFrame[i]] = enviarEntrada
	}

	return entradas, nil
}

func CargarEntradasAMemoria(pid int, entradas map[int]g.EntradaSwap) (err error) {
	tamanioPagina := g.MemoryConfig.PagSize
	err = nil

	for _, entrada := range entradas {
		dirFisica := entrada.NumeroFrame * tamanioPagina
		rta := EscribirEspacioEntrada(pid, dirFisica, string(entrada.Datos))
		if rta.Exito != nil {
			logger.Error("Error: %v", rta.Exito)
			return rta.Exito
		}

		logger.Info("## PID: <%d> - <Lectura> - Dir. Física: <%d> - Tamaño: <%d>",
			pid,
			dirFisica,
			entrada.Tamanio,
		)
	}
	return
}

func DesuspensionProcesoHandler(w http.ResponseWriter, r *http.Request) {
	inicio := time.Now()
	retrasoSwap := time.Duration(g.MemoryConfig.SwapDelay) * time.Second
	ignore := 0
	var mensaje g.DesuspensionProceso
	if err := data.LeerJson(w, r, &mensaje); err != nil {
		return
	}
	respuesta := g.RespuestaMemoria{
		Exito:   true,
		Mensaje: "Proceso cargado a Memoria",
	}

	g.MutexSwapBool.Lock()
	estaProcesoEnMemoria := !g.EstaEnSwap[mensaje.PID]
	g.MutexSwapBool.Unlock()

	if estaProcesoEnMemoria {
		respuesta = g.RespuestaMemoria{Exito: false, Mensaje: "Ya esta en Memoria"}
		ignore = 1
	}

	if ignore != 1 {
		entradas, errEntradasSwap := CargarEntradasDeSwap(mensaje.PID)
		if errEntradasSwap != nil {
			logger.Error("Error al cargar entradas: %v", errEntradasSwap)
			http.Error(w, "error: %v", http.StatusConflict)
			respuesta = g.RespuestaMemoria{Exito: false, Mensaje: fmt.Sprintf("Error: %s", errEntradasSwap.Error())}
			return
		}

		errEntradasMem := CargarEntradasAMemoria(mensaje.PID, entradas)
		if errEntradasMem != nil {
			logger.Error("Error al cargar entradas: %v", errEntradasMem)
			http.Error(w, "error: %v", http.StatusConflict)
			respuesta = g.RespuestaMemoria{Exito: false, Mensaje: fmt.Sprintf("Error: %s", errEntradasMem.Error())}
			return
		}

		g.MutexProcesosPorPID.Lock()
		proceso := g.ProcesosPorPID[mensaje.PID]
		g.MutexProcesosPorPID.Unlock()
		IncrementarMetrica(proceso, 1, IncrementarSubidasMP)

		g.EstaEnSwap[mensaje.PID] = false

		time.Sleep(time.Duration(g.MemoryConfig.SwapDelay) * time.Second)

		tiempoTranscurrido := time.Now().Sub(inicio)
		g.CalcularEjecutarSleep(tiempoTranscurrido, retrasoSwap)

	}
	if err := json.NewEncoder(w).Encode(respuesta); err != nil {
		logger.Error("Error al serializar mock de espacio: %v", err)
	}
	w.WriteHeader(http.StatusOK)
	//w.Write([]byte("Respuesta devuelta"))
}
