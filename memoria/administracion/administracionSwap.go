package administracion

import (
	"encoding/json"
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"net/http"
	"time"
)

func PasarSwapEntradaPagina(numeroFrame int) {}

func SuspenderProceso(w http.ResponseWriter, r *http.Request) {
	inicio := time.Now()
	retrasoSwap := time.Duration(globals.MemoryConfig.SwapDelay) * time.Second

	var mensaje globals.SuspensionProceso
	err := json.NewDecoder(r.Body).Decode(&mensaje)
	if err != nil {
		http.Error(w, "Error leyendo JSON de Kernel\n", http.StatusBadRequest)
		return
	}

	// PasarSwapEntradaPagina(numeroFrame)
	// LiberarEntradaPagina(numeroFrame)

	proceso := globals.ProcesoSuspendido{}

	logger.Info("## PID: <%d> - <Lectura> - Dir. Física: <%d> - Tamaño: <%d>", mensaje.PID, proceso.DireccionFisica, proceso.TamanioProceso)

	respuesta := globals.ExitoDesuspensionProceso{}

	tiempoTranscurrido := time.Now().Sub(inicio)
	globals.CalcularEjecutarSleep(tiempoTranscurrido, retrasoSwap)

	if err := json.NewEncoder(w).Encode(respuesta); err != nil {
		logger.Error("Error al serializar mock de espacio: %v", err)
	}

	json.NewEncoder(w).Encode(respuesta)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Respuesta devuelta"))
}

// TODO: NO ES NECESARIO EL SWAPEO DE TABLAS DE PAGINAS

// TODO: SE LIBERA EN MEMORIA
// TODO: SE ESCRIBE EN SWAP LA INFO NECESARIA

func SacarEntradaPaginaSwap(numeroFrame int) {}

func LiberarEspacioEnSwap(numeroFrame int) {}

func DesuspenderProceso(w http.ResponseWriter, r *http.Request) {
	inicio := time.Now()
	retrasoSwap := time.Duration(globals.MemoryConfig.SwapDelay) * time.Second

	var mensaje globals.DesuspensionProceso
	err := json.NewDecoder(r.Body).Decode(&mensaje)
	if err != nil {
		http.Error(w, "Error leyendo JSON de Kernel\n", http.StatusBadRequest)
		return
	}
	//VerificarTamanioNecesario
	//SacarEntradaPaginaSwap(numeroFrame)
	//LiberarEspacioEnSwap(numeroFrame)
	// TODO: ActualizarEstructurasNecesarias

	time.Sleep(time.Duration(globals.MemoryConfig.SwapDelay) * time.Second)
	logger.Info("## PID: <%d>  - <Lectura> - Dir. Física: <%d> - Tamaño: <%d>")

	tiempoTranscurrido := time.Now().Sub(inicio)
	globals.CalcularEjecutarSleep(tiempoTranscurrido, retrasoSwap)

	respuesta := globals.ExitoDesuspensionProceso{}

	if err := json.NewEncoder(w).Encode(respuesta); err != nil {
		logger.Error("Error al serializar mock de espacio: %v", err)
	}

	json.NewEncoder(w).Encode(respuesta)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Respuesta devuelta"))
}

// TODO: VERIFICAR EL TAMAÑO NECESARIO

// TODO: LEER EL CONTENIDO DEL SWAP, ESCRIBIERLO EN EL FRAME ASIGNADO
// TODO: LIBERAR ESPACIO EN SWAP
// TODO: ACTUALIZAR ESTRUCTURAS NECESARIAS

// TODO: RETORNAR OK
