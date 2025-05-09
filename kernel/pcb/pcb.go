package pcb

// posibles estados de un proceso
const (
	EstadoNew         = "new"
	EstadoReady       = "ready"
	EstadoExecute     = "execute"
	EstadoBlocked     = "blocked"
	EstadoExit        = "exit"
	EstadoSuspBlocked = "suspblocked"
	EstadoSuspReady   = "suspready"
)

/*
	type PCB struct {
		PID int
		PC  int
		ME  map[string]int //asocia cada estado con la cantidad de veces que el proceso estuvo en ese estado.
		MT  map[string]int //asocia cada estado con el tiempo total que el proceso pasó en ese estado.
	}
*/

// /CHICOS ES NECESARIO AGREGAR AL PCB EL TAMAÑO Y NOMBRE DE ARCHIVO DE PESUDO COGIO TANTO PARA PLANIFICADOR DE LARGO PLAZO
// (cuando termina un proceso hay que preguntar si el pcb de NEW puede inicilizar)Y PARA SJF
type PCB struct {
	PID         int
	PC          int
	ME          map[string]int
	MT          map[string]int
	FileName    string // nombre de archivo de pesudoCodigo
	ProcessSize int
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
