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
	//ESTA FUNCION VA A ATENDER EL REQUERIMIENTO QUE PASE DE SUSṔ BLQOUEADO A SUSP READY.
	go AtenderSuspBlockedAFinIO()
	go DespacharIO()
}

// ManejadorMedianoPlazo se quedará escuchando
// bloqueos e iniciará un timer por cada PID que llegue.

func ManejadorMedianoPlazo() {
	for bp := range Utils.ChannelProcessBlocked {
		// arrancá un timer en paralelo para CADA proceso bloqueado
		//y lleno el MAP con un mutex entonces cada pid tiene un mutex individual
		ch := make(chan Utils.IOEvent, 1)
		Utils.MutexIOWaiters.Lock()
		Utils.IOWaiters[bp.PID] = ch
		Utils.MutexIOWaiters.Unlock()
		go monitorBloqueado(bp)
	}
}

// moverDeBlockedAReady quita de BLOCKED y encola en READY
func moverDeBlockedAReady(ioLibre Utils.IOEvent) bool {
	// busca en ColaBloqueado
	Utils.MutexBloqueado.Lock()

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

	// Remover de BLOCKEd

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

// pasa de SUSP BLOQUEADO A SUSP READY
func AtenderSuspBlockedAFinIO() {
	for ev := range Utils.FinIODesdeSuspBlocked {
		// Para cada evento de SUSP.BLOCKED arrancoo un hilo
		go func(ev Utils.IOEvent) {
			// Obtener el canal individual para ese PID
			Utils.MutexIOWaiters.Lock()
			ch, ok := Utils.IOWaiters[ev.PID]
			Utils.MutexIOWaiters.Unlock()
			if !ok {
				logger.Warn("AtenderSuspBlockedAFinIO: no hay canal para PID %d", ev.PID)
				return
			}

			// Espero el fin de I/O para este PID
			ioFin := <-ch

			// Mover de SUSP.BLOCKED a SUSP.READY
			Utils.MutexBloqueadoSuspendido.Lock()
			var proc *pcb.PCB
			for _, p := range algoritmos.ColaBloqueadoSuspendido.Values() {
				if p.PID == ioFin.PID {
					proc = p
					break
				}
			}
			if proc != nil {
				algoritmos.ColaBloqueadoSuspendido.Remove(proc)
				pcb.CambiarEstado(proc, pcb.EstadoSuspReady)
				Utils.MutexSuspendidoReady.Lock()
				algoritmos.ColaSuspendidoReady.Add(proc)
				Utils.MutexSuspendidoReady.Unlock()
				logger.Info("## (%d) Pasa de SUSP.BLOCKED a SUSP.READY", proc.PID)
				// Notificar a largo plazo para reintentar ingresar
				Utils.InitProcess <- struct{}{}
			} else {
				logger.Warn("AtenderSuspBlockedAFinIO: PID %d no estaba en SUSP.BLOCKED", ev.PID)
			}
			Utils.MutexBloqueadoSuspendido.Unlock()

			//  Limpiar el canal para este PID
			Utils.MutexIOWaiters.Lock()
			delete(Utils.IOWaiters, ev.PID)
			Utils.MutexIOWaiters.Unlock()
		}(ev)
	}
}

func DespacharIO() {
	for {
		ioLibre := <-Utils.NotificarIOLibre // Notificación de que una IO se liberó

		Utils.MutexBloqueado.Lock()

		// Buscar un pedido que coincida con el tipo de IO liberado
		var pedidoEncontrado *algoritmos.PedidoIO
		for _, pedido := range algoritmos.PedidosIO.Values() {
			if pedido.Nombre == ioLibre.Nombre {
				pedidoEncontrado = pedido
				break
			}
		}

		// Si no se encontró un pedido compatible, no se hace nada
		if pedidoEncontrado == nil {
			Utils.MutexBloqueado.Unlock()
			continue
		}

		// Buscar una instancia libre de ese tipo de IO
		ioList := globals.IOs[ioLibre.Nombre]
		var ioAsignada *globals.DatosIO
		var indice int = -1

		globals.IOMu.Lock()
		for i := range ioList {
			if !ioList[i].Ocupada {
				indice = i
				ioAsignada = &ioList[i]
				break
			}
		}
		globals.IOMu.Unlock()

		// Si no se encontró instancia libre (raro pero posible en concurrencia), salir
		if ioAsignada == nil {
			logger.Info("No hay instancias libres de IO tipo %s", ioLibre.Nombre)
			Utils.MutexBloqueado.Unlock()
			continue
		}

		// Marcar la IO como ocupada con el PID asignado
		globals.IOMu.Lock()
		ioList[indice].Ocupada = true
		ioList[indice].PID = pedidoEncontrado.PID
		globals.IOMu.Unlock()

		logger.Info("Asignada IO <%s> a proceso <%d>", ioLibre.Nombre, pedidoEncontrado.PID)

		// Remover el pedido de la cola de bloqueado
		algoritmos.PedidosIO.Remove(pedidoEncontrado)
		Utils.MutexBloqueado.Unlock()

		// Enviar al módulo IO correspondiente
		go comunicacion.EnviarContextoIO(*ioAsignada, pedidoEncontrado.PID, pedidoEncontrado.Duracion)
	}
}
