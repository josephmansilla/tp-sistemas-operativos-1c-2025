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
}

var MemoryConfig *Config

// Tipo de datos recibidos de1 Kernel
type EntradaDePagina struct {
	EstaPresente  bool `json:"esta_presente"`
	NumeroFrame   int  `json:"numero_frame"`
	EstaEnUso     bool `json:"esta_en_uso"`
	FueModificado bool `json:"fue_modificado"`
}

type MetricasProceso struct {
	AccesoATablasPaginas     int `json:"acceso_a_tablas_paginas"`
	InstruccionesSolicitadas int `json:"instrucciones_solicitadas"`
	BajadasSwap              int `json:"bajadas_swap"`
	SubidasMP                int `json:"subidas_mp"`
	LecturasDeMemoria        int `json:"lecturas_de_memoria"`
	EscriturasDeMemoria      int `json:"escrituras_de_memoria"`
}
type Proceso struct {
	PID        int             `json:"pid"`
	TablaNivel map[int]*int    `json:"tabla_nivel"`
	Metricas   MetricasProceso `json:"metricas_proceso"`
}

type DatosParaDump struct {
	PID       int    `json:"pid"`
	TimeStamp string `json:"timeStamp"`
} // HABRIA QUE VER QUE TIPO DE DATOS ES EL TIMESTAMP

var MemoriaPrincipal [][]byte // MP simulada
// TODO: LA INICIALIZAR LA MEMORIA SE CALCULA LA CANTIDAD DE FRAMES HACIENDO
// TODO: MEMORYSIZE / PAGSIZE SACADOS DEL CONFIG
var FramesLibres []bool //los frames van a estar en True si están libres

// EspacioDeUsuario => make([]byte, TamMemoria)

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
