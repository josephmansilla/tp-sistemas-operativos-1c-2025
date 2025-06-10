package administracion

import (
	"encoding/json"
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"log"
	"net/http"
	"time"
)

func InicializarTablaRaiz() globals.TablaPaginas {
	cantidadEntradasPorTabla := globals.MemoryConfig.EntriesPerPage
	return make(globals.TablaPaginas, cantidadEntradasPorTabla)
}

func CrearIndicePara(nroPagina int) []int {

	entradas := make([]int, globals.CantidadNiveles)
	divisor := 1

	for i := globals.CantidadNiveles - 1; i >= 0; i-- {
		entradas[i] = (nroPagina / divisor) % globals.EntradasPorPagina
		divisor *= globals.EntradasPorPagina
	}

	return entradas
}

func BuscarEntradaPagina(procesoBuscado *globals.Proceso, indices []int) *globals.EntradaPagina {
	//	cantidadNiveles := globals.MemoryConfig.NumberOfLevels TODO: DEBERIA SER LO MISMO LA LONG DE INDICES Y LA CANT DE NIVELES
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
	globals.MutexProcesosPorPID.Lock()
	procesoBuscado, err := globals.ProcesosPorPID[pid]
	globals.MutexProcesosPorPID.Unlock()

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
	tamanioMaximo := globals.MemoryConfig.MemorySize

	for numeroFrame := 0; numeroFrame < tamanioMaximo; numeroFrame++ {
		globals.MutexEstructuraFramesLibres.Lock()
		bool := globals.FramesLibres[numeroFrame]
		globals.MutexEstructuraFramesLibres.Unlock()

		if bool == true {
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
	globals.MutexEstructuraFramesLibres.Lock()
	globals.FramesLibres[numeroFrame] = false
	globals.MutexEstructuraFramesLibres.Unlock()

	globals.MutexCantidadFramesLibres.Lock()
	globals.CantidadFramesLibres--
	globals.MutexCantidadFramesLibres.Unlock()
}

func LiberarEntradaPagina(numeroFrameALiberar int) {
	globals.MutexEstructuraFramesLibres.Lock()
	globals.FramesLibres[numeroFrameALiberar] = true
	globals.MutexEstructuraFramesLibres.Unlock()

	//TODO
}

func DividirBytesEnPaginas(informacionEnBytes []byte) [][]byte {
	var paginas [][]byte
	tamanioPagina := globals.TamanioMaximoFrame
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

func AsignarDatosAPaginacion(proceso *globals.Proceso, informacionEnBytes []byte) error {

	fragmentosPaginas := DividirBytesEnPaginas(informacionEnBytes)

	for numeroPagina := range fragmentosPaginas {
		frame := AsignarNumeroEntradaPagina()
		if frame == -1 {
			logger.Error("No hay marcos libres")
			break
		}

		entrada := &globals.EntradaPagina{
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

func InsertarEntradaPaginaEnTabla(tablaRaiz globals.TablaPaginas, numeroPagina int, entrada *globals.EntradaPagina) {
	indices := CrearIndicePara(numeroPagina)

	actual := tablaRaiz[indices[0]]

	for i := 1; i < len(indices)-1; i++ {
		if actual.Subtabla == nil {
			actual.Subtabla = make(map[int]*globals.TablaPagina)
		}
		actual = actual.Subtabla[indices[i]]
	}
	if actual.EntradasPaginas == nil {
		actual.EntradasPaginas = make(map[int]*globals.EntradaPagina)
	}
	actual.EntradasPaginas[indices[len(indices)-1]] = entrada
}

func EscribirEspacioEntrada(pid int, numeroFrame int, datosEscritura string) globals.ExitoEscrituraPagina {
	stringEnBytes := LecturaPseudocodigo(datosEscritura)
	ModificarEstadoEntradaEscritura(pid, numeroFrame, stringEnBytes)

	exito := globals.ExitoEscrituraPagina{
		Exito:           true,
		DireccionFisica: numeroFrame,
		Mensaje:         "Proceso fue modificado correctamente en memoria",
	}

	return exito
}

func ModificarEstadoEntradaEscritura(numeroFrame int, pid int, datosEnBytes []byte) {
	globals.MutexMemoriaPrincipal.Lock()
	globals.MemoriaPrincipal[numeroFrame] = datosEnBytes
	globals.MutexMemoriaPrincipal.Unlock()

	globals.MutexProcesosPorPID.Lock()
	proceso := globals.ProcesosPorPID[pid]
	globals.MutexProcesosPorPID.Unlock()
	IncrementarMetrica(proceso, IncrementarEscrituraDeMemoria)
}
func LeerEspacioEntrada(pid int, numeroFrame int) globals.ExitoLecturaPagina {
	var datosLectura globals.ExitoLecturaPagina = ObtenerDatosMemoria(numeroFrame)
	ModificarEstadoEntradaLectura(pid)
	return datosLectura
}

func ModificarEstadoEntradaLectura(pid int) {
	globals.MutexProcesosPorPID.Lock()
	proceso := globals.ProcesosPorPID[pid]
	globals.MutexProcesosPorPID.Unlock()
	IncrementarMetrica(proceso, IncrementarLecturaDeMemoria)
	logger.Info("## Modificacion del estado entrada exitosa")
}

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
	var mensaje globals.LecturaPagina
	err := json.NewDecoder(r.Body).Decode(&mensaje)
	if err != nil {
		http.Error(w, "Error leyendo JSON de Kernel\n", http.StatusBadRequest)
		return
	}

	pid := mensaje.PID
	direccionFisica := mensaje.DireccionFisica
	respuesta := LeerEspacioEntrada(pid, direccionFisica)

	logger.Info("## Leer Página Completa - Dir. Física: <%d>", direccionFisica)

	time.Sleep(time.Duration(globals.DelayMemoria) * time.Second)

	if err := json.NewEncoder(w).Encode(respuesta); err != nil {
		logger.Error("Error al serializar mock de espacio: %v", err)
	}

	logger.Info("## Lectura Éxitosa")

	json.NewEncoder(w).Encode(respuesta)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Respuesta devuelta"))
}

func ActualizarPaginaCompleta(w http.ResponseWriter, r *http.Request) {
	var mensaje globals.EscrituraPagina
	err := json.NewDecoder(r.Body).Decode(&mensaje)
	if err != nil {
		http.Error(w, "Error leyendo JSON de Kernel\n", http.StatusBadRequest)
		return
	}

	tamanioNecesario := mensaje.TamanioNecesario
	pid := mensaje.PID
	datosASobreEscribir := mensaje.DatosASobreEscribir
	direccionFisica := mensaje.DireccionFisica

	if tamanioNecesario > globals.TamanioMaximoFrame {
		log.Fatal("No se puede cargar en una pagina este tamaño")
		// TODO: FATAL ...
	}
	respuesta := EscribirEspacioEntrada(pid, direccionFisica, datosASobreEscribir)

	logger.Info("## PID: <%d> - Actualizar Página Completa - Dir. Física: <%d>", pid, direccionFisica)

	time.Sleep(time.Duration(globals.DelayMemoria) * time.Second)

	if err := json.NewEncoder(w).Encode(respuesta); err != nil {
		logger.Error("Error al serializar mock de espacio: %v", err)
	}

	logger.Info("## Escritura Éxitosa")

	json.NewEncoder(w).Encode(respuesta)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Respuesta devuelta"))
}
