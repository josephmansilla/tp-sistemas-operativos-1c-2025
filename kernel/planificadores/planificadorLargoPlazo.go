package planificadores

import (
	"github.com/sisoputnfrba/tp-golang/kernel/Utils"
	"github.com/sisoputnfrba/tp-golang/kernel/comunicacion"
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/pcb"
	logger "github.com/sisoputnfrba/tp-golang/utils/logger"
	"strconv"
)

//PLANIFICADOR LARGO PLAZO:
//encargado de gestionar las peticiones a la memoria
//para la creación y eliminación de procesos e hilos.

// CREA PRIMER PROCESO (NEW)
func CrearPrimerProceso(fileName string, tamanio int) {
	// Paso 1: Crear el PCB
	pid := globals.GenerarNuevoPID()
	pcbNuevo := pcb.PCB{
		PID: pid,
		PC:  0,
		ME:  make(map[string]int),
		MT:  make(map[string]int),
	}

	logger.Info("## (<%d>) Se crea el Primer proceso - Estado: NEW", pid)

	// Paso 2: Agregar a la cola  READY
	globals.ColaReady.Add(&pcbNuevo)

}

// ///////////////////
// //////// ESTO ES MIO
func PlanificadorLargoPlazo() {
	logger.Info("Iniciando el planificador de largo plazo")
	go ManejadorCreacionProcesos()
	go ManejadorFinalizacionProcesos()
}

func ManejadorCreacionProcesos() {
	logger.Info("Esperando solicitudes de INIT_PROC para creación de procesos")
	for {
		// Recibir args [filename, size, pid]
		args := <-Utils.ChannelProcessArguments
		fileName := args[0]
		size, _ := strconv.Atoi(args[1])
		pid, _ := strconv.Atoi(args[2])
		logger.Info("Solicitud INIT_PROC recibida: filename=%s, size=%d, pid=%d", fileName, size, pid)

		// Si existiera comunicación con Memoria, se haría aquí
		// exito, err := utils.SolicitarCreacionEnMemoria(fileName, size)
		// if err != nil || !exito { ... reintento ... }

		// Mover PCB de NEW a READY
		agregarProcesoAReady(pid)
	}
}
func reintentarCreacion(pid int, fileName string, processSize int) {
	for {
		logger.Info("<PID: %v> Esperando liberacion de memoria", pid)
		<-Utils.InitProcess
		exito, err := comunicacion.SolicitarCreacionEnMemoria(fileName, processSize)
		if err == nil && exito {
			agregarProcesoAReady(pid)
			Utils.MutexPuedoCrearProceso.Unlock()
			break
		}
	}
}

func agregarProcesoAReady(pid int) {
	// 1) Buscar el PCB en NEW con mutex
	Utils.MutexNuevo.Lock()
	pcbPtr := BuscarPCBPorPID(pid)
	if pcbPtr == nil {
		Utils.MutexNuevo.Unlock()
		logger.Error("agregarProcesoAReady: PCB pid=%d no existe en NEW", pid)
		return
	}

	// 2) Agregar a READY
	Utils.MutexReady.Lock()
	globals.ColaReady.Add(pcbPtr)
	Utils.MutexReady.Unlock()

	// 3) Remover de NEW
	globals.ColaNuevo.Remove(pcbPtr)
	Utils.MutexNuevo.Unlock()

	logger.Info("PCB pid=%d movido de NEW a READY", pid)

	// 4) Señal al planificador para continuar
	Utils.SemProcessCreateOK <- struct{}{}
}

func ManejadorFinalizacionProcesos() {
	for {
		pid := <-Utils.ChannelFinishprocess

		// Paso 1: Avisar a memoria que finalizó
		request := comunicacion.ConsultaAMemoria{
			Hilo:      comunicacion.Hilo{PID: comunicacion.Pid(pid)},
			Tipo:      comunicacion.FinishProcess,
			Arguments: map[string]interface{}{},
		}

		err := comunicacion.SendMemoryRequest(request)
		if err != nil {
			logger.Error("Error al finalizar proceso PID <%d>: %v", pid, err)
			continue
		}

		logger.Debug("Proceso PID <%d> finalizado. Liberando espacio.", pid)

		// Paso 2: Sacar de Running
		Utils.MutexEjecutando.Lock()
		for _, pcb := range globals.ColaEjecutando.Values() {
			if pcb.PID == pid {
				globals.ColaEjecutando.Remove(pcb)
				break
			}
		}
		Utils.MutexEjecutando.Unlock()

		// Paso 3: Intentar inicializar proceso de SUSP.READY
		if intentarInicializarDesdeSuspReady() {
			continue
		}

		// Paso 4: Si no hay en SUSP.READY, intentar desde NEW
		intentarInicializarDesdeNew()

		// Avisar que el proceso fue terminado correctamente
		Utils.ChannelFinishProcess2 <- true
	}
}

// ESTO ES PARA MANEJAR LAS SYSCALLS, ES DISTINTO AL LA CRACION DEL PRIMER PROCESO, PORQUE A LA SYSCALL pueden
// llamarla varias CPU entonces muchas cpu pueden estar tocando la listas de NEW, y no la añade de una a ready como el primer proceso
// YA QUE ESTO DEPENDE DEL PLANIFICADOR DE LARGO PLAZO
func CrearProceso(fileName string, tamanio int) {
	// Paso 1: Crear el PCB
	pid := globals.GenerarNuevoPID()
	pcbNuevo := pcb.PCB{
		PID:         pid,
		PC:          0,
		ME:          make(map[string]int),
		MT:          make(map[string]int),
		FileName:    fileName,
		ProcessSize: tamanio,
	}

	logger.Info("## (<%d>) Se crea el proceso - Estado: NEW", pid)

	// Paso 2: Agregar a la cola NEW con protección de mutex
	Utils.MutexNuevo.Lock()
	globals.ColaNuevo.Add(&pcbNuevo)
	logger.Info("SE AÑADIO A LA COLA nuevo DESDE SYSCALL")
	Utils.MutexNuevo.Unlock()

	args := []string{pcbNuevo.FileName, strconv.Itoa(pcbNuevo.ProcessSize), strconv.Itoa(pcbNuevo.PID)}
	Utils.ChannelProcessArguments <- args

}
func intentarInicializarDesdeSuspReady() bool {
	Utils.MutexSuspendidoReady.Lock()
	defer Utils.MutexSuspendidoReady.Unlock()
	// la consigna PIDE QUE PRIMERO INTENTE PASAR A READY UN PROCESO DE LA LISTA SUSPENDIDO READY LUEGO DE READY

	if globals.ColaSuspendidoReady.IsEmpty() {
		return false
	}

	pcb := globals.ColaSuspendidoReady.First()
	exito, err := comunicacion.SolicitarCreacionEnMemoria(pcb.FileName, pcb.ProcessSize)
	if err != nil || !exito {
		logger.Info("No se pudo inicializar proceso desde SUSP.READY PID <%d>", pcb.PID)
		return false
	}
	Utils.MutexSuspendidoReady.Lock()
	globals.ColaSuspendidoReady.Remove(pcb)
	Utils.MutexSuspendidoReady.Unlock()

	Utils.MutexReady.Lock()
	globals.ColaReady.Add(pcb)
	Utils.MutexReady.Unlock()

	logger.Info("PID <%d> pasó de SUSP.READY a READY", pcb.PID)
	return true
}
func intentarInicializarDesdeNew() {
	Utils.MutexNuevo.Lock()
	defer Utils.MutexNuevo.Unlock()

	if globals.ColaNuevo.IsEmpty() {
		return
	}

	pcb := globals.ColaNuevo.First()
	exito, err := comunicacion.SolicitarCreacionEnMemoria(pcb.FileName, pcb.ProcessSize)
	if err != nil || !exito {
		logger.Info("No se pudo inicializar proceso desde NEW PID <%d>", pcb.PID)
		return
	}
	Utils.MutexNuevo.Lock()
	globals.ColaNuevo.Remove(pcb)
	Utils.MutexNuevo.Unlock()

	Utils.MutexReady.Lock()
	globals.ColaReady.Add(pcb)
	Utils.MutexReady.Unlock()

	logger.Info("PID <%d> pasó de NEW a READY", pcb.PID)
}
func BuscarPCBPorPID(pid int) *pcb.PCB {
	for _, p := range globals.ColaNuevo.Values() {
		if p.PID == pid {
			return p
		}
	}
	return nil
}
func MostrarPCBsEnNew() {
	Utils.MutexNuevo.Lock()
	defer Utils.MutexNuevo.Unlock()

	if globals.ColaReady.IsEmpty() {
		logger.Info("La cola READYYY está vacía.")
		return
	}

	for _, pcb := range globals.ColaReady.Values() {
		logger.Info("- esto son los PCB en NEW CON PID: %d", pcb.PID)
	}
}
