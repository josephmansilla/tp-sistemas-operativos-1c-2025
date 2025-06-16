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
	NumeroFrame int    `json:"numero_pagina"`
	Datos       []byte `json:"frame_swap"`
	Tamanio     int    `json:"tamanio"`
}

type PedidoKernel struct {
	PID int `json:"pid"`
}

type DesuspensionProceso struct {
	PID int `json:"pid"`
}
