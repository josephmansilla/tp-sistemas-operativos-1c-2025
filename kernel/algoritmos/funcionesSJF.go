package algoritmos

import (
	"github.com/sisoputnfrba/tp-golang/kernel/pcb"
)

func SeleccionarSJF() *pcb.PCB {
	/*
		primero = ColaReady.First().EstimadoRafaga
		for id, cpu := range ColaReady.elements {
			if primero.rafaga > segundo.rafaga {
				primero = segundo
			}
		}
	*/
	pcb := ColaReady.First()
	return pcb
}
