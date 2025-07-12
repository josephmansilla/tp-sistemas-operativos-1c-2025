package main

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/traducciones"
	"github.com/sisoputnfrba/tp-golang/cpu/utils"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		logger.Info("Uso: go run cpu.go <ID_CPU>")
		os.Exit(1)
	}
	globals.ID = os.Args[1]

	// Logger
	logFileName := fmt.Sprintf("logs/cpu_%s.log", globals.ID)
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		logger.Error("Error al crear archivo de log para CPU %s: %v\n", globals.ID, err)
		os.Exit(1)
	}
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)
	logger.ConfigureLogger(logFileName, "INFO")

	logger.Info("======== Comenzó la ejecución del CPU con ID: %s ========", globals.ID)

	// Cargar config desde "config.json" raíz
	config := utils.Config()

	// Asignar IP y puerto según ID
	switch globals.ID {
	case "1":
		config.IpSelf = config.IpCPU
		config.PortSelf = config.PortCPU1
	case "2":
		config.IpSelf = config.IpCPU
		config.PortSelf = config.PortCPU2
	case "3":
		config.IpSelf = config.IpCPU
		config.PortSelf = config.PortCPU3
	case "4":
		config.IpSelf = config.IpCPU
		config.PortSelf = config.PortCPU4
	default:
		logger.Fatal("ID de CPU inválido: %s", globals.ID)
	}

	globals.ClientConfig = config

	traducciones.Max = globals.ClientConfig.CacheEntries
	traducciones.InitTLB()
	traducciones.InitCache()

	// Configuración de Memoria
	err = utils.RecibirConfiguracionMemoria(globals.ClientConfig.IpMemory, globals.ClientConfig.PortMemory)
	if err != nil {
		logger.Fatal("Error al obtener la configuración de memoria: %v", err)
	}

	// Servidor HTTP
	mux := http.NewServeMux()
	mux.HandleFunc("/cpu/kernel", utils.RecibirContextoDeKernel)
	mux.HandleFunc("/cpu/interrupcion", utils.RecibirInterrupcion)

	fmt.Printf("Servidor escuchando en http://localhost:%d/cpu\n", globals.ClientConfig.PortSelf)
	go func() {
		logger.Info("Escuchando en %s:%d...", globals.ClientConfig.IpSelf, globals.ClientConfig.PortSelf)
		err = http.ListenAndServe(fmt.Sprintf("%s:%d", globals.ClientConfig.IpSelf, globals.ClientConfig.PortSelf), mux)
		if err != nil {
			logger.Fatal("Error al iniciar servidor HTTP: %v", err)
		}
	}()

	// Conexión con Kernel
	utils.EnviarIpPuertoIDAKernel(globals.ClientConfig.IpKernel, globals.ClientConfig.PortKernel, globals.ClientConfig.IpSelf, globals.ClientConfig.PortSelf, globals.ID)

	// Bloqueo el main
	select {}
}
