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
	go DesconexionIO()
	go DesalojarProceso()
}

func DespacharProceso() {
	for {
		// WAIT hasta que llegue un proceso a READY
		//o se libere una CPU por SYSCALL DE EXIT O I/O
		<-Utils.NotificarDespachador
		logger.Info("Arranca Despachador")

		var proceso *pcb.PCB

		switch globals.KConfig.SchedulerAlgorithm {
		case "FIFO":
			proceso = algoritmos.ColaReady.First() //Toma el primero de la cola Ready
		case "SJF":
			proceso = algoritmos.SeleccionarSJF() //Toma la estimacion mas corta
		case "SRT":
			proceso = algoritmos.SeleccionarSJF() //Toma la estimacion mas corta
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

			if globals.KConfig.SchedulerAlgorithm == "SRT" {
				algoritmos.Desalojo(proceso)
			}
			continue
		}

		proceso.CpuID = cpuID
		logger.Info("Proceso <%d> -> EXECUTE en CPU <%s>", proceso.PID, cpuID)

		Utils.MutexEjecutando.Lock()
		pcb.CambiarEstado(proceso, pcb.EstadoExecute)
		algoritmos.ColaEjecutando.Add(proceso)
		Utils.MutexEjecutando.Unlock()

		Utils.MutexReady.Lock()
		algoritmos.ColaReady.Remove(proceso)
		Utils.MutexReady.Unlock()

		cpu := globals.CPUs[cpuID]
		cpu.Ocupada = true
		globals.CPUs[cpuID] = cpu

		go comunicacion.EnviarContextoCPU(cpuID, proceso)
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
	for {
		//WAIT mensaje contexto interrupcion
		msg := <-Utils.ContextoInterrupcion
		pid := msg.PID
		pc := msg.PC
		cpuID := msg.CpuID

		logger.Info("## (<%d>) Interrumpido de CPU <%s>", pid, cpuID)

		//BUSCAR en EXECUTE y actualizar PC proveniente de CPU
		var proceso *pcb.PCB
		Utils.MutexEjecutando.Lock()
		for _, p := range algoritmos.ColaEjecutando.Values() {
			if p.PID == pid {
				algoritmos.ColaEjecutando.Remove(p)
				p.PC = pc
				proceso = p
				break
			}
		}
		Utils.MutexEjecutando.Unlock()

		liberarCPU(cpuID)

		//ENVIAR A READY
		Utils.MutexBloqueado.Lock()
		pcb.CambiarEstado(proceso, pcb.EstadoReady)
		algoritmos.ColaReady.Add(proceso)
		logger.Info("## (<%d>) Pasa del estado EXECUTE al estado READY", proceso.PID)
		Utils.MutexBloqueado.Unlock()
	}
}

func BloquearProceso() {
	for {
		//WAIT mensaje de IO (bloqueante)
		msg := <-Utils.NotificarComienzoIO
		pid := msg.PID
		pc := msg.PC
		cpuID := msg.CpuID
		tipoIO := msg.Nombre

		//BUSCAR en EXECUTE y actualizar PC proveniente de CPU
		var proceso *pcb.PCB
		Utils.MutexEjecutando.Lock()
		for _, p := range algoritmos.ColaEjecutando.Values() {
			if p.PID == pid {
				algoritmos.ColaEjecutando.Remove(p)
				p.PC = pc
				proceso = p
				cpuID = p.CpuID
				break
			}
		}
		Utils.MutexEjecutando.Unlock()

		liberarCPU(cpuID)

		//Necesito encontrar a que IO libre mandarle
		globals.IOMu.Lock()
		//Exista el tipo de IO en el map?
		for {
			_, existe := globals.IOs[tipoIO]
			if existe {
				break
			} else {
				//ENVIAR A EXIT si no existe o cola null
				Utils.MutexSalida.Lock()
				pcb.CambiarEstado(proceso, pcb.EstadoExit)
				algoritmos.ColaSalida.Add(proceso)
				logger.Info("## %s NO SE ENCUENTRA", tipoIO)
				logger.Info("## (<%d>) Pasa del estado EXECUTE al estado EXIT", proceso.PID)
				Utils.MutexSalida.Unlock()
			}
		}

		// Buscar una instancia no ocupada
		var ioData *globals.DatosIO
		for {
			for i := range globals.IOs[tipoIO] {
				if !globals.IOs[tipoIO][i].Ocupada {
					globals.IOs[tipoIO][i].Ocupada = true
					ioData = &globals.IOs[tipoIO][i]
				}
			}
			//No hay libres, va a bloqueado igual pero espera en ese estado
			logger.Info("No hay IOs libres del tipo: %s. Debe esperar", tipoIO)
		}
		globals.IOMu.Unlock()

		//ENVIAR A BLOCKED
		Utils.MutexBloqueado.Lock()
		pcb.CambiarEstado(proceso, pcb.EstadoBlocked)
		algoritmos.ColaBloqueado.Add(proceso)
		logger.Info("## (<%d>) Pasa del estado EXECUTE al estado BLOCKED", proceso.PID)
		Utils.MutexBloqueado.Unlock()

		//Cuando el corto plazo termina de bloquear al proceso en particular
		//le avisa al Mediano plazo para que empiece el Timer para ESE proceso
		//le manda PID, DURACION y NOMBRE IO como señal
		go func(p int) {
			Utils.ChannelProcessBlocked <- Utils.BlockProcess{
				PID:      pid,
				PC:       pc,
				Nombre:   ioData.Tipo,
				Duracion: msg.Duracion,
			}
		}(pid)
	}
}

func FinDeIO() {
	for {
		// wait del mediano plazo para que indique orden. el proceso esta efectivamente en la cola susp. Bloqueado.
		Utils.NotificarTimeoutBlocked <- struct{}{}

		//WAIT mensaje fin de IO (bloqueante)
		//pid := <-Utils.NotificarFinIO

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
		pcb.CambiarEstado(proceso, pcb.EstadoReady)
		algoritmos.ColaReady.Add(proceso)
		logger.Info("## (%d) finalizó IO y pasa a READY", pid)
		Utils.MutexReady.Unlock()

		//Notificar al despachador llegada a READY
		Utils.NotificarDespachador <- pid
	}
}

func DesconexionIO() {
	for {
		//WAIT mensaje desconexion de IO
		pid := <-Utils.NotificarDesconexion

		if pid < 0 {
			continue // ignorar PIDs inválidos
		}

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

		if proceso == nil {
			logger.Warn("## PID <%d> desconectado pero no estaba en BLOCKED", pid)
			continue
		}

		//MOVER A EXIT
		Utils.MutexSalida.Lock()
		pcb.CambiarEstado(proceso, pcb.EstadoExit)
		algoritmos.ColaSalida.Add(proceso)
		logger.Info("## (%d) finalizó IO y pasa a EXIT por desconexión", pid)
		Utils.MutexSalida.Unlock()

		logger.Info("## IO desconectado correctamente para PID <%d>", pid)
	}
}
