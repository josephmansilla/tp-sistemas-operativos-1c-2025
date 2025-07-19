package planificadores

import (
	"github.com/sisoputnfrba/tp-golang/kernel/Utils"
	"github.com/sisoputnfrba/tp-golang/kernel/algoritmos"
	"github.com/sisoputnfrba/tp-golang/kernel/comunicacion"
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/pcb"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"time"
)

func PlanificadorMedianoPlazo() {
	logger.Info("Iniciando el planificador de Mediano Plazo")
	go ManejadorMedianoPlazo()
	//ESTA FUNCION VA A ATENDER EL REQUERIMIENTO QUE PASE DE SUSP.BLOCKED A SUSP.READY.
	go AtenderSuspBlockedAFinIO()
	go DespacharIO()
}

// ManejadorMedianoPlazo se quedará escuchando
// bloqueos e iniciará un timer por cada PID que llegue.

func ManejadorMedianoPlazo() {
	for bp := range Utils.ChannelProcessBlocked {
		// arrancá un timer en paralelo para CADA proceso bloqueado
		//y lleno el MAP con un mutex entonces cada pid tiene un mutex individual
		/*ch := make(chan Utils.IOEvent, 1)
		Utils.MutexIOWaiters.Lock()
		Utils.IOWaiters[bp.PID] = ch
		Utils.MutexIOWaiters.Unlock()*/
		ch1 := make(chan Utils.IOEvent, 1)

		// 2) Guárdalo en el map protegido con mutex
		Utils.MutexIOFinishedWaiters.Lock()
		Utils.FinIOWaiters[bp.PID] = ch1
		logger.Error("INICIALICE CANAL DE ")
		Utils.MutexIOFinishedWaiters.Unlock()
		go monitorBloqueado(bp)
	}
}

// moverDeBlockedAReady quita de BLOCKED y encola en READY
func moverDeBlockedAReady(ioLibre Utils.IOEvent) bool {
	// busca en ColaBloqueado
	Utils.MutexBloqueado.Lock()

	// Remover de BLOCKED
	var proceso *pcb.PCB
	for _, p := range algoritmos.ColaBloqueado.Values() {
		if p.PID == ioLibre.PID {
			proceso = p
			algoritmos.ColaBloqueado.Remove(p)
			break
		}
	}
	Utils.MutexBloqueado.Unlock()

	if proceso == nil {
		// No se encontró el proceso en BLOCKED
		return false
	}

	globals.IOMu.Lock()
	instancias, ok := globals.IOs[ioLibre.Nombre]

	if ok {
		for i := range instancias {
			if instancias[i].PID == ioLibre.PID && instancias[i].Puerto == ioLibre.Puerto {
				globals.IOs[ioLibre.Nombre][i].Ocupada = false
				break
			}
		}
	}
	globals.IOMu.Unlock()

	// Agregar a READY
	Utils.MutexReady.Lock()
	pcb.CambiarEstado(proceso, pcb.EstadoReady)
	algoritmos.ColaReady.Add(proceso)
	logger.Info("## <%d> finalizó IO y pasa a READY", ioLibre.PID)
	Utils.MutexReady.Unlock()

	//Señal al corto plazo para despachar
	Utils.NotificarDespachador <- ioLibre.PID
	return true
}

// moverDeBlockedASuspBlocked quita de BLOCKED y encola en SUSP.BLOCKED
func moverDeBlockedASuspBlocked(pid int) bool {
	Utils.MutexBloqueado.Lock()

	// busca en ColaBloqueado
	var proceso *pcb.PCB
	for _, p := range algoritmos.ColaBloqueado.Values() {
		if p.PID == pid {
			proceso = p
			break
		}
	}

	if proceso == nil {
		// No se encontró el proceso en BLOCKED
		logger.Info("## MedianoPlazo: No se encontró en Blocked")
		//Utils.MutexBloqueado.Unlock()
		return false
	}

	algoritmos.ColaBloqueado.Remove(proceso)
	Utils.MutexBloqueado.Unlock()

	Utils.MutexBloqueadoSuspendido.Lock()
	pcb.CambiarEstado(proceso, pcb.EstadoSuspBlocked)
	algoritmos.ColaBloqueadoSuspendido.Add(proceso)
	Utils.MutexBloqueadoSuspendido.Unlock()
	logger.Info("## (<%d>) Pasa del estado BLOCKED al estado SUSP.BLOCKED", proceso.PID)

	return true
}

func monitorBloqueado(bp Utils.BlockProcess) {
	pid := bp.PID

	// Obtenemos ambos canales bajo sus mutex
	/*Utils.MutexIOWaiters.Lock()
	blockedCh, okB := Utils.IOWaiters[pid]
	Utils.MutexIOWaiters.Unlock()
	if !okB {
		logger.Warn("monitorBloqueado: no existe blockedCh para PID %d", pid)
		return
	}*/

	Utils.MutexIOFinishedWaiters.Lock()
	finishedCh, okF := Utils.FinIOWaiters[pid]
	Utils.MutexIOFinishedWaiters.Unlock()
	if !okF {
		logger.Warn("monitorBloqueado: no existe finishedCh para PID %d", pid)
		return
	}

	// 1) Esperamos señal de que realmente se bloqueó (del corto plazo)
	/*<-blockedCh*/

	logger.Info("Arrancó TIMER para PID <%d>", pid)
	suspensión := time.Duration(globals.KConfig.SuspensionTime) * time.Millisecond
	timer := time.NewTimer(suspensión)
	defer timer.Stop()

	// encolamos el pedido de I/O
	Utils.MutexPedidosIO.Lock()
	algoritmos.PedidosIO.Add(&algoritmos.PedidoIO{
		Nombre:   bp.Nombre,
		PID:      bp.PID,
		Duracion: bp.Duracion,
	})
	Utils.MutexPedidosIO.Unlock()
	Utils.NotificarIOLibre <- Utils.IOEvent{Nombre: bp.Nombre, PID: bp.PID}

	select {
	case ioEvt := <-finishedCh:
		// fin de I/O antes del timeout → READY
		moverDeBlockedAReady(ioEvt)

	case <-timer.C:
		// timeout → SUSP.BLOCKED
		if moverDeBlockedASuspBlocked(pid) {
			logger.Info("PID <%d> → SUSP.BLOCKED (timeout)", pid)
			if err := comunicacion.SolicitarSuspensionEnMemoria(pid); err == nil {
				Utils.InitProcess <- struct{}{}
			}
			// avisamos al subplanificador para pasarlo luego a SUSP.READY
			Utils.FinIODesdeSuspBlocked <- Utils.IOEvent{PID: pid, Nombre: bp.Nombre}
		}
	}
}

// pasa de SUSP BLOQUEADO A SUSP READY (con FIN DE IO)
func AtenderSuspBlockedAFinIO() {
	for ev := range Utils.FinIODesdeSuspBlocked {
		// Obtener el canal dedicado a este PID
		Utils.MutexIOWaiters.Lock()
		finIOChan, ok := Utils.FinIOWaiters[ev.PID]
		Utils.MutexIOWaiters.Unlock()

		if !ok {
			logger.Warn("AtenderSuspBlockedAFinIO: no hay canal para PID %d", ev.PID)
			continue
		}

		// Arrancar un goroutine que espere el verdadero fin de I/O
		go func(pid int, ch chan Utils.IOEvent) {
			logger.Info("Esperando fin de I/O para PID <%d> en SUSP.BLOCKED...", pid)
			ioEvt := <-ch
			logger.Info("ME LLEGO FIN DE IO PARA PROCESO DE PI: <%d>", ioEvt.PID)

			// Pasar de SUSP.BLOCKED a SUSP.READY
			Utils.MutexBloqueadoSuspendido.Lock()
			var proc *pcb.PCB
			for _, p := range algoritmos.ColaBloqueadoSuspendido.Values() {
				if p.PID == pid {
					proc = p
					break
				}
			}
			Utils.MutexBloqueadoSuspendido.Unlock()

			if proc != nil {
				switch globals.KConfig.ReadyIngressAlgorithm {
				case "FIFO":
					Utils.MutexSuspendidoReady.Lock()
					pcb.CambiarEstado(proc, pcb.EstadoSuspReady)
					algoritmos.ColaSuspendidoReady.Add(proc)
					Utils.MutexSuspendidoReady.Unlock()
					logger.Info("PID <%d> pasa de SUSP.BLOCKED a SUSP.READY (FIFO)", pid)
				case "PMCP":
					pcb.CambiarEstado(proc, pcb.EstadoSuspReady)
					algoritmos.AddPMCPSusp(proc)
					logger.Info("PID <%d> pasa de SUSP.BLOCKED a SUSP.READY (PMCP)", pid)
				default:
					logger.Error("Algoritmo de ingreso desconocido")
					return
				}
				logger.Info("## (%d) Pasa de SUSP.BLOCKED a SUSP.READY", pid)

				// Notificar al planificador de largo plazo
				Utils.InitProcess <- struct{}{}
			} else {
				logger.Warn("AtenderSuspBlockedAFinIO: PID %d no estaba en SUSP.BLOCKED", pid)
			}

			// Limpiar el canal dedicado a este PID
			Utils.MutexIOWaiters.Lock()
			delete(Utils.IOWaiters, pid)
			Utils.MutexIOWaiters.Unlock()
			//muestro cola suspready cada vez que llega un fin IO DEBERIA ESTAR CARGADA CON LOS PROCESOS
			MostrarColasSUSPREADY()

		}(ev.PID, finIOChan)
	}
}

func DespacharIO() {
	for {
		<-Utils.NotificarIOLibre // Esperar señal de IO libre

		Utils.MutexPedidosIO.Lock()
		pedidos := algoritmos.PedidosIO.Values()

		//Si no hay pedidos continua...
		if len(pedidos) == 0 {
			Utils.MutexPedidosIO.Unlock()
			continue
		}

		//BUSCAR PEDIDOS DE IO pendientes
		pedido := algoritmos.PedidosIO.First() // FIFO
		Utils.MutexPedidosIO.Unlock()

		// Buscar una instancia de IO LIBRE del tipo Nombre
		globals.IOMu.Lock()
		var ioAsignada *globals.DatosIO
		for i := range globals.IOs[pedido.Nombre] {
			if !globals.IOs[pedido.Nombre][i].Ocupada {
				globals.IOs[pedido.Nombre][i].Ocupada = true
				globals.IOs[pedido.Nombre][i].PID = pedido.PID
				ioAsignada = &globals.IOs[pedido.Nombre][i]
				break
			}
		}
		globals.IOMu.Unlock()

		if ioAsignada == nil {
			// No se encontró una IO libre
			logger.Warn("No se encontró IO libre: %s. PID <%d> Debe esperar", pedido.Nombre, pedido.PID)
			continue
		}

		logger.Info("Asignada IO <%s> (puerto %d) a proceso <%d>", ioAsignada.Tipo, ioAsignada.Puerto, pedido.PID)
		algoritmos.PedidosIO.Remove(pedido)
		go comunicacion.EnviarContextoIO(*ioAsignada, pedido.PID, pedido.Duracion)
	}
}

func ManejadorIO(bp Utils.BlockProcess) {
	// Buscar una instancia de IO LIBRE
	globals.IOMu.Lock()
	ioList := globals.IOs[bp.Nombre]
	var ioAsignada *globals.DatosIO
	for i := range ioList {
		if !ioList[i].Ocupada {
			ioList[i].Ocupada = true
			ioList[i].PID = bp.PID
			ioAsignada = &ioList[i]
			break
		}
	}
	globals.IOMu.Unlock()

	if ioAsignada == nil {
		//No se encontró una IO libre, el proceso debe esperar
		//(se agrega a pedidos pendientes)
		//LOOP ESPERANDO SEMAFOROS DE IOS LIBRES
		logger.Info("No hay IOs libres del tipo: %s. <%d> debe esperar", bp.Nombre, bp.PID)
		Utils.MutexPedidosIO.Lock()
		algoritmos.PedidosIO.Add(&algoritmos.PedidoIO{
			Nombre:   bp.Nombre,
			PID:      bp.PID,
			Duracion: bp.Duracion,
		})
		Utils.MutexPedidosIO.Unlock()
		return
	}

	//Lo envia a IO hallada para bloquearse
	logger.Info("Asignada IO <%s> (puerto %d) a proceso <%d>", bp.Nombre, ioAsignada.Puerto, bp.PID)
	go comunicacion.EnviarContextoIO(*ioAsignada, bp.PID, bp.Duracion)
}
func MostrarColasSUSPREADY() {
	lista := algoritmos.ColaSuspendidoReady.Values()

	if len(lista) == 0 {
		logger.Info("Cola SUSPENDIDO READY vacía")
		return
	}

	logger.Info("Contenido de la cola SUSPENDIDO READY:")
	for _, proceso := range lista {
		logger.Info(" - PCB EN COLA SUSPENDIDO READY con PID: %d, TAMAÑO: %d", proceso.PID, proceso.ProcessSize)
	}
}

/*
	func monitorBloqueado(bp Utils.BlockProcess) {
		pid := bp.PID

		//ACA se guarda la referencia ala posicion del mutex correspondiente al mismo proceso que abrio el hilo de monitor
		Utils.MutexIOWaiters.Lock()
		finIOChan, ok := Utils.IOWaiters[pid]
		Utils.MutexIOWaiters.Unlock()
		if !ok {
			logger.Warn("monitorBloqueado: no existe canal para PID %d", pid)
			return
		}

		logger.Info("Arrancó TIMER para PID <%d>", pid)
		suspensión := time.Duration(globals.KConfig.SuspensionTime) * time.Millisecond
		timer := time.NewTimer(suspensión)
		defer timer.Stop()

		//DESPACHAR A IO
		//agregar pedido a listaDepedidos
		Utils.MutexPedidosIO.Lock()
		algoritmos.PedidosIO.Add(&algoritmos.PedidoIO{
			Nombre:   bp.Nombre,
			PID:      bp.PID,
			Duracion: bp.Duracion,
		})
		Utils.MutexPedidosIO.Unlock()

		Utils.NotificarIOLibre <- Utils.IOEvent{
			Nombre: bp.Nombre,
			PID:    bp.PID,
		}

		select {
		case ioEvt := <-finIOChan:
			// fin de IO antes del timeout → READY
			moverDeBlockedAReady(ioEvt)

		case <-timer.C:
			if moverDeBlockedASuspBlocked(pid) {
				logger.Info("PID <%d> → SUSP.BLOCKED (timeout)", pid)
				if err := comunicacion.SolicitarSuspensionEnMemoria(pid); err == nil {
					Utils.InitProcess <- struct{}{}
				}
				Utils.FinIODesdeSuspBlocked <- Utils.IOEvent{PID: pid, Nombre: bp.Nombre}
			}
		}
	}
*/
