## KERNEL 

## FUNCIONALIDAD

1. LEER ARCHIVO DE CONFIGURACION -> utils.Config(filepath)
2. CARGAR SUS DATOS EN GLOBALS -> en el struct Config 
3. LISTEN en los puertos HTTP
4. RECIBIR Y GUARDAR EN GLOBALS info. de algun modulo

## 🔌 1. Endpoint expuesto

El Kernel escucha conexiones entrantes desde otros módulos en:

`http://localhost:8081/kernel/io`
`http://localhost:8081/kernel/cpu`

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

Estos mensajes se decodifican en un struct de Go como los siguientes:

```go
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

// Datos recibidos
type DatosIO struct {
	Nombre string
	Ip     string
	Puerto int
}

type DatosCPU struct {
	Ip     string
	Puerto int
}
```

## 3. Estructura

kernel/ 
├── utils/ # Funciones auxiliares (leer JSON, manejar requests) 
	│ 
	└── utils.go 
├── globals/ 
	│ 
	└── globals.go 
├── config.json # Archivo de configuración 
├── go.mod # Módulo Go 
├── kernel.go # Lógica del módulo Kernel 
└── README.md # Documentación del proyecto
