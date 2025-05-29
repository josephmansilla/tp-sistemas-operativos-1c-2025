package utils

import (
	"encoding/json"
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/data"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"net/http"
)

// --------------------------------------------------------
// --------------- FUNCIONALIDAD DE KERNEL ----------------
// --------------------------------------------------------

// RESPONDE AL KERNEL
type RespuestaMemoria struct {
	Exito   bool   `json:"exito"`
	Mensaje string `json:"mensaje"`
}

// FUNCION PARA RECIBIR LOS MENSAJES PROVENIENTES DEL KERNEL
func RecibirMensajeDeKernel(w http.ResponseWriter, r *http.Request) {
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
	respuesta := RespuestaMemoria{
		Exito:   true,
		Mensaje: "Proceso creado correctamente en memoria",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respuesta)
}

func ObtenerEspacioLibreMock(w http.ResponseWriter, r *http.Request) {
	respuesta := globals.EspacioLibreRTA{EspacioLibre: globals.MemoryConfig.MemorySize}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(respuesta); err != nil {
		logger.Error("Error al serializar mock de espacio: %v", err)
	}

	logger.Info("## Espacio libre mock devuelto - Tamaño: <%d>", respuesta.EspacioLibre)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ESPACIO DEVUELTO"))
}
