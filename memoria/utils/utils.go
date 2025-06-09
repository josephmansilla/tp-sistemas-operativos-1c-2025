package utils

import (
	"github.com/sisoputnfrba/tp-golang/utils/logger"
)

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
