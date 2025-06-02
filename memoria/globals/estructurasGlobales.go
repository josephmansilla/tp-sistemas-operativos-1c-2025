package globals

import (
	"encoding/json"
	logger "github.com/sisoputnfrba/tp-golang/utils/logger"
	"os"
)

func ConfigCheck(filepath string) *Config {
	var configCheck *Config
	configFile, err := os.Open(filepath)
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&configCheck)
	return configCheck
} // err handling

var MemoryConfig *Config

type EntradaPagina struct {
	NumeroFrame   int  `json:"numero_frame"`
	EstaPresente  bool `json:"esta_presente"`
	EstaEnUso     bool `json:"esta_en_uso"`
	FueModificado bool `json:"fue_modificado"`
}

type TablaPagina struct {
	Subtabla        map[int]*TablaPagina   `json:"subtabla"`
	EntradasPaginas map[int]*EntradaPagina `json:"entradas_pagina"`
} // las entradasPaginas se instancian en nil hasta el último nivel

type TablaPaginas map[int]*TablaPagina

type Proceso struct {
	PID       int             `json:"pid"`
	TablaRaiz TablaPaginas    `json:"tabla_paginas"`
	Metricas  MetricasProceso `json:"metricas_proceso"`
}

type ProcesosMap map[int]*Proceso

type MetricasProceso struct {
	AccesosTablasPaginas     int `json:"acceso_tablas_paginas"`
	InstruccionesSolicitadas int `json:"instrucciones_solicitadas"`
	BajadasSwap              int `json:"bajadas_swap"`
	SubidasMP                int `json:"subidas_mp"`
	LecturasDeMemoria        int `json:"lecturas_de_memoria"`
	EscriturasDeMemoria      int `json:"escrituras_de_memoria"`
}

type DatosParaDump struct {
	PID       int    `json:"pid"`
	TimeStamp string `json:"timeStamp"`
} // HABRIA QUE VER QUE TIPO DE DATOS ES EL TIMESTAMP

var MemoriaPrincipal []byte // MP simulada
var FramesLibres []bool     //los frames van a estar en True si están libres
var ProcesosMapeable ProcesosMap

// SUPER PENDIENTES
type ArgmentosCreacionProceso struct {
	NombrePseudocodigo string `json:"nombre_pseudocodigo"`
	TamanioProceso     int    `json:"tamanioProceso"`
	// PID
}

type PedidoAMemoria struct {
	Thread    Thread                 `json:"thread"`
	Type      string                 `json:"type"`
	Arguments map[string]interface{} `json:"arguments"`
}

type Thread struct {
	PID int `json:"pid"`
}
