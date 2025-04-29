package instrucciones

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"log"
	"net/http"
	"strings"
)

type MensajeInstruccion struct {
	PID int `json:"pid"`
	PC  int `json:"pc"`
}

type RespuestaInstruccion struct {
	Instruccion string `json:"instruccion"`
}

type Interrupcion struct {
	PID int `json:"pid"`
	PC  int `json:"pc"`
}

func FaseFetch(ipDestino string, puertoDestino int, pidPropio int, pcInicial int) {
	pc := pcInicial

	for {
		mensaje := MensajeInstruccion{
			PID: pidPropio,
			PC:  pc,
		}

		jsonData, err := json.Marshal(mensaje)
		if err != nil {
			log.Printf("Error codificando mensaje a JSON: %s", err)
			break
		}

		url := fmt.Sprintf("http://%s:%d/memoria/instruccion", ipDestino, puertoDestino)

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("Error haciendo POST a Memoria: %s", err)
			break
		}
		resp.Body.Close()

		var respuesta RespuestaInstruccion
		err = json.NewDecoder(resp.Body).Decode(&respuesta)
		if err != nil {
			log.Printf("Error decodificando respuesta de Memoria: %s", err)
			break
		}

		if respuesta.Instruccion == "" {
			log.Printf("No hay instruccion para PID %d (PC %d)", pidPropio, pc)
			break
		}

		log.Printf("Instrucción recibida (PC %d): %s", pc, respuesta.Instruccion)

		// Parsear y ejecutar instrucción
		if seguir := FaseDecode(respuesta.Instruccion); !seguir {
			log.Println("Se pidió un syscall, finalizando ejecución del proceso.")
			break
		}

		pc++
		globals.CurrentContext.PC = pc
	}
}

func FaseDecode(instruccion string) bool {
	partes := strings.Fields(instruccion)
	if len(partes) == 0 {
		log.Println("Instrucción vacía")
		return true
	}

	nombre := partes[0]
	args := partes[1:]

	return FaseExecute(nombre, args)
}

func FaseExecute(nombre string, args []string) bool {
	instrucFunc, existe := InstruccionSet[nombre]
	if !existe {
		log.Printf("Instrucción desconocida: %s", nombre)
		return true
	}

	err := instrucFunc(globals.CurrentContext, args)
	log.Printf("Ejecutando instrucción: %s", nombre)
	if err != nil {
		log.Printf("Error ejecutando %s: %v", nombre, err)
		return false
	}

	if FaseCheckInterrupt() {
		log.Println("Finalizando ejecución por interrupción.")
		return false
	}

	return true
}

func FaseCheckInterrupt() bool {
	globals.MutexInterrupcion.Lock()
	defer globals.MutexInterrupcion.Unlock()

	if !globals.InterrupcionPendiente {
		return false
	}

	if globals.CurrentContext == nil {
		log.Println("Interrupción recibida pero no hay contexto en ejecución.")
		globals.InterrupcionPendiente = false
		return false
	}

	if globals.PIDInterrumpido != globals.CurrentContext.PID {
		log.Printf("Interrupción recibida para PID %d, pero estoy ejecutando PID %d. Ignorando.",
			globals.PIDInterrumpido, globals.CurrentContext.PID)
		return false
	}

	pid := globals.CurrentContext.PID
	pc := globals.CurrentContext.PC

	// Preparar JSON
	body := Interrupcion{
		PID: pid,
		PC:  pc,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Printf("Error serializando contexto interrumpido: %v", err)
		return false
	}

	// Enviar al Kernel
	url := fmt.Sprintf("http://%s:%d/kernel/contexto_interrumpido", globals.ClientConfig.IpKernel, globals.ClientConfig.PortKernel)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("Error enviando contexto interrumpido al Kernel: %v", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Kernel respondió con error al recibir interrupción: %s", resp.Status)
		return false
	}

	log.Printf("Contexto interrumpido enviado a Kernel. PID: %d, PC: %d", pid, pc)

	// Limpiar la interrupción
	globals.InterrupcionPendiente = false

	return true
}
