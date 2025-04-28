package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/instrucciones"
	"github.com/sisoputnfrba/tp-golang/utils/data"
	"log"
	"net/http"
	"os"
	"strings"
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
	err := data.EnviarDatos(url, mensaje)
	//Verifico si hubo error y logue si lo hubo
	if err != nil {
		log.Printf("Error enviando IP, Puerto e ID al Kernel: %s", err.Error())
		return
	}
	//Si no hubo error, logueo que salio bien
	log.Println("IP, Puerto e ID enviados exitosamente al Kernel")
}

// Recibo PID y PC de Kernel
func RecibirContextoDeKernel(w http.ResponseWriter, r *http.Request) {
	var mensajeRecibido MensajeDeKernel
	if err := data.LeerJson(w, r, &mensajeRecibido); err != nil {
		return
	}

	log.Printf("Me llego el PID:%d y el PC:%d", mensajeRecibido.PID, mensajeRecibido.PC)
	//Con el PID y PC le pido a Memoria las instrucciones

	FaseFetch(globals.ClientConfig.IpMemory, globals.ClientConfig.PortMemory, mensajeRecibido.PID, mensajeRecibido.PC)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("STATUS OK"))

}

func FaseFetch(ipDestino string, puertoDestino int, pidPropio int, pcInicial int) {
	pc := pcInicial

	for {
		mensaje := MensajeInstruccion{
			PID: pidPropio,
			PC:  pc,
		}

		jsonData, err := json.Marshal(mensaje)
		if err != nil {
			log.Printf("Error codificando mensaje a JSON: %s", err)
			break
		}

		url := fmt.Sprintf("http://%s:%d/memoria/instruccion", ipDestino, puertoDestino)

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("Error haciendo POST a Memoria: %s", err)
			break
		}
		defer resp.Body.Close()

		var respuesta RespuestaInstruccion
		err = json.NewDecoder(resp.Body).Decode(&respuesta)
		if err != nil {
			log.Printf("Error decodificando respuesta de Memoria: %s", err)
			break
		}

		if respuesta.Instruccion == "" {
			log.Printf("No hay instruccion para PID %d (PC %d)", pidPropio, pc)
			break
		}

		log.Printf("Instrucción recibida (PC %d): %s", pc, respuesta.Instruccion)

		// Ejecuta la instrucción como goroutine
		FaseDecode(respuesta.Instruccion) // usar semaforos para variables compartidas

		pc++
	}
}

func FaseDecode(instruccion string) {
	partes := strings.Fields(instruccion)
	if len(partes) == 0 {
		log.Println("Instrucción vacía")
		return
	}

	nombreInstruccion := partes[0]
	argumentos := partes[1:]

	// Llamo a ejecutar la instrucción
	FaseExecute(nombreInstruccion, argumentos)
}

func FaseExecute(nombre string, argumentos []string) {
	// Usamos el instructionSet y el contexto globales
	instruccion, existe := instrucciones.InstructionSet[nombre]
	if !existe {
		log.Printf("Instrucción desconocida: %s", nombre)
		return
	}

	// Usamos CurrentContext directamente
	err := instruccion(globals.CurrentContext, argumentos) // ejecuto funcion
	log.Printf("Ejecutando instruccion: %s", nombre)
	if err != nil {
		log.Printf("Error ejecutando %s: %v", nombre, err)
	}
}
