package globals

var MemoryConfig *Config
var MemoriaPrincipal []byte         // MP simulada
var FramesLibres []bool             //los frames van a estar en True si están libres
var CantidadFramesLibres int        // simplemente recuenta la cantidad de frames
var ProcesosPorPID map[int]*Proceso // guardo procesos con los PID
var ProcesosSuspendidos map[int]*SwapProceso

//var PaginasEnSwap map[int][]EntradaSwap

func InstanciarEstructurasGlobales() {
	ProcesosPorPID = make(map[int]*Proceso)
	ProcesosSuspendidos = make(map[int]*SwapProceso)
	//	PaginasEnSwap = make(map[int][]EntradaSwap)
}

type EstadoMemoria struct {
	CantidadFramesLibres int
	CantidadBytesUsados  int
	CantidadBytesTotales int
} //TODO: ver si dejar
