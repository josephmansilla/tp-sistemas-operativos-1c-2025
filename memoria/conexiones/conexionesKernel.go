package conexiones

import (
	"encoding/json"
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/data"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"net/http"
)

// FUNCION PARA RECIBIR LOS MENSAJES PROVENIENTES DEL KERNEL
/*func RecibirMensajeDeKernel(w http.ResponseWriter, r *http.Request) {
	var mensaje globals.DatosRespuestaDeKernel

	data.LeerJson(w, r, &mensaje)

	globals.RespuestaKernel = globals.DatosRespuestaDeKernel{
		Pseudocodigo:   mensaje.Pseudocodigo,
		TamanioMemoria: mensaje.TamanioMemoria,
	}

	CargarInstrucciones(mensaje.Pseudocodigo)

	logger.Info("Archivo Pseudocodigo: %s\n", mensaje.Pseudocodigo)
	logger.Info("Tamanio de Memoria Pedido: %d\n", mensaje.TamanioMemoria)

	// RESPUESTA AL KERNEL
	respuesta := globals.RespuestaMemoria{
		Exito:   true,
		Mensaje: "Proceso creado correctamente en memoria",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respuesta)
}*/

func ObtenerEspacioLibre(w http.ResponseWriter, r *http.Request) {

	espacioLibre := globals.CantidadFramesLibres * globals.TamanioMaximoFrame

	respuesta := globals.EspacioLibreRTA{EspacioLibre: espacioLibre}

	logger.Info("## Espacio libre mock devuelto - Tamaño: <%d>", respuesta.EspacioLibre)

	if err := json.NewEncoder(w).Encode(respuesta); err != nil {
		logger.Error("Error al serializar mock de espacio: %v", err)
	}
	json.NewEncoder(w).Encode(respuesta)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ESPACIO DEVUELTO"))
}

// TODO: CAMBIAR CON INICIALIZACIONPROCESO
func RecibirMensajeDeKernel(w http.ResponseWriter, r *http.Request) {
	var mensaje globals.DatosRespuestaDeKernel

	data.LeerJson(w, r, &mensaje)

	globals.RespuestaKernel = globals.DatosRespuestaDeKernel{
		Pseudocodigo:   mensaje.Pseudocodigo,
		TamanioMemoria: mensaje.TamanioMemoria,
		PID:            mensaje.PID,
	}

	CargarInstrucciones(mensaje.PID, mensaje.Pseudocodigo)

	logger.Info("Archivo Pseudocodigo: %s\n", mensaje.Pseudocodigo)
	logger.Info("Tamanio de Memoria Pedido: %d\n", mensaje.TamanioMemoria)

	// RESPUESTA AL KERNEL
	respuesta := globals.RespuestaMemoria{
		Exito:   true,
		Mensaje: "Proceso creado correctamente en memoria",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respuesta)
}

// Mapa global: PID → Lista de instrucciones
var InstruccionesPorPID map[int][]string = make(map[int][]string)

// Cargar instrucción para un PID específico
func CargarInstruccionParaPID(pid int, instruccion string) {
	InstruccionesPorPID[pid] = append(InstruccionesPorPID[pid], instruccion)
	logger.Info("Se cargó una instrucción para PID %d", pid)
}

// Obtener instrucción por PID y PC
func ObtenerInstruccionB(pid int, pc int) string {
	instrucciones, existe := InstruccionesPorPID[pid]
	if !existe || pc < 0 || pc >= len(instrucciones) {
		return ""
	}
	return instrucciones[pc]
}
