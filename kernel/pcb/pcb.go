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

type PCB struct {
	PID Pid
	PC  int
	ME  map[string]int //asocia cada estado con la cantidad de veces que el proceso estuvo en ese estado.
	MT  map[string]int //asocia cada estado con el tiempo total que el proceso pasó en ese estado.
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
