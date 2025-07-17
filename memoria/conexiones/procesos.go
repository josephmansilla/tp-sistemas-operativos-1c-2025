package conexiones

import (
	"encoding/json"
	"fmt"
	adm "github.com/sisoputnfrba/tp-golang/memoria/administracion"
	g "github.com/sisoputnfrba/tp-golang/memoria/estructuras"
	"github.com/sisoputnfrba/tp-golang/utils/data"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"log"
	"net/http"
	"os"
	"time"
)

func InicializacionProcesoHandler(w http.ResponseWriter, r *http.Request) {
	var mensaje g.InitProceso
	respuesta := g.RespuestaMemoria{
		Exito:   true,
		Mensaje: "Proceso creado correctamente en memoria",
	}
	if err := data.LeerJson(w, r, &mensaje); err != nil {
		logger.Error("Error al leer JSON desde Kernel: %v", err)
		http.Error(w, "Error de parseo de JSON", http.StatusBadRequest)
		return
	}

	pid := mensaje.PID
	tamanioProceso := mensaje.TamanioMemoria
	logger.Info("## PID: <%d> - Proceso Creado - Tamaño: <%d>", pid, tamanioProceso)

	err := adm.InicializarProceso(pid, tamanioProceso, mensaje.Pseudocodigo)
	if err != nil {
		logger.Error("Error al inicializar proceso PID=%d: %v", pid, err)
		respuesta.Exito = false
		respuesta.Mensaje = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	if errEncode := json.NewEncoder(w).Encode(respuesta); errEncode != nil {
		logger.Error("Error al serializar finalizacion de espacio: %v", errEncode)
		http.Error(w, "Error al procesar la respuesta", http.StatusInternalServerError)
		return
	}
}

func FinalizacionProcesoHandler(w http.ResponseWriter, r *http.Request) {

	var mensaje g.ConsultaProceso

	if err := data.LeerJson(w, r, &mensaje); err != nil {
		return
	}

	pid := mensaje.PID

	metricas, err := adm.LiberarProceso(pid)
	if err != nil {
		logger.Error("Hubo un error al eliminar el proceso %v", err)
	}

	logger.Info("## PID: <%d>  - Proceso Destruido - "+
		"Métricas - Acc.T.Pag: <%d>; Inst.Sol.: <%d>; "+
		"SWAP: <%d>; Mem. Prin.: <%d>; Lec.Mem.: <&d>; "+
		"Esc.Mem.: <Esc.Mem.>", pid, metricas.AccesosTablasPaginas,
		metricas.InstruccionesSolicitadas, metricas.BajadasSwap, metricas.SubidasMP,
		metricas.LecturasDeMemoria, metricas.EscriturasDeMemoria)

	respuesta := g.RespuestaMemoria{
		Exito:   true,
		Mensaje: "Proceso eliminado correctamente en memoria",
	}
	if errEncode := json.NewEncoder(w).Encode(respuesta); errEncode != nil {
		logger.Error("Error al serializar finalizacion de espacio: %v", errEncode)
		http.Error(w, "Error al procesar la respuesta", http.StatusInternalServerError)
		return
	}

}

func MemoriaDumpHandler(w http.ResponseWriter, r *http.Request) {
	var dump g.ConsultaDump

	if err := data.LeerJson(w, r, &dump); err != nil {
		return
	}

	dumpFileName := fmt.Sprintf("%s%d-%s.dmp", g.MemoryConfig.DumpPath, dump.PID, dump.TimeStamp)

	dumpFile, err := os.OpenFile(dumpFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		logger.Error("Error al crear archivo de log para <%d-%s>: %v\n", dump.PID, dump.TimeStamp, err)
		os.Exit(1)
	}
	log.SetOutput(dumpFile)
	defer dumpFile.Close()

	logger.Info("## PID: <%d> - Memory Dump solicitado", dump.PID)

	contenido, err := adm.RealizarDumpMemoria(dump.PID)
	if err != nil {
		logger.Error("Error encontrando PID: %v", err)
		http.Error(w, "Error encontrando PID", http.StatusInternalServerError)
		return
	}
	adm.ParsearContenido(dumpFile, dump.PID, contenido)

	logger.Info("## Archivo Dump fue creado con EXITO")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Dump Realizado"))
	if err != nil {
		logger.Error("%% Error al serializar respuesta de realizar dump: %v", err)
		return
	}
}

func SuspensionProcesoHandler(w http.ResponseWriter, r *http.Request) {
	inicio := time.Now()
	retrasoSwap := time.Duration(g.MemoryConfig.SwapDelay) * time.Millisecond
	ignore := 0
	var mensaje g.ConsultaProceso
	if err := data.LeerJson(w, r, &mensaje); err != nil {
		return
	}
	respuesta := g.RespuestaMemoria{
		Exito:   true,
		Mensaje: "Proceso cargado a SWAP",
	}

	g.MutexProcesosPorPID.Lock()
	estaProcesoEnSwap := g.ProcesosPorPID[mensaje.PID].EstaEnSwap
	g.MutexProcesosPorPID.Unlock()

	if estaProcesoEnSwap {
		respuesta = g.RespuestaMemoria{Exito: false, Mensaje: "Ya esta en SWAP"}
		ignore = 1
	}

	if ignore != 1 {
		entradas := adm.RecolectarEntradasParaSwap(mensaje.PID)

		errSwap := adm.CargarEntradasASwap(mensaje.PID, entradas) // REQUIERE ACTUALIZAR ESTRUCTURAS
		if errSwap != nil {
			logger.Error("Error: %v", errSwap)
			http.Error(w, "error: %v", http.StatusConflict)
			respuesta = g.RespuestaMemoria{Exito: false, Mensaje: fmt.Sprintf("Error: %s", errSwap.Error())}
			return
		}

		tiempoTranscurrido := time.Now().Sub(inicio)
		g.CalcularEjecutarSleep(tiempoTranscurrido, retrasoSwap)

	}

	if errEncode := json.NewEncoder(w).Encode(respuesta); errEncode != nil {
		logger.Error("Error al serializar la suspensión del proceso: %v", errEncode)
		http.Error(w, "Error al procesar la respuesta", http.StatusInternalServerError)
		return
	}
}

func DesuspensionProcesoHandler(w http.ResponseWriter, r *http.Request) {
	inicio := time.Now()
	retrasoSwap := time.Duration(g.MemoryConfig.SwapDelay) * time.Millisecond
	ignore := 0
	var mensaje g.ConsultaProceso
	if err := data.LeerJson(w, r, &mensaje); err != nil {
		return
	}
	respuesta := g.RespuestaMemoria{
		Exito:   true,
		Mensaje: "Proceso cargado a Memoria",
	}

	g.MutexProcesosPorPID.Lock()
	estaProcesoEnMemoria := g.ProcesosPorPID[mensaje.PID].EstaEnSwap
	g.MutexProcesosPorPID.Unlock()

	if !estaProcesoEnMemoria {
		respuesta = g.RespuestaMemoria{Exito: false, Mensaje: "Ya esta en Memoria"}
		ignore = 1
	}

	if ignore != 1 {
		entradas, errEntradasSwap := adm.CargarEntradasDesdeSwap(mensaje.PID)
		if errEntradasSwap != nil {
			logger.Error("Error al cargar entradas: %v", errEntradasSwap)
			http.Error(w, "error: %v", http.StatusConflict)
			respuesta = g.RespuestaMemoria{Exito: false, Mensaje: fmt.Sprintf("Error: %s", errEntradasSwap.Error())}
			return
		}

		errEntradasMem := adm.CargarEntradasAMemoria(mensaje.PID, entradas)
		if errEntradasMem != nil {
			logger.Error("Error al cargar entradas: %v", errEntradasMem)
			http.Error(w, "error: %v", http.StatusConflict)
			respuesta = g.RespuestaMemoria{Exito: false, Mensaje: fmt.Sprintf("Error: %s", errEntradasMem.Error())}
			return
		}

		tiempoTranscurrido := time.Now().Sub(inicio)
		g.CalcularEjecutarSleep(tiempoTranscurrido, retrasoSwap)

	}
	if err := json.NewEncoder(w).Encode(respuesta); err != nil {
		logger.Error("Error al serializar la desuspension del proceso: %v", err)
		http.Error(w, "Error al procesar la respuesta", http.StatusInternalServerError)
		return
	}
}
