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
	PID           int `json:"pid"`
	NumeroTabla   int `json:"numeroTabla"`
	EntradaIndice int `json:"entrada"`
}

// Respuesta de Memoria
type RespuestaTabla struct {
	EsUltimoNivel bool `json:"esUltimoNivel"`
	NumeroTabla   int  `json:"numeroTabla"`
	NumeroMarco   int  `json:"marco"`
}

// Función principal de traducción
func Traducir(dirLogica int, tlb *TLB) int {
	tamPagina := globals.TamPag
	entradasPorNivel := globals.EntradasPorNivel
	niveles := globals.CantidadNiveles
	nroPagina := dirLogica / tamPagina
	desplazamiento := dirLogica % tamPagina

	// Consulto la tlb
	if marco, ok := tlb.Buscar(nroPagina); ok {
		return marco*tamPagina + desplazamiento
	}

	// La pagina no esta en la tlb, tengo que ir a memoria
	entradas := descomponerPagina(nroPagina, niveles, entradasPorNivel)

	var tablaActual = 0
	var marco int = -1

	for nivel := 0; nivel < niveles; nivel++ {
		resp, err := accederTabla(globals.PIDActual, tablaActual, entradas[nivel])
		if err != nil {
			log.Printf("Error accediendo a nivel %d de la tabla: %s", nivel, err.Error())
			return -1
		}

		if resp.EsUltimoNivel {
			marco = resp.NumeroMarco
			break
		} else {
			tablaActual = resp.NumeroTabla //vuelvo a ejecutar hasta llegar a la ult tabla
		}
	}

	if marco == -1 {
		log.Printf("No se pudo traducir la dirección lógica %d", dirLogica)
		return -1
	}

	// Agrego la entrada a la tlb
	tlb.AgregarEntrada(nroPagina, marco)

	return marco*tamPagina + desplazamiento
}

// Descompone el número de página en los índices para cada nivel
func descomponerPagina(nroPagina, niveles, entradasPorNivel int) []int {
	entradas := make([]int, niveles)
	divisor := 1

	for i := niveles - 1; i >= 0; i-- {
		entradas[i] = (nroPagina / divisor) % entradasPorNivel
		divisor *= entradasPorNivel
	}

	return entradas
}

// Hace una petición HTTP a Memoria para resolver una entrada de tabla
func accederTabla(pid, nroTabla, entrada int) (RespuestaTabla, error) {
	url := fmt.Sprintf("http://%s:%d/memoria/tabla",
		globals.ClientConfig.IpMemory,
		globals.ClientConfig.PortMemory,
	)

	mensaje := MensajeTabla{
		PID:           pid,
		NumeroTabla:   nroTabla,
		EntradaIndice: entrada,
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
