package administracion

import (
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"net/http"
)

func SuspenderProceso(w http.ResponseWriter, r *http.Request) {
	//TODO: se escriben sus páginas en el arhcivo
	// TODO: se liberan los frames
	logger.Info("## PID: <PID>  - <Lectura> - Dir. Física: <DIRECCIÓN_FÍSICA> - Tamaño: <TAMAÑO>")
}

func DesuspenderProceso(w http.ResponseWriter, r *http.Request) {
	//TODO: se copian las páginas desde el archivo a nuevos frames
	logger.Info("## PID: <PID>  - <Lectura> - Dir. Física: <DIRECCIÓN_FÍSICA> - Tamaño: <TAMAÑO>")
}
