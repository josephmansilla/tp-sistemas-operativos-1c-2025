package estructuras

type EntradaSwapInfo struct {
	NumeroFrame    int `json:"numero_pagina"`
	PosicionInicio int `json:"posicion_inicio"`
	Tamanio        int `json:"tamanio"`
}

type SwapProcesoInfo struct {
	Entradas     map[int]*EntradaSwapInfo `json:"entradas"`
	NumerosFrame []int                    `json:"cantidad_entradas"`
}

type EntradaSwap struct {
	NumeroFrame int    `json:"numero_pagina"`
	Datos       []byte `json:"datos"`
	Tamanio     int    `json:"tamanio"`
}

type PedidoKernel struct {
	PID int `json:"pid"`
}

type DesuspensionProceso struct {
	PID int `json:"pid"`
}
