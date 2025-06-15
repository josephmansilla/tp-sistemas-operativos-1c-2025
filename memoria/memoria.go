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
	if len(os.Args) < 1 {
		fmt.Println("Falta el parametro: identificador del config de Memoria")
		os.Exit(1)
	}

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
	configPath := fmt.Sprintf("config_%s.json", os.Args[0])
	configData, err := os.ReadFile(configPath)
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

	mux.HandleFunc("/memoria/inicializacionProceso", conex.InicializacionProcesoHandler)

	mux.HandleFunc("/memoria/obtenerInstruccion", conex.ObtenerInstruccionHandler)

	mux.HandleFunc("/memoria/espaciolibre", conex.ObtenerEspacioLibreHandler)

	mux.HandleFunc("/memoria/tabla", conex.EnviarEntradaPaginaHandler)

	mux.HandleFunc("/memoria/leerEntradaPagina", adm.LeerPaginaCompletaHandler)

	mux.HandleFunc("/memoria/actualizarEntradaPagina", adm.ActualizarPaginaCompletaHandler)

	mux.HandleFunc("/memoria/lectura", conex.LeerEspacioUsuarioHandler)

	mux.HandleFunc("/memoria/escritura", conex.EscribirEspacioUsuarioHandler)

	mux.HandleFunc("/memoria/suspension", adm.SuspensionProcesoHandler)

	mux.HandleFunc("/memoria/desuspension", adm.DesuspensionProcesoHandler)

	mux.HandleFunc("/memoria/dump", conex.MemoriaDumpHandler)

	mux.HandleFunc("/memoria/finalizacionProceso", conex.FinalizacionProcesoHandler)

	direccion := fmt.Sprintf(":%d", portMemory)
	errListenAndServe := http.ListenAndServe(direccion, mux)
	if errListenAndServe != nil {
		panic(errListenAndServe)
	}

}
