package Utils

import "sync"

// //VARIABLES PARA PLANIFICACION
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
}
