package administracion

import (
	"bufio"
	"fmt"
	g "github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"os"
	"strings"
)

func InicializarMemoriaPrincipal() {
	tamanioMemoriaPrincipal := g.MemoryConfig.MemorySize
	tamanioPagina := g.MemoryConfig.PagSize
	cantidadFrames := tamanioMemoriaPrincipal / tamanioPagina

	g.MemoriaPrincipal = make([]byte, tamanioMemoriaPrincipal)
	ConfigurarFrames(cantidadFrames)
	g.InstanciarEstructurasGlobales()
	g.InicializarSemaforos()

	logger.Info("Tamanio Memoria Principal de %d", g.MemoryConfig.MemorySize)
	logger.Info("Memoria Principal Inicializada con %d con %d frames de %d.",
		tamanioMemoriaPrincipal, cantidadFrames, tamanioPagina)
}

func ConfigurarFrames(cantidadFrames int) {
	g.FramesLibres = make([]bool, cantidadFrames)
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

func LecturaPseudocodigo(proceso *g.Proceso, direccionPseudocodigo string, tamanioMaximo int) ([]byte, error) {
	if direccionPseudocodigo == "" {
		return nil, fmt.Errorf("el string es vacio")
	}
	ruta := "../pruebas/" + direccionPseudocodigo
	file, err := os.Open(ruta)
	if err != nil {
		logger.Error("Error al abrir el archivo: %s\n", err)
		return nil, err
	}
	defer file.Close()

	logger.Info("Se leyó el archivo")
	scanner := bufio.NewScanner(file)

	stringEnBytes := make([]byte, 0, tamanioMaximo)
	cantidadInstrucciones := 0

	for scanner.Scan() {
		lineaEnString := scanner.Text()
		logger.Info("Línea leída: %s", lineaEnString)
		lineaEnBytes := []byte(lineaEnString)

		stringEnBytes = append(stringEnBytes, lineaEnBytes...)
		proceso.OffsetInstrucciones[cantidadInstrucciones] = len(stringEnBytes)
		cantidadInstrucciones++
		// TODO: si los tests cuentan al EOF como instruccion queda así
		// TODO: sino despues del if

		if strings.TrimSpace(lineaEnString) == "EOF" {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		logger.Error("Error al leer el archivo: %s", err)
	}

	IncrementarMetrica(proceso, cantidadInstrucciones, IncrementarInstruccionesSolicitadas)

	logger.Info("Total de instrucciones cargadas para PID <%d>: %d", proceso.PID, cantidadInstrucciones)

	return stringEnBytes, nil
}

func ObtenerDatosMemoria(direccionFisica int) (datosLectura g.ExitoLecturaPagina) {
	tamanioPagina := g.MemoryConfig.PagSize
	numeroPagina := direccionFisica / tamanioPagina
	offset := direccionFisica % tamanioPagina

	inicioFrame := numeroPagina * tamanioPagina
	finFrame := inicioFrame + tamanioPagina
	bytesRestantes := tamanioPagina - offset

	if direccionFisica+bytesRestantes > finFrame {
		logger.Error("Se está leyendo afuera del frame")
		// TODO:		panic("Segment Fault - Lectura fuera del marco asignado")
		// TODO: tirar error pero sin panic porque no es un caso en los tests
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

func ModificarEstadoEntradaEscritura(direccionFisica int, pid int, datosEnBytes []byte) (err error) {
	tamanioPagina := g.MemoryConfig.PagSize
	numeroPagina := direccionFisica / tamanioPagina

	inicioFrame := numeroPagina * tamanioPagina
	finFrame := inicioFrame + tamanioPagina

	if direccionFisica+len(datosEnBytes) > finFrame {
		logger.Error("Segment Fault - Escritura fuera del marco asignado")
		return logger.ErrSegmentFault
	}

	g.MutexMemoriaPrincipal.Lock()
	copy(g.MemoriaPrincipal[direccionFisica:], datosEnBytes)
	g.MutexMemoriaPrincipal.Unlock()

	g.MutexProcesosPorPID.Lock()
	proceso := g.ProcesosPorPID[pid]
	g.MutexProcesosPorPID.Unlock()

	if proceso == nil {
		logger.Error("Se intentó acceder a un proceso inexistente o nil para PID=%d", pid)
		return fmt.Errorf("proceso nil para PID=%d", pid)
	}

	indices := CrearIndicePara(numeroPagina)
	entrada, err := BuscarEntradaPagina(proceso, indices)
	if err != nil {
		logger.Error("No se pudo encontrar la entrada de pagina: %v", err)
		return err
	}
	if entrada != nil {
		entrada.FueModificado = true
		entrada.EstaEnUso = true
	}

	IncrementarMetrica(proceso, 1, IncrementarEscrituraDeMemoria)

	return nil
}

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

func SeleccionarEntradas(pid int, direccionFisica int, entradasNecesarias int) (entradas []g.EntradaPagina, err error) {
	tamanioPagina := g.MemoryConfig.PagSize
	paginaInicio := direccionFisica / tamanioPagina
	err = nil

	g.MutexProcesosPorPID.Lock()
	proceso := g.ProcesosPorPID[pid]
	g.MutexProcesosPorPID.Unlock()

	if proceso == nil {
		return nil, fmt.Errorf("no existe el proceso con el PID: %d ; %v", pid, logger.ErrNoInstance)
	}

	for i := 0; i < entradasNecesarias; i++ {
		numeroPagina := paginaInicio + i
		indices := CrearIndicePara(numeroPagina)

		entrada, err := BuscarEntradaPagina(proceso, indices)
		if err != nil {
			return nil, fmt.Errorf("error al buscar la entrada de pagina: %d de PID %d; %v", numeroPagina, pid, err)
		}
		if entrada == nil {
			return nil, fmt.Errorf("error al encontrar la entrada de pagina: %d de PID %d; %v", numeroPagina, pid, err)
		}
		if !entrada.EstaPresente {
			// TODO: BUSCAR DE SWAPPPP ====================================================
		}
		entradas = append(entradas, *entrada)
	}

	return
} //TODO: rever no se usa el tamanioALeer

func LeerEspacioMemoria(pid int, direccionFisica int, tamanioALeer int) (confirmacionLectura g.ExitoLecturaMemoria, err error) {
	confirmacionLectura = g.ExitoLecturaMemoria{Exito: err, DatosAEnviar: ""}

	entradasNecesarias, err := g.CalcularCantidadEntradasATraer(tamanioALeer)
	if err != nil {
		return confirmacionLectura, err
	}

	entradas, err := SeleccionarEntradas(pid, direccionFisica, entradasNecesarias)
	if err != nil {
		return confirmacionLectura, err
	}

	bytesRestantes := tamanioALeer
	cant := len(entradas)
	datos := make([]byte, tamanioALeer)

	for i, entrada := range entradas {
		inicioLectura, finLectura, err := LogicaRecorrerMemoria(i, cant, entrada, direccionFisica, bytesRestantes)
		if err != nil {
			return confirmacionLectura, err
		}

		g.MutexMemoriaPrincipal.Lock()
		datos = append(datos, g.MemoriaPrincipal[inicioLectura:finLectura]...)
		g.MutexMemoriaPrincipal.Unlock()

		bytesRestantes -= finLectura - inicioLectura
		if bytesRestantes <= 0 {
			break
		}
	}
	return g.ExitoLecturaMemoria{Exito: nil, DatosAEnviar: string(datos)}, nil
}

func LogicaRecorrerMemoria(i int, cantEntradas int, entrada g.EntradaPagina, dirF int, bytesRestantes int) (inicio int, limite int, err error) {
	tamTotal := g.MemoryConfig.MemorySize
	tamPag := g.MemoryConfig.PagSize
	offsetLogico := dirF % tamPag

	base := entrada.NumeroFrame * tamPag

	if i == 0 { // PARA EL PRIMER FRAME
		var delta int
		if bytesRestantes <= tamPag { // EN CASO DE QUE LOS BYTESRESTANTES SEAN MENORES A LA ENTRADA
			delta = bytesRestantes
		} else {
			delta = tamPag
		}
		inicio = base + offsetLogico // BASE DE LA ENTRADA + POSIBLE DESPLAZAMIENTO
		limite = min(base+tamPag, inicio+delta)
		// ELIJO ENTRE EL LIMITE DE LA ENTRADA O EN CASO DE QUE SE CUMPLA EL PRIMER IF: HASTA UN LIMITE MENOR
		// QUE EL LIMITE DE LA ENTRADA
	} else if i == cantEntradas-1 { // PARA EL ÙLTIMO FRAME
		inicio = base
		limite = inicio + bytesRestantes
	} else { // PARA LOS CASOS INTERMEDIOS: NO DEBERÌAN ENTRAR SI LOS BYTES RESANTE ESTÀN CORRECTAMENTE CALCULADOS
		inicio = base
		limite = base + tamPag
	}

	if inicio >= tamTotal || limite > tamTotal {
		err = logger.ErrSegmentFault
		return 0, 0, err
	}

	return inicio, limite, nil
}

func EscribirEspacioMemoria(pid int, direccionFisica int, tamanioALeer int, datosAEscribir string) (confirmacionEscritura g.ExitoEdicionMemoria, err error) {
	confirmacionEscritura = g.ExitoEdicionMemoria{Exito: err, Booleano: false}

	entradasNecesarias, err := g.CalcularCantidadEntradasATraer(tamanioALeer)
	if err != nil {
		return confirmacionEscritura, err
	}

	entradas, err := SeleccionarEntradas(pid, direccionFisica, entradasNecesarias)
	if err != nil {
		return confirmacionEscritura, err
	}

	direccionParaAcceder := direccionFisica
	datosEnBytes := []byte(datosAEscribir)
	cant := len(entradas)
	bytesRestantes := tamanioALeer
	bytesEscritos := 0

	for i, entrada := range entradas {
		inicioEscritura, finEscritura, err := LogicaRecorrerMemoria(i, cant, entrada, direccionParaAcceder, bytesRestantes)
		cantBytes := finEscritura - inicioEscritura
		if err != nil {
			return confirmacionEscritura, err
		}

		inicioDatos := bytesEscritos
		finDatos := inicioDatos + cantBytes

		g.MutexMemoriaPrincipal.Lock()
		copy(g.MemoriaPrincipal[inicioEscritura:finEscritura], datosEnBytes[inicioDatos:finDatos])
		g.MutexMemoriaPrincipal.Unlock()

		entrada.FueModificado = true

		bytesEscritos += cantBytes
		bytesRestantes -= cantBytes
		if bytesRestantes <= 0 {
			break
		}
		direccionParaAcceder += cantBytes
	}
	return g.ExitoEdicionMemoria{Exito: err, Booleano: true}, nil

}
