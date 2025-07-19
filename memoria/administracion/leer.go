package administracion

import (
	g "github.com/sisoputnfrba/tp-golang/memoria/estructuras"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
)

func LeerEspacioEntrada(pid int, direccionFisica int) (datosLectura g.RespuestaLectura) {
	datosLectura = ObtenerDatosMemoria(direccionFisica)
	err := ModificarEstadoEntradaLectura(pid)
	if err != nil {
		return g.RespuestaLectura{Exito: err}
	}
	return datosLectura
}

func ObtenerDatosMemoria(direccionFisica int) (datosLectura g.RespuestaLectura) {
	tamanioPagina := g.MemoryConfig.PagSize

	finFrame := direccionFisica + tamanioPagina
	bytesRestantes := tamanioPagina - direccionFisica%tamanioPagina

	if direccionFisica+bytesRestantes > finFrame {
		logger.Error("Out of range - Lectura fuera del marco asignado")
	}
	/* ALWAYS FALSE
	if bytesRestantes < 0 {
		logger.Error("La lectura es más grande que la página")
	}*/

	datosEnBytes := make([]byte, bytesRestantes)

	g.MutexMemoriaPrincipal.Lock()
	copy(datosEnBytes, g.MemoriaPrincipal[direccionFisica:direccionFisica+bytesRestantes])
	g.MutexMemoriaPrincipal.Unlock()

	datosLectura = g.RespuestaLectura{
		Exito: nil,
		Valor: string(datosEnBytes),
	}

	return
}

func ModificarEstadoEntradaLectura(pid int) error {
	g.MutexProcesosPorPID.Lock()
	proceso := g.ProcesosPorPID[pid]
	g.MutexProcesosPorPID.Unlock()

	if proceso == nil {
		logger.Error("Se intentó acceder a un proceso inexistente o nil para PID <%d>", pid)
		return logger.ErrProcessNil
	}

	IncrementarMetrica(proceso, 1, IncrementarLecturaDeMemoria)
	//logger.Info("## Modificacion del estado entrada exitosa")

	return nil
}
