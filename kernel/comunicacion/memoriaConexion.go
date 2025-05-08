package comunicacion

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/utils/data"
	logger "github.com/sisoputnfrba/tp-golang/utils/logger"
	"net/http"
)

// Body JSON a recibir
type MensajeAMemoria struct {
	Filename string `json:"filename"` //filename
	Tamanio  int    `json:"tamanio_memoria"`
}

type ConsultaAMemoria struct {
	Hilo      Hilo        `json:"hilo"`
	Tipo      string      `json:"tipo"`
	Arguments interface{} `json:"argumentos"` // <-- puede ser cualquier tipo ahora (map, struct, etc.)
}

type Pid int

type Hilo struct {
	PID Pid `json:"pid"`
}

type RespuestaMemoria struct {
	Exito   bool   `json:"exito"`
	Mensaje string `json:"mensaje"`
}

// ENVIAR ARCHIVO DE PSEUDOCODIGO Y TAMAÑO
func SolicitarCreacionEnMemoria(fileName string, tamanio int) (bool, error) {
	url := fmt.Sprintf("http://%s:%d/memoria/espaciolibre", globals.KConfig.MemoryAddress, globals.KConfig.MemoryPort)

	mensaje := MensajeAMemoria{
		Filename: fileName,
		Tamanio:  tamanio,
	}

	resp, err := data.EnviarDatosConRespuesta(url, mensaje)
	if err != nil {
		logger.Error("Error enviando pseudocódigo a Memoria: %s", err.Error())
		return false, err
	}
	defer resp.Body.Close()

	var rta RespuestaMemoria
	err = json.NewDecoder(resp.Body).Decode(&rta)
	if err != nil {
		logger.Error("Error al decodificar respuesta de Memoria: %s", err.Error())
		return false, err
	}

	logger.Info("Respuesta de Memoria: %s", rta.Mensaje)
	return rta.Exito, nil
}

// PARA MANJERAR LOS MENSAJES DEL ENDPOINT QUE ESTAN EN MEMORIA
// por ejemplo: http.HandleFunc("/kernel/createProcess", createProcess)
const (
	CreateProcess = "createProcess"
	FinishProcess = "finishProcess"
	CreateThread  = "createThread"
	FinishThread  = "finishThread"
	MemoryDump    = "memoryDump"
	Compactacion  = "compactar"
)

func SendMemoryRequest(request ConsultaAMemoria) error {
	logger.Debug("Enviando request a  memoria: %v para el THREAD: %v", request.Tipo, request.Hilo)

	// Serializar mensaje
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return err
	}

	// Hacer request a memoria
	url := fmt.Sprintf("http://%s:%d/memoria/%s", globals.KConfig.MemoryAddress, globals.KConfig.MemoryPort, request.Tipo)
	logger.Debug("Enviando request a memoria: %v", url)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonRequest))
	if err != nil {
		logger.Error("Error al realizar request memoria: %v", err)
		return err
	}

	err = handleMemoryResponseError(resp, request.Tipo)
	if err != nil {
		return err
	}
	return nil
}

// esta funcion es auxiliar de sendMemoryRequest
func handleMemoryResponseError(response *http.Response, TypeRequest string) error {
	logger.Debug("Memoria respondio a: %v con: %v", TypeRequest, response.StatusCode)
	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusConflict { // Conflict es compactacion.
			err := ErrorRequestType[Compactacion]
			return err
		}
		err := ErrorRequestType[TypeRequest]
		return err
	}
	return nil
}

var ErrorRequestType = map[string]error{
	CreateProcess: errors.New("memoria: No hay espacio disponible en memoria "),
	FinishProcess: errors.New("memoria: No se puedo finalizar el proceso"),
	CreateThread:  errors.New("memoria: No se puedo crear el hilo"),
	FinishThread:  errors.New("memoria: No se pudo finalizar el hilo"),
	Compactacion:  errors.New("memoria: Se debe compactar"),
}
