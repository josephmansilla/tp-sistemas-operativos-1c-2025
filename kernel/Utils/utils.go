package Utils

import (
	"sync"
)

// //VARIABLES PARA PLANIFICACION NOOO ELIMINAR NADA TODAVIA
/*
var ChannelProcessArguments chan []string
var ChannelFinishprocess chan int
var InitProcess chan struct{}
var SemProcessCreateOK chan struct{}

// Canal que usa quien espera confirmación de finalización de proceso (opcional)
var ChannelFinishProcess2 chan bool

// Mutex que bloquea la creación simultánea de procesos
var MutexPuedoCrearProceso = &sync.Mutex{}

// Mutex por cada cola
var (
	MutexNuevo               sync.Mutex
	MutexReady               sync.Mutex
	MutexBloqueado           sync.Mutex
	MutexSalida              sync.Mutex
	MutexEjecutando          sync.Mutex
	MutexBloqueadoSuspendido sync.Mutex
	MutexSuspendidoReady     sync.Mutex
)

func InicializarMutexes() {
	MutexPuedoCrearProceso = &sync.Mutex{}
	// Los demás ya están inicializados como valores por defecto de tipo `sync.Mutex`
}
func InicializarCanales() {
	ChannelProcessArguments = make(chan []string, 10) // Canal buffered (para enviar argumentos de creación)
	ChannelFinishprocess = make(chan int, 5)          // También puede tener buffer
	InitProcess = make(chan struct{})                 // Unbuffered para sincronización
	SemProcessCreateOK = make(chan struct{}, 1)       // Unbuffered, tipo semáforo

	ChannelFinishProcess2 = make(chan bool, 5) // Puede ser buffered si varios procesos notifican
}*/
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

	//Canales de señalización
	ChannelProcessArguments chan []string
	InitProcess             chan struct{}
	SemProcessCreateOK      chan struct{}
	ChannelFinishProcess2   chan bool
	ChannelFinishprocess    chan FinishProcess
	ChannelProcessBlocked   chan int

	//AVISAR AL DESPACHADOR CUANDO UN PROCESO CAMBIA SU ESTADO
	NotificarDespachador    chan int              //PASA A READY
	NotificarComienzoIO     chan MensajeIOChannel //PASA A BLOQUEADO
	NotificarFinIO          chan int              //FIN DE IO
	NotificarDesconexion    chan int              //Desconexion DE IO
	ContextoInterrupcion    chan InterruptProcess //FIN DE EXECUTE
	NotificarTimeoutBlocked chan struct{}
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
	ChannelProcessArguments = make(chan []string, 10) // buffer para hasta 10 peticiones
	ChannelFinishprocess = make(chan FinishProcess, 5)
	InitProcess = make(chan struct{})           // sin buffer para sincronización exacta
	SemProcessCreateOK = make(chan struct{}, 1) // semáforo de 1 slot
	ChannelFinishProcess2 = make(chan bool, 5)

	NotificarDespachador = make(chan int, 10) // buffer 10 procesos listos
	NotificarComienzoIO = make(chan MensajeIOChannel, 10)
	NotificarFinIO = make(chan int, 10)
	NotificarDesconexion = make(chan int, 10)
	ContextoInterrupcion = make(chan InterruptProcess, 10)
	ChannelProcessBlocked = make(chan int, 10)
	NotificarTimeoutBlocked = make(chan struct{}, 1)
}

type MensajeIOChannel struct {
	PID      int
	PC       int
	Nombre   string
	Duracion int
	CpuID    string
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
