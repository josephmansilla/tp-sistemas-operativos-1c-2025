package traducciones

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
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
	time.Sleep(time.Millisecond * time.Duration(globals.ClientConfig.CacheDelay))
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

func EscribirEnCache(nroPagina int, datos string) error {
	time.Sleep(time.Millisecond * time.Duration(globals.ClientConfig.CacheDelay))

	for i := range Cache.Entradas {
		if Cache.Entradas[i].NroPagina == nroPagina {
			Cache.Entradas[i].Contenido = datos
			Cache.Entradas[i].Modificado = true
			Cache.Entradas[i].Usado = true
			//logger.Info("Se escribió %s en página %d", datos, nroPagina)
			return nil
		}
	}
	err := fmt.Errorf("no se encontró la página %d en la caché", nroPagina)
	//logger.Error("Error: %v", err)
	return err
}

func (c *CachePaginas) reemplazarEntrada(nueva EntradaCache) {
	time.Sleep(time.Millisecond * time.Duration(globals.ClientConfig.CacheDelay))
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
			if entrada.Modificado {
				// Si la página fue modificada, escribir en memoria
				tamPagina := globals.TamanioPagina
				dirLogica := entrada.NroPagina * tamPagina
				dirFisica := Traducir(dirLogica)

				if dirFisica != -1 {
					err := EscribirEnMemoria(dirFisica, entrada.Contenido)
					if err != nil {
						logger.Error("Error al escribir página modificada %d en dirección física %d: %v", entrada.NroPagina, dirFisica, err)
					} else {
						logger.Info("PID: %d - Memory Update - Pagina: %d - Frame: %d", globals.PIDActual, entrada.NroPagina, dirFisica)
					}
				}
			}

			logger.Info("Reemplazo CLOCK - Página %d reemplazada por Página %d", entrada.NroPagina, nueva.NroPagina)
			c.Entradas[c.Puntero] = nueva
			return
		}

		// Marcar la entrada como no usada y avanzar el puntero
		entrada.Usado = false
		c.Puntero = (c.Puntero + 1) % c.MaxEntradas
	}
}

func (c *CachePaginas) reemplazoClockM(nueva EntradaCache) {
	for {
		// primero busco U=0 y M=0
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

		// no encontre, busco U=0 y M=1
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
						logger.Info("PID: %d - Memory Update - Pagina: %d - Frame: %d", globals.PIDActual, entrada.NroPagina, dirFisica)
					}
				}
				logger.Info("Reemplazo CLOCK-M - Página %d reemplazada por Página %d", entrada.NroPagina, nueva.NroPagina)
				c.Entradas[indice] = nueva
				c.Puntero = (indice + 1) % c.MaxEntradas
				return
			}
		}

		// no encontre nada, seteo u = 0
		for i := 0; i < c.MaxEntradas; i++ {
			c.Entradas[i].Usado = false
		}
		// vuelvo a arrancar el loop
	}
}

func (c *CachePaginas) LimpiarCache() {
	time.Sleep(time.Millisecond * time.Duration(globals.ClientConfig.CacheDelay))
	if c == nil {
		return
	}

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

			logger.Info("PID: %d - Memory Update - Pagina: %d - Frame: %d", globals.PIDActual, entrada.NroPagina, dirFisica)
		}
	}

	//Elimino las entradas
	c.Entradas = make([]EntradaCache, 0, c.MaxEntradas)
	c.Puntero = 0

	//logger.Info("Caché vaciada correctamente")
}
