package utils

import (
	"encoding/json"
	"fmt"
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/instrucciones"
	"github.com/sisoputnfrba/tp-golang/utils/data"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
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

type ConsultaConfigMemoria struct {
	TamanioPagina    int `json:"tamanioPagina"`
	EntradasPorNivel int `json:"entradasPorNivel"`
	CantidadNiveles  int `json:"cantidadNiveles"`
}

func Config(filepath string) *globals.Config {
	//Recibe un string filepath (ruta al archivo de configuración).
	var config *globals.Config

	//Abrir archivo en la ruta filepath
	configFile, err := os.Open(filepath)

	if err != nil {
		logger.Fatal(err.Error())
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
		logger.Error("Error enviando IP, Puerto e ID al Kernel: %s", err.Error())
		return
	}
	//Si no hubo error, logueo que salio bien
	logger.Info("IP, Puerto e ID enviados exitosamente al Kernel")
}

// Recibo PCB de Kernel
func RecibirContextoDeKernel(w http.ResponseWriter, r *http.Request) {
	var msg MensajeDeKernel

	// Intentar decodificar el JSON del request
	if err := data.LeerJson(w, r, &msg); err != nil {
		logger.Error("Error al recibir JSON: %v", err)
		http.Error(w, "Error procesando datos del Kernel", http.StatusInternalServerError)
		return
	}

	// Actualizar solo PID y PC del PCB global
	globals.PIDActual = msg.PID
	globals.PCActual = msg.PC

	logger.Info("Me llegó el contexto con PID: %d, PC: %d", globals.PIDActual, globals.PCActual)

	// Validar que ClientConfig esté inicializado
	if globals.ClientConfig == nil {
		logger.Info("ClientConfig no está inicializado.")
		http.Error(w, "Configuración del cliente no inicializada", http.StatusInternalServerError)
		return
	}

	// Pedir a Memoria las instrucciones usando PID y PC
	instrucciones.FaseFetch(globals.ClientConfig.IpMemory, globals.ClientConfig.PortMemory)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("STATUS OK"))
}

func RecibirInterrupcion(w http.ResponseWriter, r *http.Request) {
	var interrumpido struct {
		PID int `json:"pid"`
	}

	if err := data.LeerJson(w, r, &interrumpido); err != nil {
		logger.Error("Error leyendo interrupción: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	logger.Info("## Llega interrupción al puerto Interrupt")

	globals.MutexInterrupcion.Lock()
	globals.InterrupcionPendiente = true //aseguro la mutua exclusion
	globals.PIDInterrumpido = interrumpido.PID
	globals.MutexInterrupcion.Unlock()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Interrupción registrada"))
}

func RecibirConfiguracionMemoria(ipDestino string, puertoDestino int) error {
	url := fmt.Sprintf("http://%s:%d/memoria/configuracion", ipDestino, puertoDestino)

	// Mensaje que CPU envía
	mensaje := struct {
		PID int `json:"pid"`
	}{
		PID: globals.PIDActual,
	}

	// Recibiremos esto
	var msg ConsultaConfigMemoria

	err := data.EnviarDatosYRecibirRespuesta(url, mensaje, &msg)
	if err != nil {
		logger.Error("Error al consultar la configuración de memoria: %s", err.Error())
		return err
	}

	globals.TamanioPagina = msg.TamanioPagina
	globals.EntradasPorNivel = msg.EntradasPorNivel
	globals.CantidadNiveles = msg.CantidadNiveles

	logger.Info("Configuración de Memoria: Tamaño de Página: %d, Entradas por Página: %d, Niveles: %d",
		msg.TamanioPagina, msg.EntradasPorNivel, msg.CantidadNiveles)

	return nil
}
