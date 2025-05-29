package utils

import (
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"net/http"
)

// --------------------------------------------------------------------
// ---------- FORMA PARTE DEL ACCESO A LAS TABLAS DE PÁGINAS ----------
// --------------------------------------------------------------------

func AsignarFrame() int {
	indiceLibre := -1
	return indiceLibre
}

func LiberarFrame(frameALiberar int) {}

func AsignarProceso(PID int, cantidadPaginas int) {

}

func BuscarPagina(proceso *globals.Proceso, indices []int) *globals.EntradaPagina {
	posActual := proceso.TablaRaiz[indices[0]]

	for i := 0; i < len(indices); i++ {
		if posActual == nil {
			return nil
		}
		if i == (len(indices) - 1) {
			return posActual.Paginas[indices[i]]
		}
		posActual = posActual.Subtabla
	}
	return nil
}

// -.-.--...

func AccesoTablaPaginas(w http.ResponseWriter, r *http.Request) int {

	//TODO

	esTablaIntermedia := false
	numeroTablaSgteNivel := 0
	esTablaUltNivel := false
	numeroFramePagina := 0

	if esTablaIntermedia {
		logger.Info("## Acceso a Tabla intermedia - Núm. Tabla Siguiente: <%d>", numeroTablaSgteNivel)
		return numeroTablaSgteNivel
	}
	if esTablaUltNivel {
		logger.Info("## Acceso a última Tabla - Núm. Frame: <%d>", numeroFramePagina)
		return numeroFramePagina
	}

	return (-1) // EN CASO DE ERROR
}
func LeerPaginaCompleta(w http.ResponseWriter, r *http.Request) {
	// TODO: RETORNAR EL CONTENIDO DESDE LA PAGINA A PARTIR DEL BYTE ENVIADO DE DIRECC FIS. DE LA USER MEMORY
	// TODO: ESTE DEBERÁ COINCIDEIR CON LA POS DEL BYTE 0 DE LA PAGINA
	logger.Info("## Leer Página Completa - Dir. Física: <DIR>")
}
func ActualizarPaginaCompleta(w http.ResponseWriter, r *http.Request) bool {
	err := 1 // valor basura

	// TODO: SE ESCRIBE LA PAGINA COMPLETO A PARTIR DEL BYTE 0 DE LA DIRECCION FISICA ENVIADA

	if err != 0 { // err != nil
		logger.Error("Error al actualizar la página - %s", err)
		return false
	}

	// TODO: RETONAR OK SI ES SATISFACTORIO
	logger.Info("## PID: <PID> - Actualizar Página Completa - Dir. Física: <DIR>")
	return true
}
