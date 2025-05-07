package globals

import (
	"encoding/json"
	"os"

	logger "github.com/sisoputnfrba/tp-golang/utils/logger"
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

type DatosConsultaDeKernel struct {
	PID            int `json:"pid"`
	TamanioMemoria int `json:"tamanio_memoria"`
	// con el tamaño de memoria consulta si es posible ejecutarlo en memoria
}

type DatosRespuestaDeKernel struct {
	Pseudocodigo   string `json:"filename"`
	TamanioMemoria int    `json:"tamanio_memoria"`
}

// Tipo de datos recibidos de la CPU

type DatosDeCPU struct {
	PID int `json:"pid"`
	PC  int `json:"pc"`
}

type DatosParaCPU struct {
	// TODO
}

type ContextoDeCPU struct {
	PID int `json:"pid"`
	PC  int `json:"pc"`
}

type InstruccionCPU struct {
	Instruccion string `json:"instruccion"`
}

type EspacioLibreRTA struct {
	EspacioLibre int `json:"espacio_libre"`
}

var RespuestaKernel DatosRespuestaDeKernel
var Kernel DatosConsultaDeKernel
var CPU DatosDeCPU

// EspacioDeUsuario => make([]byte, TamMemoria)

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
