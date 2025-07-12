package globals

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	IpKernel   string `json:"ip_kernel"`
	PortKernel int    `json:"puerto_kernel"`
	IpIo       string `json:"ip_io"`
	PuertoIo1  int    `json:"puerto_io1"`
	PuertoIo2  int    `json:"puerto_io2"`
	PuertoIo3  int    `json:"puerto_io3"`
	PuertoIo4  int    `json:"puerto_io4"`
	LogLevel   string `json:"log_level"`
	Type       string `json:"type"`
}

var IoConfig *Config
var Nombre string
var Puerto int

func (cfg Config) Validate() error {
	if cfg.IpKernel == "" {
		return errors.New("falta el campo 'ip_kernel'")
	}
	if cfg.PortKernel <= 0 {
		return errors.New("falta el campo 'port_kernel' o es inválido")
	}
	if cfg.PuertoIo1 <= 0 {
		return errors.New("falta el campo 'port_io' o es inválido")
	}
	if cfg.PuertoIo2 <= 0 {
		return errors.New("falta el campo 'port_io' o es inválido")
	}
	if cfg.PuertoIo3 <= 0 {
		return errors.New("falta el campo 'port_io' o es inválido")
	}
	if cfg.PuertoIo4 <= 0 {
		return errors.New("falta el campo 'port_io' o es inválido")
	}
	if cfg.IpIo == "" {
		return errors.New("falta el campo 'ip_io'")
	}
	if cfg.LogLevel == "" {
		return errors.New("falta el campo 'log_level'")
	}
	if cfg.Type == "" {
		return errors.New("falta el campo 'type'")
	}
	return nil
}

func CargarConfig() *Config {
	configFile, err := os.Open("../config.json")
	if err != nil {
		fmt.Printf("No se pudo abrir config: %v\n", err)
		os.Exit(1)
	}
	defer configFile.Close()

	var cfg Config
	decoder := json.NewDecoder(configFile)
	if err := decoder.Decode(&cfg); err != nil {
		fmt.Printf("Error parseando config: %v\n", err)
		os.Exit(1)
	}
	if err := cfg.Validate(); err != nil {
		fmt.Printf("Config inválida: %v\n", err)
		os.Exit(1)
	}
	return &cfg
}

func (cfg *Config) PuertoPorNombre(nombre string) {
	n := strings.ToUpper(nombre)
	switch n {
	case "DISCO1":
		Puerto = cfg.PuertoIo1
	case "DISCO2":
		Puerto = cfg.PuertoIo2
	case "DISCO3":
		Puerto = cfg.PuertoIo3
	case "DISCO4":
		Puerto = cfg.PuertoIo4
	default:
		Puerto = -1
	}
}
