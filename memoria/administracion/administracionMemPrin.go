package administracion

import (
	"fmt"
	g "github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/data"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"log"
	"net/http"
	"os"
	"sort"
	"sync"
)

func InicializarMemoriaPrincipal() {
	tamanioMemoriaPrincipal := g.MemoryConfig.MemorySize
	tamanioPagina := g.MemoryConfig.PagSize
	cantidadFrames := tamanioMemoriaPrincipal / tamanioPagina

	g.MemoriaPrincipal = make([]byte, tamanioMemoriaPrincipal)

	g.FramesLibres = make([]bool, cantidadFrames)
	ConfigurarFrames(cantidadFrames)

	logger.Info("Tamanio Memoria Principal de %d", g.MemoryConfig.MemorySize)
	logger.Info("Memoria Principal Inicializada con %d con %d frames de %d.",
		tamanioMemoriaPrincipal, cantidadFrames, tamanioPagina)
}

func ConfigurarFrames(cantidadFrames int) { //TODO: MAS O MENOS OBSOLETO
	g.MutexEstructuraFramesLibres.Lock()
	for i := 0; i < cantidadFrames; i++ {
		g.FramesLibres[i] = true
	}
	g.MutexEstructuraFramesLibres.Unlock()
	g.CantidadFramesLibres = cantidadFrames
	logger.Info("Todos los frames están libres.")
}

func TieneTamanioNecesario(tamanioProceso int) (resultado bool) {
	var framesNecesarios = float64(tamanioProceso) / float64(g.MemoryConfig.PagSize)

	g.MutexCantidadFramesLibres.Lock()
	resultado = framesNecesarios <= float64(g.CantidadFramesLibres)
	g.MutexCantidadFramesLibres.Unlock()
	return
} //TODO: testear

func LecturaPseudocodigo(archivoPseudocodigo string) (stringEnBytes []byte, err error) {
	err = nil
	string := archivoPseudocodigo
	if string == "" {
		return nil, fmt.Errorf("El string es vacio: %w")
	}
	stringEnBytes = []byte(string)
	return
} //TODO: testear

func ObtenerDatosMemoria(direccionFisica int) (datosLectura g.ExitoLecturaPagina) {
	tamanioPagina := g.MemoryConfig.PagSize
	numeroPagina := direccionFisica / tamanioPagina
	offset := direccionFisica % tamanioPagina

	inicioFrame := numeroPagina * tamanioPagina
	finFrame := inicioFrame + tamanioPagina
	bytesRestantes := tamanioPagina - offset

	if direccionFisica+bytesRestantes > finFrame {
		logger.Fatal("¡¡¡¡¡¡¡¡¡¡Segment Fault!!!!!!!!!!!!")
		panic("Segment Fault - Lectura fuera del marco asignado")
	}

	pseudocodigoEnBytes := make([]byte, bytesRestantes)

	g.MutexMemoriaPrincipal.Lock()
	copy(pseudocodigoEnBytes, g.MemoriaPrincipal[direccionFisica:direccionFisica+bytesRestantes])
	g.MutexMemoriaPrincipal.Unlock()

	pseudocodigoEnString := string(pseudocodigoEnBytes)

	datosLectura = g.ExitoLecturaPagina{
		Exito:           nil,
		PseudoCodigo:    pseudocodigoEnString,
		DireccionFisica: direccionFisica,
	}

	return
}

func ModificarEstadoEntradaEscritura(direccionFisica int, pid int, datosEnBytes []byte) {
	tamanioPagina := g.MemoryConfig.PagSize
	numeroPagina := direccionFisica / tamanioPagina

	inicioFrame := numeroPagina * tamanioPagina
	finFrame := inicioFrame + tamanioPagina

	if direccionFisica+len(datosEnBytes) > finFrame {
		logger.Fatal("¡¡¡¡¡¡¡¡¡¡Segment Fault!!!!!!!!!!!!")
		panic("Segment Fault - Lectura fuera del marco asignado")
	}

	g.MutexMemoriaPrincipal.Lock()
	copy(g.MemoriaPrincipal[direccionFisica:], datosEnBytes)
	g.MutexMemoriaPrincipal.Unlock()

	g.MutexProcesosPorPID.Lock()
	proceso := g.ProcesosPorPID[pid]
	g.MutexProcesosPorPID.Unlock()

	indices := CrearIndicePara(numeroPagina)
	entrada, err := BuscarEntradaPagina(proceso, indices)
	if err != nil {
		logger.Error("No se pudo encontrar la entrada de pagina")
		panic("AAAAAAAAAAAAAAAAAAAAAAAAA") // TODO: ver que hacer con este error
	}
	if entrada != nil {
		entrada.FueModificado = true
		entrada.EstaEnUso = true
	}

	IncrementarMetrica(proceso, IncrementarEscrituraDeMemoria)
}

func RealizarDumpMemoria(pid int) (resultado string) {
	g.MutexProcesosPorPID.Lock()
	proceso := g.ProcesosPorPID[pid]
	g.MutexProcesosPorPID.Unlock()

	if proceso == nil {
		logger.Fatal("No existe el proceso solicitado para DUMP")
		return "Proceso no encontrado 😭🙏"
		// TODO:
	}

	resultado = fmt.Sprintf("## Dump De Memoria Para PID: %d\n\n", pid)

	var entradas []g.EntradaDump

	entradas = RecolectarEntradasProceso(*proceso)

	tamanioPagina := g.MemoryConfig.PagSize
	for _, e := range entradas {
		inicio := e.NumeroFrame * tamanioPagina
		fin := inicio + tamanioPagina

		if fin > len(g.MemoriaPrincipal) {
			logger.Error("Acceso fuera de rango al hacer dump del frame %d con PID: %d", e.NumeroFrame, pid)
			continue
			// TODO: ver que hacer
		}

		g.MutexMemoriaPrincipal.Lock()
		datos := g.MemoriaPrincipal[inicio:fin]
		g.MutexMemoriaPrincipal.Unlock()

		datosEnString := string(datos)
		resultado += fmt.Sprintf("Direccion Fisica: %d | Frame: %d | Datos: %q\n", e.DireccionFisica, e.NumeroFrame, datosEnString)
	}

	return
}

func RecolectarEntradasProceso(proceso g.Proceso) (resultados []g.EntradaDump) {
	cantidadEntradas := g.MemoryConfig.EntriesPerPage
	var wg sync.WaitGroup
	canal := make(chan g.EntradaDump, cantidadEntradas)

	for _, subtabla := range proceso.TablaRaiz {
		wg.Add(1)
		go func(st *g.TablaPagina) {
			defer wg.Done()
			RecorrerTablaPaginaDeFormaConcurrente(st, canal)
		}(subtabla)
	}

	go func() {
		wg.Wait()
		close(canal)
	}()

	for entrada := range canal {
		resultados = append(resultados, entrada)
	}
	// TODO: NO ES NECESARIO Y LO PUEDO BORRAR QUEDA PENDIENTE DEJARLO O NO
	sort.Slice(resultados, func(i, j int) bool {
		return resultados[i].DireccionFisica < resultados[j].DireccionFisica
	})

	return
}

func RecorrerTablaPaginaDeFormaConcurrente(tabla *g.TablaPagina, canal chan g.EntradaDump) {

	if tabla.Subtabla != nil {
		for _, subTabla := range tabla.Subtabla {
			RecorrerTablaPaginaDeFormaConcurrente(subTabla, canal)
		}
		return
	}
	for i, entrada := range tabla.EntradasPaginas {
		if tabla.EntradasPaginas[i].EstaPresente {
			canal <- g.EntradaDump{
				DireccionFisica: g.MemoryConfig.PagSize * entrada.NumeroFrame,
				NumeroFrame:     entrada.NumeroFrame,
			}
		}
	}
}

func RecorrerTablaPagina(tabla *g.TablaPagina, resultados *[]g.EntradaDump) {

	if tabla.Subtabla != nil {
		for _, subTabla := range tabla.Subtabla {
			RecorrerTablaPagina(subTabla, resultados)
		}
		return
	}
	for i, entrada := range tabla.EntradasPaginas {
		if tabla.EntradasPaginas[i].EstaPresente {
			*resultados = append(*resultados, g.EntradaDump{
				DireccionFisica: g.MemoryConfig.PagSize * entrada.NumeroFrame,
				NumeroFrame:     entrada.NumeroFrame,
			})
		}
	}
}

// TODO: para probar
func DumpGlobal() (resultado string) {
	g.MutexProcesosPorPID.Lock()
	for pid := range g.ProcesosPorPID {
		g.MutexProcesosPorPID.Unlock()

		resultado += RealizarDumpMemoria(pid) + "\n"

		g.MutexProcesosPorPID.Lock()
	}
	g.MutexProcesosPorPID.Unlock()

	return
}

func ParsearContenido(dumpFile *os.File, contenido string) {
	_, err := dumpFile.WriteString(contenido)
	if err != nil {
		logger.Error("Error al escribir contenido en el archivo dump: %v", err)
	}
} //TODO: rever

func MemoriaDump(w http.ResponseWriter, r *http.Request) {
	var dump g.DatosParaDump

	if err := data.LeerJson(w, r, &dump); err != nil {
		logger.Error("Error al recibir JSON: %v", err)
		http.Error(w, "Error procesando datos del Kernel", http.StatusInternalServerError)
		return
	}

	dumpFileName := fmt.Sprintf("%s/<%d>-<%s>.dmp", g.MemoryConfig.DumpPath, dump.PID, dump.TimeStamp)
	logger.Info("## Se creo el file: %d ", dumpFileName)
	dumpFile, err := os.OpenFile(dumpFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Printf("Error al crear archivo de log para <%d-%s>: %v\n", dump.PID, dump.TimeStamp, err)
		os.Exit(1)
	}
	log.SetOutput(dumpFile)
	defer dumpFile.Close()

	logger.Info("## PID: <%d>  - Memory Dump solicitado", dump.PID)

	contenido := RealizarDumpMemoria(dump.PID)
	// TODO: verificacion esta vacio
	ParsearContenido(dumpFile, contenido)

	logger.Info("## Archivo Dump fue creado con EXITO")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Dump Realizado"))
}
