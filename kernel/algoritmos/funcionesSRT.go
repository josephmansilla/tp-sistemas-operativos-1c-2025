package algoritmos

import (
	"github.com/sisoputnfrba/tp-golang/kernel/Utils"
	"github.com/sisoputnfrba/tp-golang/kernel/comunicacion"
	"github.com/sisoputnfrba/tp-golang/kernel/pcb"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"time"
)

/*
SJF con Desalojo
Funciona igual que el anterior con la variante que,
al ingresar un proceso en la cola de Ready y no haber CPUs libres,
se debe evaluar si dicho proceso tiene una rafaga
más corta que los que se encuentran en ejecución.
*/
func Desalojo(procesoEntrante *pcb.PCB) {
	tiempoEntrante := procesoEntrante.EstimadoRafaga
	logger.Info("SRT: Evaluando posible desalojo por llegada de <%d> con ráfaga estimada %.0f ms", procesoEntrante.PID, tiempoEntrante)

	Utils.MutexEjecutando.Lock()
	defer Utils.MutexEjecutando.Unlock()

	if len(ColaEjecutando.Values()) == 0 {
		logger.Error("SRT: No hay procesos en ejecución")
		return
	}

	var procesoAInterrumpir *pcb.PCB
	var cpuAInterrumpir string
	var mayorTiempoRestante float64 = 0

	for _, p := range ColaEjecutando.Values() {
		// Queremos interrumpir al que tenga MAYOR tiempo restante
		tiempoEjecutado := float64(time.Since(p.TiempoEstado).Milliseconds())
		tiempoRestante := p.EstimadoRafaga - tiempoEjecutado
		logger.Debug("SRT: PID <%d> - Ejecutado: %.0f ms - Restante: %.0f ms", p.PID, tiempoEjecutado, tiempoRestante)

		if tiempoRestante > mayorTiempoRestante {
			mayorTiempoRestante = tiempoRestante
			procesoAInterrumpir = p
			cpuAInterrumpir = p.CpuID
		}
	}

	if procesoAInterrumpir == nil {
		logger.Error("SRT: No se encontró proceso con mayor tiempo restante para interrumpir.")
		return
	}

	/*
		Se debe informar a la CPU que posea al Proceso con el tiempo restante
		MAS ALTO que debe desalojar, para que pueda ser planificado el nuevo.
	*/

	if tiempoEntrante < mayorTiempoRestante {
		logger.Info("## (<%d>) - Desalojado por algoritmo SJF/SRT", procesoEntrante.PID)
		logger.Debug("SRT: Proceso <%d> interrumpe a <%d> en CPU <%s> (%.2f < %.2f)",
			procesoEntrante.PID, procesoAInterrumpir.PID, cpuAInterrumpir, tiempoEntrante, mayorTiempoRestante)
		comunicacion.AvisarDesalojoCPU(cpuAInterrumpir, procesoAInterrumpir)
	} else {
		logger.Info("SRT: Proceso <%d> NO tiene menor tiempo restante que los procesos ejecutando (%.2f > %.2f)",
			procesoEntrante.PID, tiempoEntrante, mayorTiempoRestante)
	}
}
