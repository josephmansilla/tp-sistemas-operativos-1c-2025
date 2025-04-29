package main

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/utils"
	"log"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Falta el parametro: identificador de CPU")
		os.Exit(1)
	}

	ID := os.Args[1]

	logFileName := fmt.Sprintf("cpu_%s.log", ID)
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Printf("Error al crear archivo de log para CPU %s: %v\n", ID, err)
		os.Exit(1)
	}
	log.SetOutput(logFile)

	log.Printf("Comenzo ejecucion del CPU con ID: %s", ID)

	//CPU CLIENTE
	globals.ClientConfig = utils.Config("config.json")

	if globals.ClientConfig == nil {
		log.Fatal("No se pudo cargar el archivo de configuracion")
	}

	//Solicito PID y PC para ejecutar Instrucciones
	//1. Creo el handler
	mux := http.NewServeMux()
	mux.HandleFunc("/cpu/kernel", utils.RecibirContextoDeKernel)
	mux.HandleFunc("/cpu/interrupcion", utils.RecibirInterrupcion)

	//2. Uso una goroutine para que no se bloquee el modulo
	go func() {
		log.Printf("Escuchando en %s:%d...", globals.ClientConfig.IpSelf, globals.ClientConfig.PortSelf)
		err = http.ListenAndServe(fmt.Sprintf("%s:%d", globals.ClientConfig.IpSelf, globals.ClientConfig.PortSelf), mux)
		if err != nil {
			log.Fatalf("Error al iniciar servidor HTTP: %v", err)
		}
	}()

	//Las CPUs deberán conectarse al Kernel (destino)
	//3. Envi0 su IP, su PUERTO y su ID. (self)
	utils.EnviarIpPuertoIDAKernel(globals.ClientConfig.IpKernel, globals.ClientConfig.PortKernel, globals.ClientConfig.IpSelf, globals.ClientConfig.PortSelf, ID)

	//4. Evito que el modulo termine
	select {}
}
