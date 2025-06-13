package comunicacion

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/utils/data"
	logger "github.com/sisoputnfrba/tp-golang/utils/logger"
)

// Body JSON a recibir
type MensajeAMemoria struct {
	Filename string `json:"filename"` //filename
	Tamanio  int    `json:"tamanio_memoria"`
	PID      int    `json:"pid"`
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
	EspacioLibre int `json:"espacio_libre"`
}

// ENVIAR ARCHIVO DE PSEUDOCODIGO Y TAMAÑO
func SolicitarEspacioEnMemoria(fileName string, tamanio int) int {
	url := fmt.Sprintf("http://%s:%d/memoria/espaciolibre", globals.KConfig.MemoryAddress, globals.KConfig.MemoryPort)

	mensaje := MensajeAMemoria{
		Filename: fileName,
		Tamanio:  tamanio,
	}

	resp, err := data.EnviarDatosConRespuesta(url, mensaje)
	if err != nil {
		logger.Error("Error enviando pseudocódigo a Memoria: %s", err.Error())
	}
	defer resp.Body.Close()

	var rta RespuestaMemoria
	err = json.NewDecoder(resp.Body).Decode(&rta)
	if err != nil {
		logger.Error("Error al decodificar respuesta de Memoria: %s", err.Error())
	}

	logger.Info("Memoria dice => Espacio libre: %d", rta.EspacioLibre)
	return rta.EspacioLibre
}

// ENVIAR ARCHIVO DE PSEUDOCODIGO Y TAMAÑO
func EnviarArchivoMemoria(fileName string, tamanio int, pid int) {
	url := fmt.Sprintf("http://%s:%d/memoria/kernel", globals.KConfig.MemoryAddress, globals.KConfig.MemoryPort)

	mensaje := MensajeAMemoria{
		Filename: fileName,
		Tamanio:  tamanio,
		PID:      pid,
	}

	err := data.EnviarDatos(url, mensaje)
	if err != nil {
		logger.Error("Error enviando pseudocódigo a Memoria: %s", err.Error())
	}
}

// PARA MANJERAR LOS MENSAJES DEL ENDPOINT QUE ESTAN EN MEMORIA
// por ejemplo: http.HandleFunc("/kernel/createProcess", createProcess)
const (
	CreateProcess = "createProcess"
	FinishProcess = "finishProcess"
	MemoryDump    = "memoryDump"
)

// esta funcion es auxiliar de sendMemoryRequest
/*func handleMemoryResponseError(response *http.Response, TypeRequest string) error {
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
}*/

var ErrorRequestType = map[string]error{
	CreateProcess: errors.New("memoria: No hay espacio disponible en memoria "),
	FinishProcess: errors.New("memoria: No se puedo finalizar el proceso"),
}
