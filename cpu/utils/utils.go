package utils

import (
	"encoding/json"
	"fmt"
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/instrucciones"
	"github.com/sisoputnfrba/tp-golang/utils/data"
	"log"
	"net/http"
	"os"
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
	//Verifico si hubo error y logueo si lo hubo
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

	// Intentar decodificar el JSON del request
	if err := data.LeerJson(w, r, &mensajeRecibido); err != nil {
		log.Printf("Error al recibir JSON: %v", err)
		http.Error(w, "Error procesando datos del Kernel", http.StatusInternalServerError)
		return
	}

	// Verificar que CurrentContext esté inicializado
	if globals.CurrentContext == nil {
		globals.CurrentContext = &globals.ExecutionContext{}
		log.Println("Inicializando CurrentContext.")
	}

	// Asignar el PID y PC a CurrentContext
	globals.CurrentContext.PID = mensajeRecibido.PID
	globals.CurrentContext.PC = mensajeRecibido.PC

	log.Printf("Me llegó el PID:%d y el PC:%d", mensajeRecibido.PID, mensajeRecibido.PC)

	// Verificar que ClientConfig esté inicializado
	if globals.ClientConfig == nil {
		log.Printf("ClientConfig no está inicializado.")
		http.Error(w, "Configuración del cliente no inicializada", http.StatusInternalServerError)
		return
	}

	// Con el PID y PC le pido a Memoria las instrucciones
	instrucciones.FaseFetch(globals.ClientConfig.IpMemory, globals.ClientConfig.PortMemory, mensajeRecibido.PID, mensajeRecibido.PC)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("STATUS OK"))
}

func RecibirInterrupcion(w http.ResponseWriter, r *http.Request) {
	var interrumpido struct {
		PID int `json:"pid"`
	}

	if err := data.LeerJson(w, r, &interrumpido); err != nil {
		log.Printf("Error leyendo interrupción: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	log.Printf("## Llega interrupción al puerto Interrupt")

	globals.MutexInterrupcion.Lock()
	globals.InterrupcionPendiente = true //aseguro la mutua exclusion
	globals.PIDInterrumpido = interrumpido.PID
	globals.MutexInterrupcion.Unlock()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Interrupción registrada"))
}

// TODO: PP
func ConsultarConfiguracionMemoria(ipDestino string, puertoDestino int) error {
	// Construimos la URL del endpoint donde la memoria debe proporcionar estos datos.
	url := fmt.Sprintf("http://%s:%d/memoria/configuracion", ipDestino, puertoDestino)

	// Enviamos la solicitud GET a memoria.
	var respuesta struct {
		PageSize int `json:"page_size"` // Tamaño de la página
	}

	// Usamos la función RecibirDatos del paquete `data` para obtener los datos de configuración
	err := data.RecibirDatos(url, &respuesta)
	if err != nil {
		log.Printf("Error al consultar la configuración de memoria: %s", err.Error())
		return err
	}

	// Almacenamos los valores en las variables globales
	globals.TamPag = respuesta.PageSize

	// Log de la configuración obtenida
	log.Printf("Configuración de Memoria: Tamaño de Página: %d, Entradas por Página: %d", globals.TamPag)

	return nil
}
