package administracion

import (
	"encoding/json"
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	g "github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/data"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"net/http"
	"sync"
	"time"
)

// Para pasar a suspendido:
// 1) Transferir los datos de cada entrada en bytes hacia el archivo, dejar en bytes.
// 2) Marcar cantidad frames como libres y no presente
// 3) Eliminar el contenido de los frames (conceptualmente mal, pero es posible que me ataje algún error de mis funciones)
// 4) Actualizar estructuras necesarias

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

	for _, entrada := range entradas {

		logger.Info("## PID: <%d> - <Escritura> - Dir. Física: <%d> - Tamaño: <%d>",
			pid,
			entrada.NumeroFrame*tamanioPagina,
			entrada.Tamanio,
		)
	}

	return
}

func SuspensionProcesoHandler(w http.ResponseWriter, r *http.Request) {
	inicio := time.Now()
	retrasoSwap := time.Duration(g.MemoryConfig.SwapDelay) * time.Second

	var mensaje g.SuspensionProceso
	if err := data.LeerJson(w, r, &mensaje); err != nil {
		return
	}

	entradas, errEntradas := CargarEntradasDeMemoria(mensaje.PID)
	if errEntradas != nil {
		logger.Error("Error: %v", errEntradas)
		http.Error(w, "error: %v", http.StatusNoContent)
		return
	}
	errSwap := CargarEntradasASwap(mensaje.PID, entradas) // REQUIERE ACTUALIZAR ESTRUCTURAS
	if errSwap != nil {
		logger.Error("Error: %v", errSwap)
		http.Error(w, "error: %v", http.StatusNoContent)
		return
	}

	g.MutexProcesosPorPID.Lock()
	proceso := g.ProcesosPorPID[mensaje.PID]
	g.MutexProcesosPorPID.Unlock()
	IncrementarMetrica(proceso, 1, IncrementarBajadasSwap)

	respuesta := g.RespuestaMemoria{
		Exito:   true,
		Mensaje: "Proceso cargado a SWAP",
	}

	tiempoTranscurrido := time.Now().Sub(inicio)
	g.CalcularEjecutarSleep(tiempoTranscurrido, retrasoSwap)

	if err := json.NewEncoder(w).Encode(respuesta); err != nil {
		logger.Error("Error al serializar mock de espacio: %v", err)
	}

	json.NewEncoder(w).Encode(respuesta)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Respuesta devuelta"))
}

// Para sacar de suspendido:
// 1) Verificar el tamanio del proceso y si entra en Memoria
// 2) Leer cada entrada de swap, ponerlo en memoria y marcar como presente
// 3) Liberar de SWAP el espacio
// 4) Actualizar estructuras necesarias
// 5) Retonar confirmación éxitosa o fallida

func CargarEntradasDeSwap(pid int) (entradas map[int]g.EntradaSwap) {

	return
}

func CargarEntradasAMemoria(pid int, entradas map[int]g.EntradaSwap) {
	tamanioPagina := g.MemoryConfig.PagSize

	for _, entrada := range entradas {

		logger.Info("## PID: <%d> - <Lectura> - Dir. Física: <%d> - Tamaño: <%d>",
			pid,
			entrada.NumeroFrame*tamanioPagina,
			entrada.Tamanio,
		)
	}
}

func DesuspensionProcesoHandler(w http.ResponseWriter, r *http.Request) {
	inicio := time.Now()
	retrasoSwap := time.Duration(g.MemoryConfig.SwapDelay) * time.Second

	var mensaje g.DesuspensionProceso
	if err := data.LeerJson(w, r, &mensaje); err != nil {
		return
	}

	entradas := CargarEntradasDeSwap(mensaje.PID)
	CargarEntradasAMemoria

	time.Sleep(time.Duration(g.MemoryConfig.SwapDelay) * time.Second)

	tiempoTranscurrido := time.Now().Sub(inicio)
	g.CalcularEjecutarSleep(tiempoTranscurrido, retrasoSwap)

	respuesta := g.ExitoDesuspensionProceso{}

	if err := json.NewEncoder(w).Encode(respuesta); err != nil {
		logger.Error("Error al serializar mock de espacio: %v", err)
	}

	json.NewEncoder(w).Encode(respuesta)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Respuesta devuelta"))
}
