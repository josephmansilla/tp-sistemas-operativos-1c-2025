package algoritmos

import "github.com/sisoputnfrba/tp-golang/kernel/pcb"

func SeleccionarSRT() *pcb.PCB {
	if len(ColaReady.elements) == 0 {
		return nil
	}

	// Empezar con el primero
	menor := ColaReady.elements[0]
	for _, p := range ColaReady.elements {
		if p.RafagaRestante < menor.RafagaRestante {
			menor = p
		}
	}
	return menor
}
