package globals

//No se si es correcto crear una carpeta globals
type Config struct {
	IpSelf     string `json:"ip_self"`
	PortSelf   int    `json:"port_self"`
	IpMemory   string `json:"ip_memory"`
	PortMemory int    `json:"port_memory"`
	IpKernel   string `json:"ip_kernel"`
	PortKernel int    `json:"port_kernel"`
	LogLevel   string `json:"log_level"`
}

var ClientConfig *Config
