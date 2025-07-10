package administracion

import (
	"encoding/json"
	"fmt"
	g "github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/data"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"net/http"
	"time"
)

func InicializarTablaRaiz() g.TablaPaginas {
	cantidadEntradasPorTabla := g.MemoryConfig.EntriesPerPage
	tabla := make(g.TablaPaginas, cantidadEntradasPorTabla)
	for i := 0; i < cantidadEntradasPorTabla; i++ {
		tabla[i] = &g.TablaPagina{} // inicializo cada entrada con un puntero válido
	}
	return tabla
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
	if procesoBuscado == nil {
		logger.Error("Proceso es nil en BuscarEntradaPagina")
		return nil, fmt.Errorf("proceso nil")
	}

	if procesoBuscado.TablaRaiz == nil {
		logger.Error("TablaRaiz es nil en BuscarEntradaPagina")
		return nil, fmt.Errorf("tabla raiz nil")
	}

	if len(indices) == 0 {
		logger.Error("Indices vacíos en BuscarEntradaPagina")
		return nil, fmt.Errorf("indice vacío")
	}

	tablaApuntada := procesoBuscado.TablaRaiz[indices[0]]
	if tablaApuntada == nil {
		logger.Error("La tabla no existe o nunca fue inicializada")
		return nil, fmt.Errorf("la tabla no existe o nunca fue inicializada")
	}

	for i := 1; i < len(indices)-1; i++ {
		if tablaApuntada.Subtabla == nil {
			logger.Error("La subtabla no existe o nunca fue inicializada")
			return nil, fmt.Errorf("la subtabla no existe o nunca fue inicializada")
		}
		tablaApuntada = tablaApuntada.Subtabla[indices[i]]
		if tablaApuntada == nil {
			logger.Error("La subtabla no existe en índice %d", indices[i])
			return nil, fmt.Errorf("la subtabla no existe en índice %d", indices[i])
		}
	}

	if tablaApuntada.EntradasPaginas == nil {
		logger.Error("EntradasPaginas era nil para indices=%v", indices)
		logger.Error("La entrada no fue nunca inicializada")
		return nil, fmt.Errorf("la entrada nunca fue inicializada")
	}

	entradaDeseada = tablaApuntada.EntradasPaginas[indices[len(indices)-1]]
	if entradaDeseada == nil {
		logger.Error("La entrada buscada no existe")
		return nil, fmt.Errorf("la entrada buscada no existe")
	}

	logger.Info("Se encontró la entrada de número: %d", entradaDeseada.NumeroFrame)

	if entradaDeseada.EstaPresente == false {
		logger.Error("No se encuentra presente en memoria el frame")
		// TODO: Debería sacarse de SWAP y cargarse en memoria
		return entradaDeseada, nil
	}

	IncrementarMetrica(procesoBuscado, 1, IncrementarAccesosTablasPaginas)
	return entradaDeseada, nil
} // TODO: Testear casos, pero por importancia, no porque tenga dudas

func BuscarEntradaEspecifica(tablaRaiz g.TablaPaginas, numeroEntrada int) (numeroFrameMemReal int) {
	var contador *int
	for _, tabla := range tablaRaiz {
		numeroFrameMemReal, encontrado := RecorrerTablasBuscandoEntrada(tabla, numeroEntrada, contador)
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

func ObtenerEntradaPagina(pid int, indices []int) (int, error) {
	g.MutexProcesosPorPID.Lock()
	procesoBuscado, errPro := g.ProcesosPorPID[pid]
	g.MutexProcesosPorPID.Unlock()
	if !errPro {
		logger.Error("Processo Buscado no existe")
		return -1, fmt.Errorf("el proceso no existe o nunca fue inicializada: %w", logger.ErrNoInstance)
	}
	entradaPagina, errPag := BuscarEntradaPagina(procesoBuscado, indices)
	if errPag != nil {
		logger.Error("Error al buscar la entrada de página")
		return -1, fmt.Errorf("la entrada no existe o nunca fue inicializada: %w", logger.ErrNoInstance)
	}
	return entradaPagina.NumeroFrame, nil
}

func AsignarFrameLibre() (int, error) {
	cantidadFramesTotales := g.MemoryConfig.MemorySize / g.MemoryConfig.PagSize

	for numeroFrame := 0; numeroFrame < cantidadFramesTotales; numeroFrame++ {
		g.MutexEstructuraFramesLibres.Lock()
		booleano := g.FramesLibres[numeroFrame]
		g.MutexEstructuraFramesLibres.Unlock()

		if booleano == true {
			MarcarOcupadoFrame(numeroFrame)
			logger.Info("Marco Libre encontrado: %d", numeroFrame)
			return numeroFrame, nil
		}
	}
	return -10, logger.ErrNoMemory

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

func AsignarPaginasParaPID(proceso *g.Proceso, tamanio int) error {
	logger.Info("Asignando espacio de páginas para PID <%d>...", proceso.PID)

	cantidadFrames := g.CalcularCantidadFrames(tamanio)
	for i := 1; i <= cantidadFrames; i++ {
		numeroFrame, err := AsignarFrameLibre()
		if err != nil {
			logger.Error("no hay frames libres en el sistema %v", err)
			return err
		}
		entradaPagina := &g.EntradaPagina{
			NumeroFrame:   numeroFrame,
			EstaPresente:  true,
			EstaEnUso:     false,
			FueModificado: false,
		}
		InsertarEntradaPaginaEnTabla(proceso.TablaRaiz, numeroFrame, entradaPagina)
		logger.Info("El entrada #%d para el PID: <%d> se guardó en el frame <%d>...", i, proceso.PID, numeroFrame)
	}
	logger.Info("Se reservó correctamente el espacio para el PID: <%d>.", proceso.PID)
	logger.Info("Quendan <%d> frames libes", g.CantidadFramesLibres)
	return nil
}

func InsertarEntradaPaginaEnTabla(tablaRaiz g.TablaPaginas, numeroPagina int, entrada *g.EntradaPagina) {
	indices := CrearIndicePara(numeroPagina)
	actual := tablaRaiz[indices[0]]

	if actual == nil {
		actual = &g.TablaPagina{}
		tablaRaiz[indices[0]] = actual
	}

	for i := 1; i < len(indices)-1; i++ {
		if actual.Subtabla == nil {
			actual.Subtabla = make(map[int]*g.TablaPagina)
		}
		if actual.Subtabla[indices[i]] == nil {
			actual.Subtabla[indices[i]] = &g.TablaPagina{}
		}
		actual = actual.Subtabla[indices[i]]
	}

	if actual.EntradasPaginas == nil {
		actual.EntradasPaginas = make(map[int]*g.EntradaPagina)
	}
	actual.EntradasPaginas[indices[len(indices)-1]] = entrada

	logger.Info("Insertada entrada para página %d (indices=%v)", numeroPagina, indices)
}

func EscribirEspacioEntrada(pid int, direccionFisica int, datosEscritura string) g.ExitoEscrituraPagina {
	stringEnBytes := g.ConversionEnBytes(datosEscritura)
	if len(stringEnBytes) == 0 {
		logger.Error("Los datos a escribir son vacios: %v", logger.ErrNoInstance)
	}
	err := ModificarEstadoEntradaEscritura(pid, direccionFisica, stringEnBytes)
	if err != nil {
		return g.ExitoEscrituraPagina{Exito: err, DireccionFisica: direccionFisica, Mensaje: err.Error()}
	}

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

func ObtenerInstruccion(proceso *g.Proceso, pc int) (respuesta g.InstruccionCPU, err error) {
	respuesta = g.InstruccionCPU{Exito: nil, Instruccion: ""}

	if proceso == nil {
		logger.Error("Proceso recibido es nil")
		return respuesta, fmt.Errorf("proceso nil")
	}

	cantInstrucciones := len(proceso.InstruccionesEnBytes)
	logger.Info("PID <%d> - PC: %d - Cant. Instrucciones: %d", proceso.PID, pc, cantInstrucciones)

	lineaInstruccion := proceso.InstruccionesEnBytes[pc]
	logger.Info("PC a buscar: %d, Instrucción en bytes: %v", pc, lineaInstruccion)
	logger.Info("PC a buscar: %d, Instrucción: %s", pc, string(lineaInstruccion))

	respuesta.Instruccion = string(lineaInstruccion)
	return respuesta, nil
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
	retrasoSwap := time.Duration(g.MemoryConfig.MemoryDelay) * time.Millisecond

	var mensaje g.LecturaPagina
	if err := data.LeerJson(w, r, &mensaje); err != nil {
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

	err := json.NewEncoder(w).Encode(respuesta)
	if err != nil {
		return
	}
	w.WriteHeader(http.StatusOK)
	//w.Write([]byte("Respuesta devuelta"))
}

func ActualizarPaginaCompletaHandler(w http.ResponseWriter, r *http.Request) {
	inicio := time.Now()
	retrasoMemoria := time.Duration(g.MemoryConfig.MemoryDelay) * time.Millisecond

	var mensaje g.EscrituraPagina
	if err := data.LeerJson(w, r, &mensaje); err != nil {
		return
	}

	if mensaje.TamanioNecesario > g.MemoryConfig.PagSize {
		logger.Error("No se puede cargar en una pagina este tamaño")
		http.Error(w, "No se puede cargar en una pagina este tamaño", http.StatusBadRequest)
		return
	}

	pid := mensaje.PID
	datosASobreEscribir := mensaje.DatosASobreEscribir
	direccionFisica := mensaje.DireccionFisica

	respuesta := EscribirEspacioEntrada(pid, direccionFisica, datosASobreEscribir)

	logger.Info("## PID: <%d> - Actualizar Página Completa - Dir. Física: <%d>", pid, direccionFisica)

	tiempoTranscurrido := time.Now().Sub(inicio)
	g.CalcularEjecutarSleep(tiempoTranscurrido, retrasoMemoria)

	if err := json.NewEncoder(w).Encode(respuesta); err != nil {
		logger.Error("Error al serializar mock de espacio: %v", err)
	}

	logger.Info("## Escritura Éxitosa")

	errr := json.NewEncoder(w).Encode(respuesta)
	if errr != nil {
		return
	}
	w.WriteHeader(http.StatusOK)
	//w.Write([]byte("Respuesta devuelta"))
}
