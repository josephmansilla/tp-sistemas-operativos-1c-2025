package utils

import (
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/pcb"
	logger "github.com/sisoputnfrba/tp-golang/utils/logger"
)

func CrearProceso(fileName string, tamanio int) {
	logger.Info("Intentando crear el proceso con pseudocódigo: %s y tamaño: %d", fileName, tamanio)

	// Paso 1: Pedirle a memoria que reserve espacio
	exito, err := SolicitarCreacionEnMemoria(fileName, tamanio)
	if err != nil {
		logger.Error("Error al intentar reservar memoria: %v", err)
		return
	}

	if !exito {
		logger.Info("Memoria rechazó la creación del proceso (no hay espacio suficiente o error interno)")
		return
	}

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
