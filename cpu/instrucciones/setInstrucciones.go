package instrucciones

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"log"
	"net/http"
	"strconv"
)

type Instruccion func(context *globals.ExecutionContext, arguments []string) error

type MensajeContexto struct {
	PID int `json:"pid"`
	PC  int `json:"pc"`
}

type MensajeIO struct {
	PID    int    `json:"pid"`
	PC     int    `json:"pc"`
	Tiempo int    `json:"tiempo"`
	Nombre string `json:"nombre"`
}

type MensajeInitProc struct {
	PID      int    `json:"pid"`
	PC       int    `json:"pc"`
	Filename string `json:"filename"` //filename
	Tamanio  int    `json:"tamanio_memoria"`
}

// Una instruccion es una funcion que recibe un puntero a una struct con el contexto de ejecucion del proceso que se esta
// ejecutando (Pc, variables, registros, etc) y una lista de strings que son los argumentos

var InstruccionSet = map[string]Instruccion{
	// Instrucciones básicas
	"SET":   setInstruccion,
	"READ":  readMemInstruccion,
	"WRITE": writeMemInstruccion,
	"GOTO":  gotoInstruccion,
	"NOOP":  noopInstruccion,
	"SUM":   sumInstruccion,
	"SUB":   subInstruccion,
	"JNZ":   jnzInstruccion,
	"LOG":   logInstruccion,

	// Syscalls
	"DUMP_MEMORY": dumpMemoryInstruccion,
	"IO":          ioInstruccion,
	"EXIT":        exitInstruccion,
	"INIT_PROC":   iniciarProcesoInstruccion,
}

func dumpMemoryInstruccion(context *globals.ExecutionContext, arguments []string) error {
	// Validar que no hayan args
	if err := checkArguments(arguments, 0); err != nil {
		log.Printf("Error en los argumentos de la Instruccion: %s", err)
		return err
	}

	mensaje := MensajeContexto{
		PID: context.PID,
		PC:  context.PC,
	}

	// Lo codifico a JSON
	jsonData, err := json.Marshal(mensaje)
	if err != nil {
		log.Printf("Error al codificar mensaje de syscall de DUMP_MEMORY: %s", err)
		return err
	}

	// Armo la URL de syscall hacia Kernel
	url := fmt.Sprintf("http://%s:%d/kernel/dump_memory", globals.ClientConfig.IpKernel, globals.ClientConfig.PortKernel)

	// Envío el POST al Kernel con el contexto y el tiempo de IO
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error al hacer syscall de DUMP_MEMORY a Kernel: %s", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Syscall DUMP_MEMORY al Kernel respondió status: %d", resp.StatusCode)
		return fmt.Errorf("syscall DUMP_MEMORY fallida con status %d", resp.StatusCode)
	}

	log.Println("Syscall DUMP_MEMORY realizada exitosamente")
	return nil
}

func ioInstruccion(context *globals.ExecutionContext, arguments []string) error {
	// Validar que haya exactamente 1 argumento (el tiempo de IO)
	if err := checkArguments(arguments, 2); err != nil {
		log.Printf("Error en los argumentos de la Instruccion: %s", err)
		return err
	}

	// Parsear el tiempo desde el argumento
	nombreIO := arguments[0]
	tiempoIO, err := strconv.Atoi(arguments[1])
	if err != nil {
		log.Printf("Error al convertir el tiempo de IO: %s", err)
		return err
	}

	// Crear una instancia de MensajeIO con el PID, PC y el tiempo
	mensaje := MensajeIO{
		PID:    context.PID,
		PC:     context.PC,
		Tiempo: tiempoIO,
		Nombre: nombreIO,
	}

	// Lo codifico a JSON
	jsonData, err := json.Marshal(mensaje)
	if err != nil {
		log.Printf("Error al codificar mensaje de syscall de IO: %s", err)
		return err
	}

	// Armo la URL de syscall hacia Kernel
	url := fmt.Sprintf("http://%s:%d/kernel/io", globals.ClientConfig.IpKernel, globals.ClientConfig.PortKernel)

	// Envío el POST al Kernel con el contexto y el tiempo de IO
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error al hacer syscall de IO a Kernel: %s", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Syscall IO al Kernel respondió status: %d", resp.StatusCode)
		return fmt.Errorf("syscall IO fallida con status %d", resp.StatusCode)
	}

	log.Println("Syscall IO realizada exitosamente")
	return nil
}

func exitInstruccion(context *globals.ExecutionContext, arguments []string) error {
	// Validar que no hayan args
	if err := checkArguments(arguments, 0); err != nil {
		log.Printf("Error en los argumentos de la Instruccion: %s", err)
		return err
	}

	mensaje := MensajeContexto{
		PID: context.PID,
		PC:  context.PC,
	}

	// Lo codifico a JSON
	jsonData, err := json.Marshal(mensaje)
	if err != nil {
		log.Printf("Error al codificar mensaje de syscall de EXIT: %s", err)
		return err
	}

	// Armo la URL de syscall hacia Kernel
	url := fmt.Sprintf("http://%s:%d/kernel/exit", globals.ClientConfig.IpKernel, globals.ClientConfig.PortKernel)

	// Envío el POST al Kernel con el contexto y el tiempo de IO
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error al hacer syscall de EXIT a Kernel: %s", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Syscall EXIT al Kernel respondió status: %d", resp.StatusCode)
		return fmt.Errorf("syscall EXIT fallida con status %d", resp.StatusCode)
	}

	log.Println("Syscall EXIT realizada exitosamente")
	return nil
}

func iniciarProcesoInstruccion(context *globals.ExecutionContext, arguments []string) error {
	//Valido tener 2 argumentos (nombreArchivo y tamProceso)
	if err := checkArguments(arguments, 2); err != nil {
		log.Printf("Error en los argumentos de la Instruccion: %s", err)
		return err
	}

	filename := arguments[0]
	// Parsear el tamaño
	tamanio, err := strconv.Atoi(arguments[1])
	if err != nil {
		log.Printf("Error al convertir el tamaño del proceso: %s", err)
		return err
	}

	//Instancio el mensaje
	mensaje := MensajeInitProc{
		PID:      context.PID,
		PC:       context.PC,
		Filename: filename,
		Tamanio:  tamanio,
	}

	// Lo codifico a JSON
	jsonData, err := json.Marshal(mensaje)
	if err != nil {
		log.Printf("Error al codificar mensaje de syscall de IO: %s", err)
		return err
	}

	// Armo la URL de syscall hacia Kernel
	url := fmt.Sprintf("http://%s:%d/kernel/init_proceso", globals.ClientConfig.IpKernel, globals.ClientConfig.PortKernel)

	// Envío el POST al Kernel con el contexto y el tiempo de IO
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error al hacer syscall de INIT_PROC a Kernel: %s", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Syscall INIT_PROC al Kernel respondió status: %d", resp.StatusCode)
		return fmt.Errorf("Syscall INIT_PROC fallida con status %d", resp.StatusCode)
	}

	log.Println("Syscall IO realizada exitosamente")
	return nil
}

func gotoInstruccion(context *globals.ExecutionContext, arguments []string) error {}

func noopInstruccion(context *globals.ExecutionContext, arguments []string) error {}

func writeMemInstruccion(context *globals.ExecutionContext, arguments []string) error {
}

func readMemInstruccion(context *globals.ExecutionContext, arguments []string) error {

}

func jnzInstruccion(context *globals.ExecutionContext, arguments []string) error {
	if err := checkArguments(arguments, 2); err != nil {
		return err
	}

	register, err := context.GetRegister(arguments[0])
	if err != nil {
		return err
	}

	jump, err := strconv.Atoi(arguments[1])
	if err != nil {
		return err
	}

	if *register != 0 {
		context.PC = int(jump)
		log.Printf("Actualizando PC: %v", context.PC)
	}

	return nil
}

func sumInstruccion(context *globals.ExecutionContext, args []string) error {
	if err := checkArguments(args, 2); err != nil {
		return err
	}

	firstRegister, err := context.GetRegister(args[0])
	if err != nil {
		return err
	}

	secondRegister, err := context.GetRegister(args[1])
	if err != nil {
		return err
	}

	*firstRegister = *firstRegister + *secondRegister
	return nil
}

func subInstruccion(context *globals.ExecutionContext, args []string) error {
	if err := checkArguments(args, 2); err != nil {
		return err
	}

	firstRegister, err := context.GetRegister(args[0])
	if err != nil {
		return err
	}

	secondRegister, err := context.GetRegister(args[1])
	if err != nil {
		return err
	}

	*firstRegister = *firstRegister - *secondRegister
	return nil
}

func setInstruccion(ctx *globals.ExecutionContext, args []string) error {
	if err := checkArguments(args, 2); err != nil {
		return err
	}

	reg, err := ctx.GetRegister(args[0])
	if err != nil {
		return err
	}

	i, err := strconv.Atoi(args[1])
	if err != nil {
		reg2, err := ctx.GetRegister(args[1])
		if err != nil {
			return errors.New("no se pudo parsear '" + args[1] + "' como un entero o un registro")
		}
		*reg = *reg2
	} else {
		*reg = uint32(i)
	}

	return nil
}

func logInstruccion(ctx *globals.ExecutionContext, args []string) error {
	if err := checkArguments(args, 1); err != nil {
		return err
	}

	register, err := ctx.GetRegister(args[0])
	if err != nil {
		return err
	}

	log.Printf("Logging register '%v': %v", args[0], *register)
	fmt.Println(*register)
	return nil
}

func checkArguments(args []string, correctNumberOfArgs int) error {
	if len(args) != correctNumberOfArgs {
		return errors.New("se recibió una cantidad de argumentos no válida")
	}
	return nil
}
