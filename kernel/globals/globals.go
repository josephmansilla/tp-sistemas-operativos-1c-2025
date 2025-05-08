package globals

import (
	"sync"

	"github.com/sisoputnfrba/tp-golang/kernel/algoritmos"
	"github.com/sisoputnfrba/tp-golang/kernel/pcb"
)

// Datos recibidos por el Kernel
type DatosIO struct {
	Nombre string
	Ip     string
	Puerto int
}

type DatosCPU struct {
	ID     string
	Ip     string
	Puerto int
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

// ESTAS SON VARIABLES GLOBALES OJO¡¡¡¡
var ColaNuevo algoritmos.Queue[*pcb.PCB]
var NewStateQueue algoritmos.Queue[*pcb.PCB]
var ColaBLoqueado algoritmos.Queue[*pcb.PCB]
var ColaSalida algoritmos.Queue[*pcb.PCB]
var ColaEjecutando algoritmos.Queue[*pcb.PCB]
var ColaReady algoritmos.Queue[*pcb.PCB]
var ColaBloqueadoSuspendido algoritmos.Queue[*pcb.PCB]
var ColaSuspendidoReady algoritmos.Queue[*pcb.PCB]

//crear nodos a punteros PCB, para instanciarlas en main (sera un puntero a pcb creado previamente)
