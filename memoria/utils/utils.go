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

func InfomarMetricasProceso(metricasDelProceso globals.MetricasProceso) {

	logger.Info("## Final proceso: %d", metricasDelProceso.AccesosTablasPaginas)
}
