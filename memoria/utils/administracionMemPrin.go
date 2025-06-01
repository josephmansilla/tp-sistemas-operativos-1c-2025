package utils

import (
	"fmt"
	globalData "github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/data"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"log"
	"net/http"
	"os"
)

// TODO: mutex para crear la MP
func InicializarMemoriaPrincipal() {
	cantidadFrames := CalcularCantidadFrames()

	globalData.MemoriaPrincipal = make([]byte, cantidadFrames)
	globalData.FramesLibres = make([]bool, cantidadFrames)
	ConfigurarFrames(cantidadFrames)

	logger.Info("Tamanio Memoria Principal de %d", globalData.MemoryConfig.MemorySize)
	logger.Info("Memoria Principal Inicializada con %d frames de %d cada una.", cantidadFrames, globalData.MemoryConfig.PagSize)
}

func CalcularCantidadFrames() int {
	tamanioMemoriaPrincipal := globalData.MemoryConfig.MemorySize
	tamanioPagina := globalData.MemoryConfig.PagSize

	return tamanioMemoriaPrincipal / tamanioPagina
}

func ConfigurarFrames(cantidadFrames int) {
	for i := 0; i <= cantidadFrames; i++ {
		globalData.FramesLibres[i] = true
	}
	logger.Info("Todos los frames están libres.")
}

func AsignarProceso(PID int, cantidadPaginas int) {

}

// ------------------------------------------------------------------
// ---------- FORMA PARTE DEL ACCESO A ESPACIO DE USUARIO ----------
// ------------------------------------------------------------------

func EscrituraEspacio(w http.ResponseWriter, r *http.Request) {
	// TODO: ESCRIBIR LO INDICADO EN LA DIRECCION PEDIDA

	logger.Info("## PID: <PID>  - <Escritura> - Dir. Física: <DIRECCIÓN_FÍSICA> - Tamaño: <TAMAÑO>")

	// TODO: RESPONDER CON OK
}

func LecturaEspacio(w http.ResponseWriter, r *http.Request) {

	// TODO: DEBO CREAR UNA STRUCT PARA QUE ME ENVIEN LA DIRECCION FISICA
	// TODO: DEVOLVER EL VALOR QUE SE ENCUENTRA EN LA DIRECCION PEDIDA
	logger.Info("## PID: <PID>  - <Lectura> - Dir. Física: <DIRECCIÓN_FÍSICA> - Tamaño: <TAMAÑO>")
}

func MemoriaDump(w http.ResponseWriter, r *http.Request) {
	var dump globalData.DatosParaDump

	if err := data.LeerJson(w, r, &dump); err != nil {
		logger.Error("Error al recibir JSON: %v", err)
		http.Error(w, "Error procesando datos del Kernel", http.StatusInternalServerError)
		return
	} // err handling

	dumpFileName := fmt.Sprintf("%s/<%d>-<%s>.dmp", globalData.MemoryConfig.DumpPath, dump.PID, dump.TimeStamp)
	logger.Info("EL NOMBRE DEL DUMPFILE ES: " + dumpFileName)
	dumpFile, err := os.OpenFile(dumpFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Printf("Error al crear archivo de log para <%d-%s>: %v\n", dump.PID, dump.TimeStamp, err)
		os.Exit(1)
	} // err handling y se asigna el nombre del dumpfile
	log.SetOutput(dumpFile)
	defer dumpFile.Close()

	logger.Info("## PID: <%d>  - Memory Dump solicitado", dump.PID) // se logea

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Dump Realizado"))
}
