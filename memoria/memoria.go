package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/memoria/utils"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
)

func main() {
	// ----------------------------------------------------
	// ----------- CARGO LOGS DE MEMORIA EN TXT ------------
	// ----------------------------------------------------
	var err = logger.ConfigureLogger("memoria.log", "INFO")
	if err != nil {
		fmt.Println("No se pudo crear el logger -", err.Error())
		os.Exit(1)
	}
	logger.Debug("Logger creado")

	// ----------------------------------------------------
	// ---------- PARTE CARGA DEL CONFIG ------------------
	// ----------------------------------------------------
	configData, err := os.ReadFile("config.json")
	if err != nil {
		logger.Fatal("No se pudo leer el archivo de configuración - %v", err.Error())
	}

	err = json.Unmarshal(configData, &globals.MemoryConfig)
	if err != nil {
		logger.Fatal("No se pudo parsear el archivo de configuración - %v", err.Error())
	}

	if err = globals.MemoryConfig.Validate(); err != nil {
		logger.Fatal("La configuración no es válida - %v", err.Error())
	}

	err = logger.SetLevel(globals.MemoryConfig.LogLevel)
	if err != nil {
		logger.Fatal("No se pudo leer el log-level - %v", err.Error())
	}

	var portMemory = globals.MemoryConfig.PortMemory
	logger.Info("======== Comenzo la ejecucion de Memoria ========")

	fmt.Printf("Servidor escuchando en http://localhost:%d/memoria\n", portMemory)

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

	//mux.HandleFunc("/memoria/cpu", utils.CreacionProceso)

	direccion := fmt.Sprintf(":%d", portMemory)
	errListenAndServe := http.ListenAndServe(direccion, mux)
	if errListenAndServe != nil {
		panic(errListenAndServe)
	}

	logger.Info("======== Final de Ejecución memoria ========")

}
