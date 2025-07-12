package main

import (
	"bufio"
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
	// Crear logger
	err := logger.ConfigureLogger("memoria.log", "INFO")
	if err != nil {
		fmt.Println("No se pudo crear el logger -", err.Error())
		os.Exit(1)
	}

	// Cargar y validar config desde config.json
	g.MemoryConfig = g.ConfigMemoria()

	err = logger.SetLevel(g.MemoryConfig.LogLevel)
	if err != nil {
		logger.Fatal("No se pudo establecer el log-level: %v", err)
	}

	port := g.MemoryConfig.PortMemory
	logger.Info("======== Comenzó la ejecución de Memoria ========")
	logger.Info("Servidor escuchando en http://localhost:%d/memoria", port)

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

	mux.HandleFunc("/memoria/suspension", conex.SuspensionProcesoHandler)

	mux.HandleFunc("/memoria/desuspension", conex.DesuspensionProcesoHandler)

	mux.HandleFunc("/memoria/dump", conex.MemoriaDumpHandler)

	mux.HandleFunc("/memoria/finalizacionProceso", conex.FinalizacionProcesoHandler)

	direccion := fmt.Sprintf(":%d", port)
	errListenAndServe := http.ListenAndServe(direccion, mux)
	if errListenAndServe != nil {
		panic(errListenAndServe)
	}

}
