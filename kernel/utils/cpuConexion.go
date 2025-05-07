package utils

import (
	"fmt"
	"net/http"
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/utils/data"
	logger "github.com/sisoputnfrba/tp-golang/utils/logger"
)

// Body JSON a recibir
type MensajeDeCPU struct {
	Ip     string `json:"ip"`
	Puerto int    `json:"puerto"`
	ID     string `json:"id"`
}

type MensajeACPU struct {
	Pid int `json:"pid"`
	Pc  int `json:"pc"`
}

func RecibirMensajeDeCPU(w http.ResponseWriter, r *http.Request) {
	var mensajeRecibido MensajeDeCPU
	if err := data.LeerJson(w, r, &mensajeRecibido); err != nil {
		return //hubo error
	}

	//Cargar en
	globals.CPU = globals.DatosCPU{
		Ip:     mensajeRecibido.Ip,
		Puerto: mensajeRecibido.Puerto,
		ID:     mensajeRecibido.ID,
	}

	logger.Info("Se ha recibido CPU: Ip: %s Puerto: %d ID: %s",
		globals.CPU.Ip, globals.CPU.Puerto, globals.CPU.ID)

	//Asignar PID al CPU
	EnviarContextoCPU(globals.CPU.Ip, globals.CPU.Puerto)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("STATUS OK"))
}

// Enviar PID y PC al CPU
func EnviarContextoCPU(ipDestino string, puertoDestino int) {
	//Construye la URL del endpoint(url + path) a donde se va a enviar el mensaje.
	url := fmt.Sprintf("http://%s:%d/cpu/kernel", ipDestino, puertoDestino)

	mensaje := MensajeACPU{
		Pid: 0, //PEDIR AL PCB
		Pc:  0, //PEDIR A MEMORIA
	}

	//Hace el POST a CPU
	err := data.EnviarDatos(url, mensaje)
	//Verifico si hubo error y logue si lo hubo
	if err != nil {
		logger.Info("Error enviando PID y PC a CPU: %s", err.Error())
		return
	}
	//Si no hubo error, logueo que salio bien
	logger.Info("PID: %d y PC: %d enviados exitosamente a CPU", mensaje.Pid, mensaje.Pc)
}