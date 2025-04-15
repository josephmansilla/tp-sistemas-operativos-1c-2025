package globals

type Config struct {
	PortMemory     int    `json:"port_memory"`
	MemorySize     int    `json:"memory_size"`
	PagSize        int    `json:"pag_size"`
	EntriesPerPage int    `json:"entries_per_page"`
	NumberOfLevels int    `json:"number_of_levels"`
	MemoryDelay    int    `json:"memory_delay"`
	SwapfilePath   string `json:"swapfile_path"`
	SwapDelay      int    `json:"swap_delay"`
	LogLevel       string `json:"log_level"`
}

var MemoryConfig *Config

// Tipo de datos recibidos de1 Kernel

type DatosDeKernel struct {
	PID            int `json:"pid"`
	TamanioMemoria int `json:"tamanio_memoria"` // Placeholder
	// toDO....
}

// Tipo de datos recibidos de la CPU

type DatosDeCPU struct {
	PID int `json:"pid"`
	PC  int `json:"pc"`
}

type DatosParaCPU struct {
	// TODO
}

var Kernel DatosDeKernel
var CPU DatosDeCPU
