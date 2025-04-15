## MEMORIA



## 🔌 1. Endpoint expuesto

La memoria se encarga de tener conexiones entrantes para los módulos:

`http://localhost:8083/memoria/kernel`
`http://localhost:8083/memoria/cpu`

## 📬 2. Formato del mensaje recibido

El cuerpo del mensaje (`body`) debe ser un JSON con una estructura dependiendo de cada Modulo:

```json
//CPU
{
  "ip": "127.0.0.1",
  "puerto": 8000
}

//IO
{
  "nombre":"impresora",
  "ip": "127.0.0.1",
  "puerto": 8000
}
```

Estos mensajes se decodifican en un struct de GO como los siguientes:

```go
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

// Datos recibidos
type DatosDeKernel struct {
	TamanioMemoria int `json:"tamanio_memoria"` 
	// Placeholder quedaría ver que tipos de datos podría recibir de kernel
}

// Tipo de datos recibidos de la CPU

type DatosDeCPU struct {
	PID int `json:"pid"`
	PC  int `json:"pc"`
}
```