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

var CPU DatosDeCPU
