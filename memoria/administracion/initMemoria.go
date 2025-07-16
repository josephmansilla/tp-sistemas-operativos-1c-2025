package administracion

import (
	g "github.com/sisoputnfrba/tp-golang/memoria/estructuras"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"sync"
)

func InicializarMemoriaPrincipal() {

	cantidadFrames := g.MemoryConfig.MemorySize / g.MemoryConfig.PagSize

	g.MemoriaPrincipal = make([]byte, g.MemoryConfig.MemorySize)
	ConfigurarFrames(cantidadFrames)
	InstanciarEstructurasGlobales()

	logger.Info("Memoria Principal Inicializada con %d bytes de tamaño con %d frames de %d bytes.",
		g.MemoryConfig.MemorySize, cantidadFrames, g.MemoryConfig.PagSize)
}

func InstanciarEstructurasGlobales() {
	g.ProcesosPorPID = make(map[int]*g.Proceso)
	g.SwapIndex = make(map[int]*g.SwapProcesoInfo)
	g.MutexMetrica = make(map[int]*sync.Mutex, g.MemoryConfig.MemorySize)
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
