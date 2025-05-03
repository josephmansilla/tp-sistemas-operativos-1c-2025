package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/kernel/pcb"
	"github.com/sisoputnfrba/tp-golang/kernel/syscalls"
	"github.com/sisoputnfrba/tp-golang/kernel/utils"
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
	//NACHITO ACA HAGO UN REFACTOR DE LOS LOGs
	err1 := utils.ConfigureLogger("kernel.log", "INFO")
	if err1 != nil {
		fmt.Println("No se pudo crear el logger -", err.Error())
		os.Exit(1)
	}
	utils.Debug("Logger creado")

	// ----------------------------------------------------
	// ---------- PARTE CARGA DEL CONFIG ------------------
	// ----------------------------------------------------
	//NACHITO ACA HAGO UN REFACTOR DE LOS CONFIG
	configData, err := os.ReadFile("config.json")
	if err != nil {
		utils.Fatal("No se pudo leer el archivo de configuración - %v", err.Error())
	}

	err = json.Unmarshal(configData, &utils.Config)
	if err != nil {
		utils.Fatal("No se pudo parsear el archivo de configuración - %v", err.Error())
	}

	if err = utils.Config.Validate(); err != nil {
		utils.Fatal("La configuración no es válida - %v", err.Error())
	}

	err = utils.SetLevel(utils.Config.LogLevel)
	if err != nil {
		utils.Fatal("No se pudo leer el log-level - %v", err.Error())
	}

	// ----------------------------------------------------
	// ---------- ENVIAR PSEUDOCODIGO A MEMORIA -----------
	// ----------------------------------------------------
	InitFirstProcess(archivoPseudocodigo, tamanioProceso)

	// ------------------------------------------------------
	// ---------- ESCUCHO REQUESTS DE CPU E IO --------------
	// ------------------------------------------------------
	mux := http.NewServeMux()
	mux.HandleFunc("/kernel/io", utils.RecibirMensajeDeIO)
	mux.HandleFunc("/kernel/cpu", utils.RecibirMensajeDeCPU)

	//SYSCALLS
	mux.HandleFunc("/kernel/contexto_interrumpido", syscalls.ContextoInterrumpido)
	mux.HandleFunc("/kernel/init_proc", syscalls.InitProcess)
	mux.HandleFunc("/kernel/exit", syscalls.Exit)
	mux.HandleFunc("/kernel/dump_memory", syscalls.DumpMemory)
	mux.HandleFunc("/kernel/syscallIO", syscalls.Io)

	fmt.Printf("Servidor escuchando en http://localhost:%d/kernel\n", utils.Config.KernelPort)

	address := fmt.Sprintf(":%d", utils.Config.KernelPort)
	err = http.ListenAndServe(address, mux)
	if err != nil {
		panic(err)
	}

	utils.ColaNuevo = utils.Queue[*pcb.PCB]{}
	utils.ColaBLoqueado = utils.Queue[*pcb.PCB]{}
	utils.ColaSalida = utils.Queue[*pcb.PCB]{}
	utils.ColaEjecutando = utils.Queue[*pcb.PCB]{}
	utils.ColaReady = utils.Queue[*pcb.PCB]{}
	utils.ColaBloqueadoSuspendido = utils.Queue[*pcb.PCB]{}
	utils.ColaSuspendidoReady = utils.Queue[*pcb.PCB]{}

	//TODO
	//1.funcion que cree primer proceso desde los argumentos del main
	//2.inicilizar todas las colas vacias, tipo de dato punteros a PCB y TCB(hilos)
	//3.fucncion que inicie planificacion largo plazo inicialmente parada esperando un enter desde la consola
	//4.inicialiar colas que representen los estados new, ready, bloqueado, suspendido blog, suspendido ready, ejecutando.

	fmt.Printf("Termine de Ejecutar")
}

func InitFirstProcess(fileName string, processSize int) {
	// Crear el PCB para el proceso inicial
	pid := 0 // Asignar el primer PID como 0 (puedes cambiar según la lógica de PID en tu sistema)
	pcb1 := pcb.PCB{
		PID: pid,
		PC:  0,
		ME:  make(map[string]int),
		MT:  make(map[string]int), //ACA LAS LISTAS PARA LA TRAZABILIDAD LAS INICIALIZO VAcias
	}

	log.Printf("## (<%d>:0) Se crea el proceso - Estado: NEW", pid)

	// Agregar el PCB a la lista de PCBs en el kernel
	utils.ColaNuevo.Add(&pcb1)
	// LE AVISO A MEMORIA QUE SE CREO UN NUEVO PROCESO
	request := utils.RequestToMemory{
		Thread:    utils.Thread{PID: utils.Pid(pid)},
		Type:      utils.CreateProcess,
		Arguments: []string{fileName, strconv.Itoa(processSize)}, // aca le envio como argumentos el nombre del archivo y el tamaño del proceso como strings
	}
	for {
		err := utils.SendMemoryRequest(request)
		if err != nil {
			utils.Error("Error al enviar request a memoria: %v", err)
			//<-kernelsync.InitProcess // Espera a que finalice otro proceso antes de intentar de nuevo
		} else {
			utils.Debug("Hay espacio disponible en memoria")
			break
		}
	}
}
