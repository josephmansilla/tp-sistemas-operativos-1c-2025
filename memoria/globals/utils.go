package globals

import (
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"time"
)

func CalcularEjecutarSleep(tiempoTranscurrido time.Duration, retraso time.Duration) {
	tiempoRestante := retraso - tiempoTranscurrido
	if tiempoRestante < retraso {
		logger.Info("Se duerme por %f", tiempoRestante)
		time.Sleep(tiempoRestante)
		logger.Info("Resucité...")
	}
}
