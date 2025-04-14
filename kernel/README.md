#KERNEL 

## 🔌 1. Endpoint expuesto

El Kernel escucha conexiones entrantes desde otros módulos en:

`http://localhost:8081/kernel/mensaje`

## 📬 2. Formato del mensaje recibido

El cuerpo del mensaje (`body`) debe ser un JSON con esta estructura:

```json
{
  "ip": "127.0.0.1",
  "puerto": 8000
}
```

Este mensaje se decodifica en un struct de Go como el siguiente:

```go
type Mensaje struct {
    //Nombre string `json:"nombre"` // (opcional)
    Ip     string `json:"ip"`
    Puerto int    `json:"puerto"`
}
```

## 3. Estructura

kernel/
├── main.go
├── utils/
│   └── utils.go
├── config.json
├── go.mod
├── kernel.go
└── README.md 
