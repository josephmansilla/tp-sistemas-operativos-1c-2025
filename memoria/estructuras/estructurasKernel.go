package estructuras

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

type FinalizacionProceso struct {
	PID int `json:"pid"`
}

type ExitoEdicionMemoria struct {
	Exito    error `json:"exito"`
	Booleano bool  `json:"booleano"`
}
