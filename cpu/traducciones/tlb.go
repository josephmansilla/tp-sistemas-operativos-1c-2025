package traducciones

import (
	"container/list"
	"log"
	"sync"
	"time"
)

type EntradaTLB struct {
	Pagina       int
	Marco        int
	UltimoAcceso time.Time
}

type TLB struct {
	entradas    map[int]*list.Element // clave: página
	orden       *list.List            // para FIFO(Back()--> Elimino el ultimo)/LRU(MoveToFront())
	maxEntradas int
	algoritmo   string
	mutex       sync.Mutex
}

func NuevaTLB(maxEntradas int, algoritmo string) *TLB {
	return &TLB{
		entradas:    make(map[int]*list.Element),
		orden:       list.New(),
		maxEntradas: maxEntradas,
		algoritmo:   algoritmo,
	}
}

// busco la pagina y devuelvo el marco, si esta es un hit
// (tlb *TLB) para trabajar y modificar directamente la TLB y no una copia
func (tlb *TLB) Buscar(pagina int) (int, bool) {
	tlb.mutex.Lock()
	defer tlb.mutex.Unlock()

	if elem, ok := tlb.entradas[pagina]; ok {
		if tlb.algoritmo == "LRU" {
			tlb.orden.MoveToFront(elem)
		}
		log.Printf("TLB Hit")
		return elem.Value.(EntradaTLB).Marco, true
	}
	log.Printf("TLB Miss")

	return -1, false
}

func (tlb *TLB) AgregarEntrada(pagina int, marco int) {
	tlb.mutex.Lock()
	defer tlb.mutex.Unlock()

	// por si un marco cambia de pag
	if elem, ok := tlb.entradas[pagina]; ok {
		elem.Value = EntradaTLB{Pagina: pagina, Marco: marco}
		if tlb.algoritmo == "LRU" {
			tlb.orden.MoveToFront(elem)
		}
		return
	}

	// si ya esta llena la tlb, reemplazo una entrada
	if len(tlb.entradas) >= tlb.maxEntradas {
		var victima *list.Element
		if tlb.algoritmo == "FIFO" || tlb.algoritmo == "LRU" {
			victima = tlb.orden.Back()
		}

		if victima != nil {
			entrada := victima.Value.(EntradaTLB)
			delete(tlb.entradas, entrada.Pagina)
			tlb.orden.Remove(victima)
		}
	}

	nuevaEntrada := EntradaTLB{Pagina: pagina, Marco: marco}
	elem := tlb.orden.PushFront(nuevaEntrada)
	tlb.entradas[pagina] = elem
}
