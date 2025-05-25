package algoritmos

import (
	"github.com/sisoputnfrba/tp-golang/kernel/pcb"
	"sync"
)

type Cola[T Nulleable[T]] struct {
	elements []T
	mutex    sync.Mutex
	Priority int
}

// ESTAS SON VARIABLES GLOBALES OJO¡¡¡¡
var ColaNuevo Cola[*pcb.PCB]
var NewStateQueue Cola[*pcb.PCB]
var ColaBloqueado Cola[*pcb.PCB]
var ColaSalida Cola[*pcb.PCB]
var ColaEjecutando Cola[*pcb.PCB]
var ColaReady Cola[*pcb.PCB]
var ColaBloqueadoSuspendido Cola[*pcb.PCB]
var ColaSuspendidoReady Cola[*pcb.PCB]
