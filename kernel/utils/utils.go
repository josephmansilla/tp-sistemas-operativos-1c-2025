package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

// Body JSON a recibir
type Mensaje struct {
	//Nombre string `json:"nombre"`
	Ip     string `json:"ip"`
	Puerto int    `json:"puerto"`
}

// w http.ResponseWriter. Se usa para escribir la respuesta al Cliente
// r *http.Request es la peticion que se recibio
func RecibirMensaje(w http.ResponseWriter, r *http.Request) {
	//decodificador JSON que lee directamente desde el body de la petición HTTP
	decoder := json.NewDecoder(r.Body)

	//Interpretar como si fuera un objeto de tipo Mensaje. Se guarda en variable mensaje.
	var mensaje Mensaje
	err := decoder.Decode(&mensaje)

	if err != nil {
		log.Printf("Error al decodificar el mensaje: %s", err.Error())

		//Devolver un HTTP 400 (Bad Request) al cliente.
		w.WriteHeader(http.StatusBadRequest)

		//Escribir un mensaje de error en el body de la respuesta.
		w.Write([]byte("Error al decodificar mensaje"))
		return
	}

	log.Println("Me llego un mensaje:")
	//Imprimir el contenido del struct mensaje
	log.Printf("%+v\n", mensaje)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

