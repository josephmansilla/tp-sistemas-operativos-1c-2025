package main

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/memoria/utils"
	logger "github.com/sisoputnfrba/tp-golang/utils/logger"
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
		logger.Fatal("No se pudo cargar el archivo de configuración", nil)
	}
	var portMemory = globals.MemoryConfig.PortMemory

	// ----------------------------------------------------
	// ----------- CARGO LOGS DE MEMORIA EN TXT -----------
	// ----------------------------------------------------

	logFileName := "memoria.log"
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		logger.Error("Error al crear archivo de log para la Memoria: %v\n", err)
		os.Exit(1)
	}
	log.SetOutput(logFile)
	err = logger.SetLevel(globals.MemoryConfig.LogLevel)
	if err != nil {
		logger.Fatal("No se pudo leer el log-level - %v", err.Error())
	}

	log.Printf("=================================================")
	log.Printf("======== Comenzo la ejecucion de Memoria ========")
	log.Printf("=================================================\n")
	fmt.Printf("Servidor escuchando en http://localhost:%d/memoria\n", portMemory)
	log.Printf("Servidor escuchando en http://localhost:%d/memoria\n", portMemory)

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
	//mux.HandleFunc("memoria/configuracion", utils.algo)
	mux.HandleFunc("/memoria/createProcess", utils.CreateProcess)

	//mux.HandleFunc("/memoria/cpu", utils.CreacionProceso)

	direccion := fmt.Sprintf(":%d", portMemory)
	errListenAndServe := http.ListenAndServe(direccion, mux)
	if errListenAndServe != nil {
		panic(errListenAndServe)
	}

	log.Println("======== Final de Ejecución memoria ========")

}
