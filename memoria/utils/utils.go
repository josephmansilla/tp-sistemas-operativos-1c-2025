package utils

import (
	"encoding/json"
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"log"
	"net/http"
	"os"
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
