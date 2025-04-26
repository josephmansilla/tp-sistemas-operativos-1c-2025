package data

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// Helper para enviar datos a un endpoint (POST) --> Mando un struct como JSON
func EnviarDatos(url string, data any) error {
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
func RecibirDatos(url string, data any) error {
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

// Leer Body JSON recibidos por POST o PUT
// Deserializa JSON en un struct de Go.
func LeerJson(w http.ResponseWriter, r *http.Request, mensaje any) error {
	err := json.NewDecoder(r.Body).Decode(mensaje)
	if err != nil {
		log.Printf("Error al decodificar el mensaje: %s", err.Error())
		http.Error(w, "Error al decodificar mensaje", http.StatusBadRequest)
		return err
	}
	log.Printf("Me llegó un mensaje: %+v", mensaje)

	//w.WriteHeader(http.StatusOK)
	//w.Write([]byte("STATUS OK"))
	return nil
}
