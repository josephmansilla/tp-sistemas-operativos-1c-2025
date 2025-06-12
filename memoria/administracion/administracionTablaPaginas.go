package administracion

import (
	"encoding/json"
	"fmt"
	g "github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"log"
	"net/http"
	"time"
)

func InicializarTablaRaiz() g.TablaPaginas {
	cantidadEntradasPorTabla := g.MemoryConfig.EntriesPerPage
	return make(g.TablaPaginas, cantidadEntradasPorTabla)
}

func CrearIndicePara(nroPagina int) (indices []int) {

	indices = make([]int, g.CantidadNiveles)
	divisor := 1

	for i := g.CantidadNiveles - 1; i >= 0; i-- {
		indices[i] = (nroPagina / divisor) % g.EntradasPorPagina
		divisor *= g.EntradasPorPagina
	}
	return
}

func BuscarEntradaPagina(procesoBuscado *g.Proceso, indices []int) (entradaDeseada *g.EntradaPagina, err error) {
	err = nil

	tamanioIndices := len(indices)
	if tamanioIndices == 0 {
		logger.Error("Índice vacío")
		return nil, fmt.Errorf("el indice indicado es vacío: %w", logger.ErrIsEmpty)
	}

	tablaApuntada := procesoBuscado.TablaRaiz[indices[0]]
	if tablaApuntada == nil {
		logger.Fatal("La tabla no existe o nunca fue inicializada")
		return nil, fmt.Errorf("la tabla no existe o nunca fue inicializada: %w", logger.ErrNoInstance)
	}
	// TODO: optaria por dejar cantidad niveles
	for i := 1; i <= tamanioIndices-1; i++ {
		if tablaApuntada.Subtabla == nil {
			logger.Error("La subtabla no existe o nunca fue inicializada")
			// TODO: buscar de swap la tabla
			return nil, fmt.Errorf("la subtabla no existe o nunca fue inicializada: %w", logger.ErrNoInstance)
		}
		tablaApuntada = tablaApuntada.Subtabla[indices[i]]
	}
	if tablaApuntada == nil {
		logger.Error("La tabla no exite o no fue nunca inicializada")
		return nil, fmt.Errorf("la tabla no existe o nunca fue inicializada: %w", logger.ErrNoInstance)
	}
	if tablaApuntada.EntradasPaginas == nil {
		logger.Error("La entrada no fue nunca inicializada")
		return nil, fmt.Errorf("la entrada nunca fue inicializada %w", logger.ErrNoInstance)
	}

	entradaDeseada = tablaApuntada.EntradasPaginas[indices[tamanioIndices-1]]
	logger.Info("Se encontró la entrada de número: %d", entradaDeseada.NumeroFrame)

	if entradaDeseada.EstaPresente == false {
		logger.Error("No se encuentra presente en memoria el frame")
		// TODO: Debería sacarse de SWAP y cargarse en memoria
		return entradaDeseada, nil
	}

	IncrementarMetrica(procesoBuscado, IncrementarAccesosTablasPaginas)
	return
}

func ObtenerEntradaPagina(pid int, indices []int) int {
	g.MutexProcesosPorPID.Lock()
	procesoBuscado, errPro := g.ProcesosPorPID[pid]
	g.MutexProcesosPorPID.Unlock()
	if !errPro {
		logger.Error("Processo Buscado no existe")
		return -1
	}
	entradaPagina, errPag := BuscarEntradaPagina(procesoBuscado, indices)
	if errPag != nil {
		logger.Error("Error al buscar la entrada de página")
		return -1
	}
	return entradaPagina.NumeroFrame
}

func AsignarNumeroEntradaPagina() int {
	numeroEntradaLibre := -1
	tamanioMaximo := g.MemoryConfig.MemorySize

	for numeroFrame := 0; numeroFrame < tamanioMaximo; numeroFrame++ {
		g.MutexEstructuraFramesLibres.Lock()
		bool := g.FramesLibres[numeroFrame]
		g.MutexEstructuraFramesLibres.Unlock()

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
	g.MutexEstructuraFramesLibres.Lock()
	g.FramesLibres[numeroFrame] = false
	g.MutexEstructuraFramesLibres.Unlock()

	g.MutexCantidadFramesLibres.Lock()
	g.CantidadFramesLibres--
	g.MutexCantidadFramesLibres.Unlock()
}

func LiberarEntradaPagina(numeroFrameALiberar int) {
	g.MutexEstructuraFramesLibres.Lock()
	g.FramesLibres[numeroFrameALiberar] = true
	g.MutexEstructuraFramesLibres.Unlock()

	g.MutexCantidadFramesLibres.Lock()
	g.CantidadFramesLibres++
	g.MutexCantidadFramesLibres.Unlock()
}

func AsignarDatosAPaginacion(proceso *g.Proceso, informacionEnBytes []byte) error {
	tamanioPagina := g.MemoryConfig.PagSize
	totalBytes := len(informacionEnBytes)

	for offset := 0; offset < totalBytes; offset += tamanioPagina {
		end := offset + tamanioPagina
		if end > totalBytes {
			end = totalBytes
		}

		fragmentoACargar := informacionEnBytes[offset:end]
		numeroPagina := AsignarNumeroEntradaPagina()
		if numeroPagina == -1 {
			logger.Error("No hay marcos libres")
			break
		}

		entradaPagina := &g.EntradaPagina{
			NumeroFrame:   numeroPagina,
			EstaPresente:  true,
			EstaEnUso:     true,
			FueModificado: false,
		}

		direccionFisica := numeroPagina * tamanioPagina
		ModificarEstadoEntradaEscritura(direccionFisica, proceso.PID, fragmentoACargar)
		InsertarEntradaPaginaEnTabla(proceso.TablaRaiz, numeroPagina, entradaPagina)
	}
	return nil
}

func InsertarEntradaPaginaEnTabla(tablaRaiz g.TablaPaginas, numeroPagina int, entrada *g.EntradaPagina) {
	indices := CrearIndicePara(numeroPagina)

	actual := tablaRaiz[indices[0]]

	for i := 1; i < len(indices)-1; i++ {
		if actual.Subtabla == nil {
			actual.Subtabla = make(map[int]*g.TablaPagina)
		}
		actual = actual.Subtabla[indices[i]]
	}
	if actual.EntradasPaginas == nil {
		actual.EntradasPaginas = make(map[int]*g.EntradaPagina)
	}
	actual.EntradasPaginas[indices[len(indices)-1]] = entrada
}

func EscribirEspacioEntrada(pid int, direccionFisica int, datosEscritura string) g.ExitoEscrituraPagina {
	stringEnBytes, err := LecturaPseudocodigo(datosEscritura)
	if err != nil {
		logger.Error("Los datos a escribir son vacios: %v", err)

	}
	ModificarEstadoEntradaEscritura(pid, direccionFisica, stringEnBytes)

	exito := g.ExitoEscrituraPagina{
		Exito:           nil,
		DireccionFisica: direccionFisica,
		Mensaje:         "Proceso fue modificado correctamente en memoria",
	}

	return exito
}

func LeerEspacioEntrada(pid int, direccionFisica int) (datosLectura g.ExitoLecturaPagina) {
	datosLectura = ObtenerDatosMemoria(direccionFisica)
	ModificarEstadoEntradaLectura(pid)
	return datosLectura
}

func ModificarEstadoEntradaLectura(pid int) {
	g.MutexProcesosPorPID.Lock()
	proceso := g.ProcesosPorPID[pid]
	g.MutexProcesosPorPID.Unlock()
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
	inicio := time.Now()
	retrasoSwap := time.Duration(g.MemoryConfig.MemoryDelay) * time.Second

	var mensaje g.LecturaPagina
	err := json.NewDecoder(r.Body).Decode(&mensaje)
	if err != nil {
		http.Error(w, "Error leyendo JSON de Kernel\n", http.StatusBadRequest)
		return
	}

	pid := mensaje.PID
	direccionFisica := mensaje.DireccionFisica
	respuesta := LeerEspacioEntrada(pid, direccionFisica)

	logger.Info("## Leer Página Completa - Dir. Física: <%d>", direccionFisica)

	time.Sleep(time.Duration(g.MemoryConfig.MemoryDelay) * time.Second)

	tiempoTranscurrido := time.Now().Sub(inicio)
	g.CalcularEjecutarSleep(tiempoTranscurrido, retrasoSwap)

	if err := json.NewEncoder(w).Encode(respuesta); err != nil {
		logger.Error("Error al serializar mock de espacio: %v", err)
	}

	logger.Info("## Lectura Éxitosa")

	json.NewEncoder(w).Encode(respuesta)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Respuesta devuelta"))
}

func ActualizarPaginaCompleta(w http.ResponseWriter, r *http.Request) {
	inicio := time.Now()
	retrasoMemoria := time.Duration(g.MemoryConfig.MemoryDelay) * time.Second

	var mensaje g.EscrituraPagina
	err := json.NewDecoder(r.Body).Decode(&mensaje)
	if err != nil {
		http.Error(w, "Error leyendo JSON de Kernel\n", http.StatusBadRequest)
		return
	}

	tamanioNecesario := mensaje.TamanioNecesario
	pid := mensaje.PID
	datosASobreEscribir := mensaje.DatosASobreEscribir
	direccionFisica := mensaje.DireccionFisica

	if tamanioNecesario > g.MemoryConfig.PagSize {
		log.Fatal("No se puede cargar en una pagina este tamaño")
		// TODO: FATAL ...
	}
	respuesta := EscribirEspacioEntrada(pid, direccionFisica, datosASobreEscribir)

	logger.Info("## PID: <%d> - Actualizar Página Completa - Dir. Física: <%d>", pid, direccionFisica)

	time.Sleep(time.Duration(g.MemoryConfig.SwapDelay) * time.Second)

	tiempoTranscurrido := time.Now().Sub(inicio)
	g.CalcularEjecutarSleep(tiempoTranscurrido, retrasoMemoria)

	if err := json.NewEncoder(w).Encode(respuesta); err != nil {
		logger.Error("Error al serializar mock de espacio: %v", err)
	}

	logger.Info("## Escritura Éxitosa")

	json.NewEncoder(w).Encode(respuesta)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Respuesta devuelta"))
}
