package globals

var RespuestaKernel DatosRespuestaDeKernel

// var Kernel DatosConsultaDeKernel

type DatosConsultaDeKernel struct {
	PID            int `json:"pid"`
	TamanioMemoria int `json:"tamanio_memoria"`
}

type DatosRespuestaDeKernel struct {
	PID            int    `json:"pid"`
	TamanioMemoria int    `json:"tamanio_memoria"`
	Pseudocodigo   string `json:"filename"`
}

type RespuestaMemoria struct {
	Exito   bool   `json:"exito"` // TODO: cambiar a error
	Mensaje string `json:"mensaje"`
}

type RespuestaEspacioLibre struct {
	EspacioLibre int `json:"espacio_libre"`
}

type DatosFinalizacionProceso struct {
	PID int `json:"pid"`
}

type LecturaMemoria struct {
	PID              int `json:"pid"`
	DireccionFisica  int `json:"direccion_fisica"`
	TamanioARecorrer int `json:"tamanio_a_recorrer"`
}

type ExitoLecturaMemoria struct {
	Exito        error  `json:"exito"`
	DatosAEnviar string `json:"datos_a_enviar"`
}

type EscrituraMemoria struct {
	PID              int `json:"pid"`
	DireccionFisica  int `json:"direccion_fisica"`
	TamanioARecorrer int `json:"tamanio_a_recorrer"`
}

type ExitoEscrituraMemoria struct {
	Exito    error `json:"exito"`
	Booleano bool  `json:"booleano"`
}
