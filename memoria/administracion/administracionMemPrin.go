package administracion

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/data"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"log"
	"net/http"
	"os"
)

func InicializarMemoriaPrincipal() {
	tamanioMemoriaPrincipal := globals.MemoryConfig.MemorySize
	tamanioPagina := globals.MemoryConfig.PagSize
	cantidadFrames := tamanioMemoriaPrincipal / tamanioPagina

	globals.MemoriaPrincipal = make([]byte, tamanioMemoriaPrincipal)

	globals.FramesLibres = make([]bool, cantidadFrames)
	ConfigurarFrames(cantidadFrames)

	logger.Info("Tamanio Memoria Principal de %d", globals.MemoryConfig.MemorySize)
	logger.Info("Memoria Principal Inicializada con %d con %d frames de %d.",
		tamanioMemoriaPrincipal, cantidadFrames, tamanioPagina)
}

func ConfigurarFrames(cantidadFrames int) { //TODO: OBSOLETO
	globals.MutexEstructuraFramesLibres.Lock()
	for i := 0; i < cantidadFrames; i++ {
		globals.FramesLibres[i] = true
	}
	globals.MutexEstructuraFramesLibres.Unlock()
	globals.CantidadFramesLibres = cantidadFrames
	logger.Info("Todos los frames están libres.")
}

func TieneTamanioNecesario(tamanioProceso int) bool {
	var framesNecesarios = float64(tamanioProceso) / float64(globals.MemoryConfig.PagSize)

	globals.MutexCantidadFramesLibres.Lock()
	bool := framesNecesarios <= float64(globals.CantidadFramesLibres)
	globals.MutexCantidadFramesLibres.Unlock()
	return bool
} //TODO: testear

func LecturaPseudocodigo(archivoPseudocodigo string) []byte {
	string := archivoPseudocodigo
	stringEnBytes := []byte(string)
	return stringEnBytes
} //TODO: testear

func ObtenerDatosMemoria(direccionFisica int) globals.ExitoLecturaPagina {
	tamanioPagina := globals.MemoryConfig.PagSize
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

	globals.MutexMemoriaPrincipal.Lock()
	copy(pseudocodigoEnBytes, globals.MemoriaPrincipal[direccionFisica:direccionFisica+bytesRestantes])
	globals.MutexMemoriaPrincipal.Unlock()

	pseudocodigoEnString := string(pseudocodigoEnBytes)

	datosLectura := globals.ExitoLecturaPagina{
		Exito:           true,
		PseudoCodigo:    pseudocodigoEnString,
		DireccionFisica: direccionFisica,
	}

	return datosLectura
}

func ModificarEstadoEntradaEscritura(direccionFisica int, pid int, datosEnBytes []byte) {
	tamanioPagina := globals.MemoryConfig.PagSize
	numeroPagina := direccionFisica / tamanioPagina

	inicioFrame := numeroPagina * tamanioPagina
	finFrame := inicioFrame + tamanioPagina

	if direccionFisica+len(datosEnBytes) > finFrame {
		logger.Fatal("¡¡¡¡¡¡¡¡¡¡Segment Fault!!!!!!!!!!!!")
		panic("Segment Fault - Lectura fuera del marco asignado")
	}

	globals.MutexMemoriaPrincipal.Lock()
	copy(globals.MemoriaPrincipal[direccionFisica:], datosEnBytes)
	globals.MutexMemoriaPrincipal.Unlock()

	globals.MutexProcesosPorPID.Lock()
	proceso := globals.ProcesosPorPID[pid]
	globals.MutexProcesosPorPID.Unlock()

	indices := CrearIndicePara(numeroPagina)
	entrada := BuscarEntradaPagina(proceso, indices)
	if entrada != nil {
		entrada.FueModificado = true
		entrada.EstaEnUso = true
	}

	IncrementarMetrica(proceso, IncrementarEscrituraDeMemoria)
}

func RealizarDumpMemoria(pid int) string {
	globals.MutexProcesosPorPID.Lock()
	proceso := globals.ProcesosPorPID[pid]
	globals.MutexProcesosPorPID.Unlock()
	if proceso == nil {
		logger.Fatal("No existe el proceso solicitado para DUMP")
		// TODO:
	}

	resultado := fmt.Sprintf("## Dump De Memoria Para PID: %d\n\n", pid)

	tamanioProceso := 10000 // tamanioMaximoProceso / globals.TamanioMaximoFrame
	for numeroPagina := 0; numeroPagina < tamanioProceso; numeroPagina++ {
		indices := CrearIndicePara(numeroPagina)
		entrada := BuscarEntradaPagina(proceso, indices)

		if entrada == nil || !entrada.EstaPresente {
			continue
		}
		frame := entrada.NumeroFrame
		globals.MutexMemoriaPrincipal.Lock()
		datos := globals.MemoriaPrincipal[frame]
		globals.MutexMemoriaPrincipal.Unlock()
		datosEnString := string(datos)
		resultado += fmt.Sprintf("Pagina: %d | Frame: %d | Datos: %s\n", numeroPagina, frame, datosEnString)
	}
	return resultado
}

func ParsearContenido(dumpFile *os.File, contenido string) {
	_, err := dumpFile.WriteString(contenido)
	if err != nil {
		logger.Error("Error al escribir contenido en el archivo dump: %v", err)
	}
} //TODO: rever

func MemoriaDump(w http.ResponseWriter, r *http.Request) {
	var dump globals.DatosParaDump

	if err := data.LeerJson(w, r, &dump); err != nil {
		logger.Error("Error al recibir JSON: %v", err)
		http.Error(w, "Error procesando datos del Kernel", http.StatusInternalServerError)
		return
	}

	dumpFileName := fmt.Sprintf("%s/<%d>-<%s>.dmp", globals.MemoryConfig.DumpPath, dump.PID, dump.TimeStamp)
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
