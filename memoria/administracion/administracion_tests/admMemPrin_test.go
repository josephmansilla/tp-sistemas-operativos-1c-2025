package administracion_tests

/* import (
	adm "github.com/sisoputnfrba/tp-golang/memoria/administracion"
	g "github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTieneTamanioNecesario(t *testing.T) {
	g.MemoryConfig.MemorySize = 1024
	g.MemoryConfig.PagSize = 128
	g.CantidadFramesLibres = 4 // 4 * 128 = 512

	tamanio := 500
	ok := adm.TieneTamanioNecesario(tamanio)
	assert.True(t, ok)

	tamanio = 700
	ok = adm.TieneTamanioNecesario(tamanio)
	assert.False(t, ok)
}

func TestObtenerDatosMemoria(t *testing.T) {
	g.MemoryConfig.PagSize = 4
	g.MemoryConfig.MemorySize = 16
	g.MemoriaPrincipal = make([]byte, g.MemoryConfig.MemorySize)
	copy(g.MemoriaPrincipal[4:8], []byte("hola"))
	datos := adm.ObtenerDatosMemoria(4)
	assert.Equal(t, "hola", datos.PseudoCodigo)
}
*/
