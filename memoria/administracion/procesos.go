package administracion

import (
	"fmt"
	g "github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
)

func InicializarProceso(pid int, tamanioProceso int, nombreArchPseudocodigo string) (err error) {
	logger.Info("Inicializando proceso PID <%d>, Tamaño: <%d>, Pseudocódigo <%s>", pid, tamanioProceso, nombreArchPseudocodigo)

	if !TieneTamanioNecesario(tamanioProceso) {
		logger.Error("No hay memoria suficiente para proceso PID <%d>", pid)
		return fmt.Errorf("no hay memoria disponible para el proceso: %v", logger.ErrNoMemory)
	}

	g.MutexProcesosPorPID.Lock()
	if g.ProcesosPorPID[pid] != nil {
		g.MutexProcesosPorPID.Unlock()
		logger.Error("El proceso PID <%d> ya existe", pid)
		return fmt.Errorf("el proceso PID <%d> ya existe", pid)
	} else {
		g.MutexProcesosPorPID.Unlock()
	}
	nuevoProceso := &g.Proceso{
		PID:                  pid,
		TablaRaiz:            InicializarTablaRaiz(),
		Metricas:             InicializarMetricas(),
		InstruccionesEnBytes: make(map[int][]byte),
	}
	logger.Info("## Proceso creado en memoria para PID <%d>", pid)

	g.MutexProcesosPorPID.Lock()
	g.ProcesosPorPID[pid] = nuevoProceso
	g.MutexProcesosPorPID.Unlock()

	if nuevoProceso.TablaRaiz == nil {
		logger.Error("TablaRaiz es nil para proceso PID <%d>", pid)
		return logger.ErrNoTabla
	}

	err = LecturaPseudocodigo(nuevoProceso, nombreArchPseudocodigo)
	if err != nil {
		return fmt.Errorf("error al leer pseudocódigo: %v", logger.ErrBadRequest)
	}

	err = AsignarPaginasParaPID(nuevoProceso, tamanioProceso)

	logger.Info("## Datos asignados correctamente para PID <%d>", pid)

	return nil
}

func LiberarMemoriaProceso(pid int) (metricas g.MetricasProceso, err error) {
	var proceso *g.Proceso
	metricas = g.MetricasProceso{}
	err = nil

	proceso, err = DesocuparProcesoEnVectorMapeable(pid)
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

func DesocuparProcesoEnVectorMapeable(pid int) (proceso *g.Proceso, err error) {
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

func RealizarDumpMemoria(pid int) (vector []string, err error) {
	g.MutexProcesosPorPID.Lock()
	proceso := g.ProcesosPorPID[pid]
	g.MutexProcesosPorPID.Unlock()

	if proceso == nil {
		logger.Error("No existe el proceso solicitado para DUMP")
		return vector, logger.ErrProcessNil
	}

	entradas := RecolectarEntradasProcesoDump(*proceso)
	tamanioPagina := g.MemoryConfig.PagSize

	vector = make([]string, len(g.FramesLibres))

	for i := 0; i < len(entradas); i++ {
		numeroFrame := entradas[i]
		inicio := numeroFrame * tamanioPagina
		fin := inicio + tamanioPagina

		if fin > len(g.MemoriaPrincipal) {
			logger.Error("Acceso fuera de rango al hacer dump del frame %d con PID: %d", numeroFrame, pid)
			fin = len(g.MemoriaPrincipal) - 1
			continue
		}

		g.MutexMemoriaPrincipal.Lock()
		datos := g.MemoriaPrincipal[inicio:fin]
		g.MutexMemoriaPrincipal.Unlock()

		datosEnString := string(datos)
		resul := fmt.Sprintf("Direccion Fisica: %d | Frame: %d | Datos: %q\n", inicio, numeroFrame, datosEnString)

		vector[numeroFrame] = resul
	}

	return
}
func RecolectarEntradasProcesoDump(proceso g.Proceso) (resultados []int) {
	for _, subtabla := range proceso.TablaRaiz {
		RecorrerTablaPagina(subtabla, &resultados)
	}
	return
}

func RecorrerTablaPagina(tabla *g.TablaPagina, resultados *[]int) {

	if tabla.Subtabla != nil {
		for _, subTabla := range tabla.Subtabla {
			RecorrerTablaPagina(subTabla, resultados)
		}
		return
	}
	for _, entrada := range tabla.EntradasPaginas {
		if entrada.EstaPresente {
			*resultados = append(*resultados, entrada.NumeroFrame)
		}
	}
}

func InicializarMetricas() (metricas g.MetricasProceso) {
	metricas = g.MetricasProceso{
		AccesosTablasPaginas:     0,
		InstruccionesSolicitadas: 0,
		BajadasSwap:              0,
		SubidasMP:                0,
		LecturasDeMemoria:        0,
		EscriturasDeMemoria:      0,
	}
	return
}

func IncrementarMetrica(proceso *g.Proceso, cantidad int, funcMetrica g.OperacionMetrica) {
	g.MutexMetrica[proceso.PID].Lock()
	funcMetrica(&proceso.Metricas, cantidad)
	g.MutexMetrica[proceso.PID].Unlock()
}

func IncrementarAccesosTablasPaginas(metrica *g.MetricasProceso, cantidad int) {
	metrica.AccesosTablasPaginas += cantidad
}
func IncrementarInstruccionesSolicitadas(metrica *g.MetricasProceso, cantidad int) {
	metrica.InstruccionesSolicitadas += cantidad
}
func IncrementarBajadasSwap(metrica *g.MetricasProceso, cantidad int) {
	metrica.BajadasSwap += cantidad
}
func IncrementarSubidasMP(metrica *g.MetricasProceso, cantidad int) {
	metrica.SubidasMP += cantidad
}
func IncrementarLecturaDeMemoria(metrica *g.MetricasProceso, cantidad int) {
	metrica.LecturasDeMemoria += cantidad
}
func IncrementarEscrituraDeMemoria(metrica *g.MetricasProceso, cantidad int) {
	metrica.EscriturasDeMemoria += cantidad
}
