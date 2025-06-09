package conexiones

import (
	"bufio"
	"encoding/json"
	"github.com/sisoputnfrba/tp-golang/memoria/administracion"
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/data"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"net/http"
	"os"
	"strings"
)

func RecibirMensajeDeCPU(w http.ResponseWriter, r *http.Request) {
	var mensaje globals.DatosDeCPU
	data.LeerJson(w, r, &mensaje)

	globals.CPU = globals.DatosDeCPU{
		PID: mensaje.PID,
		PC:  mensaje.PC,
	}

	logger.Info("PID Pedido: %d", mensaje.PID)
	logger.Info("PC Pedido: %d", mensaje.PC)
}

func RetornarMensajeDeCPU(w http.ResponseWriter, r *http.Request) globals.DatosDeCPU {
	var mensaje globals.DatosDeCPU
	data.LeerJson(w, r, &mensaje)

	globals.CPU = globals.DatosDeCPU{
		PID: mensaje.PID,
		PC:  mensaje.PC,
	}
	// StringInstruccion = ObtenerInstruccion(mensaje.PID, mensaje.PC)
	// se debe devolver el string mediante un JSON por el ResponseWriter
	return globals.CPU
}

func ObtenerInstruccion(w http.ResponseWriter, r *http.Request) {
	var mensaje globals.ContextoDeCPU
	err := json.NewDecoder(r.Body).Decode(&mensaje)
	if err != nil {
		http.Error(w, "Error leyendo JSON del CPU\n", http.StatusBadRequest)
		return
	}

	pid := mensaje.PID
	pc := mensaje.PC

	var instruccion string
	if pc >= 0 && pc < len(InstruccionesPorPID[pid]) {
		instruccion = InstruccionesPorPID[pid][pc]
	} else {
		instruccion = "" // Esto indica fin del archivo o error de PC
	}

	respuesta := globals.InstruccionCPU{
		Instruccion: instruccion,
	}

	logger.Info("## PID: <%d>  - Obtener instrucción: <%d> - Instrucción: <%s>", mensaje.PID, mensaje.PC, instruccion)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respuesta)
}

func EnviarConfiguracionMemoria(w http.ResponseWriter, r *http.Request) {
	// Leer el PID desde el cuerpo del request
	var pidData struct {
		PID int `json:"pid"`
	}

	err := data.LeerJson(w, r, &pidData)
	if err != nil {
		// El error ya está logueado por LeerJson
		return
	}
	logger.Info("Recibí petición de configuración desde PID: %d", pidData.PID)

	mensaje := globals.ConsultaConfigMemoria{
		TamanioPagina:    globals.MemoryConfig.PagSize,
		EntradasPorNivel: globals.MemoryConfig.EntriesPerPage,
		CantidadNiveles:  globals.MemoryConfig.NumberOfLevels,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(mensaje); err != nil {
		logger.Error("Error al codificar la respuesta JSON: %v", err)
		http.Error(w, "Error al procesar la respuesta", http.StatusInternalServerError)
	}

}
func CargarInstrucciones(pid int, nombreArchivo string) {
	ruta := "../pruebas/" + nombreArchivo

	file, err := os.Open(ruta)
	if err != nil {
		logger.Error("Error al abrir el archivo: %s\n", err)
		return
	}
	defer file.Close()

	logger.Info("Se leyó el archivo")
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		linea := scanner.Text()
		logger.Info("Línea leída: %s", linea)
		CargarInstruccionParaPID(pid, linea)
		if strings.TrimSpace(linea) == "EOF" {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		logger.Error("Error al leer el archivo: %s", err)
	}

	logger.Info("Total de instrucciones cargadas para PID <%d>: %d", pid, len(InstruccionesPorPID[pid]))
}

func EnviarEntradaPagina(w http.ResponseWriter, r *http.Request) {
	var mensaje globals.MensajePedidoTablaCPU
	err := json.NewDecoder(r.Body).Decode(&mensaje)
	if err != nil {
		http.Error(w, "Error leyendo JSON del CPU\n", http.StatusBadRequest)
		return
	}

	pid := mensaje.PID
	indices := mensaje.IndicesEntrada

	marco := administracion.ObtenerEntradaPagina(pid, indices)

	respuesta := globals.RespuestaTablaCPU{
		NumeroMarco: marco,
	}

	logger.Info("## Número Frame enviado: %d ", marco)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respuesta)
	w.Write([]byte("marco devuelto"))
}
