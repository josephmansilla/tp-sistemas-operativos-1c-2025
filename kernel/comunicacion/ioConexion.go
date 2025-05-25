package comunicacion

import (
	"fmt"
	"io"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/utils/data"
	logger "github.com/sisoputnfrba/tp-golang/utils/logger"
)

// Body JSON a recibir
type MensajeDeIO struct {
	Nombre string `json:"nombre"`
	Ip     string `json:"ip"`
	Puerto int    `json:"puerto"`
}

type MensajeAIO struct {
	Pid      int `json:"pid"`
	Duracion int `json:"duracion"` //en segundos
}

// w http.ResponseWriter. Se usa para escribir la respuesta al Cliente
// r *http.Request es la peticion que se recibio
func RecibirMensajeDeIO(w http.ResponseWriter, r *http.Request) {
	var mensajeRecibido MensajeDeIO
	if err := data.LeerJson(w, r, &mensajeRecibido); err != nil {
		return //hubo error
	}

	nombre := mensajeRecibido.Nombre

	//Cargar en
	globals.IOMu.Lock()
	globals.IOs[nombre] = globals.DatosIO{
		Nombre: mensajeRecibido.Nombre,
		Ip:     mensajeRecibido.Ip,
		Puerto: mensajeRecibido.Puerto,
	}
	globals.IOCond.Broadcast() // es como un signal al wait
	globals.IOMu.Unlock()

	logger.Info("Se ha recibido IO: Nombre: %s Ip: %s Puerto: %d",
		globals.IOs[nombre].Nombre, globals.IOs[nombre].Ip, globals.IOs[nombre].Puerto)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("STATUS OK"))
}

// Enviar PID y Duracion a IO
func EnviarContextoIO(nombreIO string, pid int, duracion int) {
	//Necesito elegir a que IO mandarle
	globals.IOMu.Lock()
	ioData, ok := globals.IOs[nombreIO]
	for !ok {
		globals.IOCond.Wait()
		ioData, ok = globals.IOs[nombreIO]
	}
	globals.IOMu.Unlock()

	url := fmt.Sprintf("http://%s:%d/io/kernel", ioData.Ip, ioData.Puerto)

	mensaje := MensajeAIO{
		Pid:      pid,
		Duracion: duracion,
	}

	logger.Info("## (%d) - Bloqueado por IO: %s", mensaje.Pid, nombreIO)

	resp, err := data.EnviarDatosConRespuesta(url, mensaje)
	if err != nil {
		logger.Info("Error enviando PID y Duracion a IO: %s", err.Error())
		return
	}
	defer resp.Body.Close()

	// Leer el cuerpo del response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Info("Error leyendo respuesta de IO: %s", err.Error())
		return
	}
	logger.Info("Respuesta del módulo IO: %s", string(body))
	//TODO Pasar proceso de bloqueado a ready
	logger.Info("## (%d) finalizó IO y pasa a READY", mensaje.Pid)
}
