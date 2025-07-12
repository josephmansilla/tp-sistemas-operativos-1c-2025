package estructuras

import "sync"

var MutexProcesosPorPID sync.Mutex
var MutexMemoriaPrincipal sync.Mutex
var MutexCantidadFramesLibres sync.Mutex
var MutexEstructuraFramesLibres sync.Mutex
var MutexMetrica []sync.Mutex
var MutexDump sync.Mutex
var MutexSwapIndex sync.Mutex
var MutexSwapBool sync.Mutex
