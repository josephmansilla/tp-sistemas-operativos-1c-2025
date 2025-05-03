package globals

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

var CPU DatosCPU
var IO DatosIO
var EspacioLibreProceso EspacioLibreRTA

//crear nodos a punteros PCB, para instanciarlas en main (sera un puntero a pcb creado previamente)
