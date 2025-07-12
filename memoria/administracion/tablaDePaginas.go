package administracion

import (
	"fmt"
	g "github.com/sisoputnfrba/tp-golang/memoria/estructuras"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
)

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
		return nil, logger.ErrProcessNil
	}

	if procesoBuscado.TablaRaiz == nil {
		logger.Error("TablaRaiz es nil en BuscarEntradaPagina")
		return nil, logger.ErrNoTabla
	}

	if len(indices) == 0 {
		logger.Error("Indices vacíos en BuscarEntradaPagina")
		return nil, logger.ErrNoIndices
	}

	tablaApuntada := procesoBuscado.TablaRaiz[indices[0]]
	if tablaApuntada == nil {
		logger.Error("La tabla no existe o nunca fue inicializada")
		return nil, logger.ErrNoTabla
	}

	for i := 1; i < len(indices)-1; i++ {
		if tablaApuntada.Subtabla == nil {
			logger.Error("La subtabla no existe o nunca fue inicializada")
			return nil, logger.ErrNoTabla
		}
		tablaApuntada = tablaApuntada.Subtabla[indices[i]]
		if tablaApuntada == nil {
			logger.Error("La subtabla no existe en el índice <%d>", indices[i])
			return nil, fmt.Errorf("la subtabla no existe en índice %d", indices[i])
		}
	}

	if tablaApuntada.EntradasPaginas == nil {
		logger.Error("Las EntradasPaginas era nil para el índice <%v>", indices)
		return nil, fmt.Errorf("la entrada nunca fue inicializada")
	}

	entradaDeseada = tablaApuntada.EntradasPaginas[indices[len(indices)-1]]
	if entradaDeseada == nil {
		logger.Error("La entrada buscada no existe")
		return nil, fmt.Errorf("la entrada buscada no existe")
	}

	//logger.Info("Se encontró la entrada de número: %d", entradaDeseada.NumeroFrame)

	if entradaDeseada.EstaPresente == false {
		logger.Error("## No se encuentra presente en memoria el frame")
		return entradaDeseada, nil
	}

	IncrementarMetrica(procesoBuscado, 1, IncrementarAccesosTablasPaginas)
	return entradaDeseada, nil
}

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
