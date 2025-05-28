package planificadores

import (
	"github.com/sisoputnfrba/tp-golang/kernel/Utils"
	"github.com/sisoputnfrba/tp-golang/kernel/algoritmos"
	"github.com/sisoputnfrba/tp-golang/kernel/comunicacion"
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/pcb"
	logger "github.com/sisoputnfrba/tp-golang/utils/logger"
	"strconv"
)

// CREA PRIMER PROCESO (NEW)
func CrearPrimerProceso(fileName string, tamanio int) {
	// Paso 1: Crear el PCB
	pid := globals.GenerarNuevoPID()
	pcbNuevo := pcb.PCB{
		PID:            pid,
		PC:             0,
		ME:             make(map[string]int),
		MT:             make(map[string]int),
		EstimadoRafaga: globals.Config.InitialEstimate,
		RafagaRestante: 0,
		FileName:       fileName,
		ProcessSize:    tamanio,
	}

	// Paso 2: Agregar a la cola NEW
	algoritmos.ColaNuevo.Add(&pcbNuevo)
	pcbNuevo.ME[pcb.EstadoNew]++
	pcbNuevo.Estado = pcb.EstadoNew
	logger.Info("## (<%d>) Se crea el Primer proceso - Estado: NEW", pcbNuevo.PID)

	//PASO 3: Mandar archivo pseudocodigo a Memoria
	comunicacion.EnviarArchivoMemoria(fileName, tamanio)
}

func PlanificadorLargoPlazo() {
	logger.Info("Iniciando el planificador de largo plazo")

	//1. Obtener primer proceso de Cola NEW
	var primerProceso *pcb.PCB
	primerProceso = algoritmos.ColaNuevo.First()
	algoritmos.ColaNuevo.Remove(primerProceso)

	//2. Mandar a Ready
	primerProceso.ME[pcb.EstadoReady]++
	primerProceso.Estado = pcb.EstadoReady
	algoritmos.ColaReady.Add(primerProceso)
	Utils.NotificarDespachador <- primerProceso.PID //SIGNAL QUE PASO A READY. MANDO PID
	logger.Info("## (<%d>) Pasa de estado NEW a estado READY", primerProceso.PID)

	go ManejadorCreacionProcesos()
	go ManejadorFinalizacionProcesos()
}

/*
func creacionDeProcesoFifo(){
	//memoria dice si
	//mandar a ready FIRST() de colaNew

	//memoria dice no
}

func creacionDeProcesoSJF(){
	//memoria dice si
	//mandar a ready MASCHICO() de colaNew

}*/

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

			//DICE QUE NO
			//MANDAR NEW

			//DICE QUE SI
			//SWITCH FIFO O SJF

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
		exito, err := comunicacion.SolicitarEspacioEnMemoria(fileName, processSize)
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
	pcbPtr.ME[pcb.EstadoReady]++
	pcbPtr.Estado = pcb.EstadoReady

	Utils.MutexReady.Lock()
	algoritmos.ColaReady.Add(pcbPtr)
	Utils.MutexReady.Unlock()

	logger.Info("## (<%d>) Pasa de estado NEW a estado READY", pcbPtr.PID)
	Utils.NotificarDespachador <- pcbPtr.PID //SIGNAL QUE PASO A READY. MANDO PID

	// 3) Remover de NEW
	algoritmos.ColaNuevo.Remove(pcbPtr)
	Utils.MutexNuevo.Unlock()

	// 4) Señal al planificador para continuar
	Utils.SemProcessCreateOK <- struct{}{}
	//MUESTOR LA COLA DE READY PARA VER SI SE AGREGAN CORRECTAMENTE
	MostrarColaReady()
	//MUESTRO LA COLA NEW PARA VER SI ESTAN VACIAS
	MostrarColaNew()
}

func ManejadorFinalizacionProcesos() {
	for {
		msg := <-Utils.ChannelFinishprocess
		pid := msg.PCB.PID
		cpuID := msg.CpuID
		logger.Info("ManejadorFinalizacionProcesos: recibida finalización pid=%d", pid)

		// Avisar a Memoria para liberar recursos
		/*err := comunicacion.SolicitarFinProcesoEnMemoria(pid)
		if err != nil {
			logger.Error("Error avisando fin proceso pid=%d: %v", pid, err)
			continue
		}*/ //MEMORIA TIENE QUE RECIBIR ESTE MENSAJE

		//ACA PONER UN WAIT/TUBERIA que espere a que memoria libere el proceso, ES PARA TENER UN ORDEN

		// Remover de EXECUTING
		var pcbFinalizado *pcb.PCB = nil
		Utils.MutexEjecutando.Lock()
		for _, p := range algoritmos.ColaEjecutando.Values() {
			if p.PID == pid {
				algoritmos.ColaEjecutando.Remove(p)
				pcbFinalizado = p
				break
			}
		}
		Utils.MutexEjecutando.Unlock()

		if pcbFinalizado == nil {
			logger.Warn("No se encontró PCB para PID=%d al mover a EXIT", pid)
			continue
		}

		// Agregar a EXIT
		Utils.MutexSalida.Lock()
		pcbFinalizado.ME[pcb.EstadoExit]++
		pcbFinalizado.Estado = pcb.EstadoExit
		algoritmos.ColaSalida.Add(pcbFinalizado)
		logger.Info("## (<%d>) Pasa de estado EXECUTE a estado EXIT", pcbFinalizado.PID)
		Utils.MutexSalida.Unlock()

		logger.Info("## (<%d>) - Finaliza el proceso", pcbFinalizado.PID)

		//LIBERAR CPU
		liberarCPU(cpuID)

		logger.Info(pcbFinalizado.ImprimirMetricas())

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
	algoritmos.ColaNuevo.Add(&pcbNuevo)
	logger.Info("SE AÑADIO A LA COLA nuevo DESDE SYSCALL")
	Utils.MutexNuevo.Unlock()

	args := []string{pcbNuevo.FileName, strconv.Itoa(pcbNuevo.ProcessSize), strconv.Itoa(pcbNuevo.PID)}
	Utils.ChannelProcessArguments <- args

}

func intentarInicializarDesdeNew() {
	Utils.MutexNuevo.Lock()
	defer Utils.MutexNuevo.Unlock()

	if algoritmos.ColaNuevo.IsEmpty() {
		return
	}

	pcb := algoritmos.ColaNuevo.First()
	exito, err := comunicacion.SolicitarEspacioEnMemoria(pcb.FileName, pcb.ProcessSize)
	if err != nil || !exito {
		logger.Info("No se pudo inicializar proceso desde NEW PID <%d>", pcb.PID)
		return
	}
	Utils.MutexNuevo.Lock()
	algoritmos.ColaNuevo.Remove(pcb)
	Utils.MutexNuevo.Unlock()

	Utils.MutexReady.Lock()
	algoritmos.ColaReady.Add(pcb)
	Utils.MutexReady.Unlock()
	Utils.NotificarDespachador <- pcb.PID

	logger.Info("PID <%d> pasó de NEW a READY", pcb.PID)
}

func BuscarPCBPorPID(pid int) *pcb.PCB {
	for _, p := range algoritmos.ColaNuevo.Values() {
		if p.PID == pid {
			return p
		}
	}
	return nil
}

func MostrarPCBsEnNew() {
	Utils.MutexNuevo.Lock()
	defer Utils.MutexNuevo.Unlock()

	if algoritmos.ColaReady.IsEmpty() {
		logger.Info("La cola READYYY está vacía.")
		return
	}

	for _, pcb := range algoritmos.ColaReady.Values() {
		logger.Info("- esto son los PCB en NEW CON PID: %d", pcb.PID)
	}
}

func MostrarColaReady() {
	lista := algoritmos.ColaReady.Values()

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
	lista := algoritmos.ColaNuevo.Values()

	if len(lista) == 0 {
		logger.Info("Cola NEW vacía")
		return
	}

	logger.Info("Contenido de la cola New:")
	for _, pcb := range lista {
		logger.Info(" - PCB EN COLA New con PID: %d", pcb.PID)
	}
}
