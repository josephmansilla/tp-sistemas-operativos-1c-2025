package traducciones

import (
	"encoding/json"
	"fmt"
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/utils/data"
	"log"
)

// Estructura del mensaje hacia Memoria
type MensajeTabla struct {
	PID            int   `json:"pid"`
	IndicesEntrada []int `json:"indices_entrada"`
}

// Respuesta de Memoria
type RespuestaTabla struct {
	NumeroMarco int `json:"numero_marco"`
}

type Mensaje

// Función principal de traducción
func Traducir(dirLogica int) int {
	tamPagina := globals.TamanioPagina
	entradasPorNivel := globals.EntradasPorNivel
	niveles := globals.CantidadNiveles
	nroPagina := dirLogica / tamPagina
	desplazamiento := dirLogica % tamPagina

	//Primero verifico si la cache esta activa
	if cache.EstaActiva() {
		log.Printf("Cache Activa")
	} else {
		log.Printf("Cache Inactiva")
	}

	// Inicializo la TLB
	var tlb = NuevaTLB(globals.ClientConfig.TlbEntries, globals.ClientConfig.TlbReplacement)
	tlb.AgregarEntrada(1, 2)
	tlb.AgregarEntrada(3, 4)
	tlb.AgregarEntrada(5, 6)

	// Consulto la TLB
	if marco, ok := tlb.Buscar(nroPagina); ok {
		return marco*tamPagina + desplazamiento
	}

	// La página no está en la TLB, voy a Memoria
	entradas := descomponerPagina(nroPagina, niveles, entradasPorNivel)
	marco, err := accederTabla(globals.PIDActual, entradas)
	if err != nil {
		log.Printf("No se pudo acceder a la tabla de páginas: %s", err.Error())
		return -1
	}
	if marco == -1 {
		log.Printf("No se pudo traducir la dirección lógica %d", dirLogica)
		return -1
	}

	// Agrego entrada a la TLB
	tlb.AgregarEntrada(nroPagina, marco)
	// Agrego tmb a la cache
	if cache.EstaActiva() {
		cache.Agregar(nroPagina, "", true) // reemplazá "" con el contenido real si lo tenés
		log.Printf("PID: %d - CACHE ADD - Pagina: %d", globals.PIDActual, nroPagina)
	}
	return marco*tamPagina + desplazamiento
}

// Descompone el número de página en los índices para cada nivel
func descomponerPagina(nroPagina int, niveles int, entradasPorNivel int) []int {
	entradas := make([]int, niveles)
	divisor := 1

	for i := niveles - 1; i >= 0; i-- {
		entradas[i] = (nroPagina / divisor) % entradasPorNivel
		divisor *= entradasPorNivel
	}

	return entradas
}

// Hace una petición HTTP a Memoria para resolver una entrada de tabla
func accederTabla(pid int, indices []int) (int, error) {
	url := fmt.Sprintf("http://%s:%d/memoria/tabla",
		globals.ClientConfig.IpMemory,
		globals.ClientConfig.PortMemory,
	)

	mensaje := MensajeTabla{
		PID:            pid,
		IndicesEntrada: indices,
	}

	resp, err := data.EnviarDatosConRespuesta(url, mensaje)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	var respuesta RespuestaTabla
	if err := json.NewDecoder(resp.Body).Decode(&respuesta); err != nil {
		return -1, err
	}
	log.Printf("Marco Recibido: %d", respuesta.NumeroMarco)

	return respuesta.NumeroMarco, nil
}

func Leer(pagina int, tamanio int) (string, error) {
	if cache.EstaActiva() {
		contenido, err := LeerEnCache(pagina, tamanio)
		if err != nil {
			log.Printf("Error leyendo en cache: %v", err)
			return "", err
		}
		return contenido, nil
	} else {
		//TODO leer en memoria
		contenido := ""
		cache.Agregar(pagina, contenido, true)
		log.Printf("PID: %d - CACHE ADD - Pagina: %d", globals.PIDActual, pagina)
		return "", nil
	}
}

func Escribir(pagina int, datos string) error {
	if cache.EstaActiva() {
		if err := EscribirEnCache(pagina, datos); err != nil {
			log.Printf("Error escribiendo en cache: %v", err)
			return err
		}
		return nil
	} else {
		//TODO escribir en memoria
		cache.Agregar(pagina, datos, true)
		log.Printf("PID: %d - CACHE ADD - Pagina: %d", globals.PIDActual, pagina)
		return nil
	}
}

func LeerEnMemoria(pagina int, tamanio int) {

}

func EscribirEnMemoria(pagina int, datos string) {

}
