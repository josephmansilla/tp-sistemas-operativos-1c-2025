package algoritmos

import (
	"github.com/sisoputnfrba/tp-golang/kernel/pcb"
)

// definido por archivo de configuración:
// - rafaga inicial estimada
// - alpha
// Criterio: Se elegirá el proceso que tenga la rafaga más corta.
func SeleccionarSJF() *pcb.PCB {
	if len(ColaReady.elements) == 0 {
		return nil
	}

	masChico := ColaReady.elements[0] //Tomo el primero y empiezo a comparar rafagas
	for _, p := range ColaReady.elements {
		if p.EstimadoRafaga < masChico.EstimadoRafaga {
			masChico = p
		}
	}
	return masChico
}
