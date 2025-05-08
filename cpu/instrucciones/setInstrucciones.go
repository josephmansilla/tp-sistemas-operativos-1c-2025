package instrucciones

import (
	"errors"
	"fmt"
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/cpu/traducciones"
	"github.com/sisoputnfrba/tp-golang/utils/data"
	"log"
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
	// Instrucciones de CPU
	"NOOP":  noopInstruccion,
	"GOTO":  gotoInstruccion,
	"WRITE": writeMemInstruccion,
	"READ":  readMemInstruccion,
	//Instrucciones de Registros
	"JNZ": jnzInstruccion,
	"SUM": sumInstruccion,
	"SUB": subInstruccion,
	"SET": setInstruccion,
	"LOG": logInstruccion, //Esta no se si va
	// Syscalls
	"DUMP_MEMORY": dumpMemoryInstruccion,
	"IO":          ioInstruccion,
	"EXIT":        exitInstruccion,
	"INIT_PROC":   iniciarProcesoInstruccion,
}

// INSTRUCCIONES DE CPU
func noopInstruccion(context *globals.ExecutionContext, arguments []string) error {
	if err := checkArguments(arguments, 0); err != nil {
		log.Printf("Error en los argumentos de la instrucción: %s", err)
		return err
	}
	log.Printf("## PID: %d - Instrucción NOOP ejecutada. No se realizó ninguna acción.", context.PID)
	return nil
}

func gotoInstruccion(context *globals.ExecutionContext, arguments []string) error {
	if err := checkArguments(arguments, 1); err != nil {
		log.Printf("Error en los argumentos de la instrucción: %s", err)
		return err
	}

	nuevoPC, err := strconv.Atoi(arguments[0])
	if err != nil {
		log.Printf("Error al convertir el valor de PC en la instrucción GOTO: %s", err)
		return err
	}

	globals.MutexInterrupcion.Lock()
	context.PC = nuevoPC
	globals.MutexInterrupcion.Unlock()

	log.Printf("## PID: %d - Instrucción GOTO ejecutada. Nuevo PC: %d", context.PID, context.PC)
	return nil
}

func writeMemInstruccion(context *globals.ExecutionContext, arguments []string) error {
	if err := checkArguments(arguments, 2); err != nil {
		log.Printf("Error en los argumentos de la Instruccion: %s", err)
		return err
	}

	dirLogica, err := strconv.Atoi(arguments[0])
	if err != nil {
		log.Printf("Error al convertir la direccion logica: %s", err)
		return err
	}
	datos := arguments[1]
	dirFisica := traducciones.Traducir(context.PID, dirLogica)

	if err := traducciones.Escribir(dirFisica, datos); err != nil {
		log.Printf("Error escribiendo en Memoria: %s", err)
		return err
	}

	log.Printf("PID: %s - Acción: LEER - Dirección Física: %s - Valor: %s", context.PID, dirFisica, datos)

	return nil
}

func readMemInstruccion(context *globals.ExecutionContext, arguments []string) error {
	if err := checkArguments(arguments, 2); err != nil {
		log.Printf("Error en los argumentos de la Instruccion: %s", err)
		return err
	}

	dirLogica, err := strconv.Atoi(arguments[0])
	if err != nil {
		log.Printf("Error al convertir la direccion logica: %s", err)
		return err
	}
	tamanio, err := strconv.Atoi(arguments[1])
	if err != nil {
		log.Printf("Error al convertir el tamanio: %s", err)
		return err
	}
	dirFisica := traducciones.Traducir(context.PID, dirLogica)

	valorLeido, err := traducciones.Leer(dirFisica, tamanio)
	if err != nil {
		log.Printf("Error leyendo en Memoria: %s", err)
		return err
	}

	log.Printf("PID: %s - Acción: LEER - Dirección Física: %s - Valor: %s", context.PID, dirFisica, valorLeido)

	return nil
}

// INSTRUCCIONES DE REGISTROS
func jnzInstruccion(context *globals.ExecutionContext, arguments []string) error {
	if err := checkArguments(arguments, 2); err != nil {
		return err
	}

	register, err := context.ObtenerRegistro(arguments[0])
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

	firstRegister, err := context.ObtenerRegistro(args[0])
	if err != nil {
		return err
	}

	secondRegister, err := context.ObtenerRegistro(args[1])
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

	firstRegister, err := context.ObtenerRegistro(args[0])
	if err != nil {
		return err
	}

	secondRegister, err := context.ObtenerRegistro(args[1])
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

	reg, err := ctx.ObtenerRegistro(args[0])
	if err != nil {
		return err
	}

	i, err := strconv.Atoi(args[1])
	if err != nil {
		reg2, err := ctx.ObtenerRegistro(args[1])
		if err != nil {
			return errors.New("no se pudo parsear '" + args[1] + "' como un entero o un registro")
		}
		*reg = *reg2
	} else {
		*reg = uint32(i)
	}

	return nil
}

// SYSCALLS
func dumpMemoryInstruccion(context *globals.ExecutionContext, arguments []string) error {
	if err := checkArguments(arguments, 0); err != nil {
		log.Printf("Error en los argumentos de la Instruccion: %s", err)
		return err
	}

	mensaje := MensajeContexto{
		PID: context.PID,
		PC:  context.PC,
	}

	url := fmt.Sprintf("http://%s:%d/kernel/dump_memory", globals.ClientConfig.IpKernel, globals.ClientConfig.PortKernel)

	if err := data.EnviarDatos(url, mensaje); err != nil {
		log.Printf("Error al hacer syscall de DUMP_MEMORY a Kernel: %s", err)
		return err
	}

	log.Println("Syscall DUMP_MEMORY realizada exitosamente")
	return nil
}

func ioInstruccion(context *globals.ExecutionContext, arguments []string) error {
	if err := checkArguments(arguments, 2); err != nil {
		log.Printf("Error en los argumentos de la Instruccion: %s", err)
		return err
	}

	nombreIO := arguments[0]
	tiempoIO, err := strconv.Atoi(arguments[1])
	if err != nil {
		log.Printf("Error al convertir el tiempo de IO: %s", err)
		return err
	}

	mensaje := MensajeIO{
		PID:    context.PID,
		PC:     context.PC,
		Tiempo: tiempoIO,
		Nombre: nombreIO,
	}

	url := fmt.Sprintf("http://%s:%d/kernel/syscallIO", globals.ClientConfig.IpKernel, globals.ClientConfig.PortKernel)

	if err := data.EnviarDatos(url, mensaje); err != nil {
		log.Printf("Error al hacer syscall IO a Kernel: %s", err)
		return err
	}

	log.Println("Syscall IO realizada exitosamente")
	return globals.ErrSyscallBloqueante
}

func exitInstruccion(context *globals.ExecutionContext, arguments []string) error {
	if err := checkArguments(arguments, 0); err != nil {
		log.Printf("Error en los argumentos de la Instruccion: %s", err)
		return err
	}

	mensaje := MensajeContexto{
		PID: context.PID,
		PC:  context.PC,
	}

	url := fmt.Sprintf("http://%s:%d/kernel/exit", globals.ClientConfig.IpKernel, globals.ClientConfig.PortKernel)

	if err := data.EnviarDatos(url, mensaje); err != nil {
		log.Printf("Error al hacer syscall EXIT a Kernel: %s", err)
		return err
	}

	log.Println("Syscall EXIT realizada exitosamente")
	return nil
}

func iniciarProcesoInstruccion(context *globals.ExecutionContext, arguments []string) error {
	if err := checkArguments(arguments, 2); err != nil {
		log.Printf("Error en los argumentos de la Instruccion: %s", err)
		return err
	}

	filename := arguments[0]
	tamanio, err := strconv.Atoi(arguments[1])
	if err != nil {
		log.Printf("Error al convertir el tamaño del proceso: %s", err)
		return err
	}

	mensaje := MensajeInitProc{
		PID:      context.PID,
		PC:       context.PC,
		Filename: filename,
		Tamanio:  tamanio,
	}

	url := fmt.Sprintf("http://%s:%d/kernel/init_proceso", globals.ClientConfig.IpKernel, globals.ClientConfig.PortKernel)

	if err := data.EnviarDatos(url, mensaje); err != nil {
		log.Printf("Error al hacer syscall INIT_PROC a Kernel: %s", err)
		return err
	}

	log.Println("Syscall INIT_PROC realizada exitosamente")
	return nil
}

// Esta nose si la van a pedir
func logInstruccion(ctx *globals.ExecutionContext, args []string) error {
	if err := checkArguments(args, 1); err != nil {
		return err
	}

	register, err := ctx.ObtenerRegistro(args[0])
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
