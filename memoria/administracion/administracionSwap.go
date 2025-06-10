package administracion

import (
	"encoding/json"
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"net/http"
	"time"
)

func SuspenderProceso(w http.ResponseWriter, r *http.Request) {
	var mensaje globals.SuspensionProceso
	err := json.NewDecoder(r.Body).Decode(&mensaje)
	if err != nil {
		http.Error(w, "Error leyendo JSON de Kernel\n", http.StatusBadRequest)
		return
	}

	PasarSwapEntradaPagina(numeroFrame)
	LiberarEntradaPagina(numeroFrame)

	proceso := globals.ProcesoSuspendido{}

	time.Sleep(time.Duration(globals.DelaySwap) * time.Second)

	logger.Info("## PID: <%d> - <Lectura> - Dir. Física: <%d> - Tamaño: <%d>", mensaje.PID, proceso.DireccionFisica, proceso.TamanioProceso)

	respuesta := globals.ExitoDesuspensionProceso{}

	if err := json.NewEncoder(w).Encode(respuesta); err != nil {
		logger.Error("Error al serializar mock de espacio: %v", err)
	}

	json.NewEncoder(w).Encode(respuesta)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Respuesta devuelta"))
}

func DesuspenderProceso(w http.ResponseWriter, r *http.Request) {

	var mensaje globals.DesuspensionProceso
	err := json.NewDecoder(r.Body).Decode(&mensaje)
	if err != nil {
		http.Error(w, "Error leyendo JSON de Kernel\n", http.StatusBadRequest)
		return
	}
	VerificarTamanioNecesario
	SacarEntradaPaginaSwap(numeroFrame)
	LiberarEspacioEnSwap
	ActualizarEstructurasNecesarias

	time.Sleep(time.Duration(globals.DelaySwap) * time.Second)
	logger.Info("## PID: <%d>  - <Lectura> - Dir. Física: <%d> - Tamaño: <%d>")

	respuesta := globals.ExitoDesuspensionProceso{}

	if err := json.NewEncoder(w).Encode(respuesta); err != nil {
		logger.Error("Error al serializar mock de espacio: %v", err)
	}

	json.NewEncoder(w).Encode(respuesta)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Respuesta devuelta"))

}
