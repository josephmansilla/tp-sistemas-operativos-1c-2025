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
}

// ManejadorMedianoPlazo se quedará escuchando bloqueos e
// iniciará un timer por cada PID que llegue.
func ManejadorMedianoPlazo() {
	for {
		procesoBloqueado := <-Utils.ChannelProcessBlocked // señal llegada de BLOCKED con el PID del proceso para que arranque su TIMER

		//EL PEDIDO SE ACUMULA EN LA LISTA FIFO
		io1
		102
		103

		//TODOS SE LLAMAN TECLADO => PERO EXISTEN TECLADO1 TECLADO2 Y TECLADO3

		go monitorBloqueado(procesoBloqueado)
	}
}

//Al ingresar un proceso al estado BLOCKED se deberá iniciar un timer
//el cual se encargará de esperar un tiempo determinado por archivo de configuración,
//al terminar ese tiempo si el proceso continúa en estado BLOCKED,
//él mismo deberá transicionar al estado SUSP. BLOCKED.
//En este momento se debe informar al módulo memoria que debe ser movido de
//memoria principal a swap. Cabe aclarar que en este momento vamos a tener
//más memoria libre en el sistema por lo que se debe verificar si uno o
//más nuevos procesos pueden entrar (tanto de la cola NEW como de SUSP. READY).

// ACA el proceso ya esta en la lista bloqueado y  arranca un timer y observa si recibe fin de IO
func monitorBloqueado(proceso Utils.BlockProcess) {
	logger.Info("Arrancó el TIMER del proceso: PID <%d>", proceso.PID)
	suspension := time.Duration(globals.KConfig.SuspensionTime) * time.Millisecond

	//Enviar al módulo IO (usando los datos del mensaje recibido)
	enviarAIO(proceso)
	//busca si disco libre...
	//asigna el pedido a io

	//ACITVA IO 5... WAIT 5 SEGUNDOS MANDA FIN

	comunicacion.EnviarContextoIO(proceso.Nombre, pid, msg.Duracion)

	timer := time.NewTimer(suspension)
	defer timer.Stop()

	select {
	//Cuando un dispositivo termina la IO pedida para un Proceso
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
			//aca meter un signal/tuberia al endpoint que atiende los fin de IO para que pueda pasar el proceso a susp.ready
			<-Utils.NotificarTimeoutBlocked
			//
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
