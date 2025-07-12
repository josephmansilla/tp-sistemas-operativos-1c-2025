package administracion

import (
	"fmt"
	g "github.com/sisoputnfrba/tp-golang/memoria/estructuras"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
)

func EscribirEspacioEntrada(pid int, direccionFisica int, datosEscritura string) g.ExitoEscrituraPagina {
	stringEnBytes := []byte(datosEscritura)
	if len(datosEscritura) == 0 {
		logger.Debug("Los datos a escribir son vacios: %v", logger.ErrNoInstance)
	}
	err := ModificarEstadoEntradaEscritura(pid, direccionFisica, stringEnBytes)
	if err != nil {
		return g.ExitoEscrituraPagina{Exito: err, DireccionFisica: direccionFisica, Mensaje: err.Error()}
	}

	exito := g.ExitoEscrituraPagina{
		Exito:           nil,
		DireccionFisica: direccionFisica,
		Mensaje:         "Proceso fue modificado correctamente en memoria",
	}

	return exito
}

func ModificarEstadoEntradaEscritura(direccionFisica int, pid int, datosEnBytes []byte) (err error) {
	tamanioPagina := g.MemoryConfig.PagSize
	numeroPagina := direccionFisica / tamanioPagina

	inicioFrame := numeroPagina * tamanioPagina
	finFrame := inicioFrame + tamanioPagina

	if direccionFisica+len(datosEnBytes) > finFrame {
		logger.Error("Out of range - Escritura fuera del marco asignado")
		return logger.ErrSegmentFault
	}

	g.MutexMemoriaPrincipal.Lock()
	copy(g.MemoriaPrincipal[direccionFisica:], datosEnBytes)
	g.MutexMemoriaPrincipal.Unlock()

	logger.Error("Se escribió en memoria: %d", datosEnBytes)

	g.MutexProcesosPorPID.Lock()
	proceso := g.ProcesosPorPID[pid]
	g.MutexProcesosPorPID.Unlock()

	if proceso == nil {
		logger.Error("Se intentó acceder a un proceso inexistente o nil para PID=%d", pid)
		return fmt.Errorf("proceso nil para PID=%d", pid)
	}

	indices := CrearIndicePara(numeroPagina)
	entrada, err := BuscarEntradaPagina(proceso, indices)
	if err != nil {
		logger.Error("No se pudo encontrar la entrada de pagina para modificar informes: %v", err)
		return err
	}
	if entrada != nil {
		entrada.FueModificado = true
		entrada.EstaEnUso = true
	} else {
		logger.Error("Entrada vacia")
	}

	IncrementarMetrica(proceso, 1, IncrementarEscrituraDeMemoria)

	return nil
}
