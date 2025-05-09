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
