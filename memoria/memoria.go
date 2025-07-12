package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	adm "github.com/sisoputnfrba/tp-golang/memoria/administracion"
	conex "github.com/sisoputnfrba/tp-golang/memoria/conexiones"
	g "github.com/sisoputnfrba/tp-golang/memoria/estructuras"
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
	var err = logger.ConfigureLogger("memoria.log", "INFO")
	if err != nil {
		fmt.Println("No se pudo crear el logger -", err.Error())
		os.Exit(1)
	}
	configData, err := os.ReadFile(fmt.Sprintf("configs/%s.json", os.Args[1]))
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

	logger.Info("======== Comenzo la ejecucion de Memoria ========")
	logger.Info("Servidor escuchando en http://localhost:%d/memoria", g.MemoryConfig.PortMemory)

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

	// ========== INICIALIZO ESTRUCTURAS NECESARIAS ==========
	adm.InicializarMemoriaPrincipal()

	// ========== REQUESTS Y ENDPOINTS ==========
	mux := http.NewServeMux()

	// mux.HandleFunc("/memoria/cpu", conex.RecibirMensajeDeCPUHandler)

	// CONFIG Y CONSULTAS
	mux.HandleFunc("/memoria/configuracion", conex.EnviarConfiguracionMemoriaHandler)
	mux.HandleFunc("/memoria/espaciolibre", conex.ObtenerEspacioLibreHandler)

	// INIT
	mux.HandleFunc("/memoria/inicializacionProceso", conex.InicializacionProcesoHandler)
	// DUMP
	mux.HandleFunc("/memoria/dump", conex.MemoriaDumpHandler)
	// MUERTE PROCESO
	mux.HandleFunc("/memoria/finalizacionProceso", conex.FinalizacionProcesoHandler)

	// SWAP
	mux.HandleFunc("/memoria/suspension", conex.SuspensionProcesoHandler)
	mux.HandleFunc("/memoria/desuspension", conex.DesuspensionProcesoHandler)

	// CPU
	mux.HandleFunc("/memoria/obtenerInstruccion", conex.ObtenerInstruccionHandler)
	mux.HandleFunc("/memoria/tabla", conex.EnviarEntradaPaginaHandler)
	mux.HandleFunc("/memoria/leerEntradaPagina", conex.LeerPaginaCompletaHandler)
	mux.HandleFunc("/memoria/actualizarEntradaPagina", conex.ActualizarPaginaCompletaHandler)

	// PRONTAS A MORIR
	mux.HandleFunc("/memoria/lectura", conex.LeerEspacioUsuarioHandler)
	mux.HandleFunc("/memoria/escritura", conex.EscribirEspacioUsuarioHandler)

	direccion := fmt.Sprintf(":%d", g.MemoryConfig.PortMemory)
	errListenAndServe := http.ListenAndServe(direccion, mux)
	if errListenAndServe != nil {
		panic(errListenAndServe)
	}

}
