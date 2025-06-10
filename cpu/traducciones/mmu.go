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

type MensajeEscritura struct {
	PID       int    `json:"pid"`
	DirFisica int    `json:"dirFisica"`
	Datos     string `json:"datos"`
}

type MensajeLectura struct {
	PID       int `json:"pid"`
	DirFisica int `json:"dirFisica"`
	Tamanio   int `json:"tamanio"`
}

type RespuestaLectura struct {
	ValorLeido string `json:"valorLeido"`
}

// Función principal de traducción
func Traducir(dirLogica int) int {
	tamPagina := globals.TamanioPagina
	entradasPorNivel := globals.EntradasPorNivel
	niveles := globals.CantidadNiveles

	nroPagina := dirLogica / tamPagina
	desplazamiento := dirLogica % tamPagina

	//Primero verifico si la cache esta activa
	if Cache.EstaActiva() {
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
	if Cache.EstaActiva() {
		Cache.Agregar(nroPagina, "", true) // reemplazá "" con el contenido real si lo tenés
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

func LeerEnMemoria(dirFisica int, tamanio int) (string, error) {
	msg := MensajeLectura{
		PID:       globals.PIDActual,
		DirFisica: dirFisica,
		Tamanio:   tamanio,
	}

	url := fmt.Sprintf("http://%s:%d/memoria/lectura",
		globals.ClientConfig.IpMemory,
		globals.ClientConfig.PortMemory,
	)

	resp, err := data.EnviarDatosConRespuesta(url, msg)
	if err != nil {
		log.Printf("Error enviando Direccion Fisica y Tamanio: %s", err.Error())
		return "", err
	}
	defer resp.Body.Close()

	var respuesta RespuestaLectura
	if err := json.NewDecoder(resp.Body).Decode(&respuesta); err != nil {
		log.Printf("Error decodificando respuesta de memoria: %s", err.Error())
		return "", err
	}

	log.Printf("Direccion Fisica: %d y Tamanio: %d enviados correctamente a memoria", dirFisica, tamanio)
	log.Printf("Valor leído: %s", respuesta.ValorLeido)

	return respuesta.ValorLeido, nil
}

func EscribirEnMemoria(dirFisica int, datos string) error {
	msg := MensajeEscritura{
		PID:       globals.PIDActual,
		DirFisica: dirFisica,
		Datos:     datos,
	}

	url := fmt.Sprintf("http://%s:%d/memoria/escritura",
		globals.ClientConfig.IpMemory, globals.ClientConfig.PortMemory)

	err := data.EnviarDatos(url, msg)
	if err != nil {
		log.Printf("Error enviando Direccion Fisica y Datos: %s", err.Error())
		return err
	}

	log.Println("Direccion Fisica: %d y Datos: %s enviados correctamente a memoria", dirFisica, datos)
	return nil
}
