package administracion

import (
	g "github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"sync"
)

func InicializarMemoriaPrincipal() {
	tamanioMemoriaPrincipal := g.MemoryConfig.MemorySize
	tamanioPagina := g.MemoryConfig.PagSize
	cantidadFrames := tamanioMemoriaPrincipal / tamanioPagina

	g.MemoriaPrincipal = make([]byte, tamanioMemoriaPrincipal)
	ConfigurarFrames(cantidadFrames)
	InstanciarEstructurasGlobales()
	InicializarSemaforos()

	logger.Info("Memoria Principal Inicializada con %d bytes de tamaño con %d frames de %d.",
		tamanioMemoriaPrincipal, cantidadFrames, tamanioPagina)
}

func InstanciarEstructurasGlobales() {
	g.ProcesosPorPID = make(map[int]*g.Proceso)
	g.SwapIndex = make(map[int]*g.SwapProcesoInfo)
	g.EstaEnSwap = make(map[int]bool)
}

func InicializarSemaforos() {
	g.MutexMetrica = make([]sync.Mutex, g.MemoryConfig.MemorySize) // tamaño totalmente arbitrario
}

func ConfigurarFrames(cantidadFrames int) {
	g.FramesLibres = make([]bool, cantidadFrames)
	g.MutexEstructuraFramesLibres.Lock()
	for i := 0; i < cantidadFrames; i++ {
		g.FramesLibres[i] = true
	}
	g.MutexEstructuraFramesLibres.Unlock()
	g.CantidadFramesLibres = cantidadFrames
}
