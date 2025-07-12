package administracion

import (
	g "github.com/sisoputnfrba/tp-golang/memoria/estructuras"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
)

func ObtenerDatosMemoria(direccionFisica int) (datosLectura g.ExitoLecturaPagina) {
	tamanioPagina := g.MemoryConfig.PagSize
	numeroPagina := direccionFisica / tamanioPagina
	offset := direccionFisica % tamanioPagina

	inicioFrame := numeroPagina * tamanioPagina
	finFrame := inicioFrame + tamanioPagina
	bytesRestantes := tamanioPagina - offset

	if direccionFisica+bytesRestantes > finFrame {
		logger.Error("Out of range - Lectura fuera del marco asignado")
	}

	pseudocodigoEnBytes := make([]byte, bytesRestantes)

	g.MutexMemoriaPrincipal.Lock()
	copy(pseudocodigoEnBytes, g.MemoriaPrincipal[direccionFisica:direccionFisica+bytesRestantes])
	g.MutexMemoriaPrincipal.Unlock()

	logger.Debug("Se obtuvo el pseudocodigo de memoria: %d", pseudocodigoEnBytes)

	pseudocodigoEnString := string(pseudocodigoEnBytes)

	datosLectura = g.ExitoLecturaPagina{
		Exito: nil,
		Valor: pseudocodigoEnString,
	}

	return
}

func LeerEspacioEntrada(pid int, direccionFisica int) (datosLectura g.ExitoLecturaPagina) {
	datosLectura = ObtenerDatosMemoria(direccionFisica)
	ModificarEstadoEntradaLectura(pid)
	return datosLectura
}

func ModificarEstadoEntradaLectura(pid int) {
	g.MutexProcesosPorPID.Lock()
	proceso := g.ProcesosPorPID[pid]
	g.MutexProcesosPorPID.Unlock()
	IncrementarMetrica(proceso, 1, IncrementarLecturaDeMemoria)
	logger.Info("## Modificacion del estado entrada exitosa")
}
