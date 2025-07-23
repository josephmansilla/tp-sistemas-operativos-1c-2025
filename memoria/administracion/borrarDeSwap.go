package administracion

import (
	g "github.com/sisoputnfrba/tp-golang/memoria/estructuras"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"io"
	"os"
)

// ========== DESOCUPO ESTRUCTURAS SWAP ==========

func DesocuparProcesoDeSwap(pid int) error {

	archivo, errAbrir := os.OpenFile(g.MemoryConfig.SwapfilePath, os.O_WRONLY, 0666)
	if errAbrir != nil {
		return errAbrir
	}
	defer func(archivo *os.File) {
		errCerrar := archivo.Close()
		if errCerrar != nil {
		}
	}(archivo)

	g.MutexSwapIndex.Lock()
	defer g.MutexSwapIndex.Unlock()

	for i := 0; i < len(g.SwapIndex[pid].Entradas); i++ {

		entrada := g.SwapIndex[pid].Entradas[i]
		if entrada.Tamanio != 0 {
			errBorrar := BorrarSeccionSwap(archivo, int64(entrada.PosicionInicio), entrada.Tamanio)
			if errBorrar != nil {
				return errBorrar
			}
		}

	}

	delete(g.SwapIndex, pid)

	return nil
}

// ========== LIBERO ESPACIO EN SWAP ==========

func BorrarSeccionSwap(archivo *os.File, posicionInicial int64, tamanio int) error {
	defer func(archivo *os.File) {
		err := archivo.Close()
		if err != nil {
			logger.Error("Archivo SWAP cerrado...")
		}
	}(archivo)

	_, errSeek := archivo.Seek(posicionInicial, io.SeekStart)
	if errSeek != nil {
		return errSeek
	}

	relleno := make([]byte, tamanio)

	_, errWrite := archivo.Write(relleno)
	if errWrite != nil {
		return errWrite
	}

	return nil
}
