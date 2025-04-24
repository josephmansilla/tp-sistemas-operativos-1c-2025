package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/sisoputnfrba/tp-golang/io/globals"
	"github.com/sisoputnfrba/tp-golang/utils/data"
)

type MensajeAKernel struct {
	Ip     string `json:"ip"`
	Puerto int    `json:"puerto"`
	Nombre string `json:"nombre"`
}

type MensajeDeKernel struct {
	PID      int `json:"pid"`
	Duracion int `json:"duracion"` // en segundos
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Falta el parametro: nombre de la interfaz de io")
		os.Exit(1)
	}

	nombre := os.Args[1]

	logFileName := fmt.Sprintf("io_%s.log", nombre)
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Error al crear archivo de log para io %s: %v\n", nombre, err)
		os.Exit(1)
	}
	log.SetOutput(logFile)

	log.Printf("Nombre de la Interfaz de IO: %s\n", nombre)

	log.Println("Comenzó ejecucion del IO")

	globals.ClientConfig = Config("config.json")

	if globals.ClientConfig == nil {
		log.Fatal("No se pudo cargar el archivo de configuración")
	}

	//Instancio el mensaje a mandar a Kernel
	mensaje := MensajeAKernel{
		Ip:     globals.ClientConfig.IpIo,
		Puerto: globals.ClientConfig.PortIo,
		Nombre: nombre,
	}

	//Lo mando
	EnviarIpPuertoNombreAKernel(globals.ClientConfig.IpKernel, globals.ClientConfig.PortKernel, mensaje)

	/*
		for {
			SolicitarOperacionIO(globals.ClientConfig.IpKernel, globals.ClientConfig.PortKernel)
		}*/
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
func EnviarIpPuertoNombreAKernel(ipDestino string, puertoDestino int, mensaje MensajeAKernel) {
	// Construye la URL del endpoint (url + path) a donde se va a enviar el mensaje
	url := fmt.Sprintf("http://%s:%d/kernel/io", ipDestino, puertoDestino)

	// Solicito Operaciones: Iniciar servidor IO primero
	mux := http.NewServeMux()
	mux.HandleFunc("/io/kernel", RecibirMensajeDeKernel)

	// Inicia el servidor HTTP para escuchar las peticiones del Kernel
	direccion := fmt.Sprintf("%s:%d", globals.ClientConfig.IpIo, globals.ClientConfig.PortIo)
	log.Printf("Escuchando en %s...", direccion)

	go func() {
		if err := http.ListenAndServe(direccion, mux); err != nil {
			log.Fatalf("Error al iniciar el servidor IO: %v", err)
		}
	}()

	// Asegurarse de que el servidor IO esté completamente iniciado antes de hacer el POST
	time.Sleep(1 * time.Second) // Retraso opcional para asegurar que el servidor esté escuchando

	// Hace el POST al Kernel
	err := data.EnviarDatos(url, mensaje)
	// Verifico si hubo error y logueo si lo hubo
	if err != nil {
		log.Printf("Error enviando mensaje: %s", err.Error())
		return
	}
	// Si no hubo error, logueo que todo salió bien
	log.Printf("Mensaje enviado a Kernel")
}

// Recibir PID Y Tiempo de Kernel
func RecibirMensajeDeKernel(w http.ResponseWriter, r *http.Request) {
	var mensajeRecibido MensajeDeKernel
	if err := data.LeerJson(w, r, &mensajeRecibido); err != nil {
		return //hubo error
	}

	//Realizo la operacion
	log.Printf("## PID:%d - Inicio de IO - Tiempo:%d", mensajeRecibido.PID, mensajeRecibido.Duracion)
	time.Sleep(time.Duration(mensajeRecibido.Duracion) * time.Second)

	//IO finalizada
	if err := InformarFinalizacionIO(globals.ClientConfig.IpKernel, globals.ClientConfig.PortKernel, mensajeRecibido.PID); err != nil {
		log.Printf("Error al notificar al Kernel: %s", err.Error())
	}
	log.Printf("## PID:%d - Fin de IO", mensajeRecibido.PID)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("STATUS OK"))
}

// FIN de IO
func InformarFinalizacionIO(ipDestino string, puertoDestino int, pid int) error {
	url := fmt.Sprintf("http://%s:%d/kernel/io/finalizado", ipDestino, puertoDestino)

	mensaje := struct {
		PID int `json:"pid"`
	}{
		PID: pid,
	}

	err := data.EnviarDatos(url, mensaje)
	if err != nil {
		log.Printf("Error notificando finalización de IO al Kernel: %s", err.Error())
		return err
	} else {
		return nil
	}
}

// Recibir PID y Tiempo de duracion de la operacion de IO (alternativa)
func SolicitarOperacionIO(ipDestino string, puertoDestino int) {
	//Creo la URL
	url := fmt.Sprintf("http://%s:%d/kernel/operacion", ipDestino, puertoDestino)

	var mensaje MensajeDeKernel
	err := data.RecibirDatos(url, &mensaje)
	if err != nil {
		log.Printf("Error al recibir PID y duracion del Kernel: %s", err.Error())
		return
	}
	//Realizo la operacion
	log.Printf("## PID:%d - Inicio de IO - Tiempo:%d", mensaje.PID, mensaje.Duracion)
	time.Sleep(time.Duration(mensaje.Duracion) * time.Second)

	//IO finalizada
	if err := InformarFinalizacionIO(ipDestino, puertoDestino, mensaje.PID); err != nil {
		log.Printf("Error al notificar al Kernel: %s", err.Error())
	}
	log.Printf("## PID:%d - Fin de IO", mensaje.PID)

}
