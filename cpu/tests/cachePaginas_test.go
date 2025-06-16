package tests

import (
	"github.com/sisoputnfrba/tp-golang/cpu/traducciones"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuscarPagina(t *testing.T) {
	err := CargarConfigCPU("../configs/cpu_3config.json")
	if err != nil {
		t.Fatalf("Error cargando config: %v", err)
	}
	traducciones.Max = 4

	cache := traducciones.NuevaCachePaginas()
	cache.Agregar(1, "Test", true)

	contenido, ok := cache.Buscar(1)

	assert.True(t, ok)
	assert.Equal(t, "Test", contenido)
}

func TestActivacionCache(t *testing.T) {
	err := CargarConfigCPU("../configs/cpu_3config.json")
	if err != nil {
		t.Fatalf("Error cargando config: %v", err)
	}
	traducciones.Max = 4

	cache := traducciones.NuevaCachePaginas()

	bool := cache.EstaActiva()

	assert.True(t, bool)
}

func TestMarcarUso(t *testing.T) {
	//TODO
}

func TestLeerCache(t *testing.T) {
	//TODO
}

func TestEscribirCache(t *testing.T) {
	//TODO
}

func TestReemplazoClock(t *testing.T) {
	//TODO
}

func TestReemplazoClockM(t *testing.T) {
	//TODO
}

func TestLimpiarCache(t *testing.T) {
	//TODO
}

func TestAgregarEntrada(t *testing.T) {
	//TODO
}

func TestEliminarEntrada(t *testing.T) {
	//TODO
}
