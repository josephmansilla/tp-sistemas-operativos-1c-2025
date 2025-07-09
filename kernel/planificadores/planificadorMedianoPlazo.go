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

// ManejadorMedianoPlazo se quedará escuchando
// bloqueos e iniciará un timer por cada PID que llegue.
func ManejadorMedianoPlazo() {
	for {
		procesoBloqueado := <-Utils.ChannelProcessBlocked // señal llegada de BLOCKED con el PID del proceso para que arranque su TIMER
		go monitorBloqueado(procesoBloqueado)
	}
}

//Al ingresar un proceso al estado BLOCKED se deberá iniciar un timer
//el cual se encargará de esperar un tiempo determinado por archivo de configuración,

func monitorBloqueado(proceso Utils.BlockProcess) {
	logger.Info("Arrancó el TIMER del proceso: PID <%d>", proceso.PID)
	suspension := time.Duration(globals.KConfig.SuspensionTime) * time.Millisecond

	timer := time.NewTimer(suspension)
	defer timer.Stop()

	pid := proceso.PID

	//Al terminar ese tiempo, si el proceso continúa en estado BLOCKED,
	//deberá transicionar al estado SUSP. BLOCKED.

	select {
	//Cuando un dispositivo termina la IO pedida para un Proceso
	case ioFin := <-Utils.NotificarFinIO:
		if ioFin.PID == pid {
			// llegó fin de IO antes del timeout: pasa a READY directo
			if moverDeBlockedAReady(ioFin) {
				logger.Info("PID <%d> terminó IO antes del timeout → READY", pid)
			}
		}
	case <-timer.C:
		// el timer expiró y sigue en BLOCKED → pasa a SUSP.BLOCKED
		//En este momento se debe informar al módulo memoria que debe ser movido de
		//memoria principal a swap. Cabe aclarar que en este momento vamos a tener
		//más memoria libre en el sistema por lo que se debe verificar si uno o
		//más nuevos procesos pueden entrar (tanto de la cola NEW como de SUSP. READY).

		if moverDeBlockedASuspBlocked(pid) {
			logger.Info("PID <%d> timeout expired → SUSP.BLOCKED", pid)
			// avisar a Memoria que lo saque a swap
			// TODO comunicacion.SolicitarSuspenderEnMemoria(pid)
			// señal al largo plazo para reintentar NEW/SUSP.READY
			Utils.InitProcess <- struct{}{}
			//aca meter un signal/tuberia al endpoint que atiende los fin de IO para que pueda pasar el proceso a susp.ready
			<-Utils.NotificarTimeoutBlocked
		}
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
