package globals

var CPU DatosDeCPU

type DatosDeCPU struct {
	PID int `json:"pid"`
	PC  int `json:"pc"`
}

type DatosParaCPU struct {
	// TODO
}

type ContextoDeCPU struct {
	PID int `json:"pid"`
	PC  int `json:"pc"`
}

type InstruccionCPU struct {
	Instruccion string `json:"instruccion"`
}

type ConsultaConfigMemoria struct {
	TamanioPagina    int `json:"tamanioPagina"`
	EntradasPorNivel int `json:"entradasPorNivel"`
	CantidadNiveles  int `json:"cantidadNiveles"`
}

type MensajePedidoTablaCPU struct {
	PID            int   `json:"pid"`
	IndicesEntrada []int `json:"indices_entrada"`
}

type RespuestaTablaCPU struct {
	NumeroMarco int `json:"numero_marco"`
}

// TODO: conexiones con cpu
type MensajeTabla struct {
	PID           int `json:"pid"`
	NumeroTabla   int `json:"numeroTabla"`
	EntradaIndice int `json:"entrada"`
}

type RespuestaTabla struct {
	EsUltimoNivel bool `json:"esUltimoNivel"`
	NumeroTabla   int  `json:"numeroTabla"`
	NumeroMarco   int  `json:"marco"`
}
