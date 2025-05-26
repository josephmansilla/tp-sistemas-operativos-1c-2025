package planificadores

import (
	"github.com/sisoputnfrba/tp-golang/kernel/algoritmos"
	"github.com/sisoputnfrba/tp-golang/kernel/comunicacion"
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/pcb"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
)

func PlanificarCortoPlazo() {
	for {
		var proceso *pcb.PCB

		switch globals.KConfig.SchedulerAlgorithm {
		case "FIFO":
			proceso = algoritmos.ColaReady.First()
		case "SJF":
			proceso = algoritmos.SeleccionarSJF()
		case "SRT":
			proceso = algoritmos.SeleccionarSRT()
		default:
			logger.Error("Algoritmo de planificación desconocido")
			return
		}

		if proceso == nil {
			logger.Info("No hay proceso listo para planificar")
			return
		}

		// Buscar CPU disponible (simplificadamente, agarramos la primera)
		var cpuID string
		for id, cpu := range globals.CPUs {
			if !cpu.Ocupada {
				cpu.Ocupada = true
				globals.CPUs[id] = cpu
				cpuID = id
				break
			}
		}

		if cpuID == "" {
			logger.Info("No hay CPU disponible para ejecutar el proceso <%d>", proceso.PID)
			return
		}

		comunicacion.EnviarContextoCPU(cpuID, proceso)
		proceso.ME[pcb.EstadoExecute]++
		proceso.Estado = pcb.EstadoExecute

		cpu := globals.CPUs[cpuID] // Paso 1: obtener copia
		cpu.Ocupada = true         // Paso 2: modificar la copia
		globals.CPUs[cpuID] = cpu  // Paso 3: volver a guardar la copia modificada

		logger.Info("Proceso <%d> -> EXECUTE en CPU <%s>", proceso.PID, cpuID)
	}
}
