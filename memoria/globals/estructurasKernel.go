package globals

var RespuestaKernel DatosRespuestaDeKernel

// var Kernel DatosConsultaDeKernel

type DatosConsultaDeKernel struct {
	PID            int `json:"pid"`
	TamanioMemoria int `json:"tamanio_memoria"`
	// con el tamaño de memoria consulta si es posible ejecutarlo en memoria
}

type DatosRespuestaDeKernel struct {
	Pseudocodigo   string `json:"filename"`
	TamanioMemoria int    `json:"tamanio_memoria"`
	PID            int    `json:"pid"`
}

type RespuestaMemoria struct {
	Exito   bool   `json:"exito"`
	Mensaje string `json:"mensaje"`
}

type EspacioLibreRTA struct {
	EspacioLibre int `json:"espacio_libre"`
}

type RespuestaEspacioLibre struct {
	MemoriaDisponible int `json:"memoria_disponible"`
	FramesLibres      int `json:"frames_libres"`
}
