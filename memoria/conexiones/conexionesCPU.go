package conexiones

import (
	"encoding/json"
	adm "github.com/sisoputnfrba/tp-golang/memoria/administracion"
	g "github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/data"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"net/http"
)

func RecibirMensajeDeCPUHandler(w http.ResponseWriter, r *http.Request) {
	var mensaje g.DatosDeCPU
	err := data.LeerJson(w, r, &mensaje)
	if err != nil {
		return
	}

	g.CPU = g.DatosDeCPU{
		PID: mensaje.PID,
		PC:  mensaje.PC,
	}

	logger.Info("PID Pedido: %d ; PC Pedido: %d", mensaje.PID, mensaje.PC)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(mensaje); err != nil {
		logger.Error("Error al codificar la respuesta JSON: %v", err)
		http.Error(w, "Error al procesar la respuesta", http.StatusInternalServerError)
	}
	//w.Write([]byte("Instruccion devuelta"))
}

func ObtenerInstruccionHandler(w http.ResponseWriter, r *http.Request) {
	var mensaje g.ContextoDeCPU
	err := json.NewDecoder(r.Body).Decode(&mensaje)
	if err != nil {
		http.Error(w, "Error leyendo JSON del CPU\n", http.StatusBadRequest)
		return
	}

	pid := mensaje.PID
	pc := mensaje.PC

	logger.Info("Petición de instrucción para PID: %d - PC: %d", mensaje.PID, mensaje.PC)

	g.MutexProcesosPorPID.Lock()
	proceso, ok := g.ProcesosPorPID[pid]
	g.MutexProcesosPorPID.Unlock()

	if !ok || proceso == nil {
		logger.Error("Proceso con PID %d no existe o es nil", mensaje.PID)
		http.Error(w, "Proceso no encontrado", http.StatusNotFound)
		return
	}

	respuesta, err := adm.ObtenerInstruccion(proceso, pc)
	if err != nil {
		logger.Error("Error al obtener instrucción: %v", err)
		http.Error(w, "Error al obtener instrucción", http.StatusInternalServerError)
		return
	}

	logger.Info("## PID: <%d>  - Obtener instrucción: <%d> - Instrucción: <%s>", mensaje.PID, mensaje.PC, respuesta.Instruccion)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(mensaje); err != nil {
		logger.Error("Error al codificar la respuesta JSON: %v", err)
		http.Error(w, "Error al procesar la respuesta", http.StatusInternalServerError)
	}
	//w.Write([]byte("Instruccion devuelta"))
}

func EnviarConfiguracionMemoriaHandler(w http.ResponseWriter, r *http.Request) {
	// Leer el PID desde el cuerpo del request
	var pidData struct {
		PID int `json:"pid"`
	}

	err := data.LeerJson(w, r, &pidData)
	if err != nil {
		return
	}
	logger.Info("Recibí petición de configuración desde PID: %d", pidData.PID)

	mensaje := g.ConsultaConfigMemoria{
		TamanioPagina:    g.MemoryConfig.PagSize,
		EntradasPorNivel: g.MemoryConfig.EntriesPerPage,
		CantidadNiveles:  g.MemoryConfig.NumberOfLevels,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(mensaje); err != nil {
		logger.Error("Error al codificar la respuesta JSON: %v", err)
		http.Error(w, "Error al procesar la respuesta", http.StatusInternalServerError)
	}
}

func EnviarEntradaPaginaHandler(w http.ResponseWriter, r *http.Request) {
	var mensaje g.MensajePedidoTablaCPU
	err := json.NewDecoder(r.Body).Decode(&mensaje)
	if err != nil {
		return
	}

	pid := mensaje.PID
	indices := mensaje.IndicesEntrada
	var marco int
	marco, err = adm.ObtenerEntradaPagina(pid, indices)
	if err != nil {
		logger.Error("Error: %v", err)
		http.Error(w, "Error al Leer espacio de Memoria \n", http.StatusInternalServerError)
	}

	respuesta := g.RespuestaTablaCPU{
		NumeroMarco: marco,
	}

	logger.Info("## Número Frame enviado: %d ", marco)

	w.Header().Set("Content-Type", "application/json")
	errEncode := json.NewEncoder(w).Encode(respuesta)
	if errEncode != nil {
		return
	}
	//w.Write([]byte("marco devuelto"))
}
