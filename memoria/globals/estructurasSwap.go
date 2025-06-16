package globals

type SwapProceso struct {
	PID                   int `json:"pid"`
	Inicio                int `json:"inicio"`
	Fin                   int `json:"fin"`
	EntradasSwap          int `json:"entradas_swap"`
	EntradasTotales       int `json:"entradas_totales"`
	TamanioTotalBytes     int `json:"tamanio_total_bytes"`
	TamanioTotalSwapBytes int `json:"tamanio_total_swap_bytes"`
	PunteroLectura        int `json:"puntero_lectura"`
	PunteroEstructura     int `json:"puntero_estructura"`
	UltimoAcceso          int `json:"ultimo_acceso"`
}

type EntradaSwap struct {
	NumeroPagina int    `json:"numero_pagina"`
	OffsetSwap   int    `json:"offset_swap"`
	Datos        []byte `json:"frame_swap"`
	Presente     bool   `json:"presente"`
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
