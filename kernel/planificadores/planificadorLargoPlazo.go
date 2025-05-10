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
		go func(fn string, sz, p int) {

			// Intentar crear en Memoria
			/*ok, err := comunicacion.SolicitarCreacionEnMemoria(fn, sz)
			if err != nil {
				logger.Warn("Error al consultar Memoria para PID %d: %v", p, err)
			}
			if !ok {
				logger.Info("Memoria sin espacio, pid=%d queda pendiente", p)
				// Esperar señal de que un proceso finalizó
				<-Utils.InitProcess
				logger.Info("Recibida señal de espacio libre, reintentando pid=%d", p)
				ok, err = comunicacion.SolicitarCreacionEnMemoria(fn, sz)
				if err != nil || !ok {
					logger.Error("Reintento falló para pid=%d, abortando", p)
					return
				}
			}*/ //CUANDO PEPE HAGA ESO YA SE PUEDE DESCOMENTAR

			// Pasar de NEW a READY
			agregarProcesoAReady(p)
			return //este return hay que sacarlo cuando pepe complete lo suyo
		}(fileName, size, pid)
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
	//MUESTOR LA COLA DE READY PARA VER SI SE AGREGAN CORRECTAMENTE
	MostrarColaReady()
	//MUESTRO LA COLA NEW PARA VER SI ESTAN VACIAS
	MostrarColaNew()

}

func ManejadorFinalizacionProcesos() {
	for {
		pid := <-Utils.ChannelFinishprocess
		logger.Info("ManejadorFinalizacionProcesos: recibida finalización pid=%d", pid)

		// Avisar a Memoria para liberar recursos
		/*err := comunicacion.SolicitarFinProcesoEnMemoria(pid)
		if err != nil {
			logger.Error("Error avisando fin proceso pid=%d: %v", pid, err)
			continue
		}*/ //MEMORIA TIENE QUE RECIBIR ESTE MENSAJE

		// Remover de EXECUTING
		Utils.MutexEjecutando.Lock()
		for _, p := range globals.ColaEjecutando.Values() {
			if p.PID == pid {
				globals.ColaEjecutando.Remove(p)
				break
			}
		}
		Utils.MutexEjecutando.Unlock()

		// Agregar a EXIT
		Utils.MutexSalida.Lock()
		pcbPtr := BuscarPCBPorPID(pid)
		if pcbPtr != nil {
			globals.ColaSalida.Add(pcbPtr)
			logger.Info("PCB pid=%d movido a EXIT", pid)
		} else {
			logger.Warn("No se encontró PCB para PID=%d al mover a EXIT", pid)
		}
		Utils.MutexSalida.Unlock()

		// Señal para reintentos de creación pendientes
		Utils.InitProcess <- struct{}{}
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
func MostrarColaReady() {
	lista := globals.ColaReady.Values()

	if len(lista) == 0 {
		logger.Info("Cola READY vacía")
		return
	}

	logger.Info("Contenido de la cola READY:")
	for _, pcb := range lista {
		logger.Info(" - PCB EN COLA READY con PID: %d", pcb.PID)
	}
}
func MostrarColaNew() {
	lista := globals.ColaNuevo.Values()

	if len(lista) == 0 {
		logger.Info("Cola NEW vacía")
		return
	}

	logger.Info("Contenido de la cola New:")
	for _, pcb := range lista {
		logger.Info(" - PCB EN COLA New con PID: %d", pcb.PID)
	}
}
