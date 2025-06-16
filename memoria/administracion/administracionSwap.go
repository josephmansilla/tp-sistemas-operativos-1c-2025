package administracion

import (
	"encoding/json"
	g "github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"net/http"
	"time"
)

// Para pasar a suspendido:
// 1) Transferir los datos de cada entrada en bytes hacia el archivo, dejar en bytes.
// 2) Marcar cantidad frames como libres y no presente
// 3) Eliminar el contenido de los frames (conceptualmente mal, pero es posible que me ataje algún error de mis funciones)
// 4) Actualizar estructuras necesarias

// Para sacar de suspendido:
// 1) Verificar el tamanio del proceso y si entra en Memoria
// 2) Leer cada entrada de swap, ponerlo en memoria y marcar como presente
// 3) Liberar de SWAP el espacio
// 4) Actualizar estructuras necesarias
// 5) Retonar confirmación éxitosa o fallida

func PasarSwapEntradaPagina(numeroFrame int) {}

func SuspensionProcesoHandler(w http.ResponseWriter, r *http.Request) {
	inicio := time.Now()
	retrasoSwap := time.Duration(g.MemoryConfig.SwapDelay) * time.Second

	var mensaje g.SuspensionProceso
	err := json.NewDecoder(r.Body).Decode(&mensaje)
	if err != nil {
		http.Error(w, "Error leyendo JSON de Kernel\n", http.StatusBadRequest)
		return
	}

	//PasarSwapEntradaPagina(numeroFrame)
	//LiberarEntradaPagina(numeroFrame)

	// TODO cambiar
	// proceso := g.ProcesoSuspendido{}

	//logger.Info("## PID: <%d> - <Lectura> - Dir. Física: <%d> - Tamaño: <%d>", mensaje.PID, proceso.DireccionFisica, proceso.TamanioProceso)

	respuesta := g.ExitoDesuspensionProceso{}

	tiempoTranscurrido := time.Now().Sub(inicio)
	g.CalcularEjecutarSleep(tiempoTranscurrido, retrasoSwap)

	if err := json.NewEncoder(w).Encode(respuesta); err != nil {
		logger.Error("Error al serializar mock de espacio: %v", err)
	}

	json.NewEncoder(w).Encode(respuesta)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Respuesta devuelta"))
}

func SacarEntradaPaginaSwap(numeroFrame int) {}

func LiberarEspacioEnSwap(numeroFrame int) {}

func DesuspensionProcesoHandler(w http.ResponseWriter, r *http.Request) {
	inicio := time.Now()
	retrasoSwap := time.Duration(g.MemoryConfig.SwapDelay) * time.Second

	var mensaje g.DesuspensionProceso
	err := json.NewDecoder(r.Body).Decode(&mensaje)
	if err != nil {
		http.Error(w, "Error leyendo JSON de Kernel\n", http.StatusBadRequest)
		return
	}
	//VerificarTamanioNecesario
	//SacarEntradaPaginaSwap(numeroFrame)
	//LiberarEspacioEnSwap(numeroFrame)
	//ActualizarEstructurasNecesarias

	time.Sleep(time.Duration(g.MemoryConfig.SwapDelay) * time.Second)
	logger.Info("## PID: <%d>  - <Lectura> - Dir. Física: <%d> - Tamaño: <%d>")

	tiempoTranscurrido := time.Now().Sub(inicio)
	g.CalcularEjecutarSleep(tiempoTranscurrido, retrasoSwap)

	respuesta := g.ExitoDesuspensionProceso{}

	if err := json.NewEncoder(w).Encode(respuesta); err != nil {
		logger.Error("Error al serializar mock de espacio: %v", err)
	}

	json.NewEncoder(w).Encode(respuesta)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Respuesta devuelta"))
}
