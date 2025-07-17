package administracion

/*
func SeleccionarEntradas(pid int, direccionFisica int, entradasNecesarias int) (entradas []g.EntradaPagina, err error) {
	paginaInicio := direccionFisica / g.MemoryConfig.PagSize
	err = nil

	g.MutexProcesosPorPID.Lock()
	proceso := g.ProcesosPorPID[pid]
	g.MutexProcesosPorPID.Unlock()

	if proceso == nil {
		return nil, fmt.Errorf("no existe el proceso con el PID: %d ; %v", pid, logger.ErrProcessNil)
	}

	for i := 0; i < entradasNecesarias; i++ {
		numeroPagina := paginaInicio + i
		indices := CrearIndicePara(numeroPagina)

		entrada, err := BuscarEntradaPagina(proceso, indices)
		if err != nil {
			return nil, fmt.Errorf("error al buscar la entrada de pagina: %d de PID %d; %v", i, pid, err)
		}
		if entrada == nil {
			return nil, fmt.Errorf("error al encontrar la entrada de pagina: %d de PID %d; %v", i, pid, logger.ErrNoInstance)
		}
		if !entrada.EstaPresente {
			return nil, fmt.Errorf("error al buscar la entrada de pagina: %d de PID %d; %v", i, pid, logger.ErrNotPresent)
		}
		entradas = append(entradas, *entrada)
	}

	return
}

func LeerEspacioMemoria(pid int, direccionFisica int, tamanioALeer int) (confirmacionLectura g.ExitoLecturaMemoria, err error) {
	confirmacionLectura = g.ExitoLecturaMemoria{Exito: err, ValorLeido: ""}

	entradasNecesarias, err := g.CalcularCantidadEntradas(tamanioALeer)
	if err != nil {
		return confirmacionLectura, err
	}

	entradas, err := SeleccionarEntradas(pid, direccionFisica, entradasNecesarias)
	if err != nil {
		return confirmacionLectura, err
	}

	bytesRestantes := tamanioALeer
	cant := len(entradas)
	var datos []byte

	for i, entrada := range entradas {
		inicioLectura, finLectura, err := LogicaRecorrerMemoria(i, cant, entrada, direccionFisica, bytesRestantes)
		if err != nil {
			return confirmacionLectura, err
		}

		g.MutexMemoriaPrincipal.Lock()
		datos = append(datos, g.MemoriaPrincipal[inicioLectura:finLectura]...)
		g.MutexMemoriaPrincipal.Unlock()

		logger.Debug("Datos en bytes leidos: %d", datos)

		bytesRestantes -= finLectura - inicioLectura
		if bytesRestantes <= 0 {
			break
		}
	}
	return g.ExitoLecturaMemoria{Exito: nil, ValorLeido: string(datos)}, nil
}

func LogicaRecorrerMemoria(i int, cantEntradas int, entrada g.EntradaPagina, dirF int, bytesRestantes int) (inicio int, limite int, err error) {
	tamTotal := g.MemoryConfig.MemorySize
	tamPag := g.MemoryConfig.PagSize
	offsetLogico := dirF % tamPag

	base := entrada.NumeroPagina * tamPag

	if i == 0 { // PARA EL PRIMER FRAME
		var delta int
		if bytesRestantes <= tamPag { // EN CASO DE QUE LOS BYTESRESTANTES SEAN MENORES A LA ENTRADA
			delta = bytesRestantes
		} else {
			delta = tamPag
		}
		inicio = base + offsetLogico // BASE DE LA ENTRADA + POSIBLE DESPLAZAMIENTO
		limite = min(base+tamPag, inicio+delta)
		// ELIJO ENTRE EL LIMITE DE LA ENTRADA O EN CASO DE QUE SE CUMPLA EL PRIMER IF: HASTA UN LIMITE MENOR
		// QUE EL LIMITE DE LA ENTRADA
	} else if i == cantEntradas-1 { // PARA EL ÙLTIMO FRAME
		inicio = base
		limite = inicio + bytesRestantes
	} else { // PARA LOS CASOS INTERMEDIOS: NO DEBERÌAN ENTRAR SI LOS BYTES RESANTE ESTÀN CORRECTAMENTE CALCULADOS
		inicio = base
		limite = base + tamPag
	}

	if inicio >= tamTotal || limite > tamTotal {
		err = logger.ErrSegmentFault
		return 0, 0, err
	}
	logger.Debug("Inicio de lectura: %d", inicio)
	logger.Debug("Limite de lectura: %d", limite)
	return inicio, limite, nil
}

func EscribirEspacioMemoria(pid int, direccionFisica int, tamanioALeer int, datosAEscribir string) (confirmacionEscritura g.ExitoEdicionMemoria, err error) {
	confirmacionEscritura = g.ExitoEdicionMemoria{Exito: err, Booleano: false}

	entradasNecesarias, err := g.CalcularCantidadEntradas(tamanioALeer)
	if err != nil {
		return confirmacionEscritura, err
	}

	entradas, err := SeleccionarEntradas(pid, direccionFisica, entradasNecesarias)
	if err != nil {
		return confirmacionEscritura, err
	}

	direccionParaAcceder := direccionFisica
	datosEnBytes := []byte(datosAEscribir)
	cant := len(entradas)
	bytesRestantes := tamanioALeer
	bytesEscritos := 0

	for i, entrada := range entradas {
		inicioEscritura, finEscritura, err := LogicaRecorrerMemoria(i, cant, entrada, direccionParaAcceder, bytesRestantes)
		cantBytes := finEscritura - inicioEscritura
		if err != nil {
			return confirmacionEscritura, err
		}

		inicioDatos := bytesEscritos
		finDatos := inicioDatos + cantBytes

		g.MutexMemoriaPrincipal.Lock()
		copy(g.MemoriaPrincipal[inicioEscritura:finEscritura], datosEnBytes[inicioDatos:finDatos])
		g.MutexMemoriaPrincipal.Unlock()

		bytesEscritos += cantBytes
		bytesRestantes -= cantBytes
		if bytesRestantes <= 0 {
			break
		}
		direccionParaAcceder += cantBytes
	}
	return g.ExitoEdicionMemoria{Exito: err, Booleano: true}, nil

}
*/
