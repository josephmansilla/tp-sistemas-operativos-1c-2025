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
	var mensaje g.ConsultaContextCPU
	err := json.NewDecoder(r.Body).Decode(&mensaje)
	if err != nil {
		logger.Error("Error leyendo JSON del CPU\n", err)
		http.Error(w, "Error al decodear mensaje del JSON", http.StatusInternalServerError)
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

	g.CalcularEjecutarSleep(time.Duration(g.MemoryConfig.MemoryDelay) * time.Millisecond)

	logger.Info("## PID: <%d> - Obtener instrucción: <%d> - Instrucción: <%s>", mensaje.PID, mensaje.PC, respuesta.Instruccion)

	if err := json.NewEncoder(w).Encode(respuesta); err != nil {
		logger.Error("Error al serializar la obtencion de instruccion: %v", err)
		http.Error(w, "Error al procesar la respuesta", http.StatusInternalServerError)
		return
	}
}

func EnviarEntradaPaginaHandler(w http.ResponseWriter, r *http.Request) {
	g.MutexOperacionMemoria.Lock()
	defer g.MutexOperacionMemoria.Unlock()

	var mensaje g.ConsultaMarco
	err := json.NewDecoder(r.Body).Decode(&mensaje)
	if err != nil {
		logger.Error("Error leyendo JSON del CPU\n", err)
		http.Error(w, "Error al decodear mensaje del JSON", http.StatusInternalServerError)
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

	respuesta := g.RespuestaMarco{
		NumeroMarco: marco,
	}

	logger.Info("## Número Frame enviado: <%d>", marco)

	if err := json.NewEncoder(w).Encode(respuesta); err != nil {
		logger.Error("Error al serializar enviar entrada pagina handler: %v", err)
		http.Error(w, "Error al procesar la respuesta", http.StatusInternalServerError)
		return
	}
}

func LeerEspacioUsuarioHandler(w http.ResponseWriter, r *http.Request) {
	g.MutexOperacionMemoria.Lock()
	defer g.MutexOperacionMemoria.Unlock()

	var mensaje g.LecturaProceso
	err := data.LeerJson(w, r, &mensaje)
	if err != nil {
		logger.Error("Error leyendo JSON del CPU\n", err)
		http.Error(w, "Error al decodear mensaje del JSON", http.StatusInternalServerError)
		return
	}

	pid := mensaje.PID
	direccionFisica := mensaje.DireccionFisica
	tamanioALeer := mensaje.TamanioARecorrer
	respuesta := adm.LeerEspacioEntrada(pid, direccionFisica)
	respuesta = g.RespuestaLectura{
		Exito: respuesta.Exito,
		Valor: respuesta.Valor[:tamanioALeer],
	}

	logger.Info("## PID: <%d> - <Lectura> - Dir. Física: <%d> - Tamaño: <%d>", pid, direccionFisica, tamanioALeer)

	g.CalcularEjecutarSleep(time.Duration(g.MemoryConfig.MemoryDelay) * time.Millisecond)

	logger.Info("## Lectura en espacio del PID <%d> de memoria éxitosa", mensaje.PID)

	if err := json.NewEncoder(w).Encode(respuesta); err != nil {
		logger.Error("Error al serializar la lectura de pagina: %v", err)
		http.Error(w, "Error al procesar la respuesta", http.StatusInternalServerError)
		return
	}
}

func EscribirEspacioUsuarioHandler(w http.ResponseWriter, r *http.Request) {
	g.MutexOperacionMemoria.Lock()
	defer g.MutexOperacionMemoria.Unlock()

	var mensaje g.EscrituraProceso
	if err := data.LeerJson(w, r, &mensaje); err != nil {
		logger.Error("Error leyendo JSON del CPU\n", err)
		http.Error(w, "Error al decodear mensaje del JSON", http.StatusInternalServerError)
		return
	}

	pid := mensaje.PID
	direccionFisica := mensaje.DireccionFisica
	datos := []byte(mensaje.DatosAEscribir)
	tamanioALeer := len(mensaje.DatosAEscribir)

	respuesta := adm.EscribirEspacioEntrada(pid, direccionFisica, datos)
	if respuesta.Exito != nil {
		logger.Error("Escritura con error: %v", respuesta.Exito)
		return
	}

	logger.Info("## PID: <%d> - <Escritura> - Dir. Física: <%d> - Tamaño: <%d>", pid, direccionFisica, tamanioALeer)

	g.CalcularEjecutarSleep(time.Duration(g.MemoryConfig.MemoryDelay) * time.Millisecond)

	logger.Info("## Escritura en espacio del PID <%d> de memoria éxitosa", mensaje.PID)

	if err := json.NewEncoder(w).Encode(respuesta); err != nil {
		logger.Error("Error al serializar la escritura de la pagina: %v", err)
		http.Error(w, "Error al procesar la respuesta", http.StatusInternalServerError)
		return
	}
}

/*
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


	logger.Info("## Lectura Éxitosa")

	if err := json.NewEncoder(w).Encode(respuesta); err != nil {
		logger.Error("Error al serializar la lectura de pagina completa: %v", err)
		http.Error(w, "Error al procesar la respuesta", http.StatusInternalServerError)
		return
	}

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
	datosASobreEscribir := []byte(mensaje.DatosASobreEscribir)
	direccionFisica := mensaje.DireccionFisica

	respuesta := adm.EscribirEspacioEntrada(pid, direccionFisica, datosASobreEscribir)

	logger.Info("## PID: <%d> - Actualizar Página Completa - Dir. Física: <%d>", pid, direccionFisica)

	tiempoTranscurrido := time.Now().Sub(inicio)
	g.CalcularEjecutarSleep(tiempoTranscurrido, retrasoMemoria)


	logger.Info("## Escritura Éxitosa")

	if err := json.NewEncoder(w).Encode(respuesta); err != nil {
		logger.Error("Error al serializar la escritura de pagina completa: %v", err)
		http.Error(w, "Error al procesar la respuesta", http.StatusInternalServerError)
		return
	}
}
*/
