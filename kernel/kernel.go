package main

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/utils"
	"log"
	"net/http"
)

func main() {
	//Cargar configuracion inicial
	globals.KernelConfig = utils.Config("config.json")

	if globals.KernelConfig == nil {
		log.Fatal("No se pudo cargar el archivo de configuración")
	}

	var portKernel = globals.KernelConfig.PortKernel
	//var ipMemory = globals.KernelConfig.IpMemory
	//var portMemory = globals.KernelConfig.PortMemory

	log.Println("Comenzó ejecucion del Kernel")

	mux := http.NewServeMux()
	//SERVER DE LOS OTROS MODULOS: Escuchar sus mensajes
	mux.HandleFunc("/kernel/io", utils.RecibirMensajeDeIO)
	mux.HandleFunc("/kernel/cpu", utils.RecibirMensajeDeCPU)
	mux.HandleFunc("/kernel/contexto", utils.EnviarContextoACPU)

	fmt.Printf("Servidor escuchando en http://localhost:%d/kernel\n", portKernel)

	address := fmt.Sprintf(":%d", portKernel)
	err := http.ListenAndServe(address, mux)
	if err != nil {
		panic(err)
	}

	//TODO FUNCIONES DE CLIENTE. CONEXION CON OTROS MODULOS:
	//enviar mensaje
}
