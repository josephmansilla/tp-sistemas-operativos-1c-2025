package utils

import (
	"net/http"

	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/data"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
)

// --------------------------------------------------------
// --------------- FUNCIONALIDAD DE KERNEL ----------------
// --------------------------------------------------------

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
}
