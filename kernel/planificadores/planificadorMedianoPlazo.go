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
}

// ManejadorMedianoPlazo se quedará escuchando
// bloqueos e iniciará un timer por cada PID que llegue.
func ManejadorMedianoPlazo() {
	for {
		procesoBloqueado := <-Utils.ChannelProcessBlocked // señal llegada de BLOCKED con el PID del proceso para que arranque su TIMER
		go monitorBloqueado(procesoBloqueado)
	}
}

// moverDeBlockedAReady quita de BLOCKED y encola en READY
func moverDeBlockedAReady(ioLibre Utils.IODesconexion) bool {
	// busca en ColaBloqueado
	Utils.MutexBloqueado.Lock()
	defer Utils.MutexBloqueado.Unlock()

	var proceso *pcb.PCB
	for _, p := range algoritmos.ColaBloqueado.Values() {
		if p.PID == ioLibre.PID {
			proceso = p
			break
		}
	}

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

	// Remover de BLOCKED
	algoritmos.ColaBloqueado.Remove(proceso)

	// Agregar a READY
	Utils.MutexReady.Lock()
	defer Utils.MutexReady.Unlock()
	pcb.CambiarEstado(proceso, pcb.EstadoReady)
	algoritmos.ColaReady.Add(proceso)
	logger.Info("## (%d) finalizó IO y pasa a READY", ioLibre.PID)

	//Señal al corto plazo para despachar
	Utils.NotificarDespachador <- ioLibre.PID
	return true

}

// moverDeBlockedASuspBlocked quita de BLOCKED y encola en SUSP.BLOCKED
func moverDeBlockedASuspBlocked(pid int) bool {
	Utils.MutexBloqueado.Lock()
	defer Utils.MutexBloqueado.Unlock()

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
		return false
	}

	algoritmos.ColaBloqueado.Remove(proceso)
	Utils.MutexBloqueado.Unlock()

	Utils.MutexBloqueadoSuspendido.Lock()
	defer Utils.MutexBloqueadoSuspendido.Unlock()
	pcb.CambiarEstado(proceso, pcb.EstadoSuspBlocked)
	algoritmos.ColaBloqueadoSuspendido.Add(proceso)
	logger.Info("## (<%d>) Pasa del estado BLOCKED al estado SUSP.BLOCKED", proceso.PID)
	return true
}

func monitorBloqueado(procesoAsuspender Utils.BlockProcess) {
	pid := procesoAsuspender.PID
	logger.Info("Arrancó TIMER para PID <%d>", pid)
	duración := time.Duration(globals.KConfig.SuspensionTime) * time.Millisecond
	timer := time.NewTimer(duración)
	defer timer.Stop()

	select {
	case ioEvt := <-Utils.NotificarFinIO:
		if ioEvt.PID == pid {
			// fin de IO antes del timeout → READY
			moverDeBlockedAReady(ioEvt)
		}
	case <-timer.C:
		// timeout expired → SUSP.BLOCKED
		if moverDeBlockedASuspBlocked(pid) {
			logger.Info("PID <%d> → SUSP.BLOCKED (timeout)", pid)

			if err := comunicacion.SolicitarSuspensionEnMemoria(pid); err == nil {
				// ahora esperamos la confirmación DESDE MEMORIA VIA ENDPOINT...
				go func() {
					//aca te manda el PID aunque no es necesario ya que solo funciona como semaforo
					<-Utils.NotificarTimeoutBlocked
					// señal al largo plazo para que intente replanificar
					Utils.InitProcess <- struct{}{}
				}()
			}
			//ACA LE AVISO a la  funcion AtenderSuspBlockedAFinIO() PARA QUE PASE EL PROCESO A SUSP.READY ANQUE EN REALIDAD LO PASO YO
			Utils.FinIODesdeSuspBlocked <- Utils.IOEvent{
				PID:    pid,
				Nombre: procesoAsuspender.Nombre,
			}
		}
	}
}

// pasa de SUSP BLOQUEADO A SUSP READY
func AtenderSuspBlockedAFinIO() {
	for {
		// 1) ACA ESPERO LA PRIMERA SEÑAL QUE EL PROCESO YA ESTA EN SUSPENDIDO BLOQUEADO
		ev := <-Utils.FinIODesdeSuspBlocked

		// EN ESTE HILO ESPERO LA SEÑAL DE FIN DE IO QUE
		//INDICA QUE EL PROCESO QUE ESTA EN SUSPENDIDO BLOQUEADO PASE A SUSPENDIDO READY
		go func(ev Utils.IOEvent) {
			// 2) ahora esperamos la señal real de fin de IO para ese PID
			for {
				ioFin := <-Utils.NotificarFinIO
				if ioFin.PID == ev.PID {
					break
				}
				// si no coincide, se devolveee al canal para que
				// otros handlers puedan leerla
				go func(other Utils.IODesconexion) {
					Utils.NotificarFinIO <- other
				}(ioFin)
			}

			// ahora sí, sacamos de SUSP.BLOCKED y pasamos a SUSP.READY
			Utils.MutexBloqueadoSuspendido.Lock()
			var p *pcb.PCB
			for _, x := range algoritmos.ColaBloqueadoSuspendido.Values() {
				if x.PID == ev.PID {
					p = x
					break
				}
			}
			if p != nil {
				algoritmos.ColaBloqueadoSuspendido.Remove(p)
				pcb.CambiarEstado(p, pcb.EstadoSuspReady)

				Utils.MutexSuspendidoReady.Lock()
				algoritmos.ColaSuspendidoReady.Add(p)
				Utils.MutexSuspendidoReady.Unlock()

				logger.Info("## (<%d>) Pasa de SUSP.BLOCKED a SUSP.READY", p.PID)
				// notificamos al planificador de largo para que lo intente meter YA QUE SE ACTIVA CADA VEZ QUE UN PROCESO LLEGA A NEW O
				//SUSPENDIDO READY
				Utils.InitProcess <- struct{}{}
			} else {
				logger.Warn("AtenderSuspBlockedAFinIO: PID %d no estaba en SUSP.BLOCKED", ev.PID)
			}
			Utils.MutexBloqueadoSuspendido.Unlock()
		}(ev)
	}
}
