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

var Instrucciones []string = []string{}

// función auxiliar para cargar el slice de instrucciones
func CargarListaDeInstrucciones(str string) {
	Instrucciones = append(Instrucciones, str)
	logger.Info("Se cargó una instrucción al Slice")
}

// ------------------------------------------------------------------
// ----------- FORMA PARTE DE LA MODIFICACIÓN DE PROCESOS -----------
// ------------------------------------------------------------------

func inicializacionProceso(w http.ResponseWriter, r *http.Request) {}

func CreacionProceso(w http.ResponseWriter, r *http.Request) {
	logger.Info(">>> Entró a utils.CreateProcess")
	var request globals.PedidoAMemoria
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Error decodificando el request", http.StatusBadRequest)
		logger.Error("Error decodificando request: %v", err)
		return
	}

	// Desempaquetar los argumentos
	var args globals.ArgmentosCreacionProceso
	argsBytes, _ := json.Marshal(request.Arguments)
	err = json.Unmarshal(argsBytes, &args)
	if err != nil {
		http.Error(w, "Error en los argumentos del proceso", http.StatusBadRequest)
		logger.Error("Error en los argumentos: %v", err)
		return
	}
	// Log para verificar lo recibido
	// log.Printf(">> [Memoria] Creando proceso: %s - Tamaño: %d", args.FileName, args.ProcessSize)
	logger.Info("## PID: <%d>  - Proceso Creado - Tamaño: <%d>")

	// se debe retornar el número de página de 1er nivel de ese proceso

	w.WriteHeader(http.StatusOK)
}

func FinalizacionProceso(w http.ResponseWriter, r *http.Request) {
	//toDO
	logger.Info("## PID: <PID>  - Proceso Destruido - Métricas - Acc.T.Pag: <ATP>; Inst.Sol.: <Inst.Sol>; SWAP: <SWAP>; Mem. Prin.: <Mem.Prin.>; Lec.Mem.: <Lec.Mem.>; Esc.Mem.: <Esc.Mem.>")
}

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
	//toDO

	var valorQueSeEncuentraLaDireccionPedida int = 0

	logger.Info("## PID: <PID>  - <Escritura> - Dir. Física: <DIRECCIÓN_FÍSICA> - Tamaño: <TAMAÑO>")

	return valorQueSeEncuentraLaDireccionPedida
}

func LecturaEspacio(w http.ResponseWriter, r *http.Request) {
	//toDO
	logger.Info("## PID: <PID>  - <Lectura> - Dir. Física: <DIRECCIÓN_FÍSICA> - Tamaño: <TAMAÑO>")
}

func MemoriaDump(w http.ResponseWriter, r *http.Request) {
	var dump globals.DatosParaDump

	if err := data.LeerJson(w, r, &dump); err != nil {
		log.Printf("Error al recibir JSON: %v", err)
		http.Error(w, "Error procesando datos del Kernel", http.StatusInternalServerError)
		return
	}

	globals.DatosDump = globals.DatosParaDump{
		PID:       dump.PID,
		TimeStamp: dump.TimeStamp,
	}

	dumpFileName := fmt.Sprintf(globals.MemoryConfig.DumpPath+"<%d>-<%s>.dmp", globals.DatosDump.PID, globals.DatosDump.TimeStamp)
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

// --------------------------------------------------------------------
// ---------- FORMA PARTE DEL ACCESO A LAS TABLAS DE PÁGINAS ----------
// --------------------------------------------------------------------

func AccesoTablaPaginas(w http.ResponseWriter, r *http.Request) int {

	//TODO

	esTablaIntermedia := false
	numeroTablaSgteNivel := 0
	esTablaUltNivel := false
	numeroFramePagina := 0

	if esTablaIntermedia {
		logger.Info("## Acceso a Tabla intermedia - Núm. Tabla Siguiente: <%d>", numeroTablaSgteNivel)
		return numeroTablaSgteNivel
	}
	if esTablaUltNivel {
		logger.Info("## Acceso a última Tabla - Núm. Frame: <%d>", numeroFramePagina)
		return numeroFramePagina
	}

	return (-1) // EN CASO DE ERROR
}
func LeerPaginaCompleta(w http.ResponseWriter, r *http.Request) {
	//toDO RETORNAR EL CONTENIDO DESDE LA PAGINA A PARTIR DEL BYTE ENVIADO DE DIRECC FIS. DE LA USER MEMORY
	//todo ESTE DEBERÁ COINCIDEIR CON LA POS DEL BYTE 0 DE LA PAGINA
	logger.Info("## Leer Página Completa - Dir. Física: <DIR>")
}
func ActualizarPaginaCompleta(w http.ResponseWriter, r *http.Request) bool {

	if err != nil {
		logger.Error("Error al actualizar la página - %s", err)
		return false
	}
	logger.Info("## PID: <PID> - Actualizar Página Completa - Dir. Física: <DIR>")
	return true
}
