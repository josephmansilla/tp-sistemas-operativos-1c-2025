package planificadores

import (
	"github.com/sisoputnfrba/tp-golang/kernel/Utils"
	"github.com/sisoputnfrba/tp-golang/kernel/algoritmos"
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/pcb"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"time"
)

func PlanificadorMedianoPlazo() {
	logger.Info("Iniciando el planificador de Mediano Plazo")
	go ManejadorMedianoPlazo()
}

// ManejadorMedianoPlazo se quedará escuchando bloqueos e
// iniciará un timer por cada PID que llegue.
func ManejadorMedianoPlazo() {
	for {
		pid := <-Utils.ChannelProcessBlocked // señal llegada de BLOCKED con el PID del proceso para que arranque su TIMER
		go monitorBloqueado(pid)
	}
}

// ACA el proceso ya esta en la lista bloqueado y  arranca un timer y observa si recibe fin de IO
func monitorBloqueado(pid int) {
	logger.Info("Arrancó el TIMER del proceso: PID <%d>", pid)
	suspensión := time.Duration(globals.KConfig.SuspensionTime) * time.Millisecond

	timer := time.NewTimer(suspensión)
	defer timer.Stop()

	select {
	// ESte Case vendria desde el Endpoint que atiende a IO cuando  un dispisitivo termino la IO pedida para un Proceso
	case pidIO := <-Utils.NotificarFinIO:
		if pidIO == pid {
			// llegó fin de IO antes del timeout: pasa a READY directo
			if moverDeBlockedAReady(pid) {
				logger.Info("PID=%d: IO terminó antes del timeout → READY", pid)
			}
		}
	case <-timer.C:
		// el timer expiró y sigue en BLOCKED → pasa a SUSP.BLOCKED
		if moverDeBlockedASuspBlocked(pid) {
			logger.Info("PID=%d: timeout expired → SUSP.BLOCKED", pid)
			// avisar a Memoria que lo saque a swap
			// TODO comunicacion.SolicitarSuspenderEnMemoria(pid)
			// señal al largo plazo para reintentar NEW/SUSP.READY
			Utils.InitProcess <- struct{}{}
		}
	}
}

// moverDeBlockedAReady quita de BLOCKED y encola en READY
func moverDeBlockedAReady(pid int) bool {
	Utils.MutexBloqueado.Lock()
	defer Utils.MutexBloqueado.Unlock()

	p := BuscarPCBPorPID(pid) // busca en ColaBloqueado
	if p == nil {
		return false
	}
	algoritmos.ColaBloqueado.Remove(p)

	Utils.MutexReady.Lock()
	defer Utils.MutexReady.Unlock()
	pcb.CambiarEstado(p, pcb.EstadoReady)
	algoritmos.ColaReady.Add(p)
	// señal al corto plazo para despachar
	Utils.NotificarDespachador <- pid
	return true
}

// moverDeBlockedASuspBlocked quita de BLOCKED y encola en SUSP.BLOCKED
func moverDeBlockedASuspBlocked(pid int) bool {
	Utils.MutexBloqueado.Lock()
	p := BuscarPCBPorPID(pid)
	if p == nil {
		Utils.MutexBloqueado.Unlock()
		return false
	}
	algoritmos.ColaBloqueado.Remove(p)
	Utils.MutexBloqueado.Unlock()

	Utils.MutexBloqueadoSuspendido.Lock()
	defer Utils.MutexBloqueadoSuspendido.Unlock()
	pcb.CambiarEstado(p, pcb.EstadoSuspBlocked)
	algoritmos.ColaBloqueadoSuspendido.Add(p)
	return true
}
