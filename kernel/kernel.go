package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/sisoputnfrba/tp-golang/kernel/Utils"
	"net/http"
	"os"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/kernel/algoritmos"
	"github.com/sisoputnfrba/tp-golang/kernel/comunicacion"
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/pcb"
	"github.com/sisoputnfrba/tp-golang/kernel/planificadores"
	"github.com/sisoputnfrba/tp-golang/kernel/syscalls"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
)

func main() {
	// ----------------------------------------------------
	// ---------- PARTE CARGA DE PARAMETROS ---------------
	// ----------------------------------------------------
	if len(os.Args) < 2 {
		fmt.Println("Falta el parametro: nombre del archivo de pseudocodigo")
		os.Exit(1)
	} else if len(os.Args) < 3 {
		fmt.Println("Falta el parametro: tamaño del proceso")
		os.Exit(1)
	}

	archivoPseudocodigo := os.Args[1]
	tamanioStr := os.Args[2] //Convertir a Double

	tamanioProceso, err := strconv.Atoi(tamanioStr)
	if err != nil {
		fmt.Printf("Tamaño del proceso inválido: %s\n", tamanioStr)
		os.Exit(1)
	}

	// ----------------------------------------------------
	// ----------- CARGO LOGS DE KERNEL EN TXT ------------
	// ----------------------------------------------------
	err = logger.ConfigureLogger("kernel.log", "INFO")
	if err != nil {
		fmt.Println("No se pudo crear el logger -", err.Error())
		os.Exit(1)
	}
	logger.Debug("Logger creado")

	logger.Info("Comenzo la ejecucion del Kernel")

	// ----------------------------------------------------
	// ---------- PARTE CARGA DEL CONFIG ------------------
	// ----------------------------------------------------
	configData, err := os.ReadFile("config.json")
	if err != nil {
		logger.Fatal("No se pudo leer el archivo de configuración - %v", err.Error())
	}

	err = json.Unmarshal(configData, &globals.KConfig)
	if err != nil {
		logger.Fatal("No se pudo parsear el archivo de configuración - %v", err.Error())
	}

	if err = globals.KConfig.Validate(); err != nil {
		logger.Fatal("La configuración no es válida - %v", err.Error())
	}

	err = logger.SetLevel(globals.KConfig.LogLevel)
	if err != nil {
		logger.Fatal("No se pudo leer el log-level - %v", err.Error())
	}

	//Inicilizar todas las colas vacias, tipo de dato punteros a PCB y TCB(hilos)
	algoritmos.ColaNuevo = algoritmos.Cola[*pcb.PCB]{}
	algoritmos.ColaBloqueado = algoritmos.Cola[*pcb.PCB]{}
	algoritmos.ColaSalida = algoritmos.Cola[*pcb.PCB]{}
	algoritmos.ColaEjecutando = algoritmos.Cola[*pcb.PCB]{}
	algoritmos.ColaReady = algoritmos.Cola[*pcb.PCB]{}
	algoritmos.ColaBloqueadoSuspendido = algoritmos.Cola[*pcb.PCB]{}
	algoritmos.ColaSuspendidoReady = algoritmos.Cola[*pcb.PCB]{}

	// ----------------------------------------------------
	// ---------- ENVIAR PSEUDOCODIGO A MEMORIA -----------
	// ----------------------------------------------------
	//1. Crear primer proceso desde los argumentos del main
	planificadores.CrearPrimerProceso(archivoPseudocodigo, tamanioProceso)
	//planificadores.PlanificarCortoPlazo()
	//planificadores.PlanificadorMedianoPlazo()

	// Inicializar recursos compartidos
	Utils.InicializarMutexes()
	Utils.InicializarCanales()

	// ------------------------------------------------------
	// ---------- ESCUCHO REQUESTS DE CPU E IO (Puertos) ----
	// ------------------------------------------------------
	mux := http.NewServeMux()
	mux.HandleFunc("/kernel/io", comunicacion.RecibirMensajeDeIO)
	mux.HandleFunc("/kernel/cpu", comunicacion.RecibirMensajeDeCPU)

	// ------------------------------------------------------
	// --------------------- SYSCALLS -----------------------
	// ------------------------------------------------------
	mux.HandleFunc("/kernel/contexto_interrumpido", syscalls.ContextoInterrumpido)
	mux.HandleFunc("/kernel/init_proceso", syscalls.InitProcess)
	mux.HandleFunc("/kernel/exit", syscalls.Exit)
	mux.HandleFunc("/kernel/dump_memory", syscalls.DumpMemory)
	mux.HandleFunc("/kernel/syscallIO", syscalls.Io)

	fmt.Printf("Servidor escuchando en http://localhost:%d/kernel\n", globals.KConfig.KernelPort)

	// ------------------------------------------------------
	// ---------- INICIAR PLANIFICADOR DE LARGO PLAZO  ------
	// ------------------------------------------------------

	// Esperar que el usuario presione Enter
	go iniciarLargoPlazo()

	address := fmt.Sprintf(":%d", globals.KConfig.KernelPort)
	err = http.ListenAndServe(address, mux)
	if err != nil {
		panic(err)
	}

	fmt.Printf("FIN DE EJECUCION")
}

func iniciarLargoPlazo() {
	fmt.Println("Presione ENTER para iniciar el Planificador de Largo Plazo...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	planificadores.PlanificadorLargoPlazo()
}
