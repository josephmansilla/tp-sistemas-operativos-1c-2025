package globals

import (
	"sync"
)

// Datos recibidos por el Kernel
type DatosIO struct {
	Nombre string
	Ip     string
	Puerto int
}

type DatosCPU struct {
	ID      string
	Ip      string
	Puerto  int
	Ocupada bool
}

type EspacioLibreRTA struct {
	EspacioLibre int `json:"espacio_libre"`
}

var KConfig *KernelConfig

var CPU DatosCPU
var CPUs map[string]DatosCPU = make(map[string]DatosCPU) // clave: ID del CPU
var CPUMu sync.Mutex
var CPUCond = sync.NewCond(&CPUMu)

var IO DatosIO
var IOs map[string]DatosIO = make(map[string]DatosIO) // clave: nombre del IO
var IOMu sync.Mutex
var IOCond = sync.NewCond(&IOMu)

var EspacioLibreProceso EspacioLibreRTA
var UltimoPID int = 0
var PidMutex sync.Mutex

func GenerarNuevoPID() int {
	PidMutex.Lock()
	defer PidMutex.Unlock()

	UltimoPID++
	return UltimoPID
}

//crear nodos a punteros PCB, para instanciarlas en main (sera un puntero a pcb creado previamente)
