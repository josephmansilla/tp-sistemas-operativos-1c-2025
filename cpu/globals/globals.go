package globals

import (
	"errors"
	"sync"
)

// No se si es correcto crear una carpeta globals
type Config struct {
	IpCPU            string `json:"ip_cpu"`
	PortCPU1         int    `json:"puerto_cpu1"`
	PortCPU2         int    `json:"puerto_cpu2"`
	PortCPU3         int    `json:"puerto_cpu3"`
	PortCPU4         int    `json:"puerto_cpu4"`
	IpMemory         string `json:"ip_memoria"`
	PortMemory       int    `json:"puerto_memoria"`
	IpKernel         string `json:"ip_kernel"`
	PortKernel       int    `json:"puerto_kernel"`
	TlbEntries       int    `json:"tlb_entries"`
	TlbReplacement   string `json:"tlb_replacement"`
	CacheEntries     int    `json:"cache_entries"`
	CacheDelay       int    `json:"cache_delay"`
	CacheReplacement string `json:"cache_replacement"`
	LogLevel         string `json:"log_level"`
	IpSelf           string
	PortSelf         int
}

var PIDActual int
var PCActual int
var ClientConfig *Config
var InterrupcionPendiente bool
var PIDInterrumpido int
var MutexInterrupcion sync.Mutex
var ErrSyscallBloqueante = errors.New("proceso bloqueado por syscall IO")
var SaltarIncrementoPC bool
var ID string
var TamanioPagina int
var EntradasPorNivel int
var CantidadNiveles int
