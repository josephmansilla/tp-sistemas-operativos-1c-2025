package globals

import "sync"

// Datos recibidos por el Kernel
type DatosIO struct {
	Nombre string
	Ip     string
	Puerto int
}

type DatosCPU struct {
	Ip     string
	Puerto int
	ID     string
}

type EspacioLibreRTA struct {
	EspacioLibre int `json:"espacio_libre"`
}

var KConfig *KernelConfig

var CPU DatosCPU
var IO DatosIO
var EspacioLibreProceso EspacioLibreRTA
var IOMu sync.Mutex
var IOCond = sync.NewCond(&IOMu)
var UltimoPID int = 0
var PidMutex sync.Mutex

func GenerarNuevoPID() int {
	PidMutex.Lock()
	defer PidMutex.Unlock()

	UltimoPID++
	return UltimoPID
}

//crear nodos a punteros PCB, para instanciarlas en main (sera un puntero a pcb creado previamente)
