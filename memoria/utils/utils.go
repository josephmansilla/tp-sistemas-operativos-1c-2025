package utils

import (
	"bufio"
	"encoding/json"
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"log"
	"net/http"
	"os"
	"strings"
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

func LeerJson(w http.ResponseWriter, r *http.Request, mensaje any) {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&mensaje)

	if err != nil {
		log.Printf("Error al decodificar el mensaje: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error al decodificar mensaje"))
		return
	}

	log.Println("Me llego un mensaje JSON:")
	log.Printf("%+v\n", mensaje)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("STATUS OK"))
}

// FUNCION PARA RECIBIR LOS MENSAJES PROVENIENTES DE LA CPU
func RecibirMensajeDeCPU(w http.ResponseWriter, r *http.Request) {
	var mensaje globals.DatosDeCPU
	LeerJson(w, r, &mensaje)

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
	LeerJson(w, r, &mensaje)

	globals.CPU = globals.DatosDeCPU{
		PID: mensaje.PID,
		PC:  mensaje.PC,
	}
	// StringInstruccion = ObtenerInstruccion(mensaje.PID, mensaje.PC)
	// se debe devolver el string mediante un JSON por el ResponseWriter
	return globals.CPU
}

var Instrucciones []string = []string{}

func CargarInstrucciones() { // TODO: (nombreArchivo string)
	var nombreArchivo string = "archivoPrueba.txt" // TODO: globals.RespuestaKernel.Pseudocodigo

	ruta := "../pruebas/" + nombreArchivo

	file, err := os.Open(ruta)
	if err != nil {
		log.Printf("Error al abrir el archivo: %s\n", err)
		return
	}
	defer file.Close()

	log.Println("Se leyó el archivo\n")
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lineaPseudocodigo := scanner.Text()
		log.Printf("Línea leída:%s\n", lineaPseudocodigo)
		if strings.TrimSpace(lineaPseudocodigo) == "" {
			continue
		}
		CargarListaDeInstrucciones(lineaPseudocodigo)

	}
	if err := scanner.Err(); err != nil {
		log.Printf("Error al leer el archivo:%s\n", err)
	}
	log.Printf("Total de instrucciones cargadas: %d\n", len(Instrucciones))
}

func CargarListaDeInstrucciones(str string) {
	Instrucciones = append(Instrucciones, str)
	log.Println("Se cargó una instrucción al Slice\n")
}

func ObtenerInstruccion(w http.ResponseWriter, r *http.Request) {
	CargarInstrucciones() // TODO: nombreArchivo por parametro

	var mensaje globals.ContextoDeCPU
	err := json.NewDecoder(r.Body).Decode(&mensaje)
	if err != nil {
		http.Error(w, "Error leyendo JSON del CPU\n", http.StatusBadRequest)
		return
	}

	pc := mensaje.PC
	var instruccion string
	if pc >= 0 && pc < len(Instrucciones) {
		instruccion = Instrucciones[pc]
	} else {
		instruccion = "" // Esto indica fin del archivo o error de PC
	}

	respuesta := globals.InstruccionCPU{
		Instruccion: instruccion,
	}

	log.Printf("## PID: <%d>  - Obtener instrucción: <%d> - Instrucción: <%s>", mensaje.PID, mensaje.PC, instruccion)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respuesta)
}

// FUNCION PARA RECIBIR LOS MENSAJES PROVENIENTES DEL KERNEL
func RecibirMensajeDeKernel(w http.ResponseWriter, r *http.Request) {
	var mensaje globals.DatosConsultaDeKernel
	LeerJson(w, r, &mensaje)

	globals.Kernel = globals.DatosConsultaDeKernel{
		PID:            mensaje.PID,
		TamanioMemoria: mensaje.TamanioMemoria,
		// nombreArchivo
		// agregar a DatosConsultaDeKernel
	}

	log.Printf("Se cargó todo bien!\n")
	log.Printf("PID Pedido: %d\n", mensaje.PID)
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
