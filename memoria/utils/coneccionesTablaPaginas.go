package utils

import (
	globalData "github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"net/http"
)

func InicializarTablas() { // TODO: REVER
	nivelMaximo := globalData.MemoryConfig.NumberOfLevels
	globalData.TablaDePaginas = make(map[int]*globalData.TablaPagina, nivelMaximo)
	globalData.TablaDePaginas[0] = CrearTabla(nivelMaximo - 1)
}

func CrearTabla(nivelActual int) *globalData.TablaPagina {
	if nivelActual == 1 {
		return &globalData.TablaPagina{
			Subtabla: nil,
			Paginas:  make(map[int]*globalData.EntradaPagina),
		}
	}
	return &globalData.TablaPagina{
		Subtabla: make(map[int]*globalData.TablaPagina),
		Paginas:  nil,
	}
}

func CalcularFrames(tamanioMemoriaPrincipal int, tamanioPagina int) int {
	return tamanioMemoriaPrincipal / tamanioPagina
}

func SerializarPagina(pagina globalData.EntradaPagina, numeroAsignado int) {
	pagina.NumeroFrame = numeroAsignado
	pagina.EstaPresente = true
	pagina.EstaEnUso = true
	pagina.FueModificado = false
}

func DesomponerPagina(numeroFrame int) []int {
	cantidadNiveles := globalData.MemoryConfig.NumberOfLevels
	entradasPorPagina := globalData.MemoryConfig.EntriesPerPage

	indice := make([]int, cantidadNiveles)
	divisor := 1

	for i := cantidadNiveles - 1; i >= 0; i-- {
		indice[i] = (numeroFrame / divisor) % entradasPorPagina
		divisor *= entradasPorPagina
	}

	return indice
}

func BuscarEntradaPagina(tablaRaiz globalData.TablaRaizPaginas, indices []int) *globalData.EntradaPagina {
	// err handling
	tamanioIndices := len(indices)
	tablaApuntada := tablaRaiz[indices[0]]
	if tamanioIndices == 0 {
		logger.Error("Índice vacío")
		return nil
	}

	for i := 1; i <= tamanioIndices-1; i++ {
		if tablaApuntada == nil {
			logger.Error("Segment Fault, la tabla no existe")
			return nil
		}
		if tablaApuntada.Subtabla == nil {
			logger.Error("Segment Fault, la subtabla no existe")
			return nil
		}
		tablaApuntada = tablaApuntada.Subtabla[indices[i]]
	}
	if tablaApuntada.Paginas == nil {
		logger.Error("Segment Fault, la entrada no existe")
		return nil
	}

	entradaDeseada := tablaApuntada.Paginas[indices[tamanioIndices-1]]
	logger.Info("Se encontró la entrada de número: %d", entradaDeseada.NumeroFrame)

	if entradaDeseada.EstaPresente == false {
		logger.Error("No se encuentra presente en memoria el frame")
		return nil
	}

	return entradaDeseada
}

// --------------------------------------------------------------------
// ---------- FORMA PARTE DEL ACCESO A LAS TABLAS DE PÁGINAS ----------
// --------------------------------------------------------------------

func AsignarFrame() (int, error) {
	indiceLibre := -1

	return indiceLibre, nil
}

func LiberarFrame(frameALiberar int) {}

func AsignarProceso(PID int, cantidadPaginas int) {

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
