package utils

import (
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
)

var Instrucciones []string = []string{}

// función auxiliar para cargar el slice de instrucciones
func CargarListaDeInstrucciones(str string) {
	Instrucciones = append(Instrucciones, str)
	logger.Info("Se cargó una instrucción al Slice")
}

func InfomarMetricasProceso(metricasDelProceso globals.MetricaProceso) {

	logger.Info("## Final proceso: %d", metricasDelProceso.CantAccesosTablasPaginas)
}

func IncrementarEscrituraMemoria() {}
func IncrementarLecturaMemoria()   {}
