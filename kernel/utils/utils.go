package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/utils/data"
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

// 1. CARGAR ARCHIVO CONFIG
func Config(filepath string) *globals.Config {
	//Recibe un string filepath (ruta al archivo de configuración).
	var config *globals.Config

	//Abrir archivo en la ruta filepath
	configFile, err := os.Open(filepath)

	if err != nil {
		log.Fatal(err.Error())
	}
	//defer se usa para asegurarse de cerrar recursos (archivos, conexiones, etc.)
	//incluso si hay errores más adelante.
	defer configFile.Close()

	//Crear decoder JSON que lee desde el archivo abierto (configFile).
	jsonParser := json.NewDecoder(configFile)

	//Deserializa el contenido del archivo JSON en una estructura Go.
	//llena el struct config con los valores que están en el archivo.
	jsonParser.Decode(&config)

	return config
}

// w http.ResponseWriter. Se usa para escribir la respuesta al Cliente
// r *http.Request es la peticion que se recibio
func RecibirMensajeDeIO(w http.ResponseWriter, r *http.Request) {
	var mensajeRecibido MensajeDeIO
	if err := data.LeerJson(w, r, &mensajeRecibido); err != nil {
		return //hubo error
	}

	//Cargar en
	globals.IO = globals.DatosIO{
		Nombre: mensajeRecibido.Nombre,
		Ip:     mensajeRecibido.Ip,
		Puerto: mensajeRecibido.Puerto,
	}

	log.Printf("Se ha recibido IO: Nombre: %s Ip: %s Puerto: %d",
		globals.IO.Nombre, globals.IO.Ip, globals.IO.Puerto)

	//Asignar PID y Duracion
	EnviarContextoIO(globals.IO.Ip, globals.IO.Puerto)
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

	log.Printf("Se ha recibido CPU: Ip: %s Puerto: %d ID: %s",
		globals.CPU.Ip, globals.CPU.Puerto, globals.CPU.ID)

	//Asignar PID al CPU
	EnviarContextoCPU(globals.CPU.Ip, globals.CPU.Puerto)
}

// Enviar PID y PC al CPU
func EnviarContextoCPU(ipDestino string, puertoDestino int) {
	//Construye la URL del endpoint(url + path) a donde se va a enviar el mensaje.
	url := fmt.Sprintf("http://%s:%d/cpu/kernel", ipDestino, puertoDestino)

	mensaje := PedirInformacion() //pedir a la memoria

	//Hace el POST a CPU
	err := data.EnviarDatos(url, mensaje)
	//Verifico si hubo error y logue si lo hubo
	if err != nil {
		log.Printf("Error enviando PID y PC a CPU: %s", err.Error())
		return
	}
	//Si no hubo error, logueo que salio bien
	log.Printf("PID: %d y PC: %d enviados exitosamente a CPU", mensaje.Pid, mensaje.Pc)
}

// Enviar PID y Duracion a IO
func EnviarContextoIO(ipDestino string, puertoDestino int) {
	//Construye la URL del endpoint(url + path) a donde se va a enviar el mensaje.
	url := fmt.Sprintf("http://%s:%d/io/kernel", ipDestino, puertoDestino)

	mensaje := PedirInformacionIO() //pedir a la memoria

	log.Printf("## (%d) - Bloqueado por IO: %s", mensaje.Pid, globals.IO.Nombre)

	//Hace el POST a IO
	err := data.EnviarDatos(url, mensaje)
	//Verifico si hubo error y logue si lo hubo
	if err != nil {
		log.Printf("Error enviando PID y Duracion a IO: %s", err.Error())
		return
	}
	//Si no hubo error, logueo que salio bien
	log.Printf("## (%d) finalizó IO y pasa a READY", mensaje.Pid)
}

// Pedir PC Y PID a la memoria
func PedirInformacion() MensajeToCPU {
	mensaje := MensajeToCPU{
		Pid: 1,
		Pc:  0,
	}
	return mensaje
}

// Pedir PC Y Duracion a la memoria
func PedirInformacionIO() MensajeToIO {
	mensaje := MensajeToIO{
		Pid:      1,
		Duracion: 10,
	}
	return mensaje
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
		log.Printf("Error enviando Pseudocodigo a Memoria: %s", err.Error())
		return
	}
	//Si no hubo error, logueo que salio bien
	log.Printf("Pseudocodigo: %s enviado exitosamente a Memoria", mensaje.Filename)
}

// NUEVA CONEXIÓN AGREGADA PARA QUE KERNEL LE CONSULTE LA DISPONIBILIDAD DE ESPACIOLIBRE A MEMORIA
func ConsultarEspacioLibreMemoria(ipDestino string, puertoDestino int) (int, error) {
	// SE LE PASA LA DIRECCIÓN DE LA MEMORIA POR LOS PARAMETROS
	// EL TIPO DE LA FUNCIÓN ES DE ENTERO Y ERROR
	// ESTOS TIPOS SERÁN USADOS PARA MANEJAR LA CONSULTA EN OTRAS FUNCIONES
	url := fmt.Sprintf("http://%s:%d/memoria/espaciolibre", ipDestino, puertoDestino)

	rta, err := http.Get(url)
	if err != nil {
		log.Printf("Error al hacer el GET a Memoria: %s", err.Error())
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
		log.Printf("Error al hacer el Decode para consultar a Memoria: %s", err.Error())
		return 0, err
	}
	log.Printf("Espacio libre reportado por Memoria: %d", respuesta.EspacioLibre)
	// SE LOGUEA EL ESPACIO LIBRE Y SE DEVUELVE, AL IGUAL QUE UN NIL PARA EL ERROR
	return respuesta.EspacioLibre, nil
}

func IntentarIniciarProceso(tamanioProceso int) {
	espacioLibre, err := /*utils.*/ ConsultarEspacioLibreMemoria(globals.KernelConfig.IpMemory, globals.KernelConfig.PortMemory)
	if err != nil {
		log.Println("No se pudo consultar a la memoria por Espacio Libre")
		return
	}

	if espacioLibre >= tamanioProceso {
		log.Println("Hay suficiente espacio libre en Memoria para el proceso")
	} else {
		log.Println("No hay suficiente espacio libre en memoria para el proceso")
	}
}

/*// Enviar PC Y PID a CPU
func EnviarContextoACPU(w http.ResponseWriter, r *http.Request) {
	mensaje := PedirInformacion() //pedir a la memoria
	log.Printf("Enviando contexto a CPU por GET: PID %d, PC %d", mensaje.Pid, mensaje.Pc)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(mensaje)
}*/
