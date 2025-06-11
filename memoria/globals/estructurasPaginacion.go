package globals

type EntradaPagina struct {
	NumeroFrame   int  `json:"numero_frame"`
	EstaPresente  bool `json:"esta_presente"`
	EstaEnUso     bool `json:"esta_en_uso"`
	FueModificado bool `json:"fue_modificado"`
}

type TablaPagina struct {
	Subtabla        map[int]*TablaPagina   `json:"subtabla"`
	EntradasPaginas map[int]*EntradaPagina `json:"entradas_pagina"`
}

type TablaPaginas map[int]*TablaPagina

type EscrituraPagina struct {
	PID                 int    `json:"pid"`
	DireccionFisica     int    `json:"direccion_fisica"`
	DatosASobreEscribir string `json:"datos_a_sobre_escribir"`
	TamanioNecesario    int    `json:"tamanio_necesario"`
}

type ExitoEscrituraPagina struct {
	Exito           bool   `json:"exito"`
	DireccionFisica int    `json:"direccion_fisica"`
	Mensaje         string `json:"mensaje"`
}

type LecturaPagina struct {
	PID             int `json:"pid"`
	DireccionFisica int `json:"direccion_fisica"`
}

type ExitoLecturaPagina struct {
	Exito           bool   `json:"exito"`
	PseudoCodigo    string `json:"pseudo_codigo"`
	DireccionFisica int    `json:"direccion_fisica"`
}
