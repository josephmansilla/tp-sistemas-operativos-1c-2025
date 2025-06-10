package globals

type Proceso struct {
	PID       int             `json:"pid"`
	TablaRaiz TablaPaginas    `json:"tabla_paginas"`
	Metricas  MetricasProceso `json:"metricas_proceso"`
}

type Ocupante struct {
	PID          int `json:"pid"`
	NumeroPagina int `json:"numero_pagina"`
}

type MetricasProceso struct {
	AccesosTablasPaginas     int `json:"acceso_tablas_paginas"`
	InstruccionesSolicitadas int `json:"instrucciones_solicitadas"`
	BajadasSwap              int `json:"bajadas_swap"`
	SubidasMP                int `json:"subidas_mp"`
	LecturasDeMemoria        int `json:"lecturas_de_memoria"`
	EscriturasDeMemoria      int `json:"escrituras_de_memoria"`
}

type OperacionMetrica func(*MetricasProceso)

type DatosParaDump struct {
	PID       int    `json:"pid"`
	TimeStamp string `json:"timeStamp"`
}
