package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/sisoputnfrba/tp-golang/io/globals"
	"github.com/sisoputnfrba/tp-golang/utils/data"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
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

type MensajeFin struct {
	PID    int    `json:"pid"`
	Nombre string `json:"nombre"` // en segundos
}

func main() {
	// ----------------------------------------------------
	// ---------- PARTE CARGA DE PARAMETROS ---------------
	// ----------------------------------------------------
	if len(os.Args) < 2 {
		fmt.Println("Falta el parametro: nombre de la interfaz de io")
		os.Exit(1)
	}

	nombre := os.Args[1]

	// ----------------------------------------------------
	// ------------- CARGO LOGS DE IO EN TXT --------------
	// ----------------------------------------------------
	logFileName := fmt.Sprintf("./logs/io_%s.log", nombre)
	var err = logger.ConfigureLogger(logFileName, "INFO")
	if err != nil {
		fmt.Println("No se pudo crear el logger -", err.Error())
		os.Exit(1)
	}
	logger.Debug("Logger creado")

	logger.Info("Comenzó ejecucion del IO")
	logger.Info("Nombre de la Interfaz de IO: %s", nombre)

	// ----------------------------------------------------
	// ---------- PARTE CARGA DEL CONFIG ------------------
	// ----------------------------------------------------
	configFilename := fmt.Sprintf("%sconfig.json", nombre)
	configPath := fmt.Sprintf("./configs/%s", configFilename)
	configData, err := os.ReadFile(configPath)
	if err != nil {
		logger.Fatal("No se pudo leer el archivo de configuración - %v", err.Error())
	}

	err = json.Unmarshal(configData, &globals.IoConfig)
	if err != nil {
		logger.Fatal("No se pudo parsear el archivo de configuración - %v", err.Error())
	}

	if err = globals.IoConfig.Validate(); err != nil {
		logger.Fatal("La configuración no es válida - %v", err.Error())
	}

	err = logger.SetLevel(globals.IoConfig.LogLevel)
	if err != nil {
		logger.Fatal("No se pudo leer el log-level - %v", err.Error())
	}

	//Instancio el mensaje a mandar a Kernel
	mensaje := MensajeAKernel{
		Ip:     globals.IoConfig.IpIo,
		Puerto: globals.IoConfig.PortIo,
		Nombre: nombre,
	}

	//Lo mando
	EnviarIpPuertoNombreAKernel(globals.IoConfig.IpKernel, globals.IoConfig.PortKernel, mensaje)

	// ------------------------------------------------------
	// ---------- ESCUCHO REQUESTS DE KERNEL ----------------
	// ------------------------------------------------------
	mux := http.NewServeMux()
	mux.HandleFunc("/io/kernel", RecibirMensajeDeKernel)

	// Inicia el servidor HTTP para escuchar las peticiones del Kernel
	direccion := fmt.Sprintf("%s:%d", globals.IoConfig.IpIo, globals.IoConfig.PortIo)
	fmt.Printf("Escuchando en %s...", direccion)

	err = http.ListenAndServe(direccion, mux)
	if err != nil {
		panic(err)
	}
}

func Config(filepath string) *globals.Config {
	//Recibe un string filepath (ruta al archivo de configuración).
	var config *globals.Config

	//Abrir archivo en la ruta filepath
	configFile, err := os.Open(filepath)

	if err != nil {
		logger.Fatal("%s", err.Error())
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

	// Hace el POST al Kernel
	err := data.EnviarDatos(url, mensaje)
	// Verifico si hubo error y logueo si lo hubo
	if err != nil {
		logger.Error("Error enviando mensaje: %s", err.Error())
		return
	}
	// Si no hubo error, logueo que todo salió bien
	logger.Info("Mensaje enviado a Kernel")
}

// Al momento de recibir una petición del Kernel,
// el módulo deberá iniciar un usleep
// por el tiempo indicado en la request.
func RecibirMensajeDeKernel(w http.ResponseWriter, r *http.Request) {
	var mensajeRecibido MensajeDeKernel
	if err := data.LeerJson(w, r, &mensajeRecibido); err != nil {
		return //hubo error
	}

	//Realizo la operacion
	logger.Info("## PID: <%d> - Inicio de IO - Tiempo: %d", mensajeRecibido.PID, mensajeRecibido.Duracion)
	time.Sleep(time.Duration(mensajeRecibido.Duracion) * time.Second)

	logger.Info("## PID: <%d> - Fin de IO", mensajeRecibido.PID)
	//IO Finalizada
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("IO finalizada correctamente"))
}

// Al finalizar deberá informar al Kernel que finalizó la solicitud de I/O
// quedará a la espera de la siguiente petición.
// ver el tema de FINALIZACION DE IO != FIN de timer de IO

//El Módulo IO, deberá notificar al Kernel de su finalización,
//para esto se deberá implementar el manejo de las señales SIGINT y SIGTERM,
//para enviar la notificación y finalizar de manera controlada.

func FinDeIO(pid int, nombre string) {
	url := fmt.Sprintf("http://%s:%d/kernel/fin_io", globals.IoConfig.IpKernel, globals.IoConfig.PortKernel)

	mensaje := MensajeFin{
		PID:    pid,
		Nombre: nombre,
	}
	logger.Info("Enviando PID <%d> y Nombre <%s> a Kernel", mensaje.PID, mensaje.Nombre)

	err := data.EnviarDatos(url, mensaje)
	if err != nil {
		logger.Info("Error enviando PID y Nombre a Kernel: %s", err.Error())
		return
	}
}
