package globals

import "sync"

var MutexProcesosPorPID sync.Mutex
var MutexMemoriaPrincipal sync.Mutex
var MutexCantidadFramesLibres sync.Mutex
var MutexEstructuraFramesLibres sync.Mutex
var MutexFrameOcupadoPorPID sync.Mutex
