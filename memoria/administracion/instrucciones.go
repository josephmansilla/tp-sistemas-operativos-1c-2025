package administracion

import (
	g "github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
)

func ObtenerInstruccion(proceso *g.Proceso, pc int) (respuesta g.InstruccionCPU, err error) {
	respuesta = g.InstruccionCPU{Exito: nil, Instruccion: ""}

	if proceso == nil {
		logger.Error("Proceso recibido es nil")
		return respuesta, logger.ErrProcessNil
	}

	lineaInstruccion := proceso.InstruccionesEnBytes[pc]

	respuesta.Instruccion = string(lineaInstruccion)
	return respuesta, nil
}
