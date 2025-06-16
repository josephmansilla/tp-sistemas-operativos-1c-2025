package globals

var MemoryConfig *Config
var MemoriaPrincipal []byte         // MP simulada
var FramesLibres []bool             //los frames van a estar en True si están libres
var CantidadFramesLibres int        // simplemente recuenta la cantidad de frames
var ProcesosPorPID map[int]*Proceso // guardo procesos con los PID
var SwapIndex map[int]*SwapProcesoInfo

func InstanciarEstructurasGlobales() {
	ProcesosPorPID = make(map[int]*Proceso)
	SwapIndex = make(map[int]*SwapProcesoInfo)
}

type EstadoMemoria struct {
	CantidadFramesLibres int
	CantidadBytesUsados  int
	CantidadBytesTotales int
} //TODO: ver si dejar
