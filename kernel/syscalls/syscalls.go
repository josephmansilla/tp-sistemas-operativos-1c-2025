package syscalls

import (
	"encoding/json"
	"fmt"
	"github.com/sisoputnfrba/tp-golang/kernel/Utils"
	"github.com/sisoputnfrba/tp-golang/kernel/algoritmos"
	"github.com/sisoputnfrba/tp-golang/kernel/pcb"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/utils/data"
	logger "github.com/sisoputnfrba/tp-golang/utils/logger"
)

// Body JSON a recibir
type MensajeInit struct {
	PID      int    `json:"pid"`
	PC       int    `json:"pc"`
	Filename string `json:"filename"`
	Tamanio  int    `json:"tamanio"`
}

type MensajeIo struct {
	PCB      *pcb.PCB `json:"pcb"`
	Duracion int      `json:"tiempo"`
	Nombre   string   `json:"nombre"`
}

type MensajeSyscall struct {
	PCB *pcb.PCB `json:"pcb"`
}

type MensajeDUMP struct {
	PID       int    `json:"pid"`
	Timestamp string `json:"timestamp"`
}

func ContextoInterrumpido(w http.ResponseWriter, r *http.Request) {
	// tu código...
}

/*
	func InitProcess(w http.ResponseWriter, r *http.Request) {
		var mensajeRecibido MensajeInit
		if err := data.LeerJson(w, r, &mensajeRecibido); err != nil {
			return
		}

		pid := mensajeRecibido.PID //
		pc := mensajeRecibido.PC   //
		filename := mensajeRecibido.Filename
		tamanio := mensajeRecibido.Tamanio

		logger.Info("## (<%d>) - Solicitó syscall: <INIT_PROC>", pid)
		logger.Info("Se ha recibido: PID: %d PC: %d Filename: %s Tamaño Memoria: %d", pid, pc, filename, tamanio)

		//LO QUE BUSCO ACA ES SOLAMENTE CREAR EL PCB Y PONERLO EN LA COLA DE NEW
		planificadores.CrearProceso(filename, tamanio)
		//cuando un proceso Inicia entonces tambien debe inciarse el planificador de largo PLAZO, osea debe intentar pasarse de NEW A READY

}
*/

func InitProcess(w http.ResponseWriter, r *http.Request) {
	// 1) Leer y parsear el JSON entrante (sin usar PID desde la CPU)
	var msg MensajeInit
	if err := data.LeerJson(w, r, &msg); err != nil {
		return
	}

	// 2) Generar el PID dentro del kernel
	pid := globals.GenerarNuevoPID()
	//pc := msg.PC           // opcional, vos no lo usás
	fileName := msg.Filename
	tamanio := msg.Tamanio

	logger.Info("## (<%d>) - Solicitó syscall: <INIT_PROC>", pid)
	logger.Info("Se ha recibido: Filename: %s Tamaño Memoria: %d", fileName, tamanio)

	// 3) Despachar la creación en segundo plano
	go func() {
		// Construir el PCB con el PID generado
		pcbNuevo := pcb.PCB{
			PID:         pid,
			PC:          0,
			FileName:    fileName,
			ProcessSize: tamanio,
			ME:          make(map[string]int),
			MT:          make(map[string]int),
		}
		logger.Info("## (<%d>) Se crea el proceso - Estado: NEW", pid)

		// Encolar en NEW con protección
		Utils.MutexNuevo.Lock()
		algoritmos.ColaNuevo.Add(&pcbNuevo)
		Utils.MutexNuevo.Unlock()
		logger.Debug("PCB <%d> añadido a NEW", pid)

		// Notificar al planificador de largo plazo
		args := []string{
			fileName,
			strconv.Itoa(tamanio),
			strconv.Itoa(pid), // ahora le paso el PID generado
		}
		Utils.ChannelProcessArguments <- args
		logger.Debug("Notificado planificador de largo plazo para PID <%d>", pid)
	}()
	<-Utils.SemProcessCreateOK

	// 4) Responder de inmediato
	w.WriteHeader(http.StatusOK)
}

func Exit(w http.ResponseWriter, r *http.Request) {
	var msg MensajeSyscall
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}
	pcb := msg.PCB
	logger.Info("## (<%d>) - Solicitó syscall: <EXIT>", pcb.PID)

	// Despachamos la señal en segundo plano para no bloquear el handler HTTP
	go func(p int) {
		Utils.ChannelFinishprocess <- p
	}(pcb.PID)
	// TODO desalojar
	// Respondemos de inmediato
	w.WriteHeader(http.StatusOK)
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
	pcb := mensajeRecibido.PCB
	nombre := mensajeRecibido.Nombre

	logger.Info("Syscall recibida: “## (<%d>) - Solicitó syscall: <IO>”", pcb.PID)

	// Aquí bloqueas el mutex mientras esperas a que el IO se registre
	globals.IOMu.Lock()
	ioData, ok := globals.IOs[nombre]
	for !ok {
		globals.IOCond.Wait()
		ioData, ok = globals.IOs[nombre] // reintenta obtenerlo
	}
	globals.IOMu.Unlock()

	logger.Info("Nombre IO: %s Duracion: %d", ioData.Nombre, mensajeRecibido.Duracion)

	//SIGNAL A Planif. CORTO PLAZO QUE LLEGO I/O
	go func(p int) {
		Utils.NotificarComienzoIO <- Utils.MensajeIOChannel{
			PID:      pcb.PID,
			Nombre:   ioData.Nombre,
			Duracion: mensajeRecibido.Duracion,
		}
	}(pcb.PID)

	w.WriteHeader(http.StatusOK)
}
