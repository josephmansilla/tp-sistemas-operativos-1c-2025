package globals

type DatosConsultaDeKernel struct {
	PID            int `json:"pid"`
	TamanioMemoria int `json:"tamanio_memoria"`
	// con el tamaño de memoria consulta si es posible ejecutarlo en memoria
}

type DatosRespuestaDeKernel struct {
	Pseudocodigo   string `json:"filename"`
	TamanioMemoria int    `json:"tamanio_memoria"`
}

type EspacioLibreRTA struct {
	EspacioLibre int `json:"espacio_libre"`
}

type DatosParaDump struct {
	PID       int    `json:"pid"`
	TimeStamp string `json:"timeStamp"`
} // HABRIA QUE VER QUE TIPO DE DATOS ES EL TIMESTAMP
