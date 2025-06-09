package administracion

import (
	globalData "github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"net/http"
	"time"
)

func SuspenderProceso(w http.ResponseWriter, r *http.Request) {
	//TODO: se escriben sus páginas en el arhcivo
	// TODO: se liberan los frames
	time.Sleep(time.Duration(globalData.DelaySwap) * time.Second)
	logger.Info("## PID: <PID>  - <Lectura> - Dir. Física: <DIRECCIÓN_FÍSICA> - Tamaño: <TAMAÑO>")
}

func DesuspenderProceso(w http.ResponseWriter, r *http.Request) {
	//TODO: se copian las páginas desde el archivo a nuevos frames
	time.Sleep(time.Duration(globalData.DelaySwap) * time.Second)
	logger.Info("## PID: <PID>  - <Lectura> - Dir. Física: <DIRECCIÓN_FÍSICA> - Tamaño: <TAMAÑO>")
}
