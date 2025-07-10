package globals

var CPU DatosDeCPU

type DatosDeCPU struct {
	PID int `json:"pid"`
	PC  int `json:"pc"`
}

type ContextoDeCPU struct {
	PID int `json:"pid"`
	PC  int `json:"pc"`
}

type InstruccionCPU struct {
	Exito       error  `json:"exito"`
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

type EscrituraProceso struct {
	PID             int    `json:"pid"`
	DireccionFisica int    `json:"direccion_fisica"`
	DatosAEscribir  string `json:"datos_a_escribir"`
}

type LecturaProceso struct {
	PID              int `json:"pid"`
	DireccionFisica  int `json:"direccion_fisica"`
	TamanioARecorrer int `json:"tamanio_a_recorrer"`
}

type ExitoLecturaMemoria struct {
	Exito      error  `json:"exito"`
	ValorLeido string `json:"valor_leido"`
}
