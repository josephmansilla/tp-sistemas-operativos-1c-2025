package instrucciones

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/traducciones"
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

func FaseFetch(ipDestino string, puertoDestino int) {
	for {
		log.Printf("## PID: %d - FETCH - Program Counter: %d", globals.PIDActual, globals.PCActual)

		mensaje := MensajeInstruccion{
			PID: globals.PIDActual,
			PC:  globals.PCActual,
		}

		jsonData, err := json.Marshal(mensaje)
		if err != nil {
			log.Printf("Error codificando mensaje a JSON: %s", err)
			break
		}

		url := fmt.Sprintf("http://%s:%d/memoria/obtenerInstruccion", ipDestino, puertoDestino)

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("Error haciendo POST a Memoria: %s", err)
			break
		}
		defer resp.Body.Close() // <-- mover esto después de confirmar que no hubo error

		var respuesta RespuestaInstruccion
		if err := json.NewDecoder(resp.Body).Decode(&respuesta); err != nil {
			log.Printf("Error decodificando respuesta de Memoria: %s", err)
			break
		}

		log.Printf("Instrucción recibida (PC %d): %s", globals.PCActual, respuesta.Instruccion)

		// Parsear y ejecutar instrucción
		if seguir := FaseDecode(respuesta.Instruccion); !seguir {
			log.Println("Se pidió un syscall, finalizando ejecución del proceso.")
			break
		}

		if !globals.SaltarIncrementoPC {
			globals.PCActual++
		} else {
			globals.SaltarIncrementoPC = false // reset para la próxima instrucción
		}
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

	err := instrucFunc(args)

	if err != nil {
		if err == globals.ErrSyscallBloqueante {
			log.Printf("Proceso %d bloqueado por syscall IO", globals.PIDActual)
			return false // Detener ejecución por syscall IO
		}

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

	if globals.PIDInterrumpido != globals.PIDActual {
		log.Printf("Interrupción recibida para PID %d, pero estoy ejecutando PID %d. Ignorando.",
			globals.PIDInterrumpido, globals.PIDActual)
		return false
	}

	traducciones.Cache.LimpiarCache()

	pid := globals.PIDActual
	pc := globals.PCActual

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
