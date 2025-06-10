package globals

var MemoryConfig *Config
var MemoriaPrincipal [][]byte        // MP simulada
var FramesLibres []bool              //los frames van a estar en True si están libres
var CantidadFramesLibres int         // simplemente recuenta la cantidad de frames
var FrameOcupadoPor map[int]Ocupante // guardo Ocupantes con los frame de indice
var ProcesosPorPID map[int]*Proceso  // guardo procesos con los PID

var CantidadNiveles int
var EntradasPorPagina int
var DelayMemoria int
var DelaySwap int
var TamanioMaximoFrame int
