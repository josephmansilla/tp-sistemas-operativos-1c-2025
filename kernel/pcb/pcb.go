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
	PID int
	PC  int
	ME  map[string]int //asocia cada estado con la cantidad de veces que el proceso estuvo en ese estado.
	MT  map[string]int //asocia cada estado con el tiempo total que el proceso pasó en ese estado.
}

//Ej ME: "ready": 3 → el proceso estuvo 3 veces en el estado listo.
//Ej MT: "execute": 12 → el proceso estuvo 12 unidades de tiempo en ejecución.
