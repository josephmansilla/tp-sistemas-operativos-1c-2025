package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	adm "github.com/sisoputnfrba/tp-golang/memoria/administracion"
	conex "github.com/sisoputnfrba/tp-golang/memoria/conexiones"
	g "github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"net/http"
	"os"
	"strings"
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
	err = json.Unmarshal(configData, &g.MemoryConfig)
	if err != nil {
		logger.Fatal("No se pudo parsear el archivo dgguración - %v", err.Error())
	}
	if err = g.MemoryConfig.Validate(); err != nil {
		logger.Fatal("La configuración no es válida - %v", err.Error())
	}
	err = logger.SetLevel(g.MemoryConfig.LogLevel)
	if err != nil {
		logger.Fatal("No se pudo leer el log-level - %v", err.Error())
	}

	var portMemory = g.MemoryConfig.PortMemory

	logger.Info("======== Comenzo la ejecucion de Memoria ========")

	logger.Info("Servidor escuchando en http://localhost:%d/memoria\n", portMemory)

	// ----------------------------------------------------------
	// --------- INICIALIZO LAS ESTRUCTURAS NECESARIAS  ---------
	// ----------------------------------------------------------
	logger.Error("Escribí exit para finalizar el módulo de Memoria")
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			texto := strings.ToLower(scanner.Text())
			if texto == "exit" {
				logger.Info("======== Final de Ejecución memoria ========")
				os.Exit(0)
			}
		}
	}()

	adm.InicializarMemoriaPrincipal()

	// ------------------------------------------------------
	// ---------- ESCUCHO REQUESTS DE CPU Y KERNEL ----------
	// ------------------------------------------------------

	mux := http.NewServeMux()
	mux.HandleFunc("/memoria/configuracion", conex.EnviarConfiguracionMemoriaHandler)

	mux.HandleFunc("/memoria/cpu", conex.RecibirMensajeDeCPUHandler)
	mux.HandleFunc("/memoria/kernel", conex.RecibirMensajeDeKernelHandler)
	mux.HandleFunc("/memoria/InicializacionProceso", conex.InicializacionProcesoHandler) // TODO: HACER CONEXIONES CON KERNEL DESPUES DEL MERGE

	mux.HandleFunc("/memoria/instruccion", conex.ObtenerInstruccion)
	// TODO: deberia devoler la instruccion que piden

	mux.HandleFunc("/memoria/espaciolibre", conex.ObtenerEspacioLibreHandler)
	mux.HandleFunc("/memoria/tabla", conex.EnviarEntradaPaginaHandler)

	// TODO: USTEDES DEBEN IMPLEMENTAR ESTAS FUNCIONES
	mux.HandleFunc("/memoria/LeerEntradaPagina", adm.LeerPaginaCompletaHandler)
	mux.HandleFunc("/memoria/ActualizarEntrada", adm.ActualizarPaginaCompletaHandler)
	mux.HandleFunc("/memoria/lectura", conex.LeerEspacioUsuarioHandler)
	// TODO: debe responder a CPU el valor de una dirección física con el delay indicado en Memory Delay
	mux.HandleFunc("/memoria/escritura", conex.EscribirEspacioUsuarioHandler)
	// TODO: recibe PID y tamaño, se crea escructuras, asigna frames y logear.
	// TODO: debe indicarle al CPU que fue éxitoso con el delay indicado en Memory Delay

	mux.HandleFunc("/memoria/suspension", adm.SuspensionProcesoHandler)
	mux.HandleFunc("/memoria/desuspension", adm.DesuspensionProcesoHandler)

	mux.HandleFunc("/memoria/dump", conex.MemoriaDumpHandler)
	mux.HandleFunc("/memoria/finalizacionProceso", conex.FinalizacionProcesoHandler)
	// TODO: debe liberar recursos y escructuras y logear metricas

	//mux.HandleFunc("/memoria/frame", utils.algo)
	//mux.HandleFunc("memoria/pagina", utils.algo)

	//mux.HandleFunc("/memoria/cpu", utils.CreacionProceso)

	direccion := fmt.Sprintf(":%d", portMemory)
	errListenAndServe := http.ListenAndServe(direccion, mux)
	if errListenAndServe != nil {
		panic(errListenAndServe)
	}

}
