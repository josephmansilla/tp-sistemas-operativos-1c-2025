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

	// Canales de señalización
	ChannelProcessArguments chan []string
	ChannelFinishprocess    chan int
	InitProcess             chan struct{}
	SemProcessCreateOK      chan struct{}
	ChannelFinishProcess2   chan bool

	//AVISAR CUANDO UN PROCESO LLEGA A READY
	NotificarProcesoReady chan int
	NotificarComienzoIO   chan MensajeIOChannel
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
	ChannelFinishprocess = make(chan int, 5)
	InitProcess = make(chan struct{})           // sin buffer para sincronización exacta
	SemProcessCreateOK = make(chan struct{}, 1) // semáforo de 1 slot
	ChannelFinishProcess2 = make(chan bool, 5)

	NotificarProcesoReady = make(chan int, 10) // buffer 10 procesos listos
	NotificarComienzoIO = make(chan MensajeIOChannel, 10)
}

type MensajeIOChannel struct {
	PID      int
	Nombre   string
	Duracion int
}
