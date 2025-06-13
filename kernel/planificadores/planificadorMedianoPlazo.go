package planificadores

import (
	"github.com/sisoputnfrba/tp-golang/kernel/Utils"
	"github.com/sisoputnfrba/tp-golang/kernel/algoritmos"
	"github.com/sisoputnfrba/tp-golang/kernel/comunicacion"
	"github.com/sisoputnfrba/tp-golang/kernel/pcb"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
)

func PlanificadorMedianoPlazo() {}

func intentarInicializarDesdeSuspReady() bool {
	Utils.MutexSuspendidoReady.Lock()
	defer Utils.MutexSuspendidoReady.Unlock()
	// la consigna PIDE QUE PRIMERO INTENTE PASAR A READY UN PROCESO DE LA LISTA SUSPENDIDO READY LUEGO DE READY

	if algoritmos.ColaSuspendidoReady.IsEmpty() {
		return false
	}

	proceso := algoritmos.ColaSuspendidoReady.First()
	espacio := comunicacion.SolicitarEspacioEnMemoria(proceso.FileName, proceso.ProcessSize)
	if espacio < proceso.ProcessSize {
		logger.Info("No se pudo inicializar proceso desde SUSP.READY PID <%d>", proceso.PID)
		return false
	}
	Utils.MutexSuspendidoReady.Lock()
	algoritmos.ColaSuspendidoReady.Remove(proceso)
	Utils.MutexSuspendidoReady.Unlock()

	Utils.MutexReady.Lock()
	pcb.CambiarEstado(proceso, pcb.EstadoReady)
	algoritmos.ColaReady.Add(proceso)
	Utils.MutexReady.Unlock()

	logger.Info("PID <%d> pasó de SUSP.READY a READY", proceso.PID)
	return true
}
