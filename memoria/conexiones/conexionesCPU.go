package conexiones

import (
	"bufio"
	"encoding/json"
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/memoria/utils"
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

/*func CargarInstrucciones(nombreArchivo string) {

	ruta := "../pruebas/" + nombreArchivo
	// con el (..) vuelve para atras en los directorios
	// accede a la carpeta de pruebas y abre el archivo pasado x parametro

	file, err := os.Open(ruta)
	if err != nil {
		logger.Error("Error al abrir el archivo: %s\n", err)
		return
	}
	defer file.Close() // se accede desde cualquier parte del código

	logger.Info("Se leyó el archivo")
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lineaPseudocodigo := scanner.Text()
		logger.Info("Línea leída:%s", lineaPseudocodigo)
		utils.CargarListaDeInstrucciones(lineaPseudocodigo)
		if strings.TrimSpace(lineaPseudocodigo) == "EOF" {
			break
		}

	}
	if err := scanner.Err(); err != nil {
		logger.Error("Error al leer el archivo:%s", err)
	}
	logger.Info("Total de instrucciones cargadas: %d", len(utils.Instrucciones))
}*/

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
	if pc >= 0 && pc < len(utils.InstruccionesPorPID[pid]) {
		instruccion = utils.InstruccionesPorPID[pid][pc]
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
		utils.CargarInstruccionParaPID(pid, linea)
		if strings.TrimSpace(linea) == "EOF" {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		logger.Error("Error al leer el archivo: %s", err)
	}

	logger.Info("Total de instrucciones cargadas para PID <%d>: %d", pid, len(utils.InstruccionesPorPID[pid]))
}
