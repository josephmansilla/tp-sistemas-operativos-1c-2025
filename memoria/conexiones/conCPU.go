package conexiones

import (
	"encoding/json"
	adm "github.com/sisoputnfrba/tp-golang/memoria/administracion"
	g "github.com/sisoputnfrba/tp-golang/memoria/estructuras"
	"github.com/sisoputnfrba/tp-golang/utils/data"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"net/http"
	"time"
)

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
	if err := json.NewEncoder(w).Encode(respuesta); err != nil {
		logger.Error("Error al codificar la respuesta JSON: %v", err)
		http.Error(w, "Error al procesar la respuesta", http.StatusInternalServerError)
	}
	//w.Write([]byte("Instruccion devuelta"))
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

// A COMBINAR
func LeerEspacioUsuarioHandler(w http.ResponseWriter, r *http.Request) {
	inicio := time.Now()
	retrasoMemoria := time.Duration(g.MemoryConfig.MemoryDelay) * time.Millisecond

	var mensaje g.LecturaProceso
	err := data.LeerJson(w, r, &mensaje)
	if err != nil {
		return
	}

	pid := mensaje.PID
	direccionFisica := mensaje.DireccionFisica
	tamanioALeer := mensaje.TamanioARecorrer
	respuesta := adm.LeerEspacioEntrada(pid, direccionFisica)
	respuesta = g.ExitoLecturaPagina{
		Exito: respuesta.Exito,
		Valor: respuesta.Valor[:tamanioALeer],
	}
	//respuesta, err := adm.LeerEspacioMemoria(pid, direccionFisica, tamanioALeer)

	logger.Info("## PID: <%d>  - <Lectura> - Dir. Física: <%d> - Tamaño: <%d>", pid, direccionFisica, tamanioALeer)

	tiempoTranscurrido := time.Now().Sub(inicio)
	g.CalcularEjecutarSleep(tiempoTranscurrido, retrasoMemoria)

	if err != nil {
		logger.Error("Error al leer en memoria: %v", err)
		http.Error(w, "Error al leer en memoria", http.StatusInternalServerError)
		return
	}

	logger.Info("## Lectura en espacio de memoria Éxitosa")

	errr := json.NewEncoder(w).Encode(respuesta)
	if errr != nil {
		return
	}
	//w.Write([]byte("Respuesta devuelta"))
}

func EscribirEspacioUsuarioHandler(w http.ResponseWriter, r *http.Request) {
	inicio := time.Now()
	retrasoMemoria := time.Duration(g.MemoryConfig.MemoryDelay) * time.Millisecond

	var mensaje g.EscrituraProceso
	if err := data.LeerJson(w, r, &mensaje); err != nil {
		return
	}

	pid := mensaje.PID
	direccionFisica := mensaje.DireccionFisica
	datos := []byte(mensaje.DatosAEscribir)
	tamanioALeer := len(mensaje.DatosAEscribir)

	respuesta := adm.ModificarEstadoEntradaEscritura(direccionFisica, pid, datos)

	logger.Info("## PID: <%d> - <Escritura> - Dir. Física: <%d> - Tamaño: <%d>", pid, direccionFisica, tamanioALeer)

	tiempoTranscurrido := time.Now().Sub(inicio)
	g.CalcularEjecutarSleep(tiempoTranscurrido, retrasoMemoria)

	/*if err != nil {
		logger.Error("Error al escribir en memoria: %v", err)
		http.Error(w, "Error al escribir en memoria", http.StatusInternalServerError)
		return
	}*/

	logger.Info("## Escritura en espacio de memoria Éxitosa")

	errr := json.NewEncoder(w).Encode(respuesta)
	if errr != nil {
		return
	}
	//w.Write([]byte("Respuesta devuelta"))
}

func LeerPaginaCompletaHandler(w http.ResponseWriter, r *http.Request) {
	inicio := time.Now()
	retrasoSwap := time.Duration(g.MemoryConfig.MemoryDelay) * time.Millisecond

	var mensaje g.LecturaPagina
	if err := data.LeerJson(w, r, &mensaje); err != nil {
		return
	}

	pid := mensaje.PID
	direccionFisica := mensaje.DireccionFisica
	respuesta := adm.LeerEspacioEntrada(pid, direccionFisica)

	logger.Info("## Leer Página Completa - Dir. Física: <%d>", direccionFisica)

	time.Sleep(time.Duration(g.MemoryConfig.MemoryDelay) * time.Second)

	tiempoTranscurrido := time.Now().Sub(inicio)
	g.CalcularEjecutarSleep(tiempoTranscurrido, retrasoSwap)

	if err := json.NewEncoder(w).Encode(respuesta); err != nil {
		logger.Error("Error al serializar mock de espacio: %v", err)
	}

	logger.Info("## Lectura Éxitosa")

	err := json.NewEncoder(w).Encode(respuesta)
	if err != nil {
		return
	}
	//w.Write([]byte("Respuesta devuelta"))
}

func ActualizarPaginaCompletaHandler(w http.ResponseWriter, r *http.Request) {
	inicio := time.Now()
	retrasoMemoria := time.Duration(g.MemoryConfig.MemoryDelay) * time.Millisecond

	var mensaje g.EscrituraPagina
	if err := data.LeerJson(w, r, &mensaje); err != nil {
		return
	}

	if mensaje.TamanioNecesario > g.MemoryConfig.PagSize {
		logger.Error("No se puede cargar en una pagina este tamaño")
		http.Error(w, "No se puede cargar en una pagina este tamaño", http.StatusBadRequest)
		return
	}

	pid := mensaje.PID
	datosASobreEscribir := mensaje.DatosASobreEscribir
	direccionFisica := mensaje.DireccionFisica

	respuesta := adm.EscribirEspacioEntrada(pid, direccionFisica, datosASobreEscribir)

	logger.Info("## PID: <%d> - Actualizar Página Completa - Dir. Física: <%d>", pid, direccionFisica)

	tiempoTranscurrido := time.Now().Sub(inicio)
	g.CalcularEjecutarSleep(tiempoTranscurrido, retrasoMemoria)

	if err := json.NewEncoder(w).Encode(respuesta); err != nil {
		logger.Error("Error al serializar mock de espacio: %v", err)
	}

	logger.Info("## Escritura Éxitosa")

	errr := json.NewEncoder(w).Encode(respuesta)
	if errr != nil {
		return
	}
	//w.Write([]byte("Respuesta devuelta"))
}
