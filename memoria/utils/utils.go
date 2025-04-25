package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/data"
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

// FUNCION PARA RECIBIR LOS MENSAJES PROVENIENTES DE LA CPU
func RecibirMensajeDeCPU(w http.ResponseWriter, r *http.Request) {
	var mensaje globals.DatosDeCPU
	data.LeerJson(w, r, &mensaje)

	globals.CPU = globals.DatosDeCPU{
		PID: mensaje.PID,
		PC:  mensaje.PC,
	}

	log.Printf("PID Pedido: %d\n", mensaje.PID)
	log.Printf("PC Pedido: %d\n", mensaje.PC)
}

// FUNCION PARA RETORNAR LOS MENSAJES PROVENIENTES DE LA CPU
func RetornarMensajeDeCPU(w http.ResponseWriter, r *http.Request) globals.DatosDeCPU {
	var mensaje globals.DatosDeCPU
	data.LeerJson(w, r, &mensaje)

	globals.CPU = globals.DatosDeCPU{
		PID: mensaje.PID,
		PC:  mensaje.PC,
	}
	// StringInstruccion = ObtenerInstruccion(mensaje.PID, mensaje.PC)
	// se debe devolver el string mediante un JSON por el ResponseWriter
	return globals.CPU
}

func ObtenerInstruccion(PID int, PC int) string {
	// PruebaFile, err := os.Open(pruebas/nombreArchivo)
	//	if err != nil { log.Fatal(err.Error()) }

	// string debería ser desde donde se abre el archivo hasta que se
	// deteca una nueva linea -> ahí deja de tomar chars y retorna ese string
	string := ""
	log.Printf("## PID: <%d>  - Obtener instrucción: <%d> - Instrucción: <%s>", PID, PC, string)
	return string
}

// FUNCION PARA RECIBIR LOS MENSAJES PROVENIENTES DEL KERNEL
func RecibirMensajeDeKernel(w http.ResponseWriter, r *http.Request) {
	var mensaje globals.DatosRespuestaDeKernel
	data.LeerJson(w, r, &mensaje)

	/*
		globals.Kernel = globals.DatosConsultaDeKernel{
			PID:            mensaje.PID,
			TamanioMemoria: mensaje.TamanioMemoria,
			// nombreArchivo
			// agregar a DatosConsultaDeKernel
		}*/

	log.Printf("Archivo Pseudocodigo: %s\n", mensaje.Pseudocodigo)
	log.Printf("Tamanio de Memoria Pedido: %d\n", mensaje.TamanioMemoria)
}

// ------------------------------------------------------------------
// ----------- FORMA PARTE DE LA MODIFICACIÓN DE PROCESOS -----------
// ------------------------------------------------------------------

func CreacionProceso(w http.ResponseWriter, r *http.Request) {
	tamanioDeseado := 1
	var datos globals.DatosDeCPU = RetornarMensajeDeCPU(w, r)

	log.Printf("## PID: <%d>  - Proceso Creado - Tamaño: <%d>", datos.PID, tamanioDeseado)
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
