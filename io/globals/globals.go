package globals

type Config struct {
	PortIo     int    `json:"port_io"`
	IpIo       string `json:"ip_io"`
	IpKernel   string `json:"ip_kernel"`
	PortKernel int    `json:"port_kernel"`
	LogLevel   string `json:"log_level"`
}

var ClientConfig *Config
