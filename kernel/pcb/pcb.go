package pcb

import (


	"github.com/sisoputnfrba/tp-golang/kernel/utils"
)

// posibles estados de un proceso
const (
	EstadoNew         = "new"
	EstadoReady       = "ready"
	EstadoExecute     = "execute"
	EstadoBlocked     = "blocked"
	EstadoExit        = "exit"
	EstadoSuspBlocked = "suspblocked"
	EstadoSuspReady   = "suspready"
)

type PCB struct {
	PID int
	PC  int
	ME  map[string]int //asocia cada estado con la cantidad de veces que el proceso estuvo en ese estado.
	MT  map[string]int //asocia cada estado con el tiempo total que el proceso pasó en ese estado.
}

func (a *PCB) Null() *PCB {
	return nil
}

func (a *PCB) Equal(b *PCB) bool {
	return a.PID == b.PID
}

var ColaNuevo utils.Queue[*PCB]
var NewStateQueue utils.Queue[*PCB]
var ColaBLoqueado utils.Queue[*PCB]
var ColaSalida utils.Queue[*PCB]
var ColaEjecutando utils.Queue[*PCB]
var ColaReady utils.Queue[*PCB]
var ColaBloqueadoSuspendido  utils.Queue[*PCB]
var ColaSuspendidoReady utils.Queue[*PCB]

//Ej ME: "ready": 3 → el proceso estuvo 3 veces en el estado listo.
//Ej MT: "execute": 12 → el proceso estuvo 12 unidades de tiempo en ejecución.
