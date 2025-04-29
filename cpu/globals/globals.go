package globals

import (
	"errors"
	"strings"
	"sync"
)

// No se si es correcto crear una carpeta globals
type Config struct {
	IpSelf     string `json:"ip_self"`
	PortSelf   int    `json:"port_self"`
	IpMemory   string `json:"ip_memory"`
	PortMemory int    `json:"port_memory"`
	IpKernel   string `json:"ip_kernel"`
	PortKernel int    `json:"port_kernel"`
	LogLevel   string `json:"log_level"`
}

type ExecutionContext struct {
	PID int    // ID del proceso
	PC  int    // Program Counter
	Ax  uint32 // unsigned int de 32 bits
	Bx  uint32
	Cx  uint32
	Dx  uint32
	Ex  uint32
	Fx  uint32
	Gx  uint32
	Hx  uint32
}

var CurrentContext *ExecutionContext
var ClientConfig *Config
var InterrupcionPendiente bool
var PIDInterrumpido int
var MutexInterrupcion sync.Mutex

// (ectx *ExecutionContext) significa que estoy trabajando sobre la struct original y no sobre una copia, GetRegister pasa a ser un metodo
func (ectx *ExecutionContext) GetRegister(str string) (*uint32, error) {
	str = strings.ToLower(str)
	switch str {
	case "ax":
		return &ectx.Ax, nil
	case "bx":
		return &ectx.Bx, nil
	case "cx":
		return &ectx.Cx, nil
	case "dx":
		return &ectx.Dx, nil
	case "ex":
		return &ectx.Ex, nil
	case "fx":
		return &ectx.Fx, nil
	case "gx":
		return &ectx.Gx, nil
	case "hx":
		return &ectx.Hx, nil
	default:
		return nil, errors.New("'" + str + "' no constituye ningún registro conocido")
	}
}
