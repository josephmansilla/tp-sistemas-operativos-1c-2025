package main

import (
	"log"

	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/utils"
)

func main() {
	log.Println("Comenzó ejecucion del Kernel")

	//CPU CLIENTE
	globals.ClientConfig = utils.Config("config.json")

	if globals.ClientConfig == nil {
		log.Fatal("No se pudo cargar el archivo de configuración")
	}

	//Las CPUs deberán conectarse al Kernel (destino)
	//enviandole su IP y su PUERTO. (self)
	utils.EnviarMensaje(globals.ClientConfig.IpKernel, globals.ClientConfig.PortKernel, globals.ClientConfig.IpSelf, globals.ClientConfig.PortSelf)

	//Al momento de recibir un PID y un PC de parte del Kernel,
	//la CPU deberá solicitarle a la Memoria la siguiente instrucción.

}
