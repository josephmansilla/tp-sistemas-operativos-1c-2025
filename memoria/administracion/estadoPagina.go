package administracion

import (
	g "github.com/sisoputnfrba/tp-golang/memoria/globals"
)

func MarcarOcupadoFrame(numeroFrame int) {
	g.MutexEstructuraFramesLibres.Lock()
	g.FramesLibres[numeroFrame] = false
	g.MutexEstructuraFramesLibres.Unlock()

	g.MutexCantidadFramesLibres.Lock()
	g.CantidadFramesLibres--
	g.MutexCantidadFramesLibres.Unlock()
}

func MarcarLibreFrame(numeroFrameALiberar int) {
	g.MutexEstructuraFramesLibres.Lock()
	g.FramesLibres[numeroFrameALiberar] = true
	g.MutexEstructuraFramesLibres.Unlock()

	g.MutexCantidadFramesLibres.Lock()
	g.CantidadFramesLibres++
	g.MutexCantidadFramesLibres.Unlock()
}
