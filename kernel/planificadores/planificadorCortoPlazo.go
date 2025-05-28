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
	logger.Info("Iniciando el Planificador de Corto Plazo")
	go DespacharProceso()
	go BloquearProceso()
	go FinDeIO()
}

func DespacharProceso() {
	for {
		// WAIT hasta que llegue un proceso a READY o nueva CPU por SYSCALL DE EXTI O I/O
		<-Utils.NotificarDespachador
		logger.Info("Arranca Despachador")

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

func liberarCPU(cpuID string) {

	cpu := globals.CPUs[cpuID]
	cpu.Ocupada = false
	globals.CPUs[cpuID] = cpu

	logger.Info("CPU <%s> libre", cpuID)

	Utils.NotificarDespachador <- 1
}

func DesalojarProceso() {

}

func BloquearProceso() {
	for {
		//WAIT mensaje de IO (bloqueante)
		msg := <-Utils.NotificarComienzoIO

		var pid = msg.PID
		var cpuID = msg.CpuID

		//BUSCAR EN PCB y SACAR DE EXECUTE
		var proceso *pcb.PCB
		Utils.MutexEjecutando.Lock()
		for _, p := range algoritmos.ColaEjecutando.Values() {
			if p.PID == pid {
				proceso = p
				algoritmos.ColaEjecutando.Remove(proceso)
			}
		}
		Utils.MutexEjecutando.Unlock()

		liberarCPU(cpuID)

		//ENVIAR A BLOCKED
		Utils.MutexBloqueado.Lock()
		proceso.ME[pcb.EstadoBlocked]++
		proceso.Estado = pcb.EstadoBlocked
		algoritmos.ColaBloqueado.Add(proceso)
		logger.Info("## (<%d>) Pasa del estado EXECUTE al estado BLOCKED", proceso.PID)
		Utils.MutexBloqueado.Unlock()

		//Enviar al módulo IO (usando los datos del mensaje recibido)
		comunicacion.EnviarContextoIO(msg.Nombre, msg.PID, msg.Duracion)
	}
}

func FinDeIO() {
	for {
		//WAIT mensaje fin de IO (bloqueante)
		pid := <-Utils.NotificarFinIO

		//BUSCAR EN PCB BLOCKED
		var proceso *pcb.PCB
		Utils.MutexBloqueado.Lock()
		for _, p := range algoritmos.ColaBloqueado.Values() {
			if p.PID == pid {
				proceso = p
				algoritmos.ColaBloqueado.Remove(proceso)
			}
		}
		Utils.MutexBloqueado.Unlock()

		Utils.MutexReady.Lock()
		proceso.ME[pcb.EstadoReady]++
		proceso.Estado = pcb.EstadoReady
		algoritmos.ColaReady.Add(proceso)
		logger.Info("## (%d) finalizó IO y pasa a READY", pid)
		Utils.MutexReady.Unlock()

		//Notificar al despachador llegada a READY
		Utils.NotificarDespachador <- pid
	}
}
