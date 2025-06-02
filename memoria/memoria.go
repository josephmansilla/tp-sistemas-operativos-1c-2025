package main

import (
	"encoding/json"
	"fmt"
	"github.com/sisoputnfrba/tp-golang/memoria/administracion"
	"github.com/sisoputnfrba/tp-golang/memoria/conexiones"
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"net/http"
	"os"
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

	// ----------------------------------------------------------
	// --------- INICIALIZO LAS ESTRUCTURAS NECESARIAS  ---------
	// ----------------------------------------------------------

	administracion.InicializarMemoriaPrincipal()

	// ------------------------------------------------------
	// ---------- ESCUCHO REQUESTS DE CPU Y KERNEL ----------
	// ------------------------------------------------------

	mux := http.NewServeMux()
	mux.HandleFunc("memoria/configuracion", conexiones.EnviarConfiguracionMemoria)
	// ESTÁ ESPERANDO LOS MENSAJES DE LOS OTROS MODULOS
	mux.HandleFunc("/memoria/cpu", conexiones.RecibirMensajeDeCPU)
	mux.HandleFunc("/memoria/kernel", conexiones.RecibirMensajeDeKernel)
	mux.HandleFunc("/memoria/instruccion", conexiones.ObtenerInstruccion)
	// TODO: deberia devoler la instruccion que piden

	mux.HandleFunc("/memoria/espaciolibre", conexiones.ObtenerEspacioLibreMock)
	// TODO: cambiar la funcion a la que escucha ,, debería devolver la cantidad de frames libres y su tamaño total

	//mux.HandleFunc("/memoria/lectura", utils.LecturaEspacio)
	// TODO: debe responder a CPU el valor de una dirección física con el delay indicado en Memory Delay
	//mux.HandleFunc("/memoria/escritura", utils.EscrituraEspacio)
	// TODO: recibe PID y tamaño, se crea escructuras, asigna frames y logear.
	// TODO: debe indicarle al CPU que fue éxitoso con el delay indicado en Memory Delay
	//mux.HandleFunc("/memoria/suspension", utils.SuspenderProceso)
	//mux.HandleFunc("/memoria/desuspension", utils.DesuspenderProceso)

	//mux.HandleFunc("/memoria/dump", utils.MemoriaDump)
	// TODO: debe liberar recursos y escructuras y logear metricas

	//mux.HandleFunc("/memoria/frame", utils.algo)
	//mux.HandleFunc("memoria/pagina", utils.algo)

	//mux.HandleFunc("/memoria/cpu", utils.CreacionProceso)

	direccion := fmt.Sprintf(":%d", portMemory)
	errListenAndServe := http.ListenAndServe(direccion, mux)
	if errListenAndServe != nil {
		panic(errListenAndServe)
	}

	logger.Info("======== Final de Ejecución memoria ========")

}
