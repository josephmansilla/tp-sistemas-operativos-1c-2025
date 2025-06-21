package conexiones

import (
	"encoding/json"
	"fmt"
	adm "github.com/sisoputnfrba/tp-golang/memoria/administracion"
	g "github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/data"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"log"
	"net/http"
	"os"
	"time"
)

func ObtenerEspacioLibreHandler(w http.ResponseWriter, r *http.Request) {
	g.MutexCantidadFramesLibres.Lock()
	cantFramesLibres := g.CantidadFramesLibres
	g.MutexCantidadFramesLibres.Unlock()

	espacioLibre := cantFramesLibres * g.MemoryConfig.PagSize

	respuesta := g.RespuestaEspacioLibre{EspacioLibre: espacioLibre}

	logger.Info("## Espacio libre devuelto - Tamaño: <%d>", respuesta.EspacioLibre)

	if err := json.NewEncoder(w).Encode(respuesta); err != nil {
		logger.Error("Error al serializar mock de espacio: %v", err)
	}
	json.NewEncoder(w).Encode(respuesta)
	w.WriteHeader(http.StatusOK)
	//w.Write([]byte("ESPACIO DEVUELTO"))
}

func InicializacionProcesoHandler(w http.ResponseWriter, r *http.Request) {
	var mensaje g.DatosRespuestaDeKernel
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

	w.Header().Set("Content-Type", "application/json")
	if errEncode := json.NewEncoder(w).Encode(respuesta); errEncode != nil {
		logger.Error("Error al codificar respuesta JSON: %v", errEncode)
	}
}

func FinalizacionProcesoHandler(w http.ResponseWriter, r *http.Request) {

	var mensaje g.FinalizacionProceso

	if err := data.LeerJson(w, r, &mensaje); err != nil {
		return
	}

	pid := mensaje.PID

	metricas, err := adm.LiberarMemoriaProceso(pid)
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
		logger.Error("Error al serializar mock de espacio: %v", errEncode)
		return
	}

	json.NewEncoder(w).Encode(respuesta)
	w.WriteHeader(http.StatusOK)
	//w.Write([]byte("Respuesta devuelta"))
}

func LeerEspacioUsuarioHandler(w http.ResponseWriter, r *http.Request) {
	inicio := time.Now()
	retrasoMemoria := time.Duration(g.MemoryConfig.MemoryDelay) * time.Second

	var mensaje g.LecturaProceso
	err := data.LeerJson(w, r, &mensaje)
	if err != nil {
		return
	}

	pid := mensaje.PID
	direccionFisica := mensaje.DireccionFisica
	tamanioALeer := mensaje.TamanioARecorrer

	respuesta, err := adm.LeerEspacioMemoria(pid, direccionFisica, tamanioALeer)
	if err != nil {
		logger.Error("Error: %v", err)
		http.Error(w, "Error al Leer espacio de Memoria \n", http.StatusInternalServerError)
	}

	logger.Info("## PID: <%d>  - <Lectura> - Dir. Física: <%d> - Tamaño: <%d>", pid, direccionFisica, tamanioALeer)

	time.Sleep(time.Duration(g.MemoryConfig.MemoryDelay) * time.Second)

	tiempoTranscurrido := time.Now().Sub(inicio)
	g.CalcularEjecutarSleep(tiempoTranscurrido, retrasoMemoria)

	if err := json.NewEncoder(w).Encode(respuesta); err != nil {
		logger.Error("Error al serializar mock de espacio: %v", err)
		return
	}

	logger.Info("## Lectura en espacio de memoria Éxitosa")

	json.NewEncoder(w).Encode(respuesta)
	w.WriteHeader(http.StatusOK)
	//w.Write([]byte("Respuesta devuelta"))
}

func EscribirEspacioUsuarioHandler(w http.ResponseWriter, r *http.Request) {
	inicio := time.Now()
	retrasoMemoria := time.Duration(g.MemoryConfig.MemoryDelay) * time.Second

	var mensaje g.EscrituraProceso
	if err := data.LeerJson(w, r, &mensaje); err != nil {
		return
	}

	pid := mensaje.PID
	direccionFisica := mensaje.DireccionFisica
	tamanioALeer := mensaje.TamanioARecorrer
	datos := mensaje.DatosAEscribir

	respuesta, err := adm.EscribirEspacioMemoria(pid, direccionFisica, tamanioALeer, datos)
	if err != nil {
		logger.Error("Error: %v", err)
		http.Error(w, "Error al Leer espacio de Memoria \n", http.StatusInternalServerError)
	}

	logger.Info("## PID: <%d> - <Escritura> - Dir. Física: <%d> - Tamaño: <%d>", pid, direccionFisica, tamanioALeer)

	time.Sleep(time.Duration(g.MemoryConfig.MemoryDelay) * time.Second)

	tiempoTranscurrido := time.Now().Sub(inicio)
	g.CalcularEjecutarSleep(tiempoTranscurrido, retrasoMemoria)

	if err := json.NewEncoder(w).Encode(respuesta); err != nil {
		logger.Error("Error al serializar mock de espacio: %v", err)
	}

	logger.Info("## Escritura en espacio de memoria Éxitosa")

	json.NewEncoder(w).Encode(respuesta)
	w.WriteHeader(http.StatusOK)
	//w.Write([]byte("Respuesta devuelta"))
}

func MemoriaDumpHandler(w http.ResponseWriter, r *http.Request) {
	var dump g.DatosParaDump

	if err := data.LeerJson(w, r, &dump); err != nil {
		return
	}

	dumpFileName := fmt.Sprintf("%s/<%d>-<%s>.dmp", g.MemoryConfig.DumpPath, dump.PID, dump.TimeStamp)
	logger.Info("## Se creo el file: %d ", dumpFileName)
	dumpFile, err := os.OpenFile(dumpFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Printf("Error al crear archivo de log para <%d-%s>: %v\n", dump.PID, dump.TimeStamp, err)
		os.Exit(1)
	}
	log.SetOutput(dumpFile)
	defer dumpFile.Close()

	logger.Info("## PID: <%d>  - Memory Dump solicitado", dump.PID)

	contenido, err := adm.RealizarDumpMemoria(dump.PID)
	if err != nil {
		logger.Error("Error encontrando PID: %v", err)
		http.Error(w, "Error encontrando PID", http.StatusInternalServerError)
		return
	}
	g.ParsearContenido(dumpFile, dump.PID, contenido)

	logger.Info("## Archivo Dump fue creado con EXITO")
	w.WriteHeader(http.StatusOK)
	//w.Write([]byte("Dump Realizado"))
}
