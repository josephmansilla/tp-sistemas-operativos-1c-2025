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

// Mapa global: PID → Lista de instrucciones
var InstruccionesPorPID map[int][]string = make(map[int][]string)

// Cargar instrucción para un PID específico
func CargarInstruccionParaPID(pid int, instruccion string) {
	InstruccionesPorPID[pid] = append(InstruccionesPorPID[pid], instruccion)
	logger.Info("Se cargó una instrucción para PID %d", pid)
}

// Obtener instrucción por PID y PC
func ObtenerInstruccion(pid int, pc int) string {
	instrucciones, existe := InstruccionesPorPID[pid]
	if !existe || pc < 0 || pc >= len(instrucciones) {
		return ""
	}
	return instrucciones[pc]
}
