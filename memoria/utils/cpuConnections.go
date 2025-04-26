package utils

import (
	"bufio"
	"encoding/json"
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/data"
	"log"
	"net/http"
	"os"
	"strings"
)

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

// FUNCION PARA DEVOLVER/RETORNAR LOS MENSAJES PROVENIENTES DE LA CPU
// DONDE LA DE ARRIBA NO RETORNA NADA
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

func CargarInstrucciones(nombreArchivo string) {
	ruta := "../pruebas/" + nombreArchivo
	// con el (..) vuelve para atras en los directorios
	// accede a la carpeta de pruebas y abre el archivo pasado x parametro

	file, err := os.Open(ruta)
	if err != nil {
		log.Printf("Error al abrir el archivo: %s\n", err)
		return
	}
	defer file.Close() // se accede desde cualquier parte del código

	log.Println("Se leyó el archivo")
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lineaPseudocodigo := scanner.Text()
		log.Printf("Línea leída:%s\n", lineaPseudocodigo)
		CargarListaDeInstrucciones(lineaPseudocodigo)
		if strings.TrimSpace(lineaPseudocodigo) == "EOF" {
			break
		}

	}
	if err := scanner.Err(); err != nil {
		log.Printf("Error al leer el archivo:%s\n", err)
	}
	log.Printf("Total de instrucciones cargadas: %d\n", len(Instrucciones))
}

func ObtenerInstruccion(w http.ResponseWriter, r *http.Request) {
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
