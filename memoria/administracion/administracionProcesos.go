package administracion

import (
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"net/http"
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

}

func TieneTamanioNecesario(tamanioProceso int) bool {
	var framesNecesarios = float64(tamanioProceso) / float64(globals.TamanioMaximoFrame)
	return framesNecesarios <= float64(globals.CantidadFramesLibres)
}

func LecturaPseudocodigo(archivoPseudocodigo string) []byte {

	string := archivoPseudocodigo
	stringEnBytes := []byte(string)

	return stringEnBytes
}

// ------------------------------------------------------------------
// ----------- FORMA PARTE DE LA MODIFICACIÓN DE PROCESOS -----------
// ------------------------------------------------------------------

func InicializacionProceso(w http.ResponseWriter, r *http.Request) {
	// TODO: VERIFICAR EL TAMAÑO NECESARIO

	// TODO: CREAR ESTRUCTURAS ADMINISTRATIVAS NECESARIAS

	// TODO: RESPONDER CON EL NUMERO DE PAGINA DE 1ER NIVEL DEL PROCESO
	logger.Info("## PID: <%d>  - Proceso Creado - Tamaño: <%d>")
}

func FinalizacionProceso(w http.ResponseWriter, r *http.Request) {
	//toDO

	logger.Info("## PID: <PID>  - Proceso Destruido - Métricas - Acc.T.Pag: <ATP>; Inst.Sol.: <Inst.Sol>; SWAP: <SWAP>; Mem. Prin.: <Mem.Prin.>; Lec.Mem.: <Lec.Mem.>; Esc.Mem.: <Esc.Mem.>")
}

func SuspensionProceso(w http.ResponseWriter, r *http.Request) {
	// TODO: NO ES NECESARIO EL SWAPEO DE TABLAS DE PAGINAS

	// TODO: SE LIBERA EN MEMORIA
	// TODO: SE ESCRIBE EN SWAP LA INFO NECESARIA

}

func DesSuspensionProceso(w http.ResponseWriter, r *http.Request) {
	// TODO: VERIFICAR EL TAMAÑO NECESARIO

	// TODO: LEER EL CONTENIDO DEL SWAP, ESCRIBIERLO EN EL FRAME ASIGNADO
	// TODO: LIBERAR ESPACIO EN SWAP
	// TODO: ACTUALIZAR ESTRUCTURAS NECESARIAS

	// TODO: RETORNAR OK
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
	funcMetrica(&proceso.Metricas)
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
