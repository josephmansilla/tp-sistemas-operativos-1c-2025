package algoritmos

import (
	"github.com/sisoputnfrba/tp-golang/kernel/Utils"
	"github.com/sisoputnfrba/tp-golang/kernel/pcb"
)

// ready_ingress_algorithm
// por ahora solo funciona con cola NEW
// Criterio: mas chico el que menos memoria solicite (size).
func AddPMCP(p *pcb.PCB) {
	Utils.MutexNuevo.Lock()
	defer Utils.MutexNuevo.Unlock()

	insertado := false
	nuevaCola := make([]*pcb.PCB, 0, len(ColaNuevo.elements)+1)

	for _, actual := range ColaNuevo.elements {
		if !insertado && p.ProcessSize < actual.ProcessSize {
			nuevaCola = append(nuevaCola, p)
			insertado = true
		}
		nuevaCola = append(nuevaCola, actual)
	}

	if !insertado {
		nuevaCola = append(nuevaCola, p)
	}

	ColaNuevo.elements = nuevaCola
}

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
