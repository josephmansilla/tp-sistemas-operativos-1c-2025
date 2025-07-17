package administracion

import (
	g "github.com/sisoputnfrba/tp-golang/memoria/estructuras"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"io"
	"os"
)

func CargarEntradasDesdeSwap(pid int) (entradas map[int]g.EntradaSwap, err error) {

	g.MutexSwapIndex.Lock()
	info, existe := g.SwapIndex[pid]
	if !existe {
		logger.Error("No existe el PID en el SwapIndex: %v", logger.ErrNoInstance)
		return nil, logger.ErrNoInstance
	}
	g.MutexSwapIndex.Unlock()

	file, err := os.Open(g.MemoryConfig.SwapfilePath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logger.Error("Error al cerrar el swapfle: %v", err)
			return
		}
	}(file)

	entradas = make(map[int]g.EntradaSwap, len(info.NumerosFrame))

	for i, entrada := range info.Entradas {
		_, errArch := file.Seek(int64(entrada.PosicionInicio), io.SeekStart)
		if errArch != nil {
			return nil, err
		}
		datos := make([]byte, 0)
		enviarEntrada := g.EntradaSwap{
			NumeroFrame: info.NumerosFrame[i],
			Datos:       datos,
			Tamanio:     entrada.Tamanio,
		}
		entradas[info.NumerosFrame[i]] = enviarEntrada
	}

	return entradas, nil
}

func CargarEntradasAMemoria(pid int, entradas map[int]g.EntradaSwap) error {

	for _, entrada := range entradas {

		dirFisica, err := ResignarPaginasParaPID(pid, entrada.NumeroFrame)
		if err != nil {
			return err
		}

		rta := EscribirEspacioEntrada(pid, dirFisica, entrada.Datos)
		if rta.Exito != nil {
			logger.Error("Error: %v", rta.Exito)
			return rta.Exito
		}

		logger.Info("## PID: <%d> - <SWAP A MEMORIA> - Dir. Física: <%d> - Tamaño: <%d>",
			pid,
			dirFisica,
			entrada.Tamanio,
		)
	}
	return nil
}

func ResignarPaginasParaPID(pid int, numeroPagina int) (int, error) {

	frameLibre, err := AsignarFrameLibre()
	if err != nil {
		logger.Error("No hay frames libres en el sistema %v", err)
		return -100, err
	}
	errr := ActualizarEntradaPaginaEnTabla(pid, numeroPagina, frameLibre)
	if errr != nil {
		return -100, errr
	}
	logger.Info("## PID <%d> ; Entrada: <%d> ; Frame: <%d>", pid, numeroPagina, frameLibre)
	return frameLibre * g.MemoryConfig.PagSize, nil
}
