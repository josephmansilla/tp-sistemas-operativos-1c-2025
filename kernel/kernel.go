package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/utils"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Falta el parametro: nombre del archivo de pseudocodigo")
		os.Exit(1)
	} else if len(os.Args) < 3 {
		fmt.Println("Falta el parametro: tamaño del proceso")
		os.Exit(1)
	}

	archivoPseudocodigo := os.Args[1]
	tamanioStr := os.Args[2] //Convertir a Double

	tamanioProceso, err := strconv.Atoi(tamanioStr)
	if err != nil {
		fmt.Printf("Tamaño del proceso inválido: %s\n", tamanioStr)
		os.Exit(1)
	}

	logFileName := "kernel.log"
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Error al crear archivo de log para Kenrel: %v\n", err)
		os.Exit(1)
	}
	log.SetOutput(logFile)

	log.Printf("Nombre del archivo de pseudocodigo: %s\n", archivoPseudocodigo)
	log.Printf("Tamaño del proceso: %d\n", tamanioProceso)

	//Cargar configuracion inicial
	globals.KernelConfig = utils.Config("config.json")

	if globals.KernelConfig == nil {
		log.Fatal("No se pudo cargar el archivo de configuración")
	}

	var portKernel = globals.KernelConfig.PortKernel
	//var ipMemory = globals.KernelConfig.IpMemory
	//var portMemory = globals.KernelConfig.PortMemory

	log.Println("Comenzó ejecucion del Kernel")
	//TODO Al iniciar el módulo, se creará un proceso inicial para que este lo planifique...

	mux := http.NewServeMux()
	//SERVER DE LOS OTROS MODULOS: Escuchar sus mensajes
	mux.HandleFunc("/kernel/io", utils.RecibirMensajeDeIO)
	mux.HandleFunc("/kernel/cpu", utils.RecibirMensajeDeCPU)
	
	//mux.HandleFunc("/kernel/contexto", utils.EnviarContextoACPU)

	fmt.Printf("Servidor escuchando en http://localhost:%d/kernel\n", portKernel)

	address := fmt.Sprintf(":%d", portKernel)
	err = http.ListenAndServe(address, mux)
	if err != nil {
		panic(err)
	}

	//TODO FUNCIONES DE CLIENTE. CONEXION CON OTROS MODULOS:
	//enviar mensaje

	//1.funcion que cree primer proceso desde los argumentos del main
	//2.inicilizar todas las colas vacias, tipo de dato punteros a PCB y TCB(hilos)
	//3.fucncion que inicie planificacion largo plazo inicialmente parada esperando un enter desde la consola
	//4.inicialiar colas que representen los estados new, ready, bloqueado, suspendido blog, suspendido ready, ejecutando.
}
