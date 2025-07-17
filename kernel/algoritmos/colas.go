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

type PedidoIO struct {
	Nombre   string
	PID      int
	Duracion int
}

// Saber si son nulos
func (a *PedidoIO) Null() *PedidoIO {
	return nil
}

// Comparar pedidos
func (a *PedidoIO) Equal(b *PedidoIO) bool {
	return a.PID == b.PID && a.Duracion == b.Duracion && a.Nombre == b.Nombre
}

// ESTAS SON VARIABLES GLOBALES OJO¡¡¡¡
var ColaNuevo Cola[*pcb.PCB]
var ColaBloqueado Cola[*pcb.PCB]
var ColaSalida Cola[*pcb.PCB]
var ColaEjecutando Cola[*pcb.PCB]
var ColaReady Cola[*pcb.PCB]
var ColaBloqueadoSuspendido Cola[*pcb.PCB]
var ColaSuspendidoReady Cola[*pcb.PCB]
var PedidosIO Cola[*PedidoIO]
