package instrucciones

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"net/http"
	"strconv"
	"strings"
)

type Instruction func(context *globals.ExecutionContext, arguments []string) error

// Una instruccion es una funcion que recibe un puntero a una struct con el contexto de ejecucion del proceso que se esta
// ejecutando (Pc, variables, registros, etc) y una lista de strings que son los argumentos

var InstructionSet = map[string]Instruction{
	"SET":       setInstruction,
	"READ_MEM":  readMemInstruction,
	"WRITE_MEM": writeMemInstruction,
	"SUM":       sumInstruction,
	"SUB":       subInstruction,
	"JNZ":       jnzInstruction,
	"LOG":       logInstruction,
	// De acá en adelante son syscalls y las soluciona querido kernel
	"DUMP_MEMORY":    dumpMemoryInstruction,
	"IO":             ioInstruction,
	"PROCESS_CREATE": processCreateInstruction,
	"THREAD_CREATE":  threadCreateInstruction,
	"THREAD_JOIN":    threadJoinInstruction,
	"THREAD_CANCEL":  threadCancelInstruction,
	"MUTEX_CREATE":   mutexCreateInstruction,
	"MUTEX_LOCK":     mutexLockInstruction,
	"MUTEX_UNLOCK":   mutexUnlockInstruction,
	"THREAD_EXIT":    threadExitInstruction,
	"PROCESS_EXIT":   processExitInstruction,
}

func writeMemInstruction(context *globals.ExecutionContext, arguments []string) error {
	dataRegister, err := context.GetRegister(arguments[1])
	if err != nil {
		return err
	}

	virtualAddressRegister, err := context.GetRegister(arguments[0])
	if err != nil {
		return err
	}

	physicalAddress := context.MemoryBase + *virtualAddressRegister
	logger.Debug("Escribiendo en physicalAddres: %v", physicalAddress)
	if *virtualAddressRegister >= context.MemorySize {
		logger.Warn("Se trató de escribir una dirección no perteneciente al proceso! Interrumpiendo...")
		interruptionChannel <- globals.Interruption{
			Type:        globals.InterruptionSegFault,
			Description: "La dirección no forma parte del espacio del memoria del proceso"}

		logger.Debug("Intentadno liberar mutexInterruption")
		//MutexInterruption.Unlock()
		logger.Debug("Liberado mutexInterruption")
		return nil
	}

	err = memoryWrite(*currentThread, physicalAddress, *dataRegister)
	if err != nil {
		return err
	}

	logger.Info("## P%v T%v - Escribió '%v' en la dirección física <%v>",
		currentThread.PID, currentThread.TID, *dataRegister, physicalAddress)

	return nil
}

func readMemInstruction(context *globals.ExecutionContext, arguments []string) error {
	dataRegister, err := context.GetRegister(arguments[0])
	logger.Debug("DataRegister: %v", *dataRegister)
	if err != nil {
		return err
	}

	virtualAddressRegister, err := context.GetRegister(arguments[1])
	logger.Debug("VirtualAddressRegister: %v", *virtualAddressRegister)
	if err != nil {
		return err
	}
	logger.Debug("MemoryBase: %v", context.MemoryBase)
	physicalAddress := context.MemoryBase + *virtualAddressRegister
	logger.Debug("Physical Address: %v", physicalAddress)

	if *virtualAddressRegister >= context.MemorySize {
		logger.Warn("Se trató de leer una dirección no perteneciente al proceso! Interrumpiendo...")
		interruptionChannel <- globals.Interruption{
			Type:        globals.InterruptionSegFault,
			Description: "La dirección no forma parte del espacio del memoria del proceso"}

		logger.Debug("Intentadno liberar mutexInterruption")
		//MutexInterruption.Unlock()
		logger.Debug("Liberado mutexInterruption")
		return nil
	}

	if currentThread == nil {
		logger.Error("Se mando a ejecutar la instrucción readMemory pero no hay ningún hilo en ejecución ?")
		return nil
	}

	*dataRegister, err = memoryRead(*currentThread, physicalAddress)
	if err != nil {
		return err
	}
	logger.Info("## P%v T%v - Leyó '%v' de la dirección física <%v>",
		currentThread.PID, currentThread.TID, *dataRegister, physicalAddress)

	return nil
}

func jnzInstruction(context *globals.ExecutionContext, arguments []string) error {
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
		context.Pc = uint32(jump)
		logger.Trace("actualizando PC: %v", context.Pc)
	}

	return nil
}

func sumInstruction(context *globals.ExecutionContext, args []string) error {
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

func subInstruction(context *globals.ExecutionContext, args []string) error {
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

func setInstruction(ctx *globals.ExecutionContext, args []string) error {
	// Check number of arguments
	if err := checkArguments(args, 2); err != nil {
		return err
	}

	// Get the register to modify
	reg, err := ctx.GetRegister(args[0])
	if err != nil {
		return err
	}

	// Try parsing second argument as int
	i, err := strconv.Atoi(args[1])
	if err != nil {
		// Ok not int, but is it a register?
		if reg2, err := ctx.GetRegister(args[1]); err != nil {
			// It is not a register nor an int
			return errors.New("no se pudo parsear '" + args[1] + "' como un entero o un registro")
		} else {
			// If it IS a register, set it to that value
			*reg = *reg2
		}
	} else {
		// Set the register
		*reg = uint32(i)
	}

	return nil
}

func logInstruction(ctx *globals.ExecutionContext, args []string) error {
	if err := checkArguments(args, 1); err != nil {
		return err
	}

	registerString := args[0]

	register, err := ctx.GetRegister(registerString)
	if err != nil {
		return err
	}

	logger.Info("Logging register '%v': %v", registerString, *register)
	fmt.Println(*register)
	return nil
}

func checkArguments(args []string, correctNumberOfArgs int) error {
	if len(args) != correctNumberOfArgs {
		return errors.New("se recibió una cantidad de argumentos no válida")
	}
	return nil
}

// A partir de acá las syscalls
func doSyscall(ctx globals.ExecutionContext, syscall syscalls.Syscall) error {
	interruption := globals.Interruption{
		Type:        globals.InterruptionSyscall,
		Description: "Interrupción por syscall",
	}
	if len(interruptionChannel) > 0 {
		logger.Debug("Llego Interruption y Syscall => Hacemos primero syscall")
		// Si queremos hacer una syscall y el kernel ya mando desalojo o fin de quantum, atende primero la syscall
		// y agregamos a deuda la de desalojo
		desalojoInterruption := <-interruptionChannel
		interruptionChannel <- interruption

		interrupcionInsatisfecha := globals.InterrupcionInsatisfecha{
			Thread:       currentThread,
			Interruption: desalojoInterruption,
		}

		deudaInterrupciones = append(deudaInterrupciones, interrupcionInsatisfecha)
	} else {
		interruptionChannel <- interruption
	}
	url := fmt.Sprintf("http://%v:%v/kernel/syscall", config.KernelAddress, config.KernelPort)
	jsonData, err := json.Marshal(syscall)
	if err != nil {
		return fmt.Errorf("error al empaquetar syscall: %v", err)
	}

	MutexInterruption.Unlock()

	resp, err := http.Post(url, "application/json", strings.NewReader(string(jsonData)))
	if err != nil || resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error al enviar syscall al kernel: %v", err)
	}

	logger.Debug("Syscall enviada al kernel")

	return nil
}

func mutexCreateInstruction(context *globals.ExecutionContext, arguments []string) error {
	if err := doSyscall(
		*context,
		syscalls.New(syscalls.MutexCreate, arguments),
	); err != nil {
		return err
	}

	return nil
}

func processExitInstruction(context *globals.ExecutionContext, arguments []string) error {
	if err := doSyscall(
		*context,
		syscalls.New(syscalls.ProcessExit, arguments),
	); err != nil {
		return err
	}

	return nil
}

func threadExitInstruction(context *globals.ExecutionContext, arguments []string) error {
	if err := doSyscall(
		*context,
		syscalls.New(syscalls.ThreadExit, arguments),
	); err != nil {
		return err
	}

	return nil
}

func mutexLockInstruction(context *globals.ExecutionContext, arguments []string) error {
	if err := doSyscall(
		*context,
		syscalls.New(syscalls.MutexLock, arguments),
	); err != nil {
		return err
	}

	return nil
}

func mutexUnlockInstruction(context *globals.ExecutionContext, arguments []string) error {
	if err := doSyscall(
		*context,
		syscalls.New(syscalls.MutexUnlock, arguments),
	); err != nil {
		return err
	}

	return nil
}

func threadCancelInstruction(context *globals.ExecutionContext, arguments []string) error {
	if err := doSyscall(
		*context,
		syscalls.New(syscalls.ThreadCancel, arguments),
	); err != nil {
		return err
	}

	return nil
}

func threadCreateInstruction(context *globals.ExecutionContext, arguments []string) error {
	if err := doSyscall(
		*context,
		syscalls.New(syscalls.ThreadCreate, arguments),
	); err != nil {
		return err
	}

	return nil
}

func threadJoinInstruction(context *globals.ExecutionContext, arguments []string) error {
	if err := doSyscall(
		*context,
		syscalls.New(syscalls.ThreadJoin, arguments),
	); err != nil {
		return err
	}

	return nil
}

func processCreateInstruction(context *globals.ExecutionContext, arguments []string) error {
	if err := doSyscall(
		*context,
		syscalls.New(syscalls.ProcessCreate, arguments),
	); err != nil {
		return err
	}

	return nil
}

func ioInstruction(context *globals.ExecutionContext, arguments []string) error {
	if err := doSyscall(
		*context,
		syscalls.New(syscalls.IO, arguments),
	); err != nil {
		return err
	}

	return nil
}

func dumpMemoryInstruction(context *globals.ExecutionContext, arguments []string) error {
	if err := doSyscall(
		*context,
		syscalls.New(syscalls.DumpMemory, arguments),
	); err != nil {
		return err
	}

	return nil
}
