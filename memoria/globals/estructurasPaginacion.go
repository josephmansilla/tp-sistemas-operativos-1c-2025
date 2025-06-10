package globals

var EntradaGenerica EntradaPagina

type EntradaPagina struct {
	NumeroFrame   int    `json:"numero_frame"`
	Datos         []byte `json:"datos"`
	EstaPresente  bool   `json:"esta_presente"`
	EstaEnUso     bool   `json:"esta_en_uso"`
	FueModificado bool   `json:"fue_modificado"`
}

type TablaPagina struct {
	Subtabla        map[int]*TablaPagina   `json:"subtabla"`
	EntradasPaginas map[int]*EntradaPagina `json:"entradas_pagina"`
} // las entradasPaginas se instancian en nil hasta el último nivel

type TablaPaginas map[int]*TablaPagina
