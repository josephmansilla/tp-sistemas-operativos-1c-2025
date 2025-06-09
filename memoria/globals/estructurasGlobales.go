package globals

var MemoryConfig *Config
var MemoriaPrincipal [][]byte // MP simulada
var FramesLibres []bool       //los frames van a estar en True si están libres
var ProcesosMapeable ProcesosMap

var CantidadNiveles int = MemoryConfig.NumberOfLevels
var EntradasPorPagina int = MemoryConfig.EntriesPerPage
var DelayMemoria int = MemoryConfig.MemoryDelay
var DelaySwap int = MemoryConfig.SwapDelay
var TamanioMaximoFrame int

// SUPER PENDIENTES
type ArgmentosCreacionProceso struct {
	NombrePseudocodigo string `json:"nombre_pseudocodigo"`
	TamanioProceso     int    `json:"tamanioProceso"`
	// PID
}

type PedidoAMemoria struct {
	Thread    Thread                 `json:"thread"`
	Type      string                 `json:"type"`
	Arguments map[string]interface{} `json:"arguments"`
}

type Thread struct {
	PID int `json:"pid"`
}
