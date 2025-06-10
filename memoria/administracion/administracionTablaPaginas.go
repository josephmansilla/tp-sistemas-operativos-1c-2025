package administracion

import (
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	data "github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"net/http"
)

func InicializarTablaRaiz() data.TablaPaginas {
	cantidadEntradasPorTabla := data.MemoryConfig.EntriesPerPage
	return make(data.TablaPaginas, cantidadEntradasPorTabla)
}

func CrearIndicePara(nroPagina int) []int {

	entradas := make([]int, data.CantidadNiveles)
	divisor := 1

	for i := data.CantidadNiveles - 1; i >= 0; i-- {
		entradas[i] = (nroPagina / divisor) % data.EntradasPorPagina
		divisor *= data.EntradasPorPagina
	}

	return entradas
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
		logger.Fatal("La tabla no existe")
		return nil
	}
	// TODO: optaria por dejar cantidad niveles
	for i := 1; i <= tamanioIndices-1; i++ {
		if tablaApuntada.Subtabla == nil {
			logger.Error("La subtabla no está en memoria")
			// TODO: buscar de swap la tabla
			return nil
		}
		tablaApuntada = tablaApuntada.Subtabla[indices[i]]
	}
	if tablaApuntada == nil {
		logger.Error("La tabla no está en memoria")
		// TODO: buscar de swap la tabla
		return nil
	}
	if tablaApuntada.EntradasPaginas == nil {
		logger.Error("La entrada no está en memoria")
		// TODO: buscar de swap la entrada
		return nil
	}

	entradaDeseada := tablaApuntada.EntradasPaginas[indices[tamanioIndices-1]]
	logger.Info("Se encontró la entrada de número: %d", entradaDeseada.NumeroFrame)

	if entradaDeseada.EstaPresente == false {
		logger.Error("No se encuentra presente en memoria el frame")
		// TODO: Debería sacarse de SWAP
		return nil
	}

	IncrementarMetrica(procesoBuscado, IncrementarAccesosTablasPaginas)
	return entradaDeseada
}

func ObtenerEntradaPagina(pid int, indices []int) int {
	procesoBuscado, err := data.ProcesosPorPID[pid]
	if !err {
		logger.Error("Processo Buscado no existe")
		return -1
	}
	entradaPagina := BuscarEntradaPagina(procesoBuscado, indices)
	if entradaPagina == nil {
		// TODO: si se la rta de buscar en swap sigue siendo nil entonces no existe por algun error raro
		logger.Error("No se encontró la entrada de página para el PID: %d", pid)
		return -1
	}
	if !entradaPagina.EstaPresente {
		logger.Error("La entrada de página de número %d y de PID: %d no se encuentra presente ", entradaPagina.NumeroFrame, pid)
		return -1
		// TODO: aca podemos deberiamos sacarlo de swap (segunda verificacion que bueno veremos si queda)
	}
	return entradaPagina.NumeroFrame
}

func AsignarNumeroEntradaPagina() int {
	numeroEntradaLibre := -1
	tamanioMaximo := data.MemoryConfig.MemorySize

	for numeroFrame := 0; numeroFrame < tamanioMaximo; numeroFrame++ {
		if data.FramesLibres[numeroFrame] == true {
			numeroEntradaLibre = numeroFrame

			MarcarOcupadoFrame(numeroEntradaLibre)

			logger.Info("Marco Libre encontrado: %d", numeroEntradaLibre)
			return numeroEntradaLibre
		}
	}
	return numeroEntradaLibre
	// TODO
}

func MarcarOcupadoFrame(numeroFrame int) {
	data.FramesLibres[numeroFrame] = false
	data.CantidadFramesLibres--
}

func LiberarEntradaPagina(numeroFrameALiberar int) {
	framesMemoriaPrincipal := data.FramesLibres
	framesMemoriaPrincipal[numeroFrameALiberar] = true
	//TODO
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

func DividirBytesEnPaginas(informacionEnBytes []byte) [][]byte {
	var paginas [][]byte
	tamanioPagina := globals.TamanioPagina
	totalBytes := len(informacionEnBytes)

	for offset := 0; offset < totalBytes; offset += tamanioPagina {
		end := offset + tamanioPagina
		if end > totalBytes {
			end = totalBytes
		}
		paginas = append(paginas, informacionEnBytes[offset:end])
	}
	return paginas
}

func AsignarDatosAPaginacion(proceso *data.Proceso, informacionEnBytes []byte) error {

	fragmentosPaginas := DividirBytesEnPaginas(informacionEnBytes)

	for numeroPagina := range fragmentosPaginas {
		frame := AsignarNumeroEntradaPagina()
		if frame == -1 {
			logger.Error("No hay marcos libres")
			break
		}

		entrada := &data.EntradaPagina{
			NumeroFrame:   frame,
			EstaPresente:  true,
			EstaEnUso:     true,
			FueModificado: false,
		}
		CargarEntradaMemoria(frame, proceso.PID, fragmentosPaginas[numeroPagina])
		InsertarEntradaPaginaEnTabla(proceso.TablaRaiz, numeroPagina, entrada)
	}
	return nil
}

func InsertarEntradaPaginaEnTabla(tablaRaiz data.TablaPaginas, numeroPagina int, entrada *data.EntradaPagina) {
	indices := CrearIndicePara(numeroPagina)

	actual := tablaRaiz[indices[0]]

	for i := 1; i < len(indices)-1; i++ {
		if actual.Subtabla == nil {
			actual.Subtabla = make(map[int]*data.TablaPagina)
		}
		actual = actual.Subtabla[indices[i]]
	}
	if actual.EntradasPaginas == nil {
		actual.EntradasPaginas = make(map[int]*data.EntradaPagina)
	}
	actual.EntradasPaginas[indices[len(indices)-1]] = entrada
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
