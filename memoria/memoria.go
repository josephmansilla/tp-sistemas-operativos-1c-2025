package main

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/memoria/utils"
	"log"
	"net/http"
	"os"
)

func main() {

	// ----------------------------------------------------
	// ---------- PRIMERA PARTE CARGA DEL CONFIG ----------
	// ----------------------------------------------------
	globals.MemoryConfig = utils.Config("config.json")
	if globals.MemoryConfig == nil {
		log.Fatal("No se pudo cargar el archivo de configuración")
	}
	var portMemory = globals.MemoryConfig.PortMemory
	log.Println("Comenzó ejecucion de la memoria")

	// ----------------------------------------------------
	// ----------- CARGO LOGS DE MEMORIA EN TXT -----------
	// ----------------------------------------------------

	logFileName := "memoria.log"
	logFile, errLogFile := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if errLogFile != nil {
		fmt.Printf("Error al crear archivo de log para la Memoria: %v\n", errLogFile)
		os.Exit(1)
	}
	log.SetOutput(logFile)

	// ------------------------------------------------------
	// ---------- ESCUCHO REQUESTS DE CPU Y KERNEL ----------
	// ------------------------------------------------------

	mux := http.NewServeMux()
	// ESTÁ ESPERANDO LOS MENSAJES DE LOS OTROS MODULOS
	mux.HandleFunc("/memoria/cpu", utils.RecibirMensajeDeCPU)
	mux.HandleFunc("/memoria/kernel", utils.RecibirMensajeDeKernel)
	mux.HandleFunc("/memoria/instruccion", utils.ObtenerInstruccion)
	mux.HandleFunc("/memoria/espaciolibre", utils.ObtenerEspacioLibreMock)
	//mux.HandleFunc("/memoria/frame", utils.algo)
	//mux.HandleFunc("memoria/pagina", utils.algo)
	mux.HandleFunc("/memoria/createProcess", utils.CreateProcess)

	//mux.HandleFunc("/memoria/cpu", utils.CreacionProceso)

	fmt.Printf("Servidor escuchando en http://localhost:%d/memoria\n", portMemory)

	direccion := fmt.Sprintf(":%d", portMemory)
	errListenAndServe := http.ListenAndServe(direccion, mux)
	if errListenAndServe != nil {
		panic(errListenAndServe)
	}

	fmt.Printf("Termine de Ejecutar")

}
