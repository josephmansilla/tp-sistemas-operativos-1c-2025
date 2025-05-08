package planificadores

import (
	"github.com/sisoputnfrba/tp-golang/kernel/comunicacion"
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/pcb"
)

func EjecutarProceso() {
	pcb := new(pcb.PCB)

	id := "1"

	comunicacion.EnviarContextoCPU(globals.CPUs[id].ID, pcb)
}
