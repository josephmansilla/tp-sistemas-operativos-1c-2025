package estructuras

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"time"
)

func CalcularEjecutarSleep(tiempoTranscurrido time.Duration, retraso time.Duration) {
	tiempoRestante := retraso - tiempoTranscurrido
	if tiempoRestante < retraso {
		time.Sleep(tiempoRestante)
	}
}

func CalcularCantidadEntradas(tamanio int) (resultado int, err error) {
	err = nil
	resultado = 0
	if tamanio < 0 {
		return resultado, fmt.Errorf("el tamanio pedido de espacio es 0 %v", logger.ErrBadRequest)
	}

	resultado = tamanio / MemoryConfig.PagSize
	if tamanio%MemoryConfig.PagSize > 0 {
		resultado++
	}
	return
}
