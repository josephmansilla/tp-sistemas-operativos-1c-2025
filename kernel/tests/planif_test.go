package tests

import (
	"testing"
	"time"

	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/pcb"
	"github.com/stretchr/testify/assert"
)

func TestCambiarEstadoYMetricas(t *testing.T) {
	globals.Config.InitialEstimate = 1000
	globals.Config.Alpha = 0.5

	p := &pcb.PCB{
		PID:            1,
		PC:             0,
		ME:             make(map[string]int),
		MT:             make(map[string]float64),
		EstimadoRafaga: globals.Config.InitialEstimate,
		FileName:       "test.txt",
		ProcessSize:    5,
		TiempoEstado:   time.Now(), //time desde que se crea
	}

	// Simular que el proceso inicia y entra a NEW
	pcb.CambiarEstado(p, pcb.EstadoNew)
	time.Sleep(10 * time.Millisecond)

	// Luego cambia a READY
	pcb.CambiarEstado(p, pcb.EstadoReady)
	time.Sleep(10 * time.Millisecond)

	// Luego a EXIT
	pcb.CambiarEstado(p, pcb.EstadoExit)

	// Verificar que haya contado los pasos
	assert.Equal(t, 1, p.ME[pcb.EstadoNew], "Debe haber 1 paso por NEW")
	assert.Equal(t, 1, p.ME[pcb.EstadoReady], "Debe haber 1 paso por READY")
	assert.Equal(t, 1, p.ME[pcb.EstadoExit], "Debe haber 1 paso por EXIT")

	// Verificar que los tiempos sean positivos
	assert.Greater(t, p.MT[pcb.EstadoNew], 1.0, "NEW debe tener tiempo > 0")
	assert.Greater(t, p.MT[pcb.EstadoReady], 1.0, "READY debe tener tiempo > 0")

	// Verificar la impresión
	metricas := p.ImprimirMetricas()
	t.Log(metricas)
	assert.Contains(t, metricas, "NEW (1)")
	assert.Contains(t, metricas, "READY (1)")
	assert.Contains(t, metricas, "EXIT (1)")
}

func TestProcesoAInterrumpir(t *testing.T) {
	now := time.Now()

	p1 := &pcb.PCB{
		PID:            1,
		EstimadoRafaga: 20000,
		TiempoEstado:   now.Add(-5 * time.Millisecond), // lleva 5ms ejecutando
	}
	p2 := &pcb.PCB{
		PID:            2,
		EstimadoRafaga: 3000,
		TiempoEstado:   now,
	}
	procesoEntrante := &pcb.PCB{
		PID:            3,
		EstimadoRafaga: 1000, // menor que el tiempo restante de p1
	}

	ejecutando := []*pcb.PCB{p1, p2}

	interrumpir := ProcesoAInterrumpir(procesoEntrante, ejecutando)
	if interrumpir == nil {
		t.Error("Se esperaba que se interrumpa un proceso")
	} else if interrumpir.PID != 1 {
		t.Errorf("Se esperaba que se interrumpa proceso PID 1, pero fue %d", interrumpir.PID)
	}
}

func ProcesoAInterrumpir(procesoEntrante *pcb.PCB, ejecutando []*pcb.PCB) *pcb.PCB {
	mayorTiempoRestante := -1.0
	var procesoAInterrumpir *pcb.PCB

	now := time.Now()

	for _, p := range ejecutando {
		tiempoEjecutado := float64(now.Sub(p.TiempoEstado).Microseconds())
		tiempoRestante := p.EstimadoRafaga - tiempoEjecutado

		if tiempoRestante > mayorTiempoRestante {
			mayorTiempoRestante = tiempoRestante
			procesoAInterrumpir = p
		}
	}

	if procesoAInterrumpir != nil && procesoEntrante.EstimadoRafaga < mayorTiempoRestante {
		return procesoAInterrumpir
	}

	return nil
}
