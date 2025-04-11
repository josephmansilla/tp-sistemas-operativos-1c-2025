package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

//Cree este utils.go para enviarMensaje()

// Body JSON a enviarse
type Mensaje struct {
	Ip     string `json:"ip"`
	Puerto int    `json:"puerto"`
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
