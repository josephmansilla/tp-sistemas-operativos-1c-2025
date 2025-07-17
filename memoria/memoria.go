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
	if len(os.Args) < 2 {
		logger.Fatal("Falta parámetro: CONFIG")
		os.Exit(1)
	}
	config := os.Args[1]
	g.MemoryConfig = g.ConfigMemoria(config)

	logger.Info("======== Comenzó la ejecución de Memoria ========")
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

	adm.InicializarMemoriaPrincipal()

	mux := http.NewServeMux()
	// ======================== CONFIGS Y CONSULTAS ========================
	mux.HandleFunc("/memoria/configuracion", conex.EnviarConfiguracionMemoriaHandler)
	mux.HandleFunc("/memoria/espaciolibre", conex.ObtenerEspacioLibreHandler)
	// ======================== PEDIDOS CPU ========================
	mux.HandleFunc("/memoria/obtenerInstruccion", conex.ObtenerInstruccionHandler)
	mux.HandleFunc("/memoria/tabla", conex.EnviarEntradaPaginaHandler)
	// ======================== INIT ========================
	mux.HandleFunc("/memoria/inicializacionProceso", conex.InicializacionProcesoHandler)
	// ======================== KILL ========================
	mux.HandleFunc("/memoria/finalizacionProceso", conex.FinalizacionProcesoHandler)
	// ======================== DUMP ========================
	mux.HandleFunc("/memoria/dump", conex.MemoriaDumpHandler)
	// ====================== SUSPENSIÓN =====================
	mux.HandleFunc("/memoria/suspension", conex.SuspensionProcesoHandler)
	mux.HandleFunc("/memoria/desuspension", conex.DesuspensionProcesoHandler)
	// ======================== LECTURA Y ESCRITURA ========================
	mux.HandleFunc("/memoria/lectura", conex.LeerEspacioUsuarioHandler)
	mux.HandleFunc("/memoria/escritura", conex.EscribirEspacioUsuarioHandler)

	errListenAndServe := http.ListenAndServe(fmt.Sprintf(":%d", g.MemoryConfig.PortMemory), mux)
	if errListenAndServe != nil {
		panic(errListenAndServe)
	}

}
