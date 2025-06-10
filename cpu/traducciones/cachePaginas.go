package traducciones

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"log"
	"sync"
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
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for i := range c.Entradas {
		if c.Entradas[i].NroPagina == nroPagina {
			c.Entradas[i].Usado = true
			log.Printf("PID: %d - Cache HIT - Página: %d", globals.PIDActual, nroPagina)
			return c.Entradas[i].Contenido, true
		}
	}
	log.Printf("PID: %d - Cache MISS - Página: %d", globals.PIDActual, nroPagina)
	return "", false
}

func (c *CachePaginas) EstaActiva() bool {
	return c != nil && c.MaxEntradas > 0
}

func (c *CachePaginas) reemplazarEntrada(nueva EntradaCache) {
	//TODO
}

func (c *CachePaginas) MarcarUso(nroPagina int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for i := range c.Entradas {
		if c.Entradas[i].NroPagina == nroPagina {
			c.Entradas[i].Usado = true
			log.Printf("PID: %d - Cache USO - Página %d marcada con bit de uso en true", globals.PIDActual, nroPagina)
			return
		}
	}
}

func LeerEnCache(nroPagina int, tamanio int) (string, error) {
	contenido, ok := Cache.Buscar(nroPagina)
	if !ok {
		err := fmt.Errorf("Página %d no encontrada en la caché", nroPagina)
		log.Printf("Error: %v", err)
		return "", err
	}
	Cache.MarcarUso(nroPagina)
	// Si querés leer solo parte del contenido: contenido = contenido[:tamanio]
	return contenido, nil
}

func EscribirEnCache(nroPagina int, datos string) error {
	for i := range Cache.Entradas {
		if Cache.Entradas[i].NroPagina == nroPagina {
			Cache.Entradas[i].Contenido = datos
			Cache.Entradas[i].Modificado = true
			Cache.Entradas[i].Usado = true
			log.Printf("Se escribió %s en página %d", datos, nroPagina)
			return nil
		}
	}
	err := fmt.Errorf("No se encontró la página %d en la caché", nroPagina)
	log.Printf("Error: %v", err)
	return err
}

/*LOGS MINIMOS RESTANTES:
Cache HIT --> Si con la dirFisica encuentra una Pagina
CACHE MISS --> Si con la dirFisica no encuentra una Pagina
CACHE ADD --> Despues de no haber encontrado la pagina, la agrega
*/
