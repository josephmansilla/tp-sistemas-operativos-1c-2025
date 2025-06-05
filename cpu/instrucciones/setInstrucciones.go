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

type Instruccion func(arguments []string) error

type MensajeDump struct {
	PID int    `json:"pid"`
	PC  int    `json:"pc"`
	ID  string `json:"id"`
}

type MensajeIO struct {
	PID    int    `json:"pid"`
	PC     int    `json:"pc"`
	Tiempo int    `json:"tiempo"`
	Nombre string `json:"nombre"`
	ID     string `json:"id"`
}

type MensajeInitProc struct {
	PID      int    `json:"pid"`
	PC       int    `json:"pc"`
	Filename string `json:"filename"` //filename
	Tamanio  int    `json:"tamanio_memoria"`
	ID       string `json:"id"`
}

type MensajeExit struct {
	PID int    `json:"pid"`
	PC  int    `json:"pc"`
	ID  string `json:"id"`
}

// Una instruccion es una funcion que recibe un puntero a una struct con el contexto de ejecucion del proceso que se esta
// ejecutando (Pc, variables, registros, etc) y una lista de strings que son los argumentos

var InstruccionSet = map[string]Instruccion{
	// Instrucciones de CPU
	"NOOP":  noopInstruccion,
	"GOTO":  gotoInstruccion,
	"WRITE": writeMemInstruccion,
	"READ":  readMemInstruccion,
	// Syscalls
	"DUMP_MEMORY": dumpMemoryInstruccion,
	"IO":          ioInstruccion,
	"EXIT":        exitInstruccion,
	"INIT_PROC":   iniciarProcesoInstruccion,
}

// INSTRUCCIONES DE CPU
func noopInstruccion(arguments []string) error {
	if err := checkArguments(arguments, 0); err != nil {
		log.Printf("Error en los argumentos de la instrucción: %s", err)
		return err
	}
	log.Printf("## PID: %d - Instrucción NOOP ejecutada. No se realizó ninguna acción.", globals.PIDActual)
	return nil
}

func gotoInstruccion(arguments []string) error {
	if err := checkArguments(arguments, 1); err != nil {
		log.Printf("Error en los argumentos de la instrucción: %s", err)
		return err
	}

	nuevoPC, err := strconv.Atoi(arguments[0])
	if err != nil {
		log.Printf("Error al convertir el valor de PC en la instrucción GOTO: %s", err)
		return err
	}

	globals.PCActual = nuevoPC
	globals.SaltarIncrementoPC = true

	log.Printf("## PID: %d - Instrucción GOTO ejecutada. Nuevo PC: %d", globals.PIDActual, globals.PCActual)
	return nil
}

func writeMemInstruccion(arguments []string) error {
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
	dirFisica := traducciones.Traducir(dirLogica)

	if err := traducciones.Escribir(dirFisica, datos); err != nil {
		log.Printf("Error escribiendo en Memoria: %s", err)
		return err
	}

	log.Printf("PID: %s - Acción: LEER - Dirección Física: %s - Valor: %s", globals.PIDActual, dirFisica, datos)

	return nil
}

func readMemInstruccion(arguments []string) error {
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
	dirFisica := traducciones.Traducir(dirLogica)

	valorLeido, err := traducciones.Leer(dirFisica, tamanio)
	if err != nil {
		log.Printf("Error leyendo en Memoria: %s", err)
		return err
	}

	log.Printf("PID: %s - Acción: LEER - Dirección Física: %s - Valor: %s", globals.PIDActual, dirFisica, valorLeido)

	return nil
}

// SYSCALLS
func dumpMemoryInstruccion(arguments []string) error {
	if err := checkArguments(arguments, 0); err != nil {
		log.Printf("Error en los argumentos de la Instruccion: %s", err)
		return err
	}

	mensaje := MensajeDump{
		PID: globals.PIDActual,
		PC:  globals.PCActual,
	}

	url := fmt.Sprintf("http://%s:%d/kernel/dump_memory", globals.ClientConfig.IpKernel, globals.ClientConfig.PortKernel)

	if err := data.EnviarDatos(url, mensaje); err != nil {
		log.Printf("Error al hacer syscall de DUMP_MEMORY a Kernel: %s", err)
		return err
	}

	log.Println("Syscall DUMP_MEMORY realizada exitosamente")
	return nil
}

func ioInstruccion(arguments []string) error {
	if err := checkArguments(arguments, 2); err != nil {
		log.Printf("Error en los argumentos de la Instruccion: %s", err)
		return err
	}

	nombreIO := arguments[0]
	tiempoIO, err := strconv.Atoi(arguments[1])
	globals.PCActual++
	if err != nil {
		log.Printf("Error al convertir el tiempo de IO: %s", err)
		return err
	}

	mensaje := MensajeIO{
		PID:    globals.PIDActual,
		PC:     globals.PCActual,
		Tiempo: tiempoIO,
		Nombre: nombreIO,
		ID:     globals.ID,
	}

	url := fmt.Sprintf("http://%s:%d/kernel/syscallIO", globals.ClientConfig.IpKernel, globals.ClientConfig.PortKernel)

	if err := data.EnviarDatos(url, mensaje); err != nil {
		log.Printf("Error al hacer syscall IO a Kernel: %s", err)
		return err
	}

	log.Println("Syscall IO realizada exitosamente")
	return globals.ErrSyscallBloqueante
}

func exitInstruccion(arguments []string) error {
	if err := checkArguments(arguments, 0); err != nil {
		log.Printf("Error en los argumentos de la Instruccion: %s", err)
		return err
	}

	mensaje := MensajeExit{
		PID: globals.PIDActual,
		PC:  globals.PCActual,
		ID:  globals.ID,
	}

	url := fmt.Sprintf("http://%s:%d/kernel/exit", globals.ClientConfig.IpKernel, globals.ClientConfig.PortKernel)

	if err := data.EnviarDatos(url, mensaje); err != nil {
		log.Printf("Error al hacer syscall EXIT a Kernel: %s", err)
		return err
	}

	log.Println("Syscall EXIT realizada exitosamente")
	return nil
}

func iniciarProcesoInstruccion(arguments []string) error {
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
		PID:      globals.PIDActual,
		PC:       globals.PCActual,
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

func checkArguments(args []string, correctNumberOfArgs int) error {
	if len(args) != correctNumberOfArgs {
		return errors.New("se recibió una cantidad de argumentos no válida")
	}
	return nil
}
