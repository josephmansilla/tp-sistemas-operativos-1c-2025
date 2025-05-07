package globals

import "errors"

type KernelConfig struct {
	MemoryAddress         string  `json:"ip_memory"`
	MemoryPort            int     `json:"port_memory"`
	SchedulerAlgorithm    string  `json:"scheduler_algorithm"`
	ReadyIngressAlgorithm string  `json:"ready_ingress_algorithm"`
	Alpha                 float64 `json:"alpha"`
	InitialEstimate       int     `json:"initial_estimate"`
	SuspensionTime        int     `json:"suspension_time"`
	LogLevel              string  `json:"log_level"`
	KernelPort            int     `json:"port_kernel"`
	KernelAddress         string  `json:"ip_kernel"`
}

func (cfg KernelConfig) Validate() error {
	if cfg.MemoryAddress == "" {
		return errors.New("falta el campo 'ip_memory'")
	}
	if cfg.MemoryPort == 0 {
		return errors.New("falta el campo 'port_memory' o es inválido")
	}
	if cfg.SchedulerAlgorithm == "" {
		return errors.New("falta el campo 'scheduler_algorithm'")
	}
	if cfg.ReadyIngressAlgorithm == "" {
		return errors.New("falta el campo 'ready_ingress_algorithm'")
	}
	if cfg.Alpha <= 0 {
		return errors.New("falta el campo 'alpha' o es inválido")
	}
	if cfg.InitialEstimate <= 0 {
		return errors.New("falta el campo 'initial_estimate' o es inválido")
	}
	if cfg.SuspensionTime < 0 {
		return errors.New("falta el campo 'suspension_time' o es inválido")
	}
	if cfg.LogLevel == "" {
		return errors.New("falta el campo 'log_level'")
	}
	if cfg.KernelPort == 0 {
		return errors.New("falta el campo 'port_kernel' o es inválido")
	}
	if cfg.KernelAddress == "" {
		return errors.New("falta el campo 'ip_kernel'")
	}
	return nil
}
