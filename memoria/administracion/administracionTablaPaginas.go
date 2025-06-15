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
	cantidadNiveles := g.MemoryConfig.NumberOfLevels
	cantidadEntradasPorTabla := g.MemoryConfig.EntriesPerPage

	indices = make([]int, cantidadNiveles)
	divisor := 1

	for i := cantidadNiveles - 1; i >= 0; i-- {
		indices[i] = (nroPagina / divisor) % cantidadEntradasPorTabla
		divisor *= cantidadEntradasPorTabla
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

	for i := 1; i <= tamanioIndices-1; i++ {
		if tablaApuntada.Subtabla == nil {
			logger.Error("La subtabla no existe o nunca fue inicializada")
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

	IncrementarMetrica(procesoBuscado, 1, IncrementarAccesosTablasPaginas)
	return
} // TODO: Testear casos, pero por importancia, no porque tenga dudas

func BuscarEntradaEspecifica(tablaRaiz g.TablaPaginas, numeroEntrada int) (numeroFrameMemReal int) {
	var contador int
	for _, tabla := range tablaRaiz {
		numeroFrameMemReal, encontrado := RecorrerTablasBuscandoEntrada(tabla, numeroEntrada, &contador)
		if encontrado {
			return numeroFrameMemReal
		}
	}
	return -1
}

func RecorrerTablasBuscandoEntrada(tabla *g.TablaPagina, numeroEntrada int, contador *int) (int, bool) {
	if tabla.Subtabla != nil {
		for _, subTabla := range tabla.Subtabla {
			numeroFrame, encontrado := RecorrerTablasBuscandoEntrada(subTabla, numeroEntrada, contador)
			if encontrado {
				return numeroFrame, true
			}
		}
		return -1, false
	}
	for _, entrada := range tabla.EntradasPaginas {
		if *contador == entrada.NumeroFrame {
			return entrada.NumeroFrame, true
		}
		*contador++
	}
	return -1, false
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
} // TODO: HACER ERROR HANDLING

func AsignarNumeroEntradaPagina() int {
	numeroEntradaLibre := -1
	tamanioMaximo := g.MemoryConfig.MemorySize / g.MemoryConfig.PagSize

	for numeroFrame := 0; numeroFrame < tamanioMaximo; numeroFrame++ {
		g.MutexEstructuraFramesLibres.Lock()
		booleano := g.FramesLibres[numeroFrame]
		g.MutexEstructuraFramesLibres.Unlock()

		if booleano == true {
			numeroEntradaLibre = numeroFrame
			MarcarOcupadoFrame(numeroEntradaLibre)

			logger.Info("Marco Libre encontrado: %d", numeroEntradaLibre)
			return numeroEntradaLibre
		}
	}
	return numeroEntradaLibre

} // TODO: ERR HANDLING

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
			// raro caso que no debería pasar pero bue,
			// por las re dudas y que no rompa nada
		}

		fragmentoACargar := informacionEnBytes[offset:end]
		numeroPagina := AsignarNumeroEntradaPagina()
		if numeroPagina == -1 {
			logger.Error("No hay marcos libres")
			break
		} // TODO: not enough for error handling and pretty fucking vage

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
} // HACER ERR HANDLING

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
	stringEnBytes := g.ConversionEnBytes(datosEscritura)
	if len(stringEnBytes) == 0 {
		logger.Error("Los datos a escribir son vacios: %v", logger.ErrNoInstance)
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
	IncrementarMetrica(proceso, 1, IncrementarLecturaDeMemoria)
	logger.Info("## Modificacion del estado entrada exitosa")
}

func LiberarTablaPaginas(tabla *g.TablaPagina, pid int) (err error) {
	err = nil

	if tabla.Subtabla != nil {
		for indice, subtabla := range tabla.Subtabla {
			err := LiberarTablaPaginas(subtabla, pid)
			if err != nil {
				return err
			}
			tabla.Subtabla[indice] = nil
		}
		tabla.Subtabla = nil
	}
	if tabla.EntradasPaginas != nil {
		for _, entrada := range tabla.EntradasPaginas {
			if entrada.EstaPresente {
				tamanioPagina := g.MemoryConfig.PagSize
				direccionFisica := entrada.NumeroFrame * tamanioPagina
				err = RemoverEspacioMemoria(direccionFisica, direccionFisica+tamanioPagina)
				LiberarEntradaPagina(entrada.NumeroFrame)
				if err != nil {
					logger.Error("Error al remover espacio de memoria del frame: \"%d\" ; %v", entrada.NumeroFrame, err)
				}
			}
			// TODO : si está en swap tambien hay que remover
		}
		tabla.EntradasPaginas = nil
	}
	return
}

func LeerPaginaCompletaHandler(w http.ResponseWriter, r *http.Request) {
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

func ActualizarPaginaCompletaHandler(w http.ResponseWriter, r *http.Request) {
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
