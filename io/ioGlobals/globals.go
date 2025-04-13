package globals

type Config struct {
	IpSelf     string `json:"ip_self"`
	PortSelf   int    `json:"port_self"`
	IpKernel   string `json:"ip_kernel"`
	PortKernel int    `json:"port_kernel"`
	LogLevel   string `json:"log_level"`
}

var ClientConfig *Config
