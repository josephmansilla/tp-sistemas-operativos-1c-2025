package administracion

import (
	"fmt"
	g "github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
)

func LiberarProceso(pid int) (metricas g.MetricasProceso, err error) {
	var proceso *g.Proceso
	metricas = g.MetricasProceso{}
	err = nil

	proceso, err = DesocuparProcesoDeEstructurasGlobales(pid)
	if err != nil {
		return metricas, err
	}
	metricas = proceso.Metricas
	for _, tabla := range proceso.TablaRaiz {
		err := LiberarTablaPaginas(tabla, pid)
		if err != nil {
			return g.MetricasProceso{}, err
		}
	}
	logger.Info("## Se liberó la memoria para el PID: %d", pid)
	return
}

// ========== LIBERAR VECTORES GLOBALES ==========
func DesocuparProcesoDeEstructurasGlobales(pid int) (proceso *g.Proceso, err error) {
	err = nil
	g.MutexProcesosPorPID.Lock()
	proceso = g.ProcesosPorPID[pid]
	delete(g.ProcesosPorPID, pid)
	g.MutexProcesosPorPID.Unlock()

	if proceso == nil {
		logger.Error("El proceso no está en el Slice de procesos mapeado por PID")
		return proceso, fmt.Errorf("no hay una instancia de pid \"%d\" en el slice de procesos por PID %v", pid, logger.ErrNoInstance)
	}

	g.MutexSwapIndex.Lock()
	delete(g.SwapIndex, pid)
	g.MutexSwapIndex.Unlock()

	return
}

// ========== DEJAR NULO LOS PUNTEROS DE LA TABLA DE PAGINAS ==========
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
				MarcarLibreFrame(entrada.NumeroFrame)
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

// ========== LIBERO EL ESPACIO EN MEMORIA ==========
func RemoverEspacioMemoria(inicio int, limite int) (err error) {
	espacioVacio := make([]byte, limite-inicio)
	if inicio < 0 || limite > len(g.MemoriaPrincipal) {
		logger.Error("El inicio es menor a cero o el limite excede el tamaño de la memoria principal")
		return fmt.Errorf("el formato de las direcciones a borrar son incorrectas %v", logger.ErrBadRequest)
	}

	g.MutexMemoriaPrincipal.Lock()
	copy(g.MemoriaPrincipal[inicio:limite], espacioVacio)
	g.MutexMemoriaPrincipal.Unlock()

	return nil
}
