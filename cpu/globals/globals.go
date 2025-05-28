package globals

import (
	"errors"
	"sync"
)

// No se si es correcto crear una carpeta globals
type Config struct {
	IpSelf         string `json:"ip_self"`
	PortSelf       int    `json:"port_self"`
	IpMemory       string `json:"ip_memory"`
	PortMemory     int    `json:"port_memory"`
	IpKernel       string `json:"ip_kernel"`
	PortKernel     int    `json:"port_kernel"`
	TlbEntries     int    `json:"tlb_entries"`
	TlbReplacement string `json:"tlb_replacement"`
	CacheDelay     int    `json:"cache_delay"`
	LogLevel       string `json:"log_level"`
}

type PCB struct {
	PID            int
	PC             int
	ME             map[string]int //asocia cada estado con la cantidad de veces que el proceso estuvo en ese estado.
	MT             map[string]int //asocia cada estado con el tiempo total que el proceso pasó en ese estado.
	FileName       string         // nombre de archivo de pseudoCodigo
	ProcessSize    int
	EstimadoRafaga float64 // Para SJF/SRT
	RafagaRestante int     // Para SRT
	Estado         string
}

var Pcb *PCB
var ClientConfig *Config
var InterrupcionPendiente bool
var PIDInterrumpido int
var MutexInterrupcion sync.Mutex
var TamPag int
var ErrSyscallBloqueante = errors.New("proceso bloqueado por syscall IO")
var SaltarIncrementoPC bool
var ID string
