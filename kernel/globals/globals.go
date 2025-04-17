package globals

type Config struct {
	IpMemory           string `json:"ip_memory"`
	PortMemory         int    `json:"port_memory"`
	PortKernel         int    `json:"port_kernel"`
	SchedulerAlgorithm string `json:"scheduler_algorithm"`
	NewAlgorithm       string `json:"new_algorithm"`
	Alpha              string `json:"alpha"`
	SuspensionTime     int    `json:"suspension_time"`
	LogLevel           string `json:"log_level"`
}

var KernelConfig *Config

// Datos recibidos por el Kernel
type DatosIO struct {
	Nombre string
	Ip     string
	Puerto int
}

type DatosCPU struct {
	Ip     string
	Puerto int
	ID     string
}

var CPU DatosCPU
var IO DatosIO
