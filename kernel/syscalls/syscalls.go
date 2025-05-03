package syscalls

import (
	"github.com/sisoputnfrba/tp-golang/kernel/pcb"
	"log"
	"net/http"
	"strconv"

	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/utils"
	"github.com/sisoputnfrba/tp-golang/utils/data"
)

// Body JSON a recibir
type MensajeInit struct {
	PID      int    `json:"pid"`
	PC       int    `json:"pc"`
	Filename string `json:"filename"`
	Tamanio  int    `json:"tamanio_memoria"`
}

type MensajeIo struct {
	PID      int    `json:"pid"`
	PC       int    `json:"pc"`
	Nombre   string `json:"nombre"`
	Duracion int    `json:"duracion"`
}

func ContextoInterrumpido(w http.ResponseWriter, r *http.Request) {
	// tu código...
}

func InitProcess(w http.ResponseWriter, r *http.Request) {
	var mensajeRecibido MensajeInit
	if err := data.LeerJson(w, r, &mensajeRecibido); err != nil {
		return //hubo error
	}

	pid := mensajeRecibido.PID
	pc := mensajeRecibido.PC
	filename := mensajeRecibido.Filename
	tamanio := mensajeRecibido.Tamanio
	utils.Info("Se ha recibido: PID: %d PC: %d Filename: %s Tamaño Memoria: %d", pid, pc, filename, tamanio)
	utils.Info("Syscall recibida: “## (<%d>) - Solicitó syscall: <INIT_PROC>”", pid)

	// Crear el PCB para el proceso inicial
	pcb1 := pcb.PCB{
		PID: pid,
		PC:  pc,
		ME:  make(map[string]int),
		MT:  make(map[string]int), //ACA LAS LISTAS PARA LA TRAZABILIDAD LAS INICIALIZO VAcias
	}

	log.Printf("## (<%v>:0) Se crea el proceso - Estado: NEW", pid)

	// Agregar el PCB a la lista de PCBs en el kernel
	utils.ColaNuevo.Add(&pcb1)
	// LE AVISO A MEMORIA QUE SE CREO UN NUEVO PROCESO
	request := utils.RequestToMemory{
		Thread:    utils.Thread{PID: utils.Pid(pid)},
		Type:      utils.CreateProcess,
		Arguments: []string{filename, strconv.Itoa(tamanio)}, // aca le envio como argumentos el nombre del archivo y el tamaño del proceso como strings
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

func Exit(w http.ResponseWriter, r *http.Request) {

}

func DumpMemory(w http.ResponseWriter, r *http.Request) {

}

func Io(w http.ResponseWriter, r *http.Request) {
	var mensajeRecibido MensajeIo
	if err := data.LeerJson(w, r, &mensajeRecibido); err != nil {
		return //hubo error
	}
	pid := mensajeRecibido.PID

	utils.Info("Se ha recibido: Nombre: %s Duracion: %d", mensajeRecibido.Nombre, mensajeRecibido.Duracion)
	utils.Info("Syscall recibida: “## (<%d>) - Solicitó syscall: <IO>”", pid)

	//Habra que buscar en lista de IOs si existe...
	if globals.IO.Nombre == mensajeRecibido.Nombre {
		utils.EnviarContextoIO(globals.IO.Ip, globals.IO.Puerto, pid, mensajeRecibido.Duracion)
	}
}
