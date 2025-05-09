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

	logger.Info("## (<%d>) Se crea el proceso - Estado: NEW", pid)

	// Paso 2: Agregar a la cola NEW
	globals.ColaNuevo.Add(&pcbNuevo)

	// Paso 3: Mover directamente a READY (sin consultar a memoria)
	globals.ColaNuevo.Remove(&pcbNuevo)
	globals.ColaReady.Add(&pcbNuevo)

	logger.Info("## (<%d>) Pasa del estado <NEW> al estado <READY>", pid)
}

// ///////////////////
// //////// ESTO ES MIO
func PlanificadorLargoPlazo() {
	logger.Info("Iniciando el planificador de largo plazo")
	go ManejadorCreacionProcesos()
	go ManejadorFinalizacionProcesos()
}

func ManejadorCreacionProcesos() {
	for {
		args := <-Utils.ChannelProcessArguments
		fileName := args[0]
		processSize, _ := strconv.Atoi(args[1])
		pid, _ := strconv.Atoi(args[2])

		exito, err := comunicacion.SolicitarCreacionEnMemoria(fileName, processSize)
		if err != nil || !exito {
			go reintentarCreacion(pid, fileName, processSize)
		} else {
			agregarProcesoAReady(pid)
			Utils.MutexPuedoCrearProceso.Unlock()
		}
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
	Utils.MutexNuevo.Lock()
	defer Utils.MutexNuevo.Unlock()

	pcbPtr := BuscarPCBPorPID(pid) //aca busca el PCB de new para pasarlo a READY sirve para el otr tipo de algoritmo
	if pcbPtr == nil {
		logger.Error("No se encontro el PCB con PID <%d>", pid)
		return
	}

	// LO AÑADO A READY

	Utils.MutexReady.Lock()
	globals.ColaReady.Add(pcbPtr)
	Utils.MutexReady.Unlock()
	// LO SACO DE NEW
	Utils.MutexNuevo.Lock()
	globals.ColaNuevo.Remove(pcbPtr)
	Utils.MutexNuevo.Unlock()

	logger.Info("<PID: %d> agregado a READY", pid)
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
		PID: pid,
		PC:  0,
		ME:  make(map[string]int),
		MT:  make(map[string]int),
	}

	logger.Info("## (<%d>) Se crea el proceso - Estado: NEW", pid)

	// Paso 2: Agregar a la cola NEW con protección de mutex
	Utils.MutexNuevo.Lock()
	globals.ColaNuevo.Add(&pcbNuevo)
	Utils.MutexNuevo.Unlock()

	logger.Info("## (<%d>) Pasa del estado <NEW> al estado <READY>", pid)
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
	Utils.MutexNuevo.Lock()
	defer Utils.MutexNuevo.Unlock()

	for _, proceso := range globals.ColaNuevo.Values() {
		if proceso.PID == pid {
			return proceso
		}
	}

	return nil
}
