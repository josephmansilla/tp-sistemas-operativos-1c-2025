package administracion

import (
	"bufio"
	"fmt"
	g "github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"os"
	"strings"
)

func LecturaPseudocodigo(proceso *g.Proceso, direccionPseudocodigo string) error {
	if direccionPseudocodigo == "" {
		return fmt.Errorf("el string es vacio")
	}
	ruta := "../pruebas/" + direccionPseudocodigo
	file, err := os.Open(ruta)
	if err != nil {
		logger.Error("Error al abrir el archivo: %s\n", err)
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	pc := 0
	for scanner.Scan() {
		lineaEnString := scanner.Text()
		lineaEnBytes := []byte(lineaEnString)

		proceso.InstruccionesEnBytes[pc] = lineaEnBytes
		pc++

		if strings.TrimSpace(lineaEnString) == "EXIT" {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		logger.Error("Error al leer el archivo: %s", err)
	}

	IncrementarMetrica(proceso, pc, IncrementarInstruccionesSolicitadas)
	logger.Info("<%d> instrucciones cargadas para PID <%d>", proceso.PID, pc)

	return nil
}
