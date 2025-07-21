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
		//logger.Debug("Arranca Despachador")

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
			//logger.Info("No hay proceso listo para planificar")
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

	//logger.Info("CPU <%s> libre", cpuID)

	Utils.NotificarDespachador <- 1
}

func DesalojarProceso() {
	for {
		//WAIT mensaje contexto interrupcion
		msg := <-Utils.ContextoInterrupcion
		pid := msg.PID
		pc := msg.PC
		cpuID := msg.CpuID

		//logger.Info("## (<%d>) Interrumpido de CPU <%s>", pid, cpuID)

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

		//AVISAR QUE UN PROCESO LLEGA A READY
		Utils.NotificarDespachador <- proceso.PID
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

		//ENVIAR A BLOCKED
		Utils.MutexBloqueado.Lock()
		pcb.CambiarEstado(proceso, pcb.EstadoBlocked)
		algoritmos.ColaBloqueado.Add(proceso)
		logger.Info("## (<%d>) Pasa del estado EXECUTE al estado BLOCKED", proceso.PID)
		Utils.MutexBloqueado.Unlock()

		// Buscar IO disponible
		globals.IOMu.Lock()
		ioList, existe := globals.IOs[tipoIO] //Existe el tipo de IO en el map?
		if !existe || len(ioList) == 0 {
			// No existe el tipo de IO / No hay ninguna instancia
			globals.IOMu.Unlock()

			//Avisar EXIT A LARGO
			logger.Info("## No existe esa IO. (<%d>) Pasa a finalizar", proceso.PID)
			Utils.ChannelFinishprocess <- Utils.FinishProcess{
				PID: proceso.PID,
				PC:  proceso.PC,
			}
			continue
		}
		globals.IOMu.Unlock()

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
		logger.Warn("Se desconectó una IO <%s:%d> - PID activo: <%d>", io.Nombre, io.Puerto, io.PID)

		// 1. Remover instancia de IO
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

		// Si no quedan más instancias, eliminar el tipo y finalizar todos los pedidos pendientes
		if len(nuevaLista) == 0 {
			delete(globals.IOs, io.Nombre)
			logger.Info("No quedan IOs activas del tipo <%s>, eliminando del mapa y finalizando pedidos", io.Nombre)

			Utils.MutexPedidosIO.Lock()
			for _, pedido := range algoritmos.PedidosIO.Values() {
				if pedido.Nombre == io.Nombre {
					logger.Warn("Finalizando PID <%d> que esperaba IO <%s> (sin instancias activas)", pedido.PID, pedido.Nombre)

					// Liberar memoria y finalizar
					comunicacion.LiberarMemoria(pedido.PID)
					finalizarProceso(pedido.PID, 0, "")
					algoritmos.PedidosIO.Remove(pedido)
				}
			}
			Utils.MutexPedidosIO.Unlock()

		} else {
			// Aún quedan IOs de este tipo
			globals.IOs[io.Nombre] = nuevaLista
		}
		globals.IOMu.Unlock()

		//Si la IO se desconecta mientras tiene un proceso...
		var proceso *pcb.PCB

		// 1. Buscar y remover de BLOCKED
		Utils.MutexBloqueado.Lock()
		for _, p := range algoritmos.ColaBloqueado.Values() {
			if p.PID == io.PID {
				proceso = p
				algoritmos.ColaBloqueado.Remove(p)
				break
			}
		}
		Utils.MutexBloqueado.Unlock()

		// 2. Buscar y remover de BLOCKED_SUSPENDED si no se encontró
		if proceso == nil {
			//logger.Warn("## PID <%d> desconectado pero no estaba en BLOCKED", pidEjecutado)

			Utils.MutexBloqueadoSuspendido.Lock()
			for _, p := range algoritmos.ColaBloqueadoSuspendido.Values() {
				if p.PID == io.PID {
					proceso = p
					algoritmos.ColaBloqueadoSuspendido.Remove(p)
					break
				}
			}
			Utils.MutexBloqueadoSuspendido.Unlock()
		}

		// 3. Si estaba esperando una IO, mandarlo a EXIT
		if proceso != nil {
			logger.Info("Finalizando PID <%d> por desconexión de IO <%s>", proceso.PID, io.Nombre)
			Utils.ChannelFinishprocess <- Utils.FinishProcess{
				PID: proceso.PID,
				PC:  proceso.PC,
			}
		}
	}
}

func MostrarCOLABLOQUEADO() {
	lista := algoritmos.ColaBloqueado.Values()

	if len(lista) == 0 {
		logger.Info("Cola NEW vacía")
		return
	}

	logger.Info("Contenido de la cola New:")
	for _, proceso := range lista {
		logger.Info(" - PCB EN COLA New con PID: %d, TAMAÑO: %d", proceso.PID, proceso.ProcessSize)
	}
}
