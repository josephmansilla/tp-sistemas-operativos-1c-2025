package main

import (
	"github.com/sisoputnfrba/tp-golang/kernel/utils"
	"log"
)

// LO TENGO QUE VOLVER A PONER EN UTILS A ESTA FUNCION QUE DEBERÍA ESTAR ACÁ
// HAY ERRORES CON LA IMPORTACIÓN DESDE ESTE ARCHIVO HASTA KERNEL.GO
func IntentarIniciarProceso(tamanioProceso int) {
	espacioLibre, err := utils.ConsultarEspacioLibreMemoria(utils.Config.MemoryAddress,utils.Config.MemoryPort)
	if err != nil {
		log.Println("No se pudo consultar a la memoria por Espacio Libre")
		return
	}

	if espacioLibre >= tamanioProceso {
		log.Println("Hay suficiente espacio libre en Memoria para el proceso")
	} else {
		log.Println("No hay suficiente espacio libre en memoria para el proceso")
	}
}
