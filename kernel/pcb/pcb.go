package pcb

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
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

// ES NECESARIO AGREGAR AL PCB EL TAMAÑO Y NOMBRE DE ARCHIVO DE PSEUDOCODIGO PARA PLANIFICADOR DE LARGO PLAZO
// (cuando termina un proceso hay que preguntar si el pcb de NEW puede inicilizar) Y PARA SJF
type PCB struct {
	PID            int
	PC             int
	ME             map[string]int
	MT             map[string]float64 // Tiempo en milisegundos con decimales por cada estado
	FileName       string             // nombre de archivo de pseudoCodigo
	ProcessSize    int                //tamaño en memoria
	EstimadoRafaga float64            // Para SJF/SRT
	Estado         string             //Estado actual
	TiempoEstado   time.Time          //Saber cuanto estuvo en un estado reciente
	CpuID          string             //conocer CpuID
}

//Ej ME: "ready": 3 → el proceso estuvo 3 veces en el estado listo.
//Ej MT: "execute": 12 → el proceso estuvo 12 unidades de tiempo en ejecución.

// Saber si son nulos
func (a *PCB) Null() *PCB {
	return nil
}

// Comparar pcbs
func (a *PCB) Equal(b *PCB) bool {
	return a.PID == b.PID
}

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

	logger.Info(salida)
	return salida
}

// Cambia el estado y actualiza metricas
func CambiarEstado(p *PCB, nuevoEstado string) {
	estadoAnterior := p.Estado

	if estadoAnterior == EstadoExecute {
		duracion := time.Since(p.TiempoEstado)
		rafagaReal := float64(duracion.Microseconds())
		//Actualizar rafaga real de CPU si viene de Execute
		ActualizarEstimacionRafaga(p, rafagaReal)
	}

	FinalizarEstado(p, estadoAnterior)

	p.ME[nuevoEstado]++
	p.Estado = nuevoEstado
	p.TiempoEstado = time.Now()
}

// Calcula antes de irse el tiempo que estuvo en ese estado
func FinalizarEstado(p *PCB, estadoAnterior string) {
	duracion := time.Since(p.TiempoEstado) //p.TiempoEnEstado()
	ms := float64(duracion.Microseconds())
	p.MT[estadoAnterior] += ms
}

// Utilizar despues de una rafaga en CPU
func ActualizarEstimacionRafaga(proceso *PCB, rafagaReal float64) float64 {
	alpha := globals.KConfig.Alpha
	proceso.EstimadoRafaga = alpha*rafagaReal + (1-alpha)*proceso.EstimadoRafaga
	return proceso.EstimadoRafaga
}

//EJEMPLO DE USO
/*
cuando termina una ráfaga en CPU:
ActualizarEstimacionRafaga(proceso, 7.001) // 7.001ms es el tiempo real que tardó la ráfaga
*/

// Esta función no se recomienda: usa el tiempo total acumulado en EXECUTE y no solo la última ráfaga.
func CalcularEstimacionRafaga(proceso *PCB) float64 {
	alpha := globals.KConfig.Alpha
	proceso.EstimadoRafaga = alpha*float64(proceso.MT[EstadoExecute]) + (1-alpha)*proceso.EstimadoRafaga
	return proceso.EstimadoRafaga
}
