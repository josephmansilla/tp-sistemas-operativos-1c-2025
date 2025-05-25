package algoritmos

import (
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/pcb"
)

func SeleccionarSJF() *pcb.PCB {
	for id, cpu := range ColaReady.elements {
		primero = ColaReady.First().EstimadoRafaga
		if primero.rafaga > segundo.rafaga {
			primero = segundo
		}
	}
}
