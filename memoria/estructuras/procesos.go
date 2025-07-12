package estructuras

type Proceso struct {
	PID                  int             `json:"pid"`
	TablaRaiz            TablaPaginas    `json:"tabla_paginas"`
	Metricas             MetricasProceso `json:"metricas_proceso"`
	InstruccionesEnBytes map[int][]byte  `json:"instrucciones_en_bytes"`
}

type MetricasProceso struct {
	AccesosTablasPaginas     int `json:"acceso_tablas_paginas"`
	InstruccionesSolicitadas int `json:"instrucciones_solicitadas"`
	BajadasSwap              int `json:"bajadas_swap"`
	SubidasMP                int `json:"subidas_mp"`
	LecturasDeMemoria        int `json:"lecturas_de_memoria"`
	EscriturasDeMemoria      int `json:"escrituras_de_memoria"`
}

type OperacionMetrica func(*MetricasProceso, int)

type DatosParaDump struct {
	PID       int    `json:"pid"`
	TimeStamp string `json:"timeStamp"`
}
