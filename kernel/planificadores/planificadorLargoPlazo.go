package planificadores

import (
	"github.com/sisoputnfrba/tp-golang/kernel/comunicacion"
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/pcb"
	logger "github.com/sisoputnfrba/tp-golang/utils/logger"
)

// CREA Y ENVIA PROCESO A NEW
// PLANIFICADOR LARGO PLAZO
func CrearProceso(fileName string, tamanio int) {
	logger.Info("Intentando crear el proceso con pseudocódigo: %s y tamaño: %d", fileName, tamanio)

	/*
		Al llegar un nuevo proceso a esta cola (NEW) y la misma esté vacía
		se enviará un pedido a Memoria para inicializar el mismo.
	*/

	// Paso 1: Pedirle a memoria que reserve espacio
	exito, err := comunicacion.SolicitarCreacionEnMemoria(fileName, tamanio)
	if err != nil {
		logger.Error("Error al intentar reservar memoria: %v", err)
		return
	}

	/*
		Si la respuesta es negativa (ya que la Memoria no tiene espacio suficiente para inicializarlo)
		se deberá esperar la finalización de otro proceso para volver a intentar inicializarlo.
	*/

	if !exito {
		logger.Info("Memoria rechazó la creación del proceso (no hay espacio suficiente o error interno)")
		return
	}

	/*
		Si la respuesta es positiva,
		se pasará el proceso al estado READY y se sigue la misma lógica con el proceso que sigue.
	*/

	// Paso 2: Crear el PCB y encolarlo
	pid := globals.GenerarNuevoPID()
	pcbNuevo := pcb.PCB{
		PID: pid,
		PC:  0,
		ME:  make(map[string]int),
		MT:  make(map[string]int),
	}

	globals.ColaNuevo.Add(&pcbNuevo)
	logger.Info("Proceso <%v> creado y agregado a la cola NEW", pid)
}

/*
En cambio, si un proceso llega a esta cola y ya hay otros esperando
se debe tener en cuenta el algoritmo definido y verificar si corresponde o no su ingreso.
*/
