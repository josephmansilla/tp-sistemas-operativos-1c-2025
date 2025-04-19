package globals

type Config struct {
	IpMemory              string `json:"ip_memory"`
	PortMemory            int    `json:"port_memory"`
	IpKernel              string `json:"ip_kernel"`
	PortKernel            int    `json:"port_kernel"`
	SchedulerAlgorithm    string `json:"scheduler_algorithm"`
	ReadyIngressAlgorithm string `json:"ready_ingress_algorithm"`
	Alpha                 string `json:"alpha"`
	SuspensionTime        int    `json:"suspension_time"`
	LogLevel              string `json:"log_level"`
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

//crear nodos a punteros PCB, para instanciarlas en main (sera un puntero a pcb creado previamente)
