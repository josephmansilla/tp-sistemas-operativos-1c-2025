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
		// Esperar señal del CPU para bloquear por IO
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

		// Buscar IO disponible
		globals.IOMu.Lock()

		//Exista el tipo de IO en el map?
		ioList, existe := globals.IOs[tipoIO]
		if !existe || len(ioList) == 0 {
			// No existe el tipo de IO o no hay ninguna instancia
			globals.IOMu.Unlock()

			Utils.MutexSalida.Lock()
			pcb.CambiarEstado(proceso, pcb.EstadoExit)
			algoritmos.ColaSalida.Add(proceso)
			logger.Info("## %s NO SE ENCUENTRA", tipoIO)
			logger.Info("## (<%d>) Pasa del estado EXECUTE al estado EXIT", proceso.PID)
			Utils.MutexSalida.Unlock()
			continue
		}

		// Buscar una instancia libre
		var ioAsignada *globals.DatosIO
		for i := range ioList {
			if !ioList[i].Ocupada {
				ioList[i].Ocupada = true
				ioList[i].PID = pid
				ioAsignada = &ioList[i]
				break
			}
		}
		globals.IOMu.Unlock()

		//ENVIAR A BLOCKED
		Utils.MutexBloqueado.Lock()
		pcb.CambiarEstado(proceso, pcb.EstadoBlocked)
		algoritmos.ColaBloqueado.Add(proceso)
		logger.Info("## (<%d>) Pasa del estado EXECUTE al estado BLOCKED", proceso.PID)
		Utils.MutexBloqueado.Unlock()

		if ioAsignada == nil {
			// No se encontró una IO libre, el proceso debe esperar (pero igual se bloquea y se agrega a pedidos pendientes)
			logger.Info("No hay IOs libres del tipo: %s. <%d> debe esperar", tipoIO, pid)
			io := algoritmos.PedidoIO{
				Nombre:   msg.Nombre,
				PID:      pid,
				Duracion: msg.Duracion,
			}
			Utils.MutexPedidosIO.Lock()
			algoritmos.PedidosIO.Add(&io)
			Utils.MutexPedidosIO.Unlock()
		} else {
			//Lo envia a IO hallada para bloquearse
			logger.Info("Asignada IO <%s> (puerto %d) a proceso <%d>", tipoIO, ioAsignada.Puerto, pid)
			go comunicacion.EnviarContextoIO(*ioAsignada, pid, msg.Duracion)
		}

		//Cuando el corto plazo termina de bloquear al proceso en particular
		//le avisa al Mediano plazo para que empiece el Timer para ESE proceso
		//le manda PID, DURACION y NOMBRE IO como señal
		go func(p int) {
			Utils.ChannelProcessBlocked <- Utils.BlockProcess{
				PID:      pid,
				PC:       pc,
				Nombre:   msg.Nombre,
				Duracion: msg.Duracion,
			}
		}(pid)
	}
}

func DesconexionIO() {
	for {
		//WAIT mensaje desconexion de IO (Tipo, PID y Puerto)
		io := <-Utils.NotificarDesconexion

		// 1) Remover instancia de IO que tenga ese puerto
		globals.IOMu.Lock()

		// Buscar el tipo de IO directamente
		instancias, existe := globals.IOs[io.Nombre]
		if !existe {
			logger.Warn("No se encontró el tipo de IO <%s>", io.Nombre)
			globals.IOMu.Unlock()
			return
		}

		// Crear nueva lista de instancias, excluyendo la desconectada
		nuevaLista := []globals.DatosIO{}
		for _, instancia := range instancias {
			if instancia.Puerto != io.Puerto {
				nuevaLista = append(nuevaLista, instancia)
			} else {
				logger.Info("Removida instancia IO <%s> con Puerto <%d>", io.Nombre, io.Puerto)
			}
		}

		// Si no quedan más instancias, eliminar el tipo del mapa
		if len(nuevaLista) == 0 {
			delete(globals.IOs, io.Nombre)
			logger.Info("No quedan IOs activas del tipo <%s>, eliminado del mapa", io.Nombre)
		} else {
			globals.IOs[io.Nombre] = nuevaLista
		}
		globals.IOMu.Unlock()

		if io.PID < 0 {
			continue // ignorar PIDs inválidos de procesos y desconectar
		}

		//Si la IO se desconecta mientras tiene un proceso...
		//2) BUSCAR PCB en BLOCKED
		var proceso *pcb.PCB
		Utils.MutexBloqueado.Lock()
		for _, p := range algoritmos.ColaBloqueado.Values() {
			if p.PID == io.PID {
				proceso = p
				algoritmos.ColaBloqueado.Remove(proceso)
			}
		}
		Utils.MutexBloqueado.Unlock()

		if proceso == nil {
			logger.Warn("## PID <%d> desconectado pero no estaba en BLOCKED", io.PID)
			continue
		}

		//3) MOVER A PROCESO A EXIT
		Utils.MutexSalida.Lock()
		pcb.CambiarEstado(proceso, pcb.EstadoExit)
		algoritmos.ColaSalida.Add(proceso)
		logger.Info("## (%d) finalizó IO y pasa a EXIT por desconexión", io.PID)
		Utils.MutexSalida.Unlock()

		logger.Info("## IO desconectado correctamente para PID <%d>", io.PID)
	}
}
