package main

import (
	"fmt"
	"log"
	"net/http"

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
	mux := http.NewServeMux()

	mux.HandleFunc("/cpu/mensaje", utils.RecibirMensaje)

	log.Printf("CPU escuchando en http://%s:%d/cpu/mensaje\n", globals.ClientConfig.IpSelf, globals.ClientConfig.PortSelf)
	err := http.ListenAndServe(fmt.Sprintf(":%d", globals.ClientConfig.PortSelf), mux)
	if err != nil {
		log.Fatal(err)
	}

	//la CPU deberá solicitarle a la Memoria la siguiente instrucción.

}
