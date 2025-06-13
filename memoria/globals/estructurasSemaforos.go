package globals

import "sync"

var MutexProcesosPorPID sync.Mutex
var MutexMemoriaPrincipal sync.Mutex
var MutexCantidadFramesLibres sync.Mutex
var MutexEstructuraFramesLibres sync.Mutex
var MutexMetrica []sync.Mutex

func CambiarEstadoFrame(numeroFrame int) {
	MutexEstructuraFramesLibres.Lock()
	if FramesLibres[numeroFrame] == false {
		FramesLibres[numeroFrame] = true
	} else {
		FramesLibres[numeroFrame] = false
	}
	MutexEstructuraFramesLibres.Unlock()
}
