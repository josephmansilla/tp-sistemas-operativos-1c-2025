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
	"strings"
)

type Instruccion func(context *globals.ExecutionContext, arguments []string) error

type MensajeContexto struct {
	PID int    `json:"pid"`
	PC  uint32 `json:"pc"`
}

type MensajeIO struct {
	PID    int    `json:"pid"`
	PC     uint32 `json:"pc"`
	Tiempo int    `json:"tiempo"`
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
	return hacerSyscall()
}

func ioInstruccion(context *globals.ExecutionContext, arguments []string) error {
	// Suponiendo que el tiempo de la IO se pasa como primer argumento
	tiempoIO, err := strconv.Atoi(arguments[0]) // Parsear el tiempo desde el argumento
	if err != nil {
		log.Printf("Error al convertir el tiempo de IO: %s", err)
		return err
	}

	// Crear una instancia de MensajeIO con el PID, PC y el tiempo
	mensaje := MensajeIO{
		PID:    context.PID,
		PC:     context.PC,
		Tiempo: tiempoIO,
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
	return hacerSyscall()
}

func iniciarProcesoInstruccion(context *globals.ExecutionContext, arguments []string) error {
	return hacerSyscall()
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
		context.PC = uint32(jump)
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

func hacerSyscall() error {
	// Preparo el mensaje con el PID y el PC actuales
	mensaje := MensajeContexto{
		PID: globals.CurrentContext.PID,
		PC:  globals.CurrentContext.PC,
	}

	// Lo codifico a JSON
	jsonData, err := json.Marshal(mensaje)
	if err != nil {
		log.Printf("Error al codificar mensaje de syscall: %s", err)
		return err
	}

	// Armo la URL de syscall hacia Kernel
	url := fmt.Sprintf("http://%s:%d/kernel/syscall", globals.ClientConfig.IpKernel, globals.ClientConfig.PortKernel)

	// Envío el POST al Kernel
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error al hacer syscall a Kernel: %s", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Syscall al Kernel respondió status: %d", resp.StatusCode)
		return fmt.Errorf("syscall fallida con status %d", resp.StatusCode)
	}

	log.Println("Syscall realizada exitosamente")
	return nil
}

func checkArguments(args []string, correctNumberOfArgs int) error {
	if len(args) != correctNumberOfArgs {
		return errors.New("se recibió una cantidad de argumentos no válida")
	}
	return nil
}
