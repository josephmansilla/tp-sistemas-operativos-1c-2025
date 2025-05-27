package planificadores

import (
	"github.com/sisoputnfrba/tp-golang/kernel/Utils"
	"github.com/sisoputnfrba/tp-golang/kernel/algoritmos"
	"github.com/sisoputnfrba/tp-golang/kernel/comunicacion"
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/pcb"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
)

func PlanificarCortoPlazo() {
	logger.Info("Iniciando el planificador de Corto Plazo")
	for {
		// WAIT hasta que llegue un proceso a READY
		pid := <-Utils.NotificarProcesoReady
		logger.Info("## (<%d>) Llega a Corto Plazo", pid)

		var proceso *pcb.PCB

		switch globals.KConfig.SchedulerAlgorithm {
		case "FIFO":
			proceso = algoritmos.ColaReady.First()
			algoritmos.ColaReady.Remove(proceso)
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
			continue
		}

		// Buscar CPU disponible
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
			continue
		}

		logger.Info("Proceso <%d> -> EXECUTE en CPU <%s>", proceso.PID, cpuID)
		proceso.ME[pcb.EstadoExecute]++
		proceso.Estado = pcb.EstadoExecute
		algoritmos.ColaEjecutando.Add(proceso)

		cpu := globals.CPUs[cpuID]
		cpu.Ocupada = true
		globals.CPUs[cpuID] = cpu

		comunicacion.EnviarContextoCPU(cpuID, proceso)
	}
}
