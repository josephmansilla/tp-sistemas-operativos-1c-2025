package administracion_tests

/*
func TestLiberarMemoriaProceso(t *testing.T) {
	// Preparar un proceso con una entrada presente
	entrada := &g.EntradaPagina{
		NumeroFrame:   0,
		EstaPresente:  true,
		EstaEnUso:     true,
		FueModificado: false,
	}
	tabla := &g.TablaPagina{
		EntradasPaginas: map[int]*g.EntradaPagina{0: entrada},
	}
	proceso := &g.Proceso{PID: 1, TablaRaiz: g.TablaPaginas{0: tabla}}
	g.ProcesosPorPID[1] = proceso
	g.MemoriaPrincipal = make([]byte, 128)
	g.MemoryConfig.PagSize = 128
	g.FramesLibres = make([]bool, 1)
	g.CantidadFramesLibres = 1

	err := adm.LiberarTablaPaginas(tabla, 1)
	assert.Nil(t, err)
}

func TestEliminarDeSlice(t *testing.T) {

	adm.InicializarMemoriaPrincipal()

	proceso := g.Proceso{
		PID:       0,
		TablaRaiz: g.TablaPaginas{},
	}

	adm.OcuparProcesoEnVectorMapeable(proceso.PID, &proceso)

	result, err := adm.DesocuparProcesoEnVectorMapeable(proceso.PID)
	if err != nil {
		t.Error(err)
	}
	assert.Empty(t, result)
}
*/
