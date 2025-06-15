package planificadores

import (
	"github.com/sisoputnfrba/tp-golang/kernel/algoritmos"
	"github.com/sisoputnfrba/tp-golang/kernel/comunicacion"
	"github.com/sisoputnfrba/tp-golang/kernel/pcb"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"sort"
	"sync"
)

// PoliticaAdmision define el comportamiento de admisión NEW→READY
type PoliticaAdmision interface {
	AlLlegarNuevo(p *pcb.PCB)
	AlFinalizarProceso()
}

var muNEW sync.Mutex

// FIFOAdmision: el primero bloquea a los demás
type FIFOAdmision struct{}

func NewFIFOAdmision() *FIFOAdmision { return &FIFOAdmision{} }

func (f *FIFOAdmision) AlLlegarNuevo(p *pcb.PCB) {
	// si es único en NEW y no hay suspendidos
	if algoritmos.ColaSuspendidoReady.IsEmpty() && algoritmos.ColaNuevo.Size() == 1 {
		f.intentarAdmitir()
	}
}

func (f *FIFOAdmision) AlFinalizarProceso() {
	f.intentarAdmitir()
}

func (f *FIFOAdmision) intentarAdmitir() {
	if algoritmos.ColaSuspendidoReady.IsEmpty() && !algoritmos.ColaNuevo.IsEmpty() {
		muNEW.Lock()
		primero := algoritmos.ColaNuevo.First()
		// aca memoria deberia responder si puede crear el proceso o no, por eso le pasamos el tamaño del proceso
		ok := comunicacion.SolicitarEspacioEnMemoria(primero.FileName, primero.ProcessSize)
		if ok {
			algoritmos.ColaNuevo.Remove(primero)
			algoritmos.ColaReady.Add(primero)
			logger.Info("FIFO: admitido PID %d", primero.PID)
		}
		muNEW.Unlock()
	}
}

// SJFAdmision: siempre trata de meter al más chico posible
type SJFAdmision struct{}

func NewSJFAdmision() *SJFAdmision { return &SJFAdmision{} }

func (s *SJFAdmision) AlLlegarNuevo(_ *pcb.PCB) { s.intentarAdmitir() }
func (s *SJFAdmision) AlFinalizarProceso()      { s.intentarAdmitir() }

func (s *SJFAdmision) intentarAdmitir() {
	if algoritmos.ColaSuspendidoReady.IsEmpty() {
		all := algoritmos.ColaNuevo.Values()
		sort.Slice(all, func(i, j int) bool {
			return all[i].ProcessSize < all[j].ProcessSize
		})
		// aca memoria deberia responder si puede crear el proceso o no, por eso le pasamos el tamaño del proceso
		for _, p := range all {
			ok := comunicacion.SolicitarEspacioEnMemoria(p.FileName, p.ProcessSize)
			if ok {
				muNEW.Lock()
				algoritmos.ColaNuevo.Remove(p)
				muNEW.Unlock()
				algoritmos.ColaReady.Add(p)
				logger.Info("SJF: admitido PID %d", p.PID)
			}
		}
	}
}
