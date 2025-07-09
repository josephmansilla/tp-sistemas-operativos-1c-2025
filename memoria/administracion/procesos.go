package administracion

import (
	"fmt"
	g "github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"sync"
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
	}
	nuevoProceso := &g.Proceso{
		PID:                  pid,
		TablaRaiz:            InicializarTablaRaiz(),
		Metricas:             InicializarMetricas(),
		InstruccionesEnBytes: make(map[int][]byte),
	}
	logger.Info("Proceso creado en memoria para PID <%d>", pid)

	g.MutexProcesosPorPID.Lock()
	g.ProcesosPorPID[pid] = nuevoProceso
	g.MutexProcesosPorPID.Unlock()
	logger.Info("Proceso PID <%d> agregado a la lista global", pid)

	//Guardar proceso: asi cualquier función que acceda por pid
	//(como ModificarEstadoEntradaEscritura o AsignarDatosAPaginacion)
	//encuentre la estructura ya guardada.

	if nuevoProceso.TablaRaiz == nil {
		logger.Error("TablaRaiz es nil para proceso PID <%d>", pid)
		return fmt.Errorf("Tabla raíz no inicializada")
	}
	logger.Info("TablaRaiz inicializada para PID=%d", pid)

	err = LecturaPseudocodigo(nuevoProceso, nombreArchPseudocodigo, tamanioProceso)
	if err != nil {
		logger.Error("Error al leer pseudocódigo para PID=%d: %v", pid, err)
		return fmt.Errorf("error al leer pseudocódigo: %v", logger.ErrBadRequest)
	}
	logger.Info("Pseudocódigo leído correctamente para PID <%d>, Longitud en bytes: <%d>", pid)

	err = AsignarDatosAPaginacion(nuevoProceso)
	err = AsginarEspacioDeSistema(nuevoProceso)
	if err != nil {
		logger.Error("Error asignando datos a paginación para PID <%d>: %v", pid, err)
		return fmt.Errorf("error asignando datos para el proceso: %v", logger.ErrInternalFailure)
	}
	logger.Info("Datos asignados correctamente para PID <%d>", pid)

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
	logger.Info("Se liberó todo para el PID: %d", pid)
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
		return vector, logger.ErrNoInstance
	}

	var entradas []int

	entradas = RecolectarEntradasProcesoDump(*proceso)

	tamanioPagina := g.MemoryConfig.PagSize
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

		g.MutexDump.Lock()
		vector[numeroFrame] = resul
		g.MutexDump.Unlock()
	}

	return
}
func RecolectarEntradasProcesoDump(proceso g.Proceso) (resultados []int) {
	cantidadEntradas := g.MemoryConfig.EntriesPerPage
	var wg sync.WaitGroup
	canal := make(chan int, cantidadEntradas)

	for _, subtabla := range proceso.TablaRaiz {
		wg.Add(1)
		go func(st *g.TablaPagina) {
			defer wg.Done()
			RecorrerTablaPaginaDeFormaConcurrenteDump(st, canal)
		}(subtabla)
	}

	go func() {
		wg.Wait()
		close(canal)
	}()

	for numeroFrame := range canal {
		resultados = append(resultados, numeroFrame)
	}

	return
}

func RecorrerTablaPaginaDeFormaConcurrenteDump(tabla *g.TablaPagina, canal chan int) {

	if tabla.Subtabla != nil {
		for _, subTabla := range tabla.Subtabla {
			RecorrerTablaPaginaDeFormaConcurrenteDump(subTabla, canal)
		}
		return
	}
	for i, entrada := range tabla.EntradasPaginas {
		if tabla.EntradasPaginas[i].EstaPresente {
			canal <- entrada.NumeroFrame
		}
	}
}

/*func RecorrerTablaPagina(tabla *g.TablaPagina, resultados *[]*g.EntradaDump) {

	if tabla.Subtabla != nil {
		for _, subTabla := range tabla.Subtabla {
			RecorrerTablaPagina(subTabla, resultados)
		}
		return
	}
	for i, entrada := range tabla.EntradasPaginas {
		if tabla.EntradasPaginas[i].EstaPresente {
			*resultados = append(*resultados, &g.EntradaDump{
				DireccionFisica: g.MemoryConfig.PagSize * entrada.NumeroFrame,
				NumeroFrame:     entrada.NumeroFrame,
			})
		}
	}
}

func DumpGlobal() (resultado string) {
	g.MutexProcesosPorPID.Lock()
	for pid := range g.ProcesosPorPID {
		g.MutexProcesosPorPID.Unlock()

		resultado += RealizarDumpMemoria(pid) + "\n"

		g.MutexProcesosPorPID.Lock()
	}
	g.MutexProcesosPorPID.Unlock()

	return
} */

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
