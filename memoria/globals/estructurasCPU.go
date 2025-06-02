package globals

// Tipo de datos recibidos de la CPU

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

var CPU DatosDeCPU
