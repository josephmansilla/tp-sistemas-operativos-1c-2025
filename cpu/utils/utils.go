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

// Body JSON que envia a Kernel
type MensajeAKernel struct {
	Ip     string `json:"ip"`
	Puerto int    `json:"puerto"`
	ID     string `json:"id"`
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

type RespuestaInstruccion struct {
	Instruccion string `json:"instruccion"`
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
func EnviarIpPuertoIDAKernel(ipDestino string, puertoDestino int, ipPropia string, puertoPropio int, id string) {
	//Creo una instancia del struct MensajeAKernel
	mensaje := MensajeAKernel{
		Ip:     ipPropia,
		Puerto: puertoPropio,
		ID:     id,
	}
	//Construye la URL del endpoint(url + path) en el Kernel a donde se va a enviar el mensaje.
	url := fmt.Sprintf("http://%s:%d/kernel/cpu", ipDestino, puertoDestino)
	//Hace el POST a kernel
	err := enviarDatos(url, mensaje)
	//Verifico si hubo error y logue si lo hubo
	if err != nil {
		log.Printf("Error enviando IP, Puerto e ID al Kernel: %s", err.Error())
		return
	}
	//Si no hubo error, logueo que salio bien
	log.Println("IP, Puerto e ID enviados exitosamente al Kernel")
}

// Solicitar el PID y PC al Kernel sin que CPU sea servidor
func SolicitarContextoDeKernel(ipDestino string, puertoDestino int) {
	url := fmt.Sprintf("http://%s:%d/kernel/contexto", ipDestino, puertoDestino) // Este endpoint debe existir en Kernel

	var mensaje MensajeDeKernel
	err := recibirDatos(url, &mensaje)
	if err != nil {
		log.Printf("Error al recibir contexto del Kernel: %s", err.Error())
		return
	}

	log.Printf("Recibido de Kernel: PID: %d, PC: %d", mensaje.PID, mensaje.PC)

	SolicitarInstruccion(globals.ClientConfig.IpMemory, globals.ClientConfig.PortMemory, mensaje.PID, mensaje.PC)
}

func SolicitarInstruccion(ipDestino string, puertoDestino int, pidPropio int, pcPropio int) (string, error) {
	// Creo el mensaje
	mensaje := MensajeInstruccion{
		PID: pidPropio,
		PC:  pcPropio,
	}

	// Codifico a JSON
	jsonData, err := json.Marshal(mensaje)
	if err != nil {
		log.Printf("Error codificando mensaje a JSON: %s", err)
		return "", err
	}

	// Armo la URL
	url := fmt.Sprintf("http://%s:%d/memoria/cpu", ipDestino, puertoDestino)

	// Envío el POST y espero respuesta
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error haciendo POST a Memoria: %s", err)
		return "", err
	}
	defer resp.Body.Close()

	// Decodifico la respuesta
	var respuesta RespuestaInstruccion
	err = json.NewDecoder(resp.Body).Decode(&respuesta)
	if err != nil {
		log.Printf("Error decodificando respuesta de Memoria: %s", err)
		return "", err
	}

	log.Printf("Instrucción recibida desde Memoria: %s", respuesta.Instruccion)

	return respuesta.Instruccion, nil
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
