package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/kernel/algoritmos"
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/pcb"
	"github.com/sisoputnfrba/tp-golang/kernel/syscalls"
	"github.com/sisoputnfrba/tp-golang/kernel/utils"
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

	globals.ColaNuevo = algoritmos.Queue[*pcb.PCB]{}
	globals.ColaBLoqueado = algoritmos.Queue[*pcb.PCB]{}
	globals.ColaSalida = algoritmos.Queue[*pcb.PCB]{}
	globals.ColaEjecutando = algoritmos.Queue[*pcb.PCB]{}
	globals.ColaReady = algoritmos.Queue[*pcb.PCB]{}
	globals.ColaBloqueadoSuspendido = algoritmos.Queue[*pcb.PCB]{}
	globals.ColaSuspendidoReady = algoritmos.Queue[*pcb.PCB]{}

	// ----------------------------------------------------
	// ---------- ENVIAR PSEUDOCODIGO A MEMORIA -----------
	// ----------------------------------------------------
	utils.CrearProceso(archivoPseudocodigo, tamanioProceso)
	// ESTA FUNCIÓN ES LA QUE TIENE QUE TENER TODA LA LÓGICA QUE TIENE
	//

	// ESTO NO SE HACE TODAVIA PRIMERO HAY QUE CONSULTAR LA MEMORIA.
	// EN INIT FIRST PROCESS DEBERÍA ESTAR LA FUNCION DE ARRIBA
	// ES PARA EL PROXIMO CHECKPOINT
	// InitFirstProcess(archivoPseudocodigo, tamanioProceso)

	// ------------------------------------------------------
	// ---------- ESCUCHO REQUESTS DE CPU E IO (Puertos) ----
	// ------------------------------------------------------
	mux := http.NewServeMux()
	mux.HandleFunc("/kernel/io", utils.RecibirMensajeDeIO)
	mux.HandleFunc("/kernel/cpu", utils.RecibirMensajeDeCPU)

	// Falta implementaciòn en kernel
	//mux.HandleFunc("/kernel/cpu", utils.EnviarIpPuertoIDAKernel)

	// ------------------------------------------------------
	// --------------------- SYSCALLS -----------------------
	// ------------------------------------------------------
	mux.HandleFunc("/kernel/contexto_interrumpido", syscalls.ContextoInterrumpido)
	mux.HandleFunc("/kernel/init_proc", syscalls.InitProcess)
	mux.HandleFunc("/kernel/exit", syscalls.Exit)
	mux.HandleFunc("/kernel/dump_memory", syscalls.DumpMemory)
	mux.HandleFunc("/kernel/syscallIO", syscalls.Io)

	fmt.Printf("Servidor escuchando en http://localhost:%d/kernel\n", globals.KConfig.KernelPort)

	address := fmt.Sprintf(":%d", globals.KConfig.KernelPort)
	err = http.ListenAndServe(address, mux)
	if err != nil {
		panic(err)
	}

	//TODO
	//1.funcion que cree primer proceso desde los argumentos del main
	//2.inicilizar todas las colas vacias, tipo de dato punteros a PCB y TCB(hilos)
	//3.fucncion que inicie planificacion largo plazo inicialmente parada esperando un enter desde la consola
	//4.inicialiar colas que representen los estados new, ready, bloqueado, suspendido blog, suspendido ready, ejecutando.

	fmt.Printf("FIN DE EJECUCION")
}
