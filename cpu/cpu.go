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

	logFileName := fmt.Sprintf("logs/cpu_%s.log", ID)
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Printf("Error al crear archivo de log para CPU %s: %v\n", ID, err)
		os.Exit(1)
	}
	log.SetOutput(logFile)

	// "SOY UN GORDO QUE NO QUIERO ANDAR"
	/*err = logger.SetLevel(globals.ClientConfig.LogLevel)
	if err != nil {
		logger.Fatal("No se pudo leer el log-level - %v", err.Error())
	}*/

	log.Printf("=========================================================")
	log.Printf("======== Comenzo la ejecucion del CPU con ID: %s ========", ID)
	log.Printf("=========================================================\n")

	//CPU CLIENTE
	configpath := fmt.Sprintf("configs/cpu_%sconfig.json", ID)
	globals.ClientConfig = utils.Config(configpath)

	if globals.ClientConfig == nil {
		log.Fatal("No se pudo cargar el archivo de configuracion")
	}

	//Solicito la configuracion de memoria
	/*err = utils.ConsultarConfiguracionMemoria(globals.ClientConfig.IpMemory, globals.ClientConfig.PortMemory)
	if err != nil {
		log.Fatalf("Error al obtener la configuración de memoria: %v", err)
	}*/
	//Solicito PID y PC para ejecutar Instrucciones
	//1. Creo el handler
	mux := http.NewServeMux()
	mux.HandleFunc("/cpu/kernel", utils.RecibirContextoDeKernel)
	mux.HandleFunc("/cpu/interrupcion", utils.RecibirInterrupcion)

	fmt.Printf("Servidor escuchando en http://localhost:%d/cpu\n", globals.ClientConfig.PortSelf)
	utils.SimularSyscallInitProcess("127.0.0.1", 8081, 0, 0, "holiii", 64)
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
