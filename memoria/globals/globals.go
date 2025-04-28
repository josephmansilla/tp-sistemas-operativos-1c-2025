package globals

type Config struct {
	PortMemory     int    `json:"port_memory"`
	IpMemory       string `json:"ip_memory"`
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

type DatosConsultaDeKernel struct {
	PID            int `json:"pid"`
	TamanioMemoria int `json:"tamanio_memoria"`
	// con el tamaño de memoria consulta si es posible ejecutarlo en memoria
}

type DatosRespuestaDeKernel struct {
	Pseudocodigo   string `json:"filename"`
	TamanioMemoria int    `json:"tamanio_memoria"`
}

// Tipo de datos recibidos de la CPU

type DatosDeCPU struct {
	PID int `json:"pid"`
	PC  int `json:"pc"`
}

type DatosParaCPU struct {
	// TODO
}

type ContextoDeCPU struct {
	PID int `json:"pid"`
	PC  int `json:"pc"`
}

type InstruccionCPU struct {
	Instruccion string `json:"instruccion"`
}

type EspacioLibreRTA struct {
	EspacioLibre int `json:"espacio_libre"`
}

var RespuestaKernel DatosRespuestaDeKernel
var Kernel DatosConsultaDeKernel
var CPU DatosDeCPU

// EspacioDeUsuario => make([]byte, TamMemoria)
