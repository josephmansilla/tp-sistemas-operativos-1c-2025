package utils

import (
	globalData "github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"net/http"
)

func InicializarTablaRaiz() map[int]*globalData.TablaPaginasMain {
	cantidadEntradasPorTabla := globalData.MemoryConfig.EntriesPerPage
	tablaRaiz := make(map[int]*globalData.TablaPaginasMain, cantidadEntradasPorTabla)
	return tablaRaiz
} //TODO: EL RESTO DE TABLAS Y ENTRADAS DE PAGINA SE VAN INSTANCIANDO A MEDIDA QUE
// TODO: SE PIDE ASIGNAR X PROCESO EN LA MEMORIA.

func CrearTablas() {
	// TODO: DEBE CREARSE A MEDIDA QUE SE ASIGNA MEMORIA A UN PROCESO
}

func SerializarPagina(pagina globalData.EntradaPagina, numeroAsignado int) {
	pagina.NumeroFrame = numeroAsignado
	pagina.EstaPresente = true
	pagina.EstaEnUso = true
	pagina.FueModificado = false
}

func DescomponerPagina(numeroFrame int) []int {
	cantidadNiveles := globalData.MemoryConfig.NumberOfLevels
	entradasPorPagina := globalData.MemoryConfig.EntriesPerPage

	indice := make([]int, cantidadNiveles)
	divisor := 1

	for i := cantidadNiveles - 1; i >= 0; i-- {
		indice[i] = (numeroFrame / divisor) % entradasPorPagina
		divisor *= entradasPorPagina
	}

	return indice
} // lo usa cpu al final

func BuscarEntradaPagina(procesoBuscado *globalData.Proceso, indices []int) *globalData.EntradaPagina {

	tamanioIndices := len(indices)
	tablaApuntada := procesoBuscado.TablaRaiz[indices[0]]
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

func ObtenerEntradaPagina(pid int, numeroPagina int) int {
	procesoBuscado, err := globalData.ProcesosMapeable[pid]
	if !err {
		logger.Error("Processo Buscado no existe")
		return -1
	}
	indices := DescomponerPagina(numeroPagina)
	entradaPagina := BuscarEntradaPagina(procesoBuscado, indices)
	if entradaPagina == nil {
		logger.Error("No se encontró la entrada de página para el PID: %d", pid)
		return -1
	}
	if !entradaPagina.EstaPresente {
		logger.Error("La entrada de página de número %d y de PID: %d no se encuentra presente ", entradaPagina.NumeroFrame, pid)
		return -1
		// aca podemos deberiamos sacarlo de swap
	}
	return entradaPagina.NumeroFrame
}

func AsignarEntradaPagina() int {
	indiceLibre := -1
	tamanioMaximo := globalData.MemoryConfig.MemorySize
	framesMemoriaPrincipal := globalData.FramesLibres
	for i := 0; i < tamanioMaximo; i++ {
		if framesMemoriaPrincipal[i] == true {
			indiceLibre = i
			return indiceLibre
		}
	}
	return indiceLibre
}

func LiberarEntradaPagina(frameALiberar int) {
	framesMemoriaPrincipal := globalData.FramesLibres
	framesMemoriaPrincipal[frameALiberar] = true
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
