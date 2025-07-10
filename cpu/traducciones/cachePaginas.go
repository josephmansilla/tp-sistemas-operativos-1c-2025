package traducciones

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"sync"
	"time"
)

type EntradaCache struct {
	NroPagina  int
	Contenido  string
	Modificado bool
	Usado      bool
}

type CachePaginas struct {
	Entradas    []EntradaCache
	MaxEntradas int
	Algoritmo   string
	Puntero     int
	mutex       sync.Mutex
}

var Cache *CachePaginas
var Max int

func InitCache() {
	Cache = NuevaCachePaginas()
}

func NuevaCachePaginas() *CachePaginas {
	if Max <= 0 {
		return nil
	}
	return &CachePaginas{
		Entradas:    make([]EntradaCache, 0, Max),
		MaxEntradas: Max,
		Algoritmo:   globals.ClientConfig.CacheReplacement,
		Puntero:     0,
	}
}

func (c *CachePaginas) Agregar(nroPagina int, contenido string, modificado bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	nueva := EntradaCache{
		NroPagina:  nroPagina,
		Contenido:  contenido,
		Modificado: modificado,
		Usado:      true,
	}

	if len(c.Entradas) < c.MaxEntradas {
		c.Entradas = append(c.Entradas, nueva)
	} else {
		c.reemplazarEntrada(nueva)
	}
}

func (c *CachePaginas) Buscar(nroPagina int) (string, bool) {
	time.Sleep(time.Millisecond * time.Duration(globals.ClientConfig.CacheDelay))

	c.mutex.Lock()
	defer c.mutex.Unlock()

	for i := range c.Entradas {
		if c.Entradas[i].NroPagina == nroPagina {
			c.Entradas[i].Usado = true
			logger.Info("PID: %d - Cache HIT - Página: %d", globals.PIDActual, nroPagina)
			return c.Entradas[i].Contenido, true
		}
	}
	logger.Info("PID: %d - Cache MISS - Página: %d", globals.PIDActual, nroPagina)
	return "", false
}

func (c *CachePaginas) EstaActiva() bool {
	return c != nil && c.MaxEntradas > 0
}

func (c *CachePaginas) MarcarUso(nroPagina int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for i := range c.Entradas {
		if c.Entradas[i].NroPagina == nroPagina {
			c.Entradas[i].Usado = true
			logger.Info("PID: %d - Cache USO - Página %d marcada con bit de uso en true", globals.PIDActual, nroPagina)
			return
		}
	}
}

func LeerEnCache(nroPagina int, tamanio int) (string, error) {
	time.Sleep(time.Millisecond * time.Duration(globals.ClientConfig.CacheDelay))

	contenido, ok := Cache.Buscar(nroPagina)
	if !ok {
		err := fmt.Errorf("Página %d no encontrada en la caché", nroPagina)
		logger.Error("Error: %v", err)
		return "", err
	}
	Cache.MarcarUso(nroPagina)
	// Si querés leer solo parte del contenido: contenido = contenido[:tamanio]
	return contenido, nil
}

func EscribirEnCache(nroPagina int, datos string) error {
	time.Sleep(time.Millisecond * time.Duration(globals.ClientConfig.CacheDelay))

	for i := range Cache.Entradas {
		if Cache.Entradas[i].NroPagina == nroPagina {
			Cache.Entradas[i].Contenido = datos
			Cache.Entradas[i].Modificado = true
			Cache.Entradas[i].Usado = true
			logger.Info("Se escribió %s en página %d", datos, nroPagina)
			return nil
		}
	}
	err := fmt.Errorf("No se encontró la página %d en la caché", nroPagina)
	logger.Error("Error: %v", err)
	return err
}

func (c *CachePaginas) reemplazarEntrada(nueva EntradaCache) {
	switch c.Algoritmo {
	case "CLOCK":
		c.reemplazoClock(nueva)
	case "CLOCK-M":
		c.reemplazoClockM(nueva)
	default:
		logger.Info("Algoritmo de reemplazo inválido: %s", c.Algoritmo)
	}
}

func (c *CachePaginas) reemplazoClock(nueva EntradaCache) {
	for {
		entrada := &c.Entradas[c.Puntero]
		if !entrada.Usado {
			tamPagina := globals.TamanioPagina
			dirLogica := entrada.NroPagina * tamPagina
			dirFisica := Traducir(dirLogica)
			if dirFisica != -1 {
				err := EscribirEnMemoria(dirFisica, entrada.Contenido)
				if err != nil {
					logger.Error("Error al escribir página modificada %d en dirección física %d: %v", entrada.NroPagina, dirFisica, err)
				} else {
					logger.Info("Página modificada %d escrita en dirección física %d antes de reemplazo", entrada.NroPagina, dirFisica)
				}
			}
		}
		logger.Info("Reemplazo CLOCK - Página %d reemplazada por Página %d", entrada.NroPagina, nueva.NroPagina)
		c.Entradas[c.Puntero] = nueva
		c.Puntero = (c.Puntero + 1) % c.MaxEntradas // Para volver a 0 si se pasa de las entradas -> (3+1) % 4 = 0
		return

		entrada.Usado = false
		c.Puntero = (c.Puntero + 1) % c.MaxEntradas
	}
}

func (c *CachePaginas) reemplazoClockM(nueva EntradaCache) {
	//Primero busco U=0, M=0
	for i := 0; i < c.MaxEntradas; i++ {
		indice := (c.Puntero + i) % c.MaxEntradas
		entrada := &c.Entradas[indice]
		if !entrada.Usado && !entrada.Modificado {
			logger.Info("Reemplazo CLOCK-M (0,0) - Página %d reemplazada por Página %d", entrada.NroPagina, nueva.NroPagina)
			c.Entradas[indice] = nueva
			c.Puntero = (indice + 1) % c.MaxEntradas
			return
		}
	}

	// Despues busco U=0, M=1
	for i := 0; i < c.MaxEntradas; i++ {
		indice := (c.Puntero + i) % c.MaxEntradas
		entrada := &c.Entradas[indice]
		if !entrada.Usado && entrada.Modificado {
			tamPagina := globals.TamanioPagina
			dirLogica := entrada.NroPagina * tamPagina
			dirFisica := Traducir(dirLogica)
			if dirFisica != -1 {
				err := EscribirEnMemoria(dirFisica, entrada.Contenido)
				if err != nil {
					logger.Error("Error al escribir página modificada %d en dirección física %d: %v", entrada.NroPagina, dirFisica, err)
				} else {
					logger.Info("Página modificada %d escrita en dirección física %d antes de reemplazo", entrada.NroPagina, dirFisica)
				}
			}
			logger.Info("Reemplazo CLOCK-M (0,1) - Página %d reemplazada por Página %d", entrada.NroPagina, nueva.NroPagina)
			c.Entradas[indice] = nueva
			c.Puntero = (indice + 1) % c.MaxEntradas
			return
		}
	}

	// No encontre, seteo todos los U=0 para volver a buscar
	for i := 0; i < c.MaxEntradas; i++ {
		c.Entradas[i].Usado = false
	}

	// Vuelvo a buscar
	c.reemplazoClockM(nueva)
}

func (c *CachePaginas) LimpiarCache() {
	if c == nil {
		return
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	tamPagina := globals.TamanioPagina

	for _, entrada := range c.Entradas {
		if entrada.Modificado {
			dirLogica := entrada.NroPagina * tamPagina

			dirFisica := Traducir(dirLogica)
			if dirFisica == -1 {
				logger.Error("Error al traducir la dirección lógica de página %d", entrada.NroPagina)
				continue
			}

			err := EscribirEnMemoria(dirFisica, entrada.Contenido)
			if err != nil {
				logger.Error("Error al escribir página %d en dirección física %d: %v", entrada.NroPagina, dirFisica, err)
				continue
			}

			logger.Info("Página modificada %d escrita en dirección física %d", entrada.NroPagina, dirFisica)
		}
	}

	//Elimino las entradas
	c.Entradas = make([]EntradaCache, 0, c.MaxEntradas)
	c.Puntero = 0

	logger.Info("Caché vaciada correctamente")
}
