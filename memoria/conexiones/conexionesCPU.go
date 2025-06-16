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
	data.LeerJson(w, r, &mensaje)

	g.CPU = g.DatosDeCPU{
		PID: mensaje.PID,
		PC:  mensaje.PC,
	}

	logger.Info("PID Pedido: %d", mensaje.PID)
	logger.Info("PC Pedido: %d", mensaje.PC)
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

	g.MutexProcesosPorPID.Lock()
	proceso := g.ProcesosPorPID[pid]
	g.MutexProcesosPorPID.Unlock()

	respuesta, err := ObtenerInstruccion(proceso, pc)

	logger.Info("## PID: <%d>  - Obtener instrucción: <%d> - Instrucción: <%s>", mensaje.PID, mensaje.PC, respuesta.Instruccion)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respuesta)
	w.Write([]byte("Insturccion devuelta"))
}

func ObtenerInstruccion(proceso *g.Proceso, pc int) (respuesta g.InstruccionCPU, err error) {
	respuesta = g.InstruccionCPU{Exito: nil, Instruccion: ""}
	cantInstrucciones := len(proceso.OffsetInstrucciones)

	var base int
	var tamanioALeer int

	if pc == 0 {
		base = 0
		tamanioALeer = proceso.OffsetInstrucciones[pc]
	} else if pc == cantInstrucciones {
		return
	} else { // Esto indica fin del archivo o error de PC
		base = proceso.OffsetInstrucciones[pc-1]
		tamanioALeer = proceso.OffsetInstrucciones[pc] - base
	}
	tamanioPagina := g.MemoryConfig.PagSize
	numeroEntradaABuscar := base / tamanioPagina
	offsetDir := base % tamanioPagina

	direccionFisica := (adm.BuscarEntradaEspecifica(proceso.TablaRaiz, numeroEntradaABuscar) * tamanioPagina) + offsetDir
	var memoria g.ExitoLecturaMemoria
	memoria, err = adm.LeerEspacioMemoria(proceso.PID, direccionFisica, tamanioALeer)
	if err != nil {
		return
	}
	respuesta = g.InstruccionCPU{Instruccion: memoria.DatosAEnviar}
	return
}

func EnviarConfiguracionMemoriaHandler(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Error leyendo JSON del CPU\n", http.StatusBadRequest)
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
	json.NewEncoder(w).Encode(respuesta)
	w.Write([]byte("marco devuelto"))
}
