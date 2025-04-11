package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
)

//Cree este utils.go para enviarMensaje()

// Body JSON a enviarse
type Mensaje struct {
	Ip     string `json:"ip"`
	Puerto int    `json:"puerto"`
}

func Config(filepath string) *globals.Config{
	//Recibe un string filepath (ruta al archivo de configuración).
	var config * globals.Config

	//Abrir archivo en la ruta filepath
	configFile, err := os.Open(filepath)

	if err != nil{
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

func EnviarMensaje(ipDestino string, puertoDestino int, ipPropia string, puertoPropio int) {
	//Instanciar struct mensaje
	mensaje := Mensaje{Ip: ipPropia, Puerto: puertoPropio}

	//Convertir a JSON (serializar)
	body, err := json.Marshal(mensaje)

	//Si ocurrio un error al convertir el mensaje a JSON, mostrar en el log.
	if err != nil {
		log.Printf("Error codificando mensaje: %s", err.Error())
	}

	//Construir URL
	url := fmt.Sprintf("http://%s:%d/kernel/mensaje", ipDestino, puertoDestino)

	//Peticion HTTP POST, donde el body es un JSON.
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))

	if err != nil {
		log.Printf("Error enviando mensaje a ip:%s puerto:%d", ipDestino, puertoDestino)
	}

	log.Printf("Respuesta del servidor: %s", resp.Status)
}
