package planificadores

import (
	"github.com/sisoputnfrba/tp-golang/kernel/Utils"
	"github.com/sisoputnfrba/tp-golang/kernel/algoritmos"
	"github.com/sisoputnfrba/tp-golang/kernel/comunicacion"
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/pcb"
	logger "github.com/sisoputnfrba/tp-golang/utils/logger"
	"strconv"
	"time"
)

// CREA PRIMER PROCESO (NEW)
func CrearPrimerProceso(fileName string, tamanio int) {
	// Paso 1: Crear el PCB
	pid := globals.GenerarNuevoPID()
	estimado := globals.Config.InitialEstimate
	pcbNuevo := pcb.PCB{
		PID:            pid,
		PC:             0,
		ME:             make(map[string]int),
		MT:             make(map[string]float64),
		EstimadoRafaga: estimado,
		RafagaRestante: 0,
		FileName:       fileName,
		ProcessSize:    tamanio,
		TiempoEstado:   time.Now(),
		CpuID:          "",
	}

	//Paso 2: Agregar el primero a la cola NEW
	algoritmos.ColaNuevo.Add(&pcbNuevo)
	pcb.CambiarEstado(&pcbNuevo, pcb.EstadoNew)
	logger.Info("## (<%d>) Se crea el Primer proceso - Estado: <%s>", pcbNuevo.PID, pcbNuevo.Estado)

	//PASO 3: Intentar crear en Memoria
	espacio := comunicacion.SolicitarEspacioEnMemoria(fileName, tamanio)
	if espacio < tamanio {
		logger.Info("Memoria sin espacio. Abortando")
		return
	}

	//PASO 4: Mandar archivo pseudocodigo a Memoria
	comunicacion.EnviarArchivoMemoria(fileName, tamanio, pid)
}

// ARRANCAR LARGO PLAZO Y PRIMER PROCESO AL PRESIONAR ENTER
func PlanificadorLargoPlazo() {
	logger.Info("Iniciando el planificador de largo plazo")

	//1. Obtener primer proceso de Cola NEW
	var primerProceso *pcb.PCB
	primerProceso = algoritmos.ColaNuevo.First()
	algoritmos.ColaNuevo.Remove(primerProceso)

	//2. Mandar a Ready
	algoritmos.ColaReady.Add(primerProceso)
	pcb.CambiarEstado(primerProceso, pcb.EstadoReady)

	Utils.NotificarDespachador <- primerProceso.PID //SIGNAL QUE PASO A READY. MANDO PID
	logger.Info("## (<%d>) Pasa de estado NEW a estado %s", primerProceso.PID, primerProceso.Estado)

	go ManejadorCreacionProcesos()
	go ManejadorFinalizacionProcesos()
}

func ManejadorCreacionProcesos() {
	logger.Info("Esperando solicitudes de INIT_PROC para creación de procesos")
	for {
		//SIGNAL llega PROCESO a COLA NEW
		// Recibir args [filename, size, pid]
		args := <-Utils.ChannelProcessArguments
		fileName := args[0]
		size, _ := strconv.Atoi(args[1])
		pid, _ := strconv.Atoi(args[2])
		logger.Info("Solicitud INIT_PROC recibida: filename=%s, size=%d, pid=%d", fileName, size, pid)

		/*
			Al llegar un nuevo proceso a esta cola
			y la misma esté vacía
			y no se tengan procesos en la cola de SUSP READY,
			se enviará un pedido a Memoria para inicializar el mismo.
		*/

		go func(fn string, sz, p int) {
			if !algoritmos.ColaNuevo.IsEmpty() || !algoritmos.ColaSuspendidoReady.IsEmpty() {
				//CAMINO NEGATIVO
			}

			//Intentar crear en Memoria
			espacio := comunicacion.SolicitarEspacioEnMemoria(fn, sz)
			if espacio < size {
				logger.Info("Memoria sin espacio, pid=%d queda pendiente", p)
				// Esperar señal de que un proceso finalizó
				<-Utils.InitProcess
				logger.Info("Recibida señal de espacio libre, reintentando pid=%d", p)
				espacio = comunicacion.SolicitarEspacioEnMemoria(fn, sz)
				if espacio < size {
					logger.Error("Reintento falló para pid=%d, abortando", p)
					return
				}
			}

			//DICE QUE SI, HAY ESPACIO
			//MANDAR PROCESO A READY
			agregarProcesoAReady(p)
			comunicacion.EnviarArchivoMemoria(fileName, size, pid)

			//return //este return hay que sacarlo cuando pepe complete lo suyo
		}(fileName, size, pid)
	}
}

func reintentarCreacion(pid int, fileName string, processSize int) {
	for {
		logger.Info("<PID: %v> Esperando liberacion de memoria", pid)
		<-Utils.InitProcess
		espacio := comunicacion.SolicitarEspacioEnMemoria(fileName, processSize)
		if espacio > processSize {
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
	Utils.MutexNuevo.Unlock()
	if pcbPtr == nil {
		logger.Error("agregarProcesoAReady: PCB pid=%d no existe en NEW", pid)
		return
	}

	// 2) Agregar a READY
	Utils.MutexReady.Lock()
	pcb.CambiarEstado(pcbPtr, pcb.EstadoReady)
	algoritmos.ColaReady.Add(pcbPtr)
	Utils.MutexReady.Unlock()

	// 3) Remover de NEW
	Utils.MutexNuevo.Lock()
	algoritmos.ColaNuevo.Remove(pcbPtr)
	Utils.MutexNuevo.Unlock()

	logger.Info("## (<%d>) Pasa de estado NEW a estado READY", pcbPtr.PID)

	// 4) Señal al planificador de corto plazo
	Utils.NotificarDespachador <- pcbPtr.PID //MANDO PID

	// 5) Señal al planificador largo para continuar
	Utils.SemProcessCreateOK <- struct{}{}

	//MUESTRO LA COLA DE READY PARA VER SI SE AGREGAN CORRECTAMENTE
	MostrarColaReady()
	//MUESTRO LA COLA NEW PARA VER SI ESTAN VACIAS
	MostrarColaNew()
}

func ManejadorFinalizacionProcesos() {
	for {
		msg := <-Utils.ChannelFinishprocess
		pid := msg.PID
		pc := msg.PC
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
				p.PC = pc //SOBREESCRIBIR NUEVO PC PROVENINIENTE DE CPU
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
		pcb.CambiarEstado(pcbFinalizado, pcb.EstadoExit)
		algoritmos.ColaSalida.Add(pcbFinalizado)
		logger.Info("## (<%d>) Pasa de estado EXECUTE a estado EXIT", pcbFinalizado.PID)
		Utils.MutexSalida.Unlock()

		logger.Info("## (<%d>) - Finaliza el proceso", pcbFinalizado.PID)
		logger.Info(pcbFinalizado.ImprimirMetricas())

		//LIBERAR CPU
		liberarCPU(cpuID)

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
		MT:          make(map[string]float64),
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

	proceso := algoritmos.ColaNuevo.First()
	espacio := comunicacion.SolicitarEspacioEnMemoria(proceso.FileName, proceso.ProcessSize)
	if espacio < proceso.ProcessSize {
		logger.Info("No se pudo inicializar proceso desde NEW PID <%d>", proceso.PID)
		return
	}
	Utils.MutexNuevo.Lock()
	algoritmos.ColaNuevo.Remove(proceso)
	Utils.MutexNuevo.Unlock()

	Utils.MutexReady.Lock()
	algoritmos.ColaReady.Add(proceso)
	Utils.MutexReady.Unlock()
	Utils.NotificarDespachador <- proceso.PID

	logger.Info("PID <%d> pasó de NEW a READY", proceso.PID)
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

	for _, proceso := range algoritmos.ColaReady.Values() {
		logger.Info("- esto son los PCB en NEW CON PID: %d", proceso.PID)
	}
}

func MostrarColaReady() {
	lista := algoritmos.ColaReady.Values()

	if len(lista) == 0 {
		logger.Info("Cola READY vacía")
		return
	}

	logger.Info("Contenido de la cola READY:")
	for _, proceso := range lista {
		logger.Info(" - PCB EN COLA READY con PID: %d", proceso.PID)
	}
}

func MostrarColaNew() {
	lista := algoritmos.ColaNuevo.Values()

	if len(lista) == 0 {
		logger.Info("Cola NEW vacía")
		return
	}

	logger.Info("Contenido de la cola New:")
	for _, proceso := range lista {
		logger.Info(" - PCB EN COLA New con PID: %d", proceso.PID)
	}
}
