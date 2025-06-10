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
	cantidadFrames := CalcularTotalFrames()

	globals.MemoriaPrincipal = make([][]byte, cantidadFrames)
	globals.MutexMemoriaPrincipal.Lock()
	for i := range globals.MemoriaPrincipal {
		globals.MemoriaPrincipal[i] = make([]byte, globals.TamanioMaximoFrame)
	}
	globals.MutexMemoriaPrincipal.Unlock()
	globals.FramesLibres = make([]bool, cantidadFrames)
	ConfigurarFrames(cantidadFrames)

	logger.Info("Tamanio Memoria Principal de %d", globals.MemoryConfig.MemorySize)
	logger.Info("Memoria Principal Inicializada con %d frames de %d cada una.", cantidadFrames, globals.MemoryConfig.PagSize)
}
func CalcularTotalFrames() int {
	tamanioMemoriaPrincipal := globals.MemoryConfig.MemorySize
	tamanioPagina := globals.MemoryConfig.PagSize

	return tamanioMemoriaPrincipal / tamanioPagina
}
func ConfigurarFrames(cantidadFrames int) {
	globals.MutexEstructuraFramesLibres.Lock()
	for i := 0; i < cantidadFrames; i++ {
		globals.FramesLibres[i] = true
	}
	globals.MutexEstructuraFramesLibres.Unlock()
	logger.Info("Todos los frames están libres.")
}

func TieneTamanioNecesario(tamanioProceso int) bool {
	var framesNecesarios = float64(tamanioProceso) / float64(globals.TamanioMaximoFrame)

	globals.MutexCantidadFramesLibres.Lock()
	bool := framesNecesarios <= float64(globals.CantidadFramesLibres)
	globals.MutexCantidadFramesLibres.Unlock()
	return bool
}

func LecturaPseudocodigo(archivoPseudocodigo string) []byte {
	string := archivoPseudocodigo
	stringEnBytes := []byte(string)
	return stringEnBytes
}

func ObtenerDatosMemoria(numeroFrame int) globals.ExitoLecturaPagina {

	globals.MutexMemoriaPrincipal.Lock()
	pseudocodigoEnBytes := globals.MemoriaPrincipal[numeroFrame]
	globals.MutexMemoriaPrincipal.Unlock()

	pseudocodigoEnString := string(pseudocodigoEnBytes)

	datosLectura := globals.ExitoLecturaPagina{
		PseudoCodigo:    pseudocodigoEnString,
		DireccionFisica: numeroFrame,
	}

	return datosLectura
}

func MemoriaDump(w http.ResponseWriter, r *http.Request) {
	var dump globals.DatosParaDump

	if err := data.LeerJson(w, r, &dump); err != nil {
		logger.Error("Error al recibir JSON: %v", err)
		http.Error(w, "Error procesando datos del Kernel", http.StatusInternalServerError)
		return
	}

	dumpFileName := fmt.Sprintf("%s/<%d>-<%s>.dmp", globals.MemoryConfig.DumpPath, dump.PID, dump.TimeStamp)
	logger.Info("EL NOMBRE DEL DUMPFILE ES: " + dumpFileName)
	dumpFile, err := os.OpenFile(dumpFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Printf("Error al crear archivo de log para <%d-%s>: %v\n", dump.PID, dump.TimeStamp, err)
		os.Exit(1)
	}
	log.SetOutput(dumpFile)
	defer dumpFile.Close()

	logger.Info("## PID: <%d>  - Memory Dump solicitado", dump.PID) // se logea
	// TODO: se debe ubicar el proceso completo
	// InformarMetricasProceso(proceso.Metricas)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Dump Realizado"))
}
