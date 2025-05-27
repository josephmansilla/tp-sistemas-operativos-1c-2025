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

		Utils.MutexEjecutando.Lock()
		algoritmos.ColaEjecutando.Add(proceso)
		Utils.MutexEjecutando.Unlock()

		Utils.MutexReady.Lock()
		algoritmos.ColaReady.Remove(proceso)
		Utils.MutexReady.Unlock()

		cpu := globals.CPUs[cpuID]
		cpu.Ocupada = true
		globals.CPUs[cpuID] = cpu

		comunicacion.EnviarContextoCPU(cpuID, proceso)
	}
}

func DesalojarProceso() {

}

func TerminarEjecucion() {
	logger.Info("Iniciando el planificador de Corto Plazo")
	for {
		// WAIT hasta que CPU finaliza
		//cpuID := <-Utils.ChannelFinishprocess

		var proceso *pcb.PCB

		logger.Info("## (<%d>) Pasa del estado EXECUTE al estado READY", proceso.PID)

		Utils.MutexEjecutando.Lock()
		algoritmos.ColaEjecutando.Remove(proceso)
		Utils.MutexEjecutando.Unlock()

		Utils.MutexReady.Lock()
		algoritmos.ColaReady.Add(proceso)
		Utils.MutexReady.Unlock()

		/*
			//liberar cpu
			cpu := globals.CPUs[cpuID]
			cpu.Ocupada = false
			globals.CPUs[cpuID] = cpu*/
	}
}
