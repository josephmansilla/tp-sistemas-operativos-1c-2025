package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/pcb"
	"github.com/sisoputnfrba/tp-golang/utils/data"
	logger "github.com/sisoputnfrba/tp-golang/utils/logger"
	"io"
	"net/http"
)

// Body JSON a recibir
type MensajeDeIO struct {
	Nombre string `json:"nombre"`
	Ip     string `json:"ip"`
	Puerto int    `json:"puerto"`
}

type MensajeDeCPU struct {
	Ip     string `json:"ip"`
	Puerto int    `json:"puerto"`
	ID     string `json:"id"`
}

type MensajeToCPU struct {
	Pid int `json:"pid"`
	Pc  int `json:"pc"`
}

type MensajeToIO struct {
	Pid      int `json:"pid"`
	Duracion int `json:"duracion"` //en segundos
}

type MensajeToMemoria struct {
	Filename string `json:"filename"` //filename
	Tamanio  int    `json:"tamanio_memoria"`
}

// w http.ResponseWriter. Se usa para escribir la respuesta al Cliente
// r *http.Request es la peticion que se recibio
func RecibirMensajeDeIO(w http.ResponseWriter, r *http.Request) {
	var mensajeRecibido MensajeDeIO
	if err := data.LeerJson(w, r, &mensajeRecibido); err != nil {
		return //hubo error
	}

	//Cargar en
	globals.IOMu.Lock()
	globals.IO = globals.DatosIO{
		Nombre: mensajeRecibido.Nombre,
		Ip:     mensajeRecibido.Ip,
		Puerto: mensajeRecibido.Puerto,
	}
	globals.IOCond.Broadcast() // es como un signal al wait
	globals.IOMu.Unlock()

	logger.Info("Se ha recibido IO: Nombre: %s Ip: %s Puerto: %d",
		globals.IO.Nombre, globals.IO.Ip, globals.IO.Puerto)
}

func RecibirMensajeDeCPU(w http.ResponseWriter, r *http.Request) {
	var mensajeRecibido MensajeDeCPU
	if err := data.LeerJson(w, r, &mensajeRecibido); err != nil {
		return //hubo error
	}

	//Cargar en
	globals.CPU = globals.DatosCPU{
		Ip:     mensajeRecibido.Ip,
		Puerto: mensajeRecibido.Puerto,
		ID:     mensajeRecibido.ID,
	}

	logger.Info("Se ha recibido CPU: Ip: %s Puerto: %d ID: %s",
		globals.CPU.Ip, globals.CPU.Puerto, globals.CPU.ID)

	//Asignar PID al CPU
	EnviarContextoCPU(globals.CPU.Ip, globals.CPU.Puerto)
}

// Enviar PID y PC al CPU
func EnviarContextoCPU(ipDestino string, puertoDestino int) {
	//Construye la URL del endpoint(url + path) a donde se va a enviar el mensaje.
	url := fmt.Sprintf("http://%s:%d/cpu/kernel", ipDestino, puertoDestino)

	mensaje := MensajeToCPU{
		Pid: 0, //PEDIR AL PCB
		Pc:  0, //PEDIR A MEMORIA
	}

	//Hace el POST a CPU
	err := data.EnviarDatos(url, mensaje)
	//Verifico si hubo error y logue si lo hubo
	if err != nil {
		logger.Info("Error enviando PID y PC a CPU: %s", err.Error())
		return
	}
	//Si no hubo error, logueo que salio bien
	logger.Info("PID: %d y PC: %d enviados exitosamente a CPU", mensaje.Pid, mensaje.Pc)
}

// Enviar PID y Duracion a IO
func EnviarContextoIO(ipDestino string, puertoDestino int, pid int, duracion int) {
	url := fmt.Sprintf("http://%s:%d/io/kernel", ipDestino, puertoDestino)

	mensaje := MensajeToIO{
		Pid:      pid,
		Duracion: duracion,
	}

	logger.Info("## (%d) - Bloqueado por IO: %s", mensaje.Pid, globals.IO.Nombre)

	resp, err := data.EnviarDatosConRespuesta(url, mensaje)
	if err != nil {
		logger.Info("Error enviando PID y Duracion a IO: %s", err.Error())
		return
	}
	defer resp.Body.Close()

	// Leer el cuerpo del response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Info("Error leyendo respuesta de IO: %s", err.Error())
		return
	}

	logger.Info("Respuesta del módulo IO: %s", string(body))
	logger.Info("## (%d) finalizó IO y pasa a READY", mensaje.Pid)
}

func EnviarFileMemoria(ipDestino string, puertoDestino int, filename string, tamanioProceso int) {
	//Construye la URL del endpoint(url + path) a donde se va a enviar el mensaje.
	url := fmt.Sprintf("http://%s:%d/memoria/kernel", ipDestino, puertoDestino)

	mensaje := MensajeToMemoria{
		Filename: filename,
		Tamanio:  tamanioProceso,
	}

	//Hace el POST a Memoria
	err := data.EnviarDatos(url, mensaje)
	//Verifico si hubo error y logue si lo hubo
	if err != nil {
		logger.Info("Error enviando Pseudocodigo a Memoria: %s", err.Error())
		return
	}
	//Si no hubo error, logueo que salio bien
	logger.Info("Pseudocodigo: %s enviado exitosamente a Memoria", mensaje.Filename)
}

// NUEVA CONEXIÓN AGREGADA PARA QUE KERNEL LE CONSULTE LA DISPONIBILIDAD DE ESPACIOLIBRE A MEMORIA
func ConsultarEspacioLibreMemoria(ipDestino string, puertoDestino int) (int, error) {
	// SE LE PASA LA DIRECCIÓN DE LA MEMORIA POR LOS PARAMETROS
	// EL TIPO DE LA FUNCIÓN ES DE ENTERO Y ERROR
	// ESTOS TIPOS SERÁN USADOS PARA MANEJAR LA CONSULTA EN OTRAS FUNCIONES
	url := fmt.Sprintf("http://%s:%d/memoria/espaciolibre", ipDestino, puertoDestino)

	rta, err := http.Get(url)
	if err != nil {
		logger.Info("Error al hacer el GET a Memoria: %s", err.Error())
		return 0, err
		// SI HAY UN ERROR SE DEVUELVE EL ERROR, PERO TAMBIÉN ES NECESARIO INDICAR
		// QUE EL ESPACIOLIBRE ES 0
	}
	defer rta.Body.Close()

	// USANDO EL STRUCT QUE COMPARTEN MEMORIA Y KERNEL
	// POR AHORA ES LÓGICA REPETIDA
	var respuesta globals.EspacioLibreRTA
	err = json.NewDecoder(rta.Body).Decode(&respuesta)
	if err != nil {
		logger.Info("Error al hacer el Decode para consultar a Memoria: %s", err.Error())
		return 0, err
	}
	logger.Info("Espacio libre reportado por Memoria: %d", respuesta.EspacioLibre)
	// SE LOGUEA EL ESPACIO LIBRE Y SE DEVUELVE, AL IGUAL QUE UN NIL PARA EL ERROR
	return respuesta.EspacioLibre, nil
}

func IntentarIniciarProceso(tamanioProceso int) {
	espacioLibre, err := ConsultarEspacioLibreMemoria(Config.MemoryAddress, Config.MemoryPort)
	if err != nil {
		logger.Info("No se pudo consultar a la memoria por Espacio Libre")
		return
	}

	if espacioLibre >= tamanioProceso {
		logger.Info("Hay suficiente espacio libre en Memoria para el proceso")
	} else {
		logger.Info("No hay suficiente espacio libre en memoria para el proceso")
	}
}

type RequestToMemory struct {
	Thread    Thread      `json:"thread"`
	Type      string      `json:"type"`
	Arguments interface{} `json:"arguments"` // <-- puede ser cualquier tipo ahora (map, struct, etc.)
}
type Pid int

type Thread struct {
	PID Pid `json:"pid"`
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

func SendMemoryRequest(request RequestToMemory) error {
	logger.Debug("Enviando request a  memoria: %v para el THREAD: %v", request.Type, request.Thread)

	// Serializar mensaje
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return err
	}

	// Hacer request a memoria
	url := fmt.Sprintf("http://%s:%d/memoria/%s", Config.MemoryAddress, Config.MemoryPort, request.Type)
	logger.Debug("Enviando request a memoria: %v", url)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonRequest))
	if err != nil {
		logger.Error("Error al realizar request memoria: %v", err)
		return err
	}

	err = handleMemoryResponseError(resp, request.Type)
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

// ESTAS SON VARIABLES GLOBALES OJO¡¡¡¡
var Config KernelConfig
var ColaNuevo Queue[*pcb.PCB]
var NewStateQueue Queue[*pcb.PCB]
var ColaBLoqueado Queue[*pcb.PCB]
var ColaSalida Queue[*pcb.PCB]
var ColaEjecutando Queue[*pcb.PCB]
var ColaReady Queue[*pcb.PCB]
var ColaBloqueadoSuspendido Queue[*pcb.PCB]
var ColaSuspendidoReady Queue[*pcb.PCB]
