package globals

import "errors"

type MemoriaConfig struct {
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
	DumpPath       string `json:"dump_path"`
	ScriptsPath    string `json:"scripts_path"`
}

func (cfg MemoriaConfig) Validate() error {
	if cfg.IpMemory == "" {
		return errors.New("falta el campo 'ip_memory'")
	}
	if cfg.PortMemory <= 0 {
		return errors.New("falta el campo 'port_memory' o es inválido")
	}
	if cfg.MemorySize <= 0 {
		return errors.New("falta el campo 'MemorySize'")
	}
	if cfg.PagSize <= 0 {
		return errors.New("falta el campo 'Pagsize'")
	}
	if cfg.EntriesPerPage <= 0 {
		return errors.New("falta el campo 'EntriesPerPage' o es inválido")
	}
	if cfg.NumberOfLevels <= 0 {
		return errors.New("falta el campo 'NumberOfLevels' o es inválido")
	}
	if cfg.MemoryDelay <= 0 {
		return errors.New("falta el campo 'MemoryDelay' o es inválido")
	}
	if cfg.SwapfilePath == "" {
		return errors.New("falta el campo 'SwapfilePath'")
	}
	if cfg.SwapDelay <= 0 {
		return errors.New("falta el campo 'SwapDelay' o es inválido")
	}
	if cfg.LogLevel == "" {
		return errors.New("falta el campo 'LogLevel'")
	}
	if cfg.DumpPath == "" {
		return errors.New("falta el campo 'DumpPath'")
	}
	if cfg.ScriptsPath == "" {
		return errors.New("falta el campo 'ScriptsPath'")
	}
	return nil
}
