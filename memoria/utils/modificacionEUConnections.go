package utils

import (
	"encoding/json"
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

func ObtenerEspacioLibreMock(w http.ResponseWriter, r *http.Request) {
	respuesta := globals.EspacioLibreRTA{EspacioLibre: globals.MemoryConfig.MemorySize}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(respuesta); err != nil {
		logger.Error("Error al serializar mock de espacio: %v", err)
	}

	logger.Info("## Espacio libre mock devuelto - Tamaño: <%d>", respuesta.EspacioLibre)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ESPACIO DEVUELTO"))
}

func EscrituraEspacio(w http.ResponseWriter, r *http.Request) int {
	// TODO: ESCRIBIR LO INDICADO EN LA DIRECCION PEDIDA

	var valorQueSeEncuentraLaDireccionPedida int = 0

	logger.Info("## PID: <PID>  - <Escritura> - Dir. Física: <DIRECCIÓN_FÍSICA> - Tamaño: <TAMAÑO>")

	// TODO: RESPONDER CON OK
	return valorQueSeEncuentraLaDireccionPedida
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

	globals.DatosDump = globals.DatosParaDump{
		PID:       dump.PID,
		TimeStamp: dump.TimeStamp,
	}

	dumpFileName := fmt.Sprintf(globals.MemoryConfig.DumpPath+"<%d>-<%s>.dmp", globals.DatosDump.PID, globals.DatosDump.TimeStamp)
	logger.Info("EL NOMBRE DEL DUMPFILE ES: " + dumpFileName)
	dumpFile, err := os.OpenFile(dumpFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Printf("Error al crear archivo de log para <%d-%s>: %v\n", globals.DatosDump.PID, globals.DatosDump.TimeStamp, err)
		os.Exit(1)
	}
	log.SetOutput(dumpFile)

	// Llamado: "<PID>-<TIMESTAMP>.dmp" dentro del path definido por el archivo de configuración
	logger.Info("## PID: <%d>  - Memory Dump solicitado", globals.DatosDump.PID)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Dump Realizado"))
}
