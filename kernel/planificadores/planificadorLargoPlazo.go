package planificadores

import (
	"github.com/sisoputnfrba/tp-golang/kernel/Utils"
	"github.com/sisoputnfrba/tp-golang/kernel/algoritmos"
	"github.com/sisoputnfrba/tp-golang/kernel/comunicacion"
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/pcb"
	logger "github.com/sisoputnfrba/tp-golang/utils/logger"
	"time"
)

// CREA PRIMER PROCESO (NEW)
func CrearPrimerProceso(fileName string, tamanio int) {
	// Paso 1: Crear el PCB
	pid := globals.GenerarNuevoPID()
	pcbNuevo := pcb.PCB{
		PID:            pid,
		PC:             0,
		ME:             make(map[string]int),
		MT:             make(map[string]float64),
		EstimadoRafaga: globals.KConfig.InitialEstimate,
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
		logger.Warn("Memoria sin espacio. Abortando")
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
	go ManejadorInicializacionProcesos()
	go ManejadorFinalizacionProcesos()
	go Reintentador()
}

func ManejadorInicializacionProcesos() {
	for {
		//SIGNAL llega PROCESO a COLA NEW / SUSP.READY
		<-Utils.InitProcess

		logger.Info("LLEGO A LARGO")
		//Al llegar un nuevo proceso a esta cola
		//y la misma esté vacía
		//y no se tengan procesos en la cola de SUSP READY,
		//se enviará un pedido a Memoria para inicializar el mismo.

		//SE TOMA EL SIGUIENTE PROCESO A ENVIAR A READY
		//(ORDENADOS PREVIAMENTE POR ALGORITMO)
		var p *pcb.PCB
		if !algoritmos.ColaSuspendidoReady.IsEmpty() {
			//SI NO ESTA VACIA -> TIENE PRIORIDAD SUSP.READY
			p = algoritmos.ColaSuspendidoReady.First()
		} else {
			p = algoritmos.ColaNuevo.First()
		}

		if p == nil {
			logger.Warn("No hay procesos para inicializar")
			continue
		}

		//MUESTRO LA COLA READY
		MostrarColaReady()
		//MUESTRO LA COLA NEW
		MostrarColaNew()

		filename := p.FileName
		size := p.ProcessSize

		//Intentar crear en Memoria
		espacio := comunicacion.SolicitarEspacioEnMemoria(filename, size)
		if espacio < size {
			logger.Info("Memoria sin espacio, PID <%d> queda pendiente", p.PID)

			// 1) Señal al planificador de corto plazo
			Utils.NotificarDespachador <- p.PID //MANDO PID

			// 2) Señal al planificador Largo para continuar
			Utils.SemProcessCreateOK <- struct{}{}
			continue
		}

		//DICE QUE SI, HAY ESPACIO
		//MANDAR PROCESO A READY
		comunicacion.EnviarArchivoMemoria(filename, size, p.PID)
		agregarProcesoAReady(p, p.Estado)
	}
}

func Reintentador() {
	for {
		//CUANDO MEMORIA LIBERA VUELVE A
		<-Utils.LiberarMemoria
		Utils.InitProcess <- struct{}{}
	}
}

// RECIBIR SYSCALLS DE CREAR PROCESO
func ManejadorCreacionProcesos() {
	logger.Debug("Esperando solicitudes de INIT_PROC para creación de procesos")
	for {
		//SIGNAL de SYSCALL INIT_PROC
		// Recibir filename, size, pid
		msg := <-Utils.ChannelProcessArguments
		fileName := msg.Filename
		size := msg.Tamanio
		pid := msg.PID
		logger.Info("Solicitud INIT_PROC recibida: filename=%s, size=%d, pid=%d", fileName, size, pid)

		// AVISAR QUE SE CREO UN PROCESO AL LARGO PLAZO
		Utils.InitProcess <- struct{}{}
	}
}

func agregarProcesoAReady(proceso *pcb.PCB, estadoAnterior string) {
	// 1) Chequear el PCB existe?
	if proceso == nil {
		logger.Error("PCB no existe (agregarProcesoAReady)")
		return
	}

	// 2) Agregar a READY
	Utils.MutexReady.Lock()
	pcb.CambiarEstado(proceso, pcb.EstadoReady)
	algoritmos.ColaReady.Add(proceso)
	Utils.MutexReady.Unlock()

	// 3) Remover de su cola anterior
	if estadoAnterior == pcb.EstadoSuspReady {
		comunicacion.DesuspensionMemoria(proceso.PID)
		Utils.MutexSuspendidoReady.Lock()
		algoritmos.ColaSuspendidoReady.Remove(proceso)
		Utils.MutexSuspendidoReady.Unlock()

	} else if estadoAnterior == pcb.EstadoNew {
		Utils.MutexNuevo.Lock()
		algoritmos.ColaNuevo.Remove(proceso)
		Utils.MutexNuevo.Unlock()
	}

	logger.Info("## (<%d>) Pasa de estado %s a estado READY", proceso.PID, estadoAnterior)

	// 4) Señal al planificador de corto plazo
	Utils.NotificarDespachador <- proceso.PID //MANDO PID

	// 5) Señal al planificador Largo para continuar
	Utils.SemProcessCreateOK <- struct{}{}

	//MUESTRO LA COLA DE READY PARA VER SI SE AGREGAN CORRECTAMENTE
	MostrarColaReady()
	//MUESTRO LA COLA NEW PARA VER SI ESTAN VACIAS
	MostrarColaNew()
}

// RECIBIR SYSCALLS DE EXIT
func ManejadorFinalizacionProcesos() {
	for {
		//logger.Info("ManejadorFinalizacionProcesos: recibida finalización pid=%d", pid)
		msg := <-Utils.ChannelFinishprocess
		pid := msg.PID
		pc := msg.PC
		cpuID := msg.CpuID

		// Avisar a Memoria para liberar recursos
		comunicacion.LiberarMemoria(pid)

		//Enviar a EXIT con metricas
		finalizarProceso(pid, pc, cpuID)
	}

}

func finalizarProceso(pid int, pc int, cpuID string) {
	var proceso *pcb.PCB = nil

	// 1. Buscar en EXECUTE y remover
	Utils.MutexEjecutando.Lock()
	for _, p := range algoritmos.ColaEjecutando.Values() {
		if p.PID == pid {
			algoritmos.ColaEjecutando.Remove(p)
			logger.Info("## (<%d>) Pasa de estado EXECUTE a estado EXIT", p.PID)
			p.PC = pc
			proceso = p
			break
		}
	}
	Utils.MutexEjecutando.Unlock()

	// 2. Si no está en EXECUTE, buscar en BLOCKED
	if proceso == nil {
		Utils.MutexBloqueado.Lock()
		for _, p := range algoritmos.ColaBloqueado.Values() {
			if p.PID == pid {
				algoritmos.ColaBloqueado.Remove(p)
				logger.Info("## (<%d>) Pasa de estado BLOCKED a estado EXIT", p.PID)
				proceso = p
				break
			}
		}
		Utils.MutexBloqueado.Unlock()
	}

	// 3. Si no está, buscar en SUSP.BLOCKED
	if proceso == nil {
		Utils.MutexBloqueadoSuspendido.Lock()
		for _, p := range algoritmos.ColaBloqueadoSuspendido.Values() {
			if p.PID == pid {
				algoritmos.ColaBloqueadoSuspendido.Remove(p)
				logger.Info("## (<%d>) Pasa de estado SUSP.BLOCKED a estado EXIT", p.PID)
				proceso = p
				break
			}
		}
		Utils.MutexBloqueadoSuspendido.Unlock()
	}

	// 4. Si no está en ninguna, loguear error
	if proceso == nil {
		logger.Error("No se pudo finalizar PID=%d, no encontrado en ninguna cola", pid)
		return
	}

	// 5. Mover a EXIT
	pcb.CambiarEstado(proceso, pcb.EstadoExit)
	Utils.MutexSalida.Lock()
	algoritmos.ColaSalida.Add(proceso)
	Utils.MutexSalida.Unlock()

	// 6. Liberar CPU si corresponde
	if cpuID != "" {
		liberarCPU(cpuID)
	}

	// 7. Log y métricas
	logger.Info("## (<%d>) Finaliza proceso", proceso.PID)
	logger.Info(proceso.ImprimirMetricas())

	// 8. Señal para liberar memoria
	//reintentos de creación pendientes
	Utils.LiberarMemoria <- struct{}{}
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
		logger.Info(" - PCB EN COLA New con PID: %d, TAMAÑO: %d", proceso.PID, proceso.ProcessSize)
	}
}
