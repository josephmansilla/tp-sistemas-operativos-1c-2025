package traducciones

import (
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"log"
	"sync"
)

type EntradaCache struct {
	Pagina     int
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

var cache *CachePaginas
var max = globals.ClientConfig.CacheEntries

func NuevaCachePaginas() *CachePaginas {
	if max <= 0 {
		return nil
	}
	return &CachePaginas{
		Entradas:    make([]EntradaCache, 0, max),
		MaxEntradas: max,
		Algoritmo:   globals.ClientConfig.CacheReplacement,
		Puntero:     0,
	}
}

func (c *CachePaginas) Agregar(pagina int, contenido string, modificado bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	nueva := EntradaCache{
		Pagina:     pagina,
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

func (c *CachePaginas) Buscar(pagina int) (string, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for i := range c.Entradas {
		if c.Entradas[i].Pagina == pagina {
			c.Entradas[i].Usado = true
			log.Printf("PID: %d - CACHE HIT - Página: %d", globals.PIDActual, pagina)
			return c.Entradas[i].Contenido, true
		}
	}
	log.Printf("PID: %d - CACHE MISS - Página: %d", globals.PIDActual, pagina)
	return "", false
}

func (c *CachePaginas) EstaActiva() bool {
	return c != nil && c.MaxEntradas > 0
}

func (c *CachePaginas) reemplazarEntrada(nueva EntradaCache) {
	//TODO
}

func (c *CachePaginas) MarcarUso(pagina int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for i := range c.Entradas {
		if c.Entradas[i].Pagina == pagina {
			c.Entradas[i].Usado = true
			log.Printf("PID: %d - CACHE USO - Página %d marcada con bit de uso en true", globals.PIDActual, pagina)
			return
		}
	}
}

func LeerEnCache(pagina int, tamanio int) (string, error) {
	if contenido, ok := cache.Buscar(pagina); ok {
		cache.MarcarUso(pagina) // actualizás el bit de uso si usás CLOCK
		// Simulás leer solo una porción si querés
		return contenido, nil
	}
	return "", nil
}

func EscribirEnCache(pagina int, datos string) error {
	for i := range cache.Entradas {
		if cache.Entradas[i].Pagina == pagina {
			cache.Entradas[i].Contenido = datos
			cache.Entradas[i].Modificado = true
			cache.Entradas[i].Usado = true
			log.Printf("Se escribio %s en página %d", datos, pagina)
			return nil
		}
	}
	return nil
}

/*LOGS MINIMOS RESTANTES:
CACHE HIT --> Si con la dirFisica encuentra una Pagina
CACHE MISS --> Si con la dirFisica no encuentra una Pagina
CACHE ADD --> Despues de no haber encontrado la pagina, la agrega
*/
