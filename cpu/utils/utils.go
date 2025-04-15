package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
)

//Cree este utils.go para enviarMensaje()

// Body JSON que envia a Kernel
type MensajeAKernel struct {
	Ip     string `json:"ip"`
	Puerto int    `json:"puerto"`
}

// Body JSON que recibe de Kernel
type MensajeDeKernel struct {
	PID int `json:"pid"`
	PC  int `json:"pc"`
}

type MensajeInstruccion struct {
	PID int `json:"pid"`
	PC  int `json:"pc"`
}

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

// Enviar IP y Puerto al Kernel
func EnviarIpPuertoAKernel(ipDestino string, puertoDestino int, ipPropia string, puertoPropio int) {
	//Creo una instancia del struct MensajeAKernel
	mensaje := MensajeAKernel{
		Ip:     ipPropia,
		Puerto: puertoPropio,
	}
	//Construye la URL del endpoint(url + path) en el Kernel a donde se va a enviar el mensaje.
	url := fmt.Sprintf("http://%s:%d/kernel/mensaje", ipDestino, puertoDestino)
	//Hace el POST a kernel
	err := enviarDatos(url, mensaje)
	//Verifico si hubo error y logue si lo hubo
	if err != nil {
		log.Printf("Error enviando IP y Puerto al Kernel: %s", err.Error())
		return
	}
	//Si no hubo error, logueo que salio bien
	log.Println("IP y Puerto enviados exitosamente al Kernel")
}

// Enviar PID y PC a Memoria
func SolicitarInstruccion(ipDestino string, puertoDestino int, pidPropio int, pcPropio int) {
	//Creo una instancia del struct MensajeInstruccion
	mensaje := MensajeInstruccion{
		PID: pidPropio,
		PC:  pcPropio,
	}
	//Construyo la URL del endpoint(url + path) en Memoria a donde se va a enviar el mensaje
	url := fmt.Sprintf("http://%s:%d/memoria/cpu", ipDestino, puertoDestino)
	//Hago el post a memoria
	err := enviarDatos(url, mensaje)
	//Verifico si hubo error y logueo si lo hubo
	if err != nil {
		log.Printf("Error enviando PID y PC a Memoria: %s", err.Error())
		return
	}
	//Si no hubo error, logueo que salio bien
	log.Println("PID y PC enviados éxitosamente a Memoria")
}

// Recibir PID y PC del Kernel
func RecibirContextoProcesoDeKernel(w http.ResponseWriter, r *http.Request) {
	//Creo un decoder para leer el JSON
	decoder := json.NewDecoder(r.Body)

	//Declaro una variable para poder guardar la info
	var mensaje MensajeDeKernel
	//Interpreto la info como MensajeDeKernel
	err := decoder.Decode(&mensaje)

	//Verifico si hubo un error y lo logueo
	if err != nil {
		log.Printf("Error al decodificar el mensaje del Kernel: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error al decodificar mensaje"))
		return
	}

	globals.PIDActual = mensaje.PID
	globals.PCActual = mensaje.PC
	/*for {
		SolicitarInstruccion(globals.ClientConfig.IpMemory, globals.ClientConfig.PortMemory, mensaje.PID, mensaje.PC)
		RecibirInstruccion(instruccion)
		EjecutarInstruccion(instruccion)
		mensaje.PC++
	}*/

	//Si salio bien loguea
	log.Println("Me llegó contexto de proceso desde Kernel:")
	log.Printf("PID: %d, PC: %d\n", mensaje.PID, mensaje.PC)

	//Le digo a kernel que llego bien
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// Helper para enviar datos a un endpoint (POST) --> Mando un struct como JSON
func enviarDatos(url string, data any) error {
	//Convierte el struct(data) a un JSON
	jsonData, err := json.Marshal(data)

	//Si no pudo serializar, devuelvo error
	if err != nil {
		return err
	}

	//POST a la url con el JSON
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	//Verifico error
	if err != nil {
		return err
	}
	//Cierro la rta, salio bien
	defer resp.Body.Close()

	return nil
}

// Helper para recibir datos desde un endpoint (GET) --> Pasa de JSON a struct
func recibirDatos(url string, data any) error {
	//Llamo al endpoint y verifico error
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	//Leo el contenido
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	//Deserializacion del JSON y lo paso a data
	err = json.Unmarshal(body, data)
	if err != nil {
		return err
	}

	return nil
}
