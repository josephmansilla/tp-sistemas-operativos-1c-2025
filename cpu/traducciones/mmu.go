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

// Función principal de traducción
func Traducir(dirLogica int) int {
	tamPagina := globals.TamanioPagina
	entradasPorNivel := globals.EntradasPorNivel
	niveles := globals.CantidadNiveles
	nroPagina := dirLogica / tamPagina
	desplazamiento := dirLogica % tamPagina

	//Inicializo la TLB
	var tlb = NuevaTLB(globals.ClientConfig.TlbEntries, globals.ClientConfig.TlbReplacement)
	tlb.AgregarEntrada(1, 2)
	tlb.AgregarEntrada(3, 4)
	tlb.AgregarEntrada(5, 6)

	// Consulto la tlb
	if marco, ok := tlb.Buscar(nroPagina); ok {
		return marco*tamPagina + desplazamiento
	}

	// La pagina no esta en la tlb, tengo que ir a memoria
	entradas := descomponerPagina(nroPagina, niveles, entradasPorNivel)

	var marco int = -1

	resp, err := accederTabla(globals.PIDActual, entradas)
	if err != nil {
		log.Fatal("ELEGI QUE PONERLE SANTI")
	}
	marco = resp

	if marco == -1 {
		log.Printf("No se pudo traducir la dirección lógica %d", dirLogica)
		return -1
	}

	// Agrego la entrada a la tlb
	tlb.AgregarEntrada(nroPagina, marco)

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
func accederTabla(pid int, indices []int) (RespuestaTabla, error) {
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
		return RespuestaTabla{}, err
	}
	defer resp.Body.Close()

	var respuesta RespuestaTabla
	if err := json.NewDecoder(resp.Body).Decode(&respuesta); err != nil {
		return RespuestaTabla{}, err
	}

	return respuesta, nil
}
