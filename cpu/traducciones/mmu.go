package traducciones

import (
	"encoding/json"
	"fmt"
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/utils/data"
	"log"
)

type MensajeAMemoria struct {
	PID                int `json:"pid"`
	NroPagina          int `json:"nroPagina"`
	CantidadEntradasTP int `json:"cantidadEntradasTP"`
}

func Traducir(pid int, dirLogica int) int {
	nroPagina := dirLogica / globals.TamPag
	desplazamiento := dirLogica % globals.TamPag

	var tlb = NuevaTLB(globals.ClientConfig.TlbEntries, globals.ClientConfig.TlbReplacement)

	marco, ok := tlb.Buscar(nroPagina)
	if !ok {
		marco = ObtenerMarco(pid, nroPagina)
		if marco == -1 {
			log.Printf("Error al obtener el marco de memoria para la página %d", nroPagina)
			return -1 // Indicar que ocurrió un error
		}
		tlb.AgregarEntrada(nroPagina, marco)
	}

	return marco*globals.TamPag + desplazamiento
}

func ObtenerMarco(pid int, nroPagina int) int {
	mensaje := MensajeAMemoria{
		PID:       pid,
		NroPagina: nroPagina,
	}

	url := fmt.Sprintf("http://%s:%d/memoria/marco",
		globals.ClientConfig.IpMemory,
		globals.ClientConfig.PortMemory,
	)

	resp, err := data.EnviarDatosConRespuesta(url, mensaje)
	if err != nil {
		log.Printf("Error enviando PID, NroPagina y CantidadEntradasTP a Memoria: %s", err.Error())
		return -1
	}
	defer resp.Body.Close()

	var respuesta struct {
		Marco int `json:"marco"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&respuesta); err != nil {
		log.Printf("Error al decodificar la respuesta de Memoria: %s", err.Error())
		return -1
	}

	log.Printf("PID: %s - OBTENER MARCO - Página: %s - Marco: %s", globals.Pcb.PID, nroPagina, respuesta.Marco)
	return respuesta.Marco
}
