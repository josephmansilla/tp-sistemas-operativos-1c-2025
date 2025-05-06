package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/sisoputnfrba/tp-golang/memoria/globals"
)

func Config(filepath string) *globals.Config {
	var config *globals.Config
	configFile, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}

var Instrucciones []string = []string{}

// función auxiliar para cargar el slice de instrucciones
func CargarListaDeInstrucciones(str string) {
	Instrucciones = append(Instrucciones, str)
	log.Print("Se cargó una instrucción al Slice\n")
}

// ------------------------------------------------------------------
// ----------- FORMA PARTE DE LA MODIFICACIÓN DE PROCESOS -----------
// ------------------------------------------------------------------

func CreacionProceso(tamanioDeseado int) {
	PID := 1 // TEMPORAL HASTA QUE SEPAMOS COMO ASIGNAR PID'S
	log.Printf("## PID: <%d>  - Proceso Creado - Tamaño: <%d>", PID, tamanioDeseado)

	// se debe retornar el número de página de 1er nivel de ese proceso
}
func DestruccionProceso(w http.ResponseWriter, r *http.Request) {
	//toDO
	log.Printf("## PID: <PID>  - Proceso Destruido - Métricas - Acc.T.Pag: <ATP>; Inst.Sol.: <Inst.Sol>; SWAP: <SWAP>; Mem. Prin.: <Mem.Prin.>; Lec.Mem.: <Lec.Mem.>; Esc.Mem.: <Esc.Mem.>")
}
func FinalizacionProceso(w http.ResponseWriter, r *http.Request) {
	// toDO
}

// ------------------------------------------------------------------
// ---------- FORMA PARTE DEL ACCESO A ESPACIO DE USUARIO ----------
// ------------------------------------------------------------------

func ObtenerEspacioLibreMock(w http.ResponseWriter, r *http.Request) {
	respuesta := globals.EspacioLibreRTA{EspacioLibre: globals.MemoryConfig.MemorySize}
	// el EspacioLibre es un valor arbitrario
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(respuesta); err != nil {
		log.Printf("Error al serializar mock de espacio: %v", err)
	}

	log.Printf("## Espacio libre mock devuelto - Tamaño: <%d>\n", respuesta.EspacioLibre)
}

func EscrituraEspacio(w http.ResponseWriter, r *http.Request) {
	//toDO
	log.Printf("## PID: <PID>  - <Escritura> - Dir. Física: <DIRECCIÓN_FÍSICA> - Tamaño: <TAMAÑO>")
}

func LecturaEspacio(w http.ResponseWriter, r *http.Request) {
	//toDO
	log.Printf("## PID: <PID>  - <Lectura> - Dir. Física: <DIRECCIÓN_FÍSICA> - Tamaño: <TAMAÑO>")
}

func MemoryDump(w http.ResponseWriter, r *http.Request) {
	// toDO
	// Llamado: "<PID>-<TIMESTAMP>.dmp" dentro del path definido por el archivo de configuración
	log.Printf("## PID: <PID>  - Memory Dump solicitado")
}

// --------------------------------------------------------------------
// ---------- FORMA PARTE DEL ACCESO A LAS TABLAS DE PÁGINAS ----------
// --------------------------------------------------------------------

func AccesoTablaPaginas(w http.ResponseWriter, r *http.Request) {
	//toDO
}
func LeerPaginaCompleta(w http.ResponseWriter, r *http.Request) {
	//toDO
}
func ActualizarPaginaCompleta(w http.ResponseWriter, r *http.Request) {
	//toDO
}

// TODO

// Para recibir Creacion de proceso
func CreateProcess(w http.ResponseWriter, r *http.Request) {
	log.Println(">>> Entró a utils.CreateProcess")
	var request RequestToMemory
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Error decodificando el request", http.StatusBadRequest)
		log.Printf("Error decodificando request: %v", err)
		return
	}

	// Desempaquetar los argumentos
	var args CreateProcessArgs
	argsBytes, _ := json.Marshal(request.Arguments)
	err = json.Unmarshal(argsBytes, &args)
	if err != nil {
		http.Error(w, "Error en los argumentos del proceso", http.StatusBadRequest)
		log.Printf("Error en los argumentos: %v", err)
		return
	}
	// Log para verificar lo recibido
	log.Printf(">> [Memoria] Creando proceso: %s - Tamaño: %d", args.FileName, args.ProcessSize)

	// ACA VA LA DE CARGAR INSTRUCCIONES DADO EL NOMBRE DE PSUDO CODIGOañadir la lógica para manejar el proceso en memoria

	w.WriteHeader(http.StatusOK)
}

type CreateProcessArgs struct {
	FileName    string `json:"fileName"`
	ProcessSize int    `json:"processSize"`
}

type RequestToMemory struct {
	Thread    Thread                 `json:"thread"`
	Type      string                 `json:"type"`
	Arguments map[string]interface{} `json:"arguments"`
}

type Thread struct {
	PID int `json:"pid"`
}
