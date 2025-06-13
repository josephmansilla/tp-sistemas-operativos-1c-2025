package globals

type ProcesoSuspendido struct {
	// 	PID 			int `json:"pid"` TODO: este voy a usar para mapear el vector
	DireccionFisica string `json:"direccion_fisica"`
	TamanioProceso  string `json:"tamanio_proceso"`
	Base            int    `json:"base"`
	Limite          int    `json:"limite"`
	// TODO: ver que mas agrego
}

type SuspensionProceso struct {
	PID    int   `json:"pid"`
	Indice []int `json:"indice"`
}

type ExitoSuspensionProceso struct {
	Exito           error `json:"exito"`
	DireccionFisica int   `json:"direccion_fisica"`
	TamanioProceso  int   `json:"tamanio_proceso"`
}

type DesuspensionProceso struct {
	PID int `json:"pid"`
}

type ExitoDesuspensionProceso struct {
	Exito           error `json:"exito"`
	DireccionFisica int   `json:"direccion_fisica"`
	TamanioProceso  int   `json:"tamanio_proceso"`
}
