package Utils

import (
	"sync"
)

var (
	// Mutex para coordinar creaciones concurrentes
	MutexPuedoCrearProceso *sync.Mutex

	// Mutex por cada cola
	MutexNuevo               sync.Mutex
	MutexReady               sync.Mutex
	MutexBloqueado           sync.Mutex
	MutexSalida              sync.Mutex
	MutexEjecutando          sync.Mutex
	MutexBloqueadoSuspendido sync.Mutex
	MutexSuspendidoReady     sync.Mutex
	MutexPedidosIO           sync.Mutex

	//Canales de señalización
	ChannelProcessArguments chan NewProcess
	InitProcess             chan struct{}
	LiberarMemoria          chan struct{}
	SemProcessCreateOK      chan struct{}
	ChannelFinishprocess    chan FinishProcess
	ChannelProcessBlocked   chan BlockProcess

	//AVISAR AL DESPACHADOR CUANDO UN PROCESO CAMBIA SU ESTADO
	NotificarDespachador    chan int              //PASA A READY
	NotificarComienzoIO     chan MensajeIOChannel //PASA A BLOQUEADO
	NotificarFinIO          chan IOEvent          //FIN DE IO
	NotificarIOLibre        chan IOEvent
	NotificarDesconexion    chan IODesconexion    //Desconexion DE IO
	ContextoInterrupcion    chan InterruptProcess //FIN DE EXECUTE
	NotificarTimeoutBlocked chan int
	FinIODesdeSuspBlocked   chan IOEvent
)

// InicializarMutexes deja listas las variables de mutex.
// Solo MutexPuedoCrearProceso requiere puntero, el resto ya
// está listo con su valor cero.
func InicializarMutexes() {
	MutexPuedoCrearProceso = &sync.Mutex{}
	// MutexNuevo, MutexReady, ... ya funcionan sin más
}

// InicializarCanales crea y configura los canales con buffers adecuados.
func InicializarCanales() {
	ChannelProcessArguments = make(chan NewProcess, 10) // buffer para hasta 10 peticiones
	ChannelFinishprocess = make(chan FinishProcess, 5)
	InitProcess = make(chan struct{})           // sin buffer para sincronización exacta
	SemProcessCreateOK = make(chan struct{}, 1) // semáforo de 1 slot
	LiberarMemoria = make(chan struct{}, 1)

	NotificarDespachador = make(chan int, 10) // buffer 10 procesos listos
	NotificarComienzoIO = make(chan MensajeIOChannel, 10)
	NotificarFinIO = make(chan IOEvent, 10)
	NotificarIOLibre = make(chan IOEvent, 10)
	NotificarDesconexion = make(chan IODesconexion, 10)
	ContextoInterrupcion = make(chan InterruptProcess, 10)
	ChannelProcessBlocked = make(chan BlockProcess, 10)
	NotificarTimeoutBlocked = make(chan int)
	FinIODesdeSuspBlocked = make(chan IOEvent, 0)
}

type MensajeIOChannel struct {
	PID      int
	PC       int
	Nombre   string
	Duracion int
	CpuID    string
}

type IODesconexion struct {
	Nombre string
	PID    int
	Puerto int
}
type FinishProcess struct {
	PID   int
	PC    int
	CpuID string
}
type InterruptProcess struct {
	PID    int
	PC     int
	CpuID  string
	Motivo string
}
type BlockProcess struct {
	PID      int
	PC       int
	Nombre   string
	Duracion int
	CpuID    string
}
type NewProcess struct {
	Filename string
	Tamanio  int
	PID      int
}
type IOEvent struct {
	PID    int
	Nombre string
	Puerto int
}
