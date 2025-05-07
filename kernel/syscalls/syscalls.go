package syscalls

import (
	"log"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/kernel/comunicacion"
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/planificadores"
	"github.com/sisoputnfrba/tp-golang/utils/data"
	logger "github.com/sisoputnfrba/tp-golang/utils/logger"
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
	Duracion int    `json:"tiempo"`
	Nombre   string `json:"nombre"`
}

func ContextoInterrumpido(w http.ResponseWriter, r *http.Request) {
	// tu código...
}

func InitProcess(w http.ResponseWriter, r *http.Request) {
	var mensajeRecibido MensajeInit
	if err := data.LeerJson(w, r, &mensajeRecibido); err != nil {
		return
	}

	pid := mensajeRecibido.PID
	pc := mensajeRecibido.PC
	filename := mensajeRecibido.Filename
	tamanio := mensajeRecibido.Tamanio

	logger.Info("Se ha recibido: PID: %d PC: %d Filename: %s Tamaño Memoria: %d", pid, pc, filename, tamanio)
	logger.Info("Syscall recibida: “## (<%d>) - Solicitó syscall: <INIT_PROC>”", pid)

	// Usar CrearProceso del paquete utils
	planificadores.CrearProceso(filename, tamanio)
}

func Exit(w http.ResponseWriter, r *http.Request) {

}

func DumpMemory(w http.ResponseWriter, r *http.Request) {

}

func Io(w http.ResponseWriter, r *http.Request) {
	var mensajeRecibido MensajeIo
	if err := data.LeerJson(w, r, &mensajeRecibido); err != nil {
		return
	}
	pid := mensajeRecibido.PID

	// Aquí bloqueas el mutex mientras esperas a que el IO se registre
	globals.IOMu.Lock()
	for globals.IO.Ip == "" { // Asumiendo que Ip vacía significa que el IO no está conectado
		globals.IOCond.Wait() // Espera hasta que el IO se registre
	}
	globals.IOMu.Unlock()

	logger.Info("Se ha recibido: Nombre: %s Duracion: %d", mensajeRecibido.Nombre, mensajeRecibido.Duracion)
	logger.Info("Syscall recibida: “## (<%d>) - Solicitó syscall: <IO>”", pid)

	comunicacion.EnviarContextoIO(globals.IO.Ip, globals.IO.Puerto, pid, mensajeRecibido.Duracion)
	log.Printf("Operacion de IO enviada correctamente")
}
