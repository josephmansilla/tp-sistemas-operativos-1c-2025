package globals

import "errors"

type Config struct {
	IpKernel   string `json:"ip_kernel"`
	PortKernel int    `json:"port_kernel"`
	PortIo     int    `json:"port_io"`
	IpIo       string `json:"ip_io"`
	LogLevel   string `json:"log_level"`
}

var IoConfig *Config


func (cfg Config) Validate() error {
	if cfg.IpKernel == "" {
		return errors.New("falta el campo 'ip_kernel'")
	}
	if cfg.PortKernel <= 0 {
		return errors.New("falta el campo 'port_kernel' o es inválido")
	}
	if cfg.PortIo <= 0 {
		return errors.New("falta el campo 'port_io' o es inválido")
	}
	if cfg.IpIo == "" {
		return errors.New("falta el campo 'ip_io'")
	}
	if cfg.LogLevel == "" {
		return errors.New("falta el campo 'log_level'")
	}
	return nil
}
