package algoritmos

import (
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/pcb"
)

//definida por archivo de configuración:
//rafaga inicial estimada
//alpha

func SeleccionarSJF() *pcb.PCB {
	var masChico = ColaReady.First()
	for _, p := range ColaReady.elements {
		if p.EstimadoRafaga < masChico.EstimadoRafaga {
			masChico = p
		}
	}
	return masChico
}

// Utilizar despues de una rafaga en CPU
func ActualizarEstimacionRafaga(proceso *pcb.PCB, rafagaReal int) {
	alpha := globals.Config.Alpha
	proceso.EstimadoRafaga = alpha*float64(rafagaReal) + (1-alpha)*proceso.EstimadoRafaga
}

//EJEMPLO DE USO
/*
cuando termina una ráfaga
ActualizarEstimacionRafaga(proceso, 7) // 7 es el tiempo real que tardó la ráfaga
*/
