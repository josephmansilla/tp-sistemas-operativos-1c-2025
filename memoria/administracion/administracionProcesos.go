package administracion

import (
	g "github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
)

func InicializarProceso(pid int, tamanioProceso int, archivoPseudocodigo string) {
	if !TieneTamanioNecesario(tamanioProceso) {
		// TODO
		logger.Error("No hay memoria")
	}
	nuevoProceso := g.Proceso{
		PID:       pid,
		TablaRaiz: InicializarTablaRaiz(),
		Metricas:  InicializarMetricas(),
	}
	pseudo, err := LecturaPseudocodigo(archivoPseudocodigo)
	if err != nil {
		logger.Error("Error al leer el pseudocodigo: %v", err)
	}

	err = AsignarDatosAPaginacion(&nuevoProceso, pseudo)
	if err != nil {
		logger.Error("Error al asignarDatosAPaginacion %v", err)
	}
	OcuparProcesoEnVectorMapeable(pid, nuevoProceso)
}

func OcuparProcesoEnVectorMapeable(pid int, nuevoProceso g.Proceso) {
	g.MutexProcesosPorPID.Lock()
	g.ProcesosPorPID[pid] = &nuevoProceso
	g.MutexProcesosPorPID.Unlock()
}

func CargarEntradaMemoria(numeroFrame int, pid int, datosEnBytes []byte) {
	direccionFisica := numeroFrame * g.MemoryConfig.PagSize
	g.MutexMemoriaPrincipal.Lock()
	for indice := 0; indice < len(datosEnBytes); indice++ {
		g.MemoriaPrincipal[direccionFisica] = datosEnBytes[indice]
	}
	g.MutexMemoriaPrincipal.Unlock()
}

// METRICAS PROCESOS

func InicializarMetricas() (metricas g.MetricasProceso) {
	metricas = g.MetricasProceso{
		AccesosTablasPaginas:     0,
		InstruccionesSolicitadas: 0,
		BajadasSwap:              0,
		SubidasMP:                0,
		LecturasDeMemoria:        0,
		EscriturasDeMemoria:      0,
	}
	return
}

func IncrementarMetrica(proceso *g.Proceso, funcMetrica g.OperacionMetrica) {
	g.MutexMetrica[proceso.PID].Lock()
	funcMetrica(&proceso.Metricas)
	g.MutexMetrica[proceso.PID].Unlock()
}

func InformarMetricasProceso(metricasDelProceso g.MetricasProceso) {

	logger.Info("## AccesosTablasPaginas: %d", metricasDelProceso.AccesosTablasPaginas)
	logger.Info("## InstruccionesSolicitadas: %d", metricasDelProceso.InstruccionesSolicitadas)
	logger.Info("## BajadasSwap: %d", metricasDelProceso.BajadasSwap)
	logger.Info("## SubidasMP: %d", metricasDelProceso.SubidasMP)
	logger.Info("## LecturasDeMemoria: %d", metricasDelProceso.LecturasDeMemoria)
	logger.Info("## EscriturasDeMemoria: %d", metricasDelProceso.EscriturasDeMemoria)

} // TODO: borrar

func IncrementarAccesosTablasPaginas(metrica *g.MetricasProceso) {
	metrica.AccesosTablasPaginas++
}
func IncrementarInstruccionesSolicitadas(metrica *g.MetricasProceso) {
	metrica.InstruccionesSolicitadas++
}
func IncrementarBajadasSwap(metrica *g.MetricasProceso) {
	metrica.BajadasSwap++
}
func IncrementarSubidasMP(metrica *g.MetricasProceso) {
	metrica.SubidasMP++
}
func IncrementarLecturaDeMemoria(metrica *g.MetricasProceso) {
	metrica.LecturasDeMemoria++
}
func IncrementarEscrituraDeMemoria(metrica *g.MetricasProceso) {
	metrica.EscriturasDeMemoria++
}
