package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sisoputnfrba/tp-golang/io/globals"
	"log"
	"net/http"
	"os"
)

type MensajeAKernel struct {
	Ip     string `json:"ip"`
	Puerto int    `json:"puerto"`
	Nombre string `json:"nombre"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Falta el parametro: nombre de la interfaz de io")
		os.Exit(1)
	}

	nombre := os.Args[1]

	fmt.Printf("Nombre de la Interfaz de IO: %s\n", nombre)

	//El IO siempre es cliente del KERNEL
	log.Println("Comenzó ejecucion del IO")

	globals.ClientConfig = Config("config.json")

	log.Println("## PID: <PID> - Inicio de IO - Tiempo: <TIEMPO_IO>")

	if globals.ClientConfig == nil {
		log.Fatal("No se pudo cargar el archivo de configuración")
	}

	//Una vez leído el nombre,
	//se conectará al Kernel y en el handshake inicial le enviará su
	//nombre, ip y puerto
	//y quedará esperando las peticiones del mismo.

	//Creo una instancia del struct MensajeAKernel
	mensaje := MensajeAKernel{
		Ip:     globals.ClientConfig.IpIo,
		Puerto: globals.ClientConfig.PortIo,
		Nombre: nombre,
	}

	EnviarIpPuertoNombreAKernel(globals.ClientConfig.IpKernel, globals.ClientConfig.PortKernel, mensaje)

	log.Println("## PID: <PID> - Fin de IO")
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
func EnviarIpPuertoNombreAKernel(ipDestino string, puertoDestino int, mensaje any) {
	//Construye la URL del endpoint(url + path) a donde se va a enviar el mensaje.
	url := fmt.Sprintf("http://%s:%d/kernel/mensaje", ipDestino, puertoDestino)

	//Hace el POST
	err := enviarDatos(url, mensaje)
	//Verifico si hubo error y logueo si lo hubo
	if err != nil {
		log.Printf("Error enviando mensaje: %s", err.Error())
		return
	}
	//Si no hubo error, logueo que todo salio bien
	log.Println("Mensaje enviado exitosamente")
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
