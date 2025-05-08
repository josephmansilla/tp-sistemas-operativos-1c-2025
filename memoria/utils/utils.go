package utils

import (
	"github.com/sisoputnfrba/tp-golang/utils/logger"
)

var Instrucciones []string = []string{}

// función auxiliar para cargar el slice de instrucciones
func CargarListaDeInstrucciones(str string) {
	Instrucciones = append(Instrucciones, str)
	logger.Info("Se cargó una instrucción al Slice")
}

func InfomarMetricasProceso(metricasDelProceso MetricaProceso) {
	logger.Info("## Final proceso: ")
}

func IncrementarEscrituraMemoria() {}
func IncrementarLecturaMemoria()   {}
