package pcb

import (
	"fmt"
	"time"
)

// posibles estados de un proceso
const (
	EstadoNew         = "NEW"
	EstadoReady       = "READY"
	EstadoExecute     = "EXECUTE"
	EstadoBlocked     = "BLOCKED"
	EstadoExit        = "EXIT"
	EstadoSuspBlocked = "SUSP_BLOCKED"
	EstadoSuspReady   = "SUSP_READY"
)

// /CHICOS ES NECESARIO AGREGAR AL PCB EL TAMAÑO Y NOMBRE DE ARCHIVO DE PSEUDOCODIGO TANTO PARA PLANIFICADOR DE LARGO PLAZO
// (cuando termina un proceso hay que preguntar si el pcb de NEW puede inicilizar)Y PARA SJF
type PCB struct {
	PID            int
	PC             int
	ME             map[string]int
	MT             map[string]float64 // Tiempo en milisegundos con decimales
	FileName       string             // nombre de archivo de pseudoCodigo
	ProcessSize    int
	EstimadoRafaga float64 // Para SJF/SRT
	RafagaRestante int     // Para SRT
	Estado         string
	TiempoEstado   time.Time
}

func (a *PCB) Null() *PCB {
	return nil
}

func (a *PCB) Equal(b *PCB) bool {
	return a.PID == b.PID
}

type Pid int
type RequestToMemory struct {
	Thread    Pid      `json:"pid"`
	Type      string   `json:"type"` //aca le indico el el json que tipo de request es por ejemplo creacionDeProceso
	Arguments []string `json:"arguments"`
}

//Ej ME: "ready": 3 → el proceso estuvo 3 veces en el estado listo.
//Ej MT: "execute": 12 → el proceso estuvo 12 unidades de tiempo en ejecución.

// ImprimirMetricas devuelve un string con las métricas de estado del proceso en el formato:
// ## (<PID>) - Métricas de estado: NEW (COUNT) (TIME), READY (COUNT) (TIME), ...
func (p *PCB) ImprimirMetricas() string {
	estados := []string{
		EstadoNew, EstadoReady, EstadoExecute, EstadoBlocked,
		EstadoExit, EstadoSuspBlocked, EstadoSuspReady,
	}

	salida := fmt.Sprintf("## (<%d>) - Métricas de estado:", p.PID)

	for _, estado := range estados {
		count := p.ME[estado]
		tiempo := p.MT[estado]
		salida += fmt.Sprintf(" %s (%d) (%.2f ms),", estado, count, tiempo)
	}

	// Eliminar la última coma
	if len(salida) > 0 {
		salida = salida[:len(salida)-1]
	}
	return salida
}

func CambiarEstado(p *PCB, nuevoEstado string) {
	estadoAnterior := p.Estado
	FinalizarEstado(p, estadoAnterior) //medir ANTES

	p.ME[nuevoEstado]++
	p.Estado = nuevoEstado
	p.TiempoEstado = time.Now()
}

func FinalizarEstado(p *PCB, estadoAnterior string) {
	duracion := time.Since(p.TiempoEstado) //p.TiempoEnEstado()
	ms := float64(duracion.Microseconds()) / 1000.0
	p.MT[estadoAnterior] += ms
}
