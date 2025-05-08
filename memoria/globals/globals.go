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
var RespuestaKernel DatosRespuestaDeKernel
var Kernel DatosConsultaDeKernel
var CPU DatosDeCPU
var DatosDump DatosParaDump

// EspacioDeUsuario => make([]byte, TamMemoria)

type MetricaProceso struct {
	CantAccesosTablasPaginas     int `json:"cant_accesos_tablas_paginas"`
	CantInstruccionesSolicitadas int `json:"cant_instrucciones_solicitadas"`
	CantBajadasHaciaSwap         int `json:"cant_bajadas_hacia_swap"`
	CantSubidasMP                int `json:"cant_subidas_mp"`
	CantLecturaMemoria           int `json:"cant_lectura_memoria"`
	CantEscrituraMemoria         int `json:"cant_escritura_memoria"`
}

type ArgmentosCreacionProceso struct {
	NombrePseudocodigo string `json:"nombre_pseudocodigo"`
	TamanioProceso     int    `json:"tamanioProceso"`
	// PID
	//
}

type PedidoAMemoria struct {
	Thread    Thread                 `json:"thread"`
	Type      string                 `json:"type"`
	Arguments map[string]interface{} `json:"arguments"`
}

type Thread struct {
	PID int `json:"pid"`
}
