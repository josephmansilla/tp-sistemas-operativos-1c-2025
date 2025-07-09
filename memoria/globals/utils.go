package globals

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"os"
	"time"
)

func ParsearContenido(dumpFile *os.File, pid int, contenido []string) {
	comienzo := fmt.Sprintf("## Dump De Memoria Para PID: %d\n\n", pid)
	_, err := dumpFile.WriteString(comienzo)
	for i := 0; i < len(contenido); i++ {
		_, err = dumpFile.WriteString(contenido[i])
		if err != nil {
			logger.Error("Error al escribir contenido en el archivo dump: %v", err)
		}
	}
}

func ConversionEnBytes(stringazo string) []byte {
	return []byte(stringazo)
}

func CalcularCantidadFrames(tamanio int) (cantidad int) {
	cantidad = tamanio / MemoryConfig.PagSize
	if tamanio%MemoryConfig.PagSize > 0 {
		cantidad++
	}
	return
}

func CalcularEjecutarSleep(tiempoTranscurrido time.Duration, retraso time.Duration) {
	tiempoRestante := retraso - tiempoTranscurrido
	if tiempoRestante < retraso {
		logger.Info("Se duerme por %f", tiempoRestante)
		time.Sleep(tiempoRestante)
		logger.Info("Resucité...")
	}
}

func CalcularCantidadEntradasATraer(tamanio int) (resultado int, err error) {
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
