package utils

import (
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/data"
	"log"
	"net/http"
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

	log.Printf("Archivo Pseudocodigo: %s\n", mensaje.Pseudocodigo)
	log.Printf("Tamanio de Memoria Pedido: %d\n", mensaje.TamanioMemoria)
}
