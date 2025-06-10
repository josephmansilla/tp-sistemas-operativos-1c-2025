package traducciones

import (
	"container/list"
	"github.com/sisoputnfrba/tp-golang/cpu/globals"
	"log"
	"sync"
	"time"
)

type EntradaTLB struct {
	NroPagina    int
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
	if algoritmo != "FIFO" && algoritmo != "LRU" {
		log.Printf("Algoritmo TLB inválido: %s", algoritmo)
	}
	return &TLB{
		entradas:    make(map[int]*list.Element),
		orden:       list.New(),
		maxEntradas: maxEntradas,
		algoritmo:   algoritmo,
	}
}

// busco la pagina y devuelvo el marco, si esta es un hit
// (tlb *TLB) para trabajar y modificar directamente la TLB y no una copia
func (tlb *TLB) Buscar(nroPagina int) (int, bool) {
	tlb.mutex.Lock()
	defer tlb.mutex.Unlock()

	if elem, ok := tlb.entradas[nroPagina]; ok {
		if tlb.algoritmo == "LRU" {
			tlb.orden.MoveToFront(elem)
		}
		log.Printf("PID: %d - TLB HIT - Pagina: %d", globals.PIDActual, nroPagina)
		return elem.Value.(EntradaTLB).Marco, true
	}
	log.Printf("PID: %d - TLB MISS - Pagina: %d", globals.PIDActual, nroPagina)
	return -1, false
}

func (tlb *TLB) AgregarEntrada(nroPagina int, marco int) {
	tlb.mutex.Lock()
	defer tlb.mutex.Unlock()

	// Si ya existe, actualizar
	if elem, ok := tlb.entradas[nroPagina]; ok {
		elem.Value = EntradaTLB{NroPagina: nroPagina, Marco: marco}
		if tlb.algoritmo == "LRU" {
			tlb.orden.MoveToFront(elem)
		}
		return
	}

	// Si está llena, elegir víctima
	if len(tlb.entradas) >= tlb.maxEntradas {
		var victima *list.Element
		victima = tlb.orden.Back() //Para los dos algoritmos voy a borrar siempre el ultimo elemento

		if victima != nil {
			entrada := victima.Value.(EntradaTLB)
			delete(tlb.entradas, entrada.NroPagina)
			tlb.orden.Remove(victima)
		}
	}

	// Agregar nueva entrada al frente
	nuevaEntrada := EntradaTLB{NroPagina: nroPagina, Marco: marco}
	elem := tlb.orden.PushFront(nuevaEntrada)
	tlb.entradas[nroPagina] = elem
}
