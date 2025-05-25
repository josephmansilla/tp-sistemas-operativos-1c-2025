package planificadores

import (
	"github.com/sisoputnfrba/tp-golang/kernel/Utils"
	"github.com/sisoputnfrba/tp-golang/kernel/algoritmos"
	"github.com/sisoputnfrba/tp-golang/kernel/comunicacion"
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

	pcb := algoritmos.ColaSuspendidoReady.First()
	exito, err := comunicacion.SolicitarEspacioEnMemoria(pcb.FileName, pcb.ProcessSize)
	if err != nil || !exito {
		logger.Info("No se pudo inicializar proceso desde SUSP.READY PID <%d>", pcb.PID)
		return false
	}
	Utils.MutexSuspendidoReady.Lock()
	algoritmos.ColaSuspendidoReady.Remove(pcb)
	Utils.MutexSuspendidoReady.Unlock()

	Utils.MutexReady.Lock()
	algoritmos.ColaReady.Add(pcb)
	Utils.MutexReady.Unlock()

	logger.Info("PID <%d> pasó de SUSP.READY a READY", pcb.PID)
	return true
}
