package algoritmos

import (
	"github.com/sisoputnfrba/tp-golang/kernel/Utils"
	"github.com/sisoputnfrba/tp-golang/kernel/comunicacion"
	"github.com/sisoputnfrba/tp-golang/kernel/pcb"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"time"
)

/*
Si no hay CPUs libres, se debe evaluar si dicho proceso
tiene una rafaga más corta que los que se encuentran en ejecución.
*/
func Desalojo(procesoEntrante *pcb.PCB) {
	Utils.MutexEjecutando.Lock()
	defer Utils.MutexEjecutando.Unlock()

	tiempoEntrante := procesoEntrante.EstimadoRafaga

	var procesoAInterrumpir *pcb.PCB
	var cpuAInterrumpir string
	var mayorTiempoRestante float64 = -1

	/*
		se debe informar a la CPU que posea al Proceso con el tiempo restante
		MAS ALTO que debe desalojar, para que pueda ser planificado el nuevo.
	*/

	for _, p := range ColaEjecutando.Values() {
		duracion := time.Since(p.TiempoEstado)
		tiempoEjecutado := float64(duracion.Milliseconds())
		tiempoRestante := p.EstimadoRafaga - tiempoEjecutado

		// Queremos interrumpir al que tenga MAYOR tiempo restante
		if tiempoRestante > mayorTiempoRestante {
			mayorTiempoRestante = tiempoRestante
			procesoAInterrumpir = p
			cpuAInterrumpir = p.CpuID
		}
	}

	//COMPARAR TIEMPO RESTANTE CON LA RAFAGA ENTRANTE
	if procesoAInterrumpir != nil && tiempoEntrante < mayorTiempoRestante {
		logger.Info("## (<%d>) - Desalojado por algoritmo SJF/SRT", procesoEntrante.PID)
		logger.Debug("SRT: Proceso <%d> interrumpe a <%d> en CPU <%s> (%.2f < %.2f)",
			procesoEntrante.PID, procesoAInterrumpir.PID, cpuAInterrumpir, tiempoEntrante, mayorTiempoRestante)
		comunicacion.AvisarDesalojoCPU(cpuAInterrumpir, procesoAInterrumpir)
	} else {
		logger.Info("SRT: Proceso <%d> NO tiene menor tiempo restante que los procesos ejecutando (Tiempo entrante: %.2f < Mayor Restante: %.2f)",
			procesoEntrante.PID, tiempoEntrante, mayorTiempoRestante)
	}
}
