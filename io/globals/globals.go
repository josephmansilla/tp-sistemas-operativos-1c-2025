package globals

type Config struct {
	IpKernel   string `json:"ip_kernel"`
	PortKernel int    `json:"port_kernel"`
	PortIo     int    `json:"port_io"`
	IpIo       string `json:"ip_io"`
	LogLevel   string `json:"log_level"`
}

var ClientConfig *Config
