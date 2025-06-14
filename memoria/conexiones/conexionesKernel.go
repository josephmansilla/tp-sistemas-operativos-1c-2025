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
	w.Write([]byte("ESPACIO DEVUELTO"))
}

func InicializacionProcesoHandler(w http.ResponseWriter, r *http.Request) {
	var mensaje g.DatosRespuestaDeKernel

	err := json.NewDecoder(r.Body).Decode(&mensaje)
	if err != nil {
		http.Error(w, "Error leyendo JSON de Kernel\n", http.StatusBadRequest)
		return
	}

	pid := mensaje.PID
	tamanioProceso := mensaje.TamanioMemoria
	adm.InicializarProceso(pid, tamanioProceso, mensaje.Pseudocodigo)

	logger.Info("## PID: <%d> - Proceso Creado - Tamaño: <%d>", pid, tamanioProceso)

	respuesta := g.RespuestaMemoria{
		Exito:   true,
		Mensaje: "Proceso creado correctamente en memoria",
	}
	if err := json.NewEncoder(w).Encode(respuesta); err != nil {
		logger.Error("Error al serializar mock de espacio: %v", err)
	}

	json.NewEncoder(w).Encode(respuesta)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Respuesta devuelta"))
}

func FinalizacionProcesoHandler(w http.ResponseWriter, r *http.Request) {

	var mensaje g.FinalizacionProceso

	err := json.NewDecoder(r.Body).Decode(&mensaje)
	if err != nil {
		http.Error(w, "Error leyendo JSON de Kernel\n", http.StatusBadRequest)
		return
	}
	var metricas g.MetricasProceso
	pid := mensaje.PID

	metricas, err = adm.LiberarMemoriaProceso(pid)
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
		Mensaje: "Proceso creado correctamente en memoria",
	}
	if err := json.NewEncoder(w).Encode(respuesta); err != nil {
		logger.Error("Error al serializar mock de espacio: %v", err)
	}

	json.NewEncoder(w).Encode(respuesta)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Respuesta devuelta"))
}

func LeerEspacioUsuarioHandler(w http.ResponseWriter, r *http.Request) {
	inicio := time.Now()
	retrasoMemoria := time.Duration(g.MemoryConfig.MemoryDelay) * time.Second

	var mensaje g.LecturaProceso
	err := json.NewDecoder(r.Body).Decode(&mensaje)
	if err != nil {
		http.Error(w, "Error leyendo JSON de Kernel\n", http.StatusBadRequest)
		return
	}

	pid := mensaje.PID
	direccionFisica := mensaje.DireccionFisica
	tamanioALeer := mensaje.TamanioARecorrer

	respuesta, err := adm.LeerEspacioMemoria(pid, direccionFisica, tamanioALeer)
	if err != nil {
		// TODO:::::: -------------------------------------------------------------
	}

	logger.Info("## PID: <%d>  - <Lectura> - Dir. Física: <%d> - Tamaño: <%d>", pid, direccionFisica, tamanioALeer)

	time.Sleep(time.Duration(g.MemoryConfig.MemoryDelay) * time.Second)

	tiempoTranscurrido := time.Now().Sub(inicio)
	g.CalcularEjecutarSleep(tiempoTranscurrido, retrasoMemoria)

	if err := json.NewEncoder(w).Encode(respuesta); err != nil {
		logger.Error("Error al serializar mock de espacio: %v", err)
	}

	logger.Info("## Lectura en espacio de memoria Éxitosa")

	json.NewEncoder(w).Encode(respuesta)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Respuesta devuelta"))
}

func EscribirEspacioUsuarioHandler(w http.ResponseWriter, r *http.Request) {
	inicio := time.Now()
	retrasoMemoria := time.Duration(g.MemoryConfig.MemoryDelay) * time.Second

	var mensaje g.EscrituraProceso
	err := json.NewDecoder(r.Body).Decode(&mensaje)
	if err != nil {
		http.Error(w, "Error leyendo JSON de Kernel\n", http.StatusBadRequest)
		return
	}

	pid := mensaje.PID
	direccionFisica := mensaje.DireccionFisica
	tamanioALeer := mensaje.TamanioARecorrer
	datos := mensaje.DatosAEscribir

	respuesta, err := adm.EscribirEspacioMemoria(pid, direccionFisica, tamanioALeer, datos)
	if err != nil {
		// TODO : ======================================
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
	w.Write([]byte("Respuesta devuelta"))
}

func MemoriaDumpHandler(w http.ResponseWriter, r *http.Request) {
	var dump g.DatosParaDump

	if err := data.LeerJson(w, r, &dump); err != nil {
		logger.Error("Error al recibir JSON: %v", err)
		http.Error(w, "Error procesando datos del Kernel", http.StatusInternalServerError)
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

	contenido := adm.RealizarDumpMemoria(dump.PID)
	// TODO: verificacion esta vacio
	g.ParsearContenido(dumpFile, contenido)

	logger.Info("## Archivo Dump fue creado con EXITO")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Dump Realizado"))
}
