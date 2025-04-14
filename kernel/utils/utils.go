package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/sisoputnfrba/tp-golang/kernel/globals"
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

func LeerJson(w http.ResponseWriter, r *http.Request, mensaje any) {
	//decodificador JSON que lee directamente desde el body de la petición HTTP
	decoder := json.NewDecoder(r.Body)

	//Interpretar como si fuera un objeto de tipo Mensaje. Se guarda en variable mensaje.
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
	w.Write([]byte("STATUS OK"))
}

// w http.ResponseWriter. Se usa para escribir la respuesta al Cliente
// r *http.Request es la peticion que se recibio
func RecibirMensajeDeIO(w http.ResponseWriter, r *http.Request) {
	var mensaje MensajeDeIO
	LeerJson(w, r, &mensaje)

	globals.IO = globals.DatosIO{
		Nombre: mensaje.Nombre,
		Ip:     mensaje.Ip,
		Puerto: mensaje.Puerto,
	}

	log.Printf("Se ha guardado: %s\n", globals.IO.Nombre)
}

func RecibirMensajeDeCPU(w http.ResponseWriter, r *http.Request) {
	var mensaje MensajeDeCPU
	LeerJson(w, r, &mensaje)

	//Cargar en
	globals.CPU = globals.DatosCPU{
		Ip:     mensaje.Ip,
		Puerto: mensaje.Puerto,
	}

	log.Printf("CPU Port: %d\n", globals.CPU.Puerto)
}
