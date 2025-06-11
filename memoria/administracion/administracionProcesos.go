package administracion

import (
	"encoding/json"
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"net/http"
	"sync"
)

func InicializarProceso(pid int, tamanioProceso int, archivoPseudocodigo string) {
	if !TieneTamanioNecesario(tamanioProceso) {
		// TODO
		logger.Error("No hay memoria")
	}
	nuevoProceso := globals.Proceso{
		PID:       pid,
		TablaRaiz: InicializarTablaRaiz(),
		Metricas:  InicializarMetricas(),
	}
	err := AsignarDatosAPaginacion(&nuevoProceso, LecturaPseudocodigo(archivoPseudocodigo))
	if err != nil {
		logger.Error("Error al asignarDatosAPaginacion %v", err)
	}
	OcuparProcesoEnVectorMapeable(pid, nuevoProceso)
}

func OcuparProcesoEnVectorMapeable(pid int, nuevoProceso globals.Proceso) {
	globals.MutexProcesosPorPID.Lock()
	globals.ProcesosPorPID[pid] = &nuevoProceso
	globals.MutexProcesosPorPID.Unlock()
}

func CargarEntradaMemoria(numeroFrame int, pid int, datosEnBytes []byte) {
	direccionFisica := numeroFrame * globals.MemoryConfig.PagSize
	globals.MutexMemoriaPrincipal.Lock()
	for indice := 0; indice < len(datosEnBytes); indice++ {
		globals.MemoriaPrincipal[direccionFisica] = datosEnBytes[indice]
	}
	globals.MutexMemoriaPrincipal.Unlock()
}

func InicializacionProceso(w http.ResponseWriter, r *http.Request) {
	var mensaje globals.DatosRespuestaDeKernel

	err := json.NewDecoder(r.Body).Decode(&mensaje)
	if err != nil {
		http.Error(w, "Error leyendo JSON de Kernel\n", http.StatusBadRequest)
		return
	}

	pid := mensaje.PID
	tamanioProceso := mensaje.TamanioMemoria
	InicializarProceso(pid, tamanioProceso, mensaje.Pseudocodigo)

	logger.Info("## PID: <%d> - Proceso Creado - Tamaño: <%d>", pid, tamanioProceso)

	respuesta := globals.RespuestaMemoria{
		Exito:   true,
		Mensaje: "Proceso creado correctamente en memoria",
	}
	if err := json.NewEncoder(w).Encode(respuesta); err != nil {
		logger.Error("Error al serializar mock de espacio: %v", err)
	}

	json.NewEncoder(w).Encode(respuesta)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Respuesta devuelta"))
} //TODO:
// TODO: RESPONDER CON EL NUMERO DE PAGINA DE 1ER NIVEL DEL PROCESO???

func FinalizacionProceso(w http.ResponseWriter, r *http.Request) {

	var mensaje globals.DatosFinalizacionProceso

	err := json.NewDecoder(r.Body).Decode(&mensaje)
	if err != nil {
		http.Error(w, "Error leyendo JSON de Kernel\n", http.StatusBadRequest)
		return
	}

	pid := mensaje.PID
	LiberarMemoria(pid) // TODO: pendiente
	globals.MutexProcesosPorPID.Lock()
	proceso := globals.ProcesosPorPID[pid]
	delete(globals.ProcesosPorPID, pid)
	globals.MutexProcesosPorPID.Unlock()

	metricas := proceso.Metricas

	// TODO: revisar logger
	logger.Info("## PID: <%d>  - Proceso Destruido - "+
		"Métricas - Acc.T.Pag: <%d>; Inst.Sol.: <%d>; "+
		"SWAP: <ñññ%d>; Mem. Prin.: <ñññ>; Lec.Mem.: <&d>; "+
		"Esc.Mem.: <Esc.Mem.>", pid, metricas.AccesosTablasPaginas,
		metricas.InstruccionesSolicitadas, metricas.BajadasSwap+metricas.SubidasMP,
		metricas.LecturasDeMemoria, metricas.EscriturasDeMemoria)

	respuesta := globals.RespuestaMemoria{
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

// METRICAS PROCESOS

func InicializarMetricas() globals.MetricasProceso {
	metricas := globals.MetricasProceso{
		AccesosTablasPaginas:     0,
		InstruccionesSolicitadas: 0,
		BajadasSwap:              0,
		SubidasMP:                0,
		LecturasDeMemoria:        0,
		EscriturasDeMemoria:      0,
	}
	return metricas
}

func IncrementarMetrica(proceso *globals.Proceso, funcMetrica globals.OperacionMetrica) {
	var mutexMetrica sync.Mutex

	mutexMetrica.Lock()
	funcMetrica(&proceso.Metricas)
	mutexMetrica.Unlock()
}

func InformarMetricasProceso(metricasDelProceso globals.MetricasProceso) {

	logger.Info("## AccesosTablasPaginas: %d", metricasDelProceso.AccesosTablasPaginas)
	logger.Info("## InstruccionesSolicitadas: %d", metricasDelProceso.InstruccionesSolicitadas)
	logger.Info("## BajadasSwap: %d", metricasDelProceso.BajadasSwap)
	logger.Info("## SubidasMP: %d", metricasDelProceso.SubidasMP)
	logger.Info("## LecturasDeMemoria: %d", metricasDelProceso.LecturasDeMemoria)
	logger.Info("## EscriturasDeMemoria: %d", metricasDelProceso.EscriturasDeMemoria)

}

func IncrementarAccesosTablasPaginas(metrica *globals.MetricasProceso) {
	metrica.AccesosTablasPaginas++
}
func IncrementarInstruccionesSolicitadas(metrica *globals.MetricasProceso) {
	metrica.InstruccionesSolicitadas++
}
func IncrementarBajadasSwap(metrica *globals.MetricasProceso) {
	metrica.BajadasSwap++
}
func IncrementarSubidasMP(metrica *globals.MetricasProceso) {
	metrica.SubidasMP++
}
func IncrementarLecturaDeMemoria(metrica *globals.MetricasProceso) {
	metrica.LecturasDeMemoria++
}
func IncrementarEscrituraDeMemoria(metrica *globals.MetricasProceso) {
	metrica.EscriturasDeMemoria++
}
