package administracion_tests

/*
func TestAsignarNumeroEntradaPagina(t *testing.T) {
	g.FramesLibres = []bool{true, false, true}
	g.CantidadFramesLibres = 2
	frame := adm.AsignarNumeroEntradaPagina()
	assert.Equal(t, 0, frame)
}

func TestAsignarDatosAPaginacion(t *testing.T) {
	proceso := &g.Proceso{PID: 1, TablaRaiz: adm.InicializarTablaRaiz()}
	g.MemoryConfig.PagSize = 4
	g.MemoryConfig.MemorySize = 128
	g.FramesLibres = make([]bool, 32)
	for i := range g.FramesLibres {
		g.FramesLibres[i] = true
	}
	g.MemoriaPrincipal = make([]byte, 128)
	datos := []byte("abcdwxyz")
	err := adm.AsignarDatosAPaginacion(proceso, datos)
	assert.Nil(t, err)
}

func TestEscribirEspacioEntrada(t *testing.T) {
	g.MemoriaPrincipal = make([]byte, 8)
	g.MemoryConfig.PagSize = 4
	g.MemoryConfig.MemorySize = 8
	proceso := &g.Proceso{PID: 1}
	g.ProcesosPorPID[1] = proceso
	exito := adm.EscribirEspacioEntrada(1, 0, "abcd")
	assert.Equal(t, "Proceso fue modificado correctamente en memoria", exito.Mensaje)
}
*/
