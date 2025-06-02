package administracion

import (
	data "github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"net/http"
)

func InicializarTablaRaiz() data.TablaPaginas {
	cantidadEntradasPorTabla := data.MemoryConfig.EntriesPerPage
	tablaRaiz := make(data.TablaPaginas, cantidadEntradasPorTabla)
	return tablaRaiz
} //TODO: EL RESTO DE TABLAS Y ENTRADAS DE PAGINA SE VAN INSTANCIANDO A MEDIDA QUE
// TODO: SE PIDE ASIGNAR X PROCESO EN LA MEMORIA.

func UbicarPagina(numeroFrame int) []int {
	cantidadNiveles := data.MemoryConfig.NumberOfLevels
	entradasPorPagina := data.MemoryConfig.EntriesPerPage

	indice := make([]int, cantidadNiveles)
	divisor := 1

	for i := cantidadNiveles - 1; i >= 0; i-- {
		indice[i] = (numeroFrame / divisor) % entradasPorPagina
		divisor *= entradasPorPagina
	}

	return indice
}

func BuscarEntradaPagina(procesoBuscado *data.Proceso, indices []int) *data.EntradaPagina {
	//	cantidadNiveles := data.MemoryConfig.NumberOfLevels TODO: DEBERIA SER LO MISMO LA LONG DE INDICES Y LA CANT DE NIVELES
	tamanioIndices := len(indices)
	if tamanioIndices == 0 {
		logger.Error("Índice vacío")
		return nil
	}

	tablaApuntada := procesoBuscado.TablaRaiz[indices[0]]
	if tablaApuntada == nil {
		logger.Error("La tabla no existe")
		return nil
	}
	// TODO: optaria por dejar cantidad niveles
	for i := 1; i <= tamanioIndices-1; i++ {
		if tablaApuntada.Subtabla == nil {
			logger.Error("La subtabla no existe")
			// debería ser un error fatal al acceder una locacion
			// de memoria a la que no tiene acceso
			return nil
		}
		tablaApuntada = tablaApuntada.Subtabla[indices[i]]
	}
	if tablaApuntada == nil {
		logger.Error("La tabla no existe")
		return nil
	}
	if tablaApuntada.EntradasPaginas == nil {
		logger.Error("La entrada no existe")
		// debería ser un error fatal al acceder una locacion
		// de memoria a la que no tiene acceso
		return nil
	}

	entradaDeseada := tablaApuntada.EntradasPaginas[indices[tamanioIndices-1]]
	logger.Info("Se encontró la entrada de número: %d", entradaDeseada.NumeroFrame)

	if entradaDeseada.EstaPresente == false {
		logger.Error("No se encuentra presente en memoria el frame")
		// Debería sacarse de SWAP
		return nil
	}

	return entradaDeseada
}

func ObtenerEntradaPagina(pid int, numeroPagina int) int {
	procesoBuscado, err := data.ProcesosMapeable[pid]
	if !err {
		logger.Error("Processo Buscado no existe")
		return -1
	}
	indices := UbicarPagina(numeroPagina)
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
	entradaLibre := -1
	tamanioMaximo := data.MemoryConfig.MemorySize
	framesMemoriaPrincipal := data.FramesLibres
	for i := 0; i < tamanioMaximo; i++ {
		if framesMemoriaPrincipal[i] == true {
			entradaLibre = i
			return entradaLibre
		}
	}
	return entradaLibre
	// retorna la entrada libre
}

func LiberarEntradaPagina(frameALiberar int) {
	framesMemoriaPrincipal := data.FramesLibres
	framesMemoriaPrincipal[frameALiberar] = true
}

func SerializarPagina(pagina data.EntradaPagina, numeroAsignado int) {
	pagina.NumeroFrame = numeroAsignado // Se le asigna para testear
	pagina.EstaPresente = true
	pagina.EstaEnUso = true
	pagina.FueModificado = false
}

func cambiarEstadoPresentePagina(pagina data.EntradaPagina) {
	pagina.EstaPresente = !pagina.EstaPresente
}
func cambiarEstadoUsoPagina(pagina data.EntradaPagina) {
	pagina.EstaEnUso = !pagina.EstaEnUso
}
func cambiarEstadoModificacionPagina(pagina data.EntradaPagina) {
	pagina.FueModificado = !pagina.FueModificado
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
