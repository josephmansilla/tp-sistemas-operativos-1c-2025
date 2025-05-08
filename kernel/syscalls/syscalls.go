package syscalls

import (
	"fmt"
	"log"
	"net/http"
	"time"

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

type MensajeSyscall struct {
	PID int `json:"pid"`
	PC  int `json:"pc"`
}

type MensajeDUMP struct {
	PID       int    `json:"pid"`
	Timestamp string `json:"timestamp"`
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

	logger.Info("## (<%d>) - Solicitó syscall: <INIT_PROC>", pid)
	logger.Info("Se ha recibido: PID: %d PC: %d Filename: %s Tamaño Memoria: %d", pid, pc, filename, tamanio)

	//Planificador Largo Plazo
	planificadores.CrearProceso(filename, tamanio)
}

func Exit(w http.ResponseWriter, r *http.Request) {
	var mensajeRecibido MensajeSyscall

	logger.Info("## (<%d>) - Solicitó syscall: <EXIT>", mensajeRecibido.PID)

	//Planificador Largo Plazo
	planificadores.FinalizarProceso(mensajeRecibido.PID)
}

func DumpMemory(w http.ResponseWriter, r *http.Request) {
	var mensaje MensajeDUMP
	if err := data.LeerJson(w, r, &mensaje); err != nil {
		return // El error ya fue respondido por LeerJson
	}

	// Generar el timestamp en formato yyyyMMddTHHmmss
	timestamp := time.Now().Format("20060102T150405")

	// Crear mensaje para enviar a Memoria
	req := MensajeDUMP{
		PID:       mensaje.PID,
		Timestamp: timestamp,
	}

	// Armar URL del módulo Memoria
	url := fmt.Sprintf("http://%s:%d/memoria/dump", globals.Config.MemoryAddress, globals.Config.MemoryPort)

	// Usar helper para enviar datos
	if err := data.EnviarDatos(url, req); err != nil {
		log.Printf("Error enviando dump a Memoria: %v", err)
		http.Error(w, "Error comunicando con Memoria", http.StatusInternalServerError)
		return
	}

	log.Printf("Se envió correctamente el pedido de dump del PID %d a Memoria", req.PID)
	w.WriteHeader(http.StatusOK)
}

func Io(w http.ResponseWriter, r *http.Request) {
	var mensajeRecibido MensajeIo
	if err := data.LeerJson(w, r, &mensajeRecibido); err != nil {
		return
	}
	pid := mensajeRecibido.PID
	nombre := mensajeRecibido.Nombre

	logger.Info("Syscall recibida: “## (<%d>) - Solicitó syscall: <IO>”", pid)

	// Aquí bloqueas el mutex mientras esperas a que el IO se registre
	globals.IOMu.Lock()
	ioData, ok := globals.IOs[nombre]
	for !ok {
		globals.IOCond.Wait()
		ioData, ok = globals.IOs[nombre] // reintenta obtenerlo
	}
	globals.IOMu.Unlock()

	logger.Info("Nombre IO: %s Duracion: %d", ioData.Nombre, mensajeRecibido.Duracion)

	comunicacion.EnviarContextoIO(ioData.Nombre, pid, mensajeRecibido.Duracion)
}
