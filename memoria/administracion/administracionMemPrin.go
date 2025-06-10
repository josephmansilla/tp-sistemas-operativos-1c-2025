package administracion

import (
	"encoding/json"
	"fmt"
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/data"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"log"
	"net/http"
	"os"
	"time"
)

func InicializarMemoriaPrincipal() {
	cantidadFrames := CalcularTotalFrames()

	globals.MemoriaPrincipal = make([][]byte, cantidadFrames)
	for i := range globals.MemoriaPrincipal {
		globals.MemoriaPrincipal[i] = make([]byte, globals.TamanioMaximoFrame)
	}
	globals.FramesLibres = make([]bool, cantidadFrames)
	ConfigurarFrames(cantidadFrames)

	logger.Info("Tamanio Memoria Principal de %d", globals.MemoryConfig.MemorySize)
	logger.Info("Memoria Principal Inicializada con %d frames de %d cada una.", cantidadFrames, globals.MemoryConfig.PagSize)
}
func CalcularTotalFrames() int {
	tamanioMemoriaPrincipal := globals.MemoryConfig.MemorySize
	tamanioPagina := globals.MemoryConfig.PagSize

	return tamanioMemoriaPrincipal / tamanioPagina
}
func ConfigurarFrames(cantidadFrames int) {
	for i := 0; i < cantidadFrames; i++ {
		globals.FramesLibres[i] = true
	}
	logger.Info("Todos los frames están libres.")
}

func TieneTamanioNecesario(tamanioProceso int) bool {
	var framesNecesarios = float64(tamanioProceso) / float64(globals.TamanioMaximoFrame)
	return framesNecesarios <= float64(globals.CantidadFramesLibres)
}
func LecturaPseudocodigo(archivoPseudocodigo string) []byte {
	string := archivoPseudocodigo
	stringEnBytes := []byte(string)
	return stringEnBytes
}

func ObtenerDatosMemoria(numeroFrame int) globals.ExitoLecturaPagina {

	pseudocodigoEnBytes := globals.MemoriaPrincipal[numeroFrame]
	pseudocodigoEnString := string(pseudocodigoEnBytes)

	datosLectura := globals.ExitoLecturaPagina{
		PseudoCodigo:    pseudocodigoEnString,
		DireccionFisica: numeroFrame,
	}

	return datosLectura
}

// ------------------------------------------------------------------
// ---------- FORMA PARTE DEL ACCESO A ESPACIO DE USUARIO ----------
// ------------------------------------------------------------------

func EscrituraEspacio(w http.ResponseWriter, r *http.Request) {
	var mensaje globals.EscrituraPagina // TODO: defiinir otro tipo
	err := json.NewDecoder(r.Body).Decode(&mensaje)
	if err != nil {
		http.Error(w, "Error leyendo JSON de Kernel\n", http.StatusBadRequest)
		return
	}

	time.Sleep(time.Duration(globals.DelayMemoria) * time.Second)

	tamanioNecesario := mensaje.TamanioNecesario
	pid := mensaje.PID
	datosASobreEscribir := mensaje.DatosASobreEscribir
	indice := mensaje.Indice

	// TODO: CAMBIAR PORQUE DEBEN HABER VARIOS MARCOS PARA EDITAR
	if tamanioNecesario > globals.TamanioMaximoFrame {
		log.Fatal("No se puede cargar en una pagina este tamaño")
		// TODO: FATAL ...
	}
	EscribirEspacioEntrada(pid, indice, datosASobreEscribir)
	// TODO: DIRECCION FISICA
	logger.Info("## PID: <%d> - <%s> - Dir. Física: <%d> - Tamaño: <%d>", pid, datosASobreEscribir, direccionFisica, tamanioNecesario)

	respuesta := globals.ExitoEscrituraPagina{
		Exito:   true,
		Mensaje: "Proceso fue modificado correctamente en memoria",
	}
	if err := json.NewEncoder(w).Encode(respuesta); err != nil {
		logger.Error("Error al serializar mock de espacio: %v", err)
	}

	json.NewEncoder(w).Encode(respuesta)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Respuesta devuelta"))
}

func LecturaEspacio(w http.ResponseWriter, r *http.Request) {

	time.Sleep(time.Duration(globals.DelayMemoria) * time.Second)
	LeerEspacio()

	// TODO: DEBO CREAR UNA STRUCT PARA QUE ME ENVIEN LA DIRECCION FISICA
	// TODO: DEVOLVER EL VALOR QUE SE ENCUENTRA EN LA DIRECCION PEDIDA
	logger.Info("## PID: <PID>  - <Lectura> - Dir. Física: <DIRECCIÓN_FÍSICA> - Tamaño: <TAMAÑO>")
}

func MemoriaDump(w http.ResponseWriter, r *http.Request) {
	var dump globals.DatosParaDump

	if err := data.LeerJson(w, r, &dump); err != nil {
		logger.Error("Error al recibir JSON: %v", err)
		http.Error(w, "Error procesando datos del Kernel", http.StatusInternalServerError)
		return
	} // err handling

	dumpFileName := fmt.Sprintf("%s/<%d>-<%s>.dmp", globals.MemoryConfig.DumpPath, dump.PID, dump.TimeStamp)
	logger.Info("EL NOMBRE DEL DUMPFILE ES: " + dumpFileName)
	dumpFile, err := os.OpenFile(dumpFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Printf("Error al crear archivo de log para <%d-%s>: %v\n", dump.PID, dump.TimeStamp, err)
		os.Exit(1)
	} // err handling y se asigna el nombre del dumpfile
	log.SetOutput(dumpFile)
	defer dumpFile.Close()

	logger.Info("## PID: <%d>  - Memory Dump solicitado", dump.PID) // se logea
	// TODO: se debe ubicar el proceso completo
	// InformarMetricasProceso(proceso.Metricas)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Dump Realizado"))
}
