package comunicacion

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/kernel/Utils"
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/utils/data"
	logger "github.com/sisoputnfrba/tp-golang/utils/logger"
	"net/http"
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

type MensajeFin struct {
	PID         int  `json:"pid"`
	Desconexion bool `json:"desconexion"`
}

// w http.ResponseWriter. Se usa para escribir la respuesta al Cliente
// r *http.Request es la peticion que se recibio
func RecibirMensajeDeIO(w http.ResponseWriter, r *http.Request) {
	var mensajeRecibido MensajeDeIO
	if err := data.LeerJson(w, r, &mensajeRecibido); err != nil {
		return //hubo error
	}

	tipo := mensajeRecibido.Nombre

	globals.IOMu.Lock()
	instancia := globals.DatosIO{
		Tipo:   mensajeRecibido.Nombre,
		Ip:     mensajeRecibido.Ip,
		Puerto: mensajeRecibido.Puerto,
	}

	// Agrega a la lista correspondiente
	globals.IOs[tipo] = append(globals.IOs[tipo], instancia)
	globals.IOMu.Unlock()

	logger.Info("Se ha recibido IO: Nombre: %s Ip: %s Puerto: %d",
		instancia.Tipo, instancia.Ip, instancia.Puerto)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("STATUS OK"))
}

// Enviar PID y Duracion a IO
func EnviarContextoIO(instanciaIO globals.DatosIO, pid int, duracion int) {

	url := fmt.Sprintf("http://%s:%d/io/kernel", instanciaIO.Ip, instanciaIO.Puerto)

	mensaje := MensajeAIO{
		Pid:      pid,
		Duracion: duracion,
	}

	logger.Info("## (%d) - Bloqueado por IO: %s", mensaje.Pid, instanciaIO.Tipo)

	err := data.EnviarDatos(url, mensaje)
	if err != nil {
		logger.Info("Error enviando PID y Duracion a IO: %s", err.Error())
		return
	}
}

// Al momento de recibir un mensaje de una IO se deberá verificar
// que el mismo sea una confirmación de fin de IO, en caso afirmativo,
// se deberá validar si hay más procesos esperando realizar dicha IO.
// En caso de que el mensaje corresponda a una desconexión de la IO,
// el proceso que estaba ejecutando en dicha IO, se deberá pasar al estado EXIT.
func RecibirFinDeIO(w http.ResponseWriter, r *http.Request) {
	var mensajeRecibido MensajeFin
	if err := data.LeerJson(w, r, &mensajeRecibido); err != nil {
		return //hubo error
	}

	pid := mensajeRecibido.PID
	if mensajeRecibido.Desconexion {
		logger.Info("Desconexion de IO: PID %d", pid)
		Utils.NotificarDesconexion <- mensajeRecibido.PID //PID -1
	} else {
		logger.Info("FIN de IO: PID %d", pid)
		Utils.NotificarFinIO <- mensajeRecibido.PID
	}
}
