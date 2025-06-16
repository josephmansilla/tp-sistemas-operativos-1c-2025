package globals

import "sync"

var MutexProcesosPorPID sync.Mutex
var MutexMemoriaPrincipal sync.Mutex
var MutexCantidadFramesLibres sync.Mutex
var MutexEstructuraFramesLibres sync.Mutex
var MutexMetrica []sync.Mutex
var MutexDump sync.Mutex

// Podría ser un slice de Mutex por PID, es medio al pedo
// pero sería conceptualmente correcto

func InicializarSemaforos() {
	MutexMetrica = make([]sync.Mutex, MemoryConfig.MemorySize*1000) // tamaño totalmente arbitrario
	// MutexDump = make([]sync.Mutex, MemoryConfig.MemorySize)
}
