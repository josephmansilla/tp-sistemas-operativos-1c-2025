package planificadores

import (
	"github.com/sisoputnfrba/tp-golang/kernel/comunicacion"
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/pcb"
	logger "github.com/sisoputnfrba/tp-golang/utils/logger"
)

//PLANIFICADOR LARGO PLAZO:
//encargado de gestionar las peticiones a la memoria
//para la creación y eliminación de procesos e hilos.

// CREA PROCESO (NEW)
func CrearProceso(fileName string, tamanio int) {

	// Paso 1: Crear el PCB
	pid := globals.GenerarNuevoPID()
	pcbNuevo := pcb.PCB{
		PID: pid,
		PC:  0,
		ME:  make(map[string]int),
		MT:  make(map[string]int),
	}

	/*
		Al llegar un nuevo proceso a cola (NEW) y la misma esté vacía
		se enviará un pedido a Memoria para inicializar el mismo.
	*/

	if globals.ColaNuevo.IsEmpty() {

		globals.ColaNuevo.Add(&pcbNuevo)
		logger.Info("## (<%d>) Se crea el proceso - Estado: NEW", pid)

		IniciarProceso(pcbNuevo, fileName, tamanio)
	} else {
		/*
			En cambio, si un proceso llega a esta cola y ya hay otros esperando
			se debe tener en cuenta el algoritmo definido
			y verificar si corresponde o no su ingreso.
		*/
		globals.ColaNuevo.Add(&pcbNuevo)
	}
}

// INICIA (READY)
func IniciarProceso(pcb pcb.PCB, fileName string, tamanio int) {
	logger.Info("Intentando crear el proceso con pseudocódigo: %s y tamaño: %d", fileName, tamanio)

	// Paso 1: Pedirle a memoria que reserve espacio
	exito, err := comunicacion.SolicitarCreacionEnMemoria(fileName, tamanio)
	if err != nil {
		logger.Error("Error al intentar reservar memoria: %v", err)
		return
	}

	if !exito {
		logger.Info("Memoria rechazó la creación del proceso (no hay espacio suficiente o error interno)")

		/*
			Si la respuesta es negativa (ya que la Memoria no tiene espacio suficiente para inicializarlo)
			se deberá esperar la finalización de otro proceso para volver a intentar inicializarlo.
		*/
		return

	}

	/*
		Si la respuesta es positiva,
		se pasará el proceso al estado READY y se sigue la misma lógica con el proceso que sigue.
	*/

	globals.ColaReady.Add(&pcb)
	globals.ColaNuevo.Remove(&pcb)
	logger.Info("## (<%d>) Pasa del estado <NEW> al estado <READY>", pcb.PID)
}

// FIN DE PROCESO (EXIT)
func FinalizarProceso(pid int) {
	//NECESITO OBTENER EL PCB para REMOVE E IMPRIMIR METRICAS

	//globals.ColaSalida.Remove();
	logger.Info("## (<%d>) - Finaliza el proceso", pid)
	logger.Info("## (<PID>) - Métricas de estado: NEW (NEW_COUNT) (NEW_TIME), READY (READY_COUNT) (READY_TIME),")
}
