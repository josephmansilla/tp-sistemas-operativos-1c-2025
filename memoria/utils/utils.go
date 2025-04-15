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

// probando cositas
func RetornarMensajeDeCPU(w http.ResponseWriter, r *http.Request) globals.DatosDeCPU {
	var mensaje globals.DatosDeCPU
	LeerJson(w, r, &mensaje)

	globals.CPU = globals.DatosDeCPU{
		PID: mensaje.PID,
		PC:  mensaje.PC,
	}

	return globals.CPU
}

func RecibirMensajeDeKernel(w http.ResponseWriter, r *http.Request) {
	var mensaje globals.DatosDeKernel
	LeerJson(w, r, &mensaje)

	globals.Kernel = globals.DatosDeKernel{
		PID:            mensaje.PID,
		TamanioMemoria: mensaje.TamanioMemoria,
	}

	log.Printf("PID Pedido: %d\n", mensaje.PID)
	log.Printf("Tamanio de Memoria Pedido: %d\n", mensaje.TamanioMemoria)
}

// FORMA PARTE DE LA MODIFICACION DE PROCESOS
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

// FORMA PARTE DEL ACCESO A ESPACIO DE USUARIO
func EscrituraEspacio(w http.ResponseWriter, r *http.Request) {
	//toDO
	log.Printf("## PID: <PID>  - <Escritura> - Dir. Física: <DIRECCIÓN_FÍSICA> - Tamaño: <TAMAÑO>")
}

func LecturaEspacio(w http.ResponseWriter, r *http.Request) {
	//toDO
	log.Printf("## PID: <PID>  - <Lectura> - Dir. Física: <DIRECCIÓN_FÍSICA> - Tamaño: <TAMAÑO>")
}

func ObtenerInstruccion(w http.ResponseWriter, r *http.Request) {
	// toDO
	log.Printf("## PID: <PID>  - Obtener instrucción: <PC> - Instrucción: <INSTRUCCIÓN> <...ARGS>")
}

func MemoryDump(w http.ResponseWriter, r *http.Request) {
	// toDO
	// Llamado: "<PID>-<TIMESTAMP>.dmp" dentro del path definido por el archivo de configuración
	log.Printf("## PID: <PID>  - Memory Dump solicitado")
}

// FORMA PARTE DEL ACCESO A LAS TABLAS DE PÁGINAS
func AccesoTablaPaginas(w http.ResponseWriter, r *http.Request) {
	//toDO
}
func LeerPaginaCompleta(w http.ResponseWriter, r *http.Request) {
	//toDO
}
func ActualizarPaginaCompleta(w http.ResponseWriter, r *http.Request) {
	//toDO
}
