package utils

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/data"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"log"
	"net/http"
	"os"
)

// TODO: EU = ESPACIO DE USUARIO

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
	var dump globals.DatosParaDump

	if err := data.LeerJson(w, r, &dump); err != nil {
		logger.Error("Error al recibir JSON: %v", err)
		http.Error(w, "Error procesando datos del Kernel", http.StatusInternalServerError)
		return
	}

	dumpFileName := fmt.Sprintf(globals.MemoryConfig.DumpPath+"<%d>-<%s>.dmp", dump.PID, dump.TimeStamp)
	logger.Info("EL NOMBRE DEL DUMPFILE ES: " + dumpFileName)
	dumpFile, err := os.OpenFile(dumpFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Printf("Error al crear archivo de log para <%d-%s>: %v\n", dump.PID, dump.TimeStamp, err)
		os.Exit(1)
	}
	log.SetOutput(dumpFile)

	// Llamado: "<PID>-<TIMESTAMP>.dmp" dentro del path definido por el archivo de configuración
	logger.Info("## PID: <%d>  - Memory Dump solicitado", dump.PID)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Dump Realizado"))
}
