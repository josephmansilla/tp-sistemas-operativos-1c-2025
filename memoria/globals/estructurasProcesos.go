package globals

type Proceso struct {
	PID                        int             `json:"pid"`
	TablaRaiz                  TablaPaginas    `json:"tabla_paginas"`
	Metricas                   MetricasProceso `json:"metricas_proceso"`
	OffsetInstruccionesEnBytes map[int][]byte  `json:"offset_instrucciones_en_bytes"`
}

/*
Por cada PID del sistema,
se deberá leer su archivo de pseudocódigo y guardar de forma estructurada
las instrucciones del mismo para poder devolverlas una a una a pedido de la CPU.

Queda a criterio del grupo utilizar
la estructura que crea conveniente para este caso de uso.

ProcesosPorPID[{{PID}}].OffsetInstruccionesEnBytes[{{PC}}] = [78 89 65 76]
*/

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

type OperacionMetrica func(*MetricasProceso, int)

type DatosParaDump struct {
	PID       int    `json:"pid"`
	TimeStamp string `json:"timeStamp"`
}

/*type EntradaDump struct {
	DireccionFisica int `json:"direccion_fisica"`
	NumeroFrame     int `json:"numero_frame"`
}*/
