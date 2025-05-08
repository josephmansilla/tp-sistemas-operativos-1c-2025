package comunicacion

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/pcb"
	"github.com/sisoputnfrba/tp-golang/utils/data"
	logger "github.com/sisoputnfrba/tp-golang/utils/logger"
	"net/http"
)

// Body JSON a recibir
type MensajeDeCPU struct {
	Ip     string `json:"ip"`
	Puerto int    `json:"puerto"`
	ID     string `json:"id"`
}

type MensajeACPU struct {
	Pid int `json:"pid"`
	Pc  int `json:"pc"`
}

func RecibirMensajeDeCPU(w http.ResponseWriter, r *http.Request) {
	var mensajeRecibido MensajeDeCPU
	if err := data.LeerJson(w, r, &mensajeRecibido); err != nil {
		logger.Error("Error al leer JSON: %s", err.Error())
		return
	}

	id := mensajeRecibido.ID

	// Cargar en globals
	globals.CPUMu.Lock()
	globals.CPUs[id] = globals.DatosCPU{
		ID:     mensajeRecibido.ID,
		Ip:     mensajeRecibido.Ip,
		Puerto: mensajeRecibido.Puerto,
	}
	globals.CPUCond.Broadcast() // Despierta a quien espera CPUs
	globals.CPUMu.Unlock()

	logger.Info("Se ha recibido CPU: Ip: %s Puerto: %d ID: %s",
		globals.CPUs[id].Ip, globals.CPUs[id].Puerto, globals.CPUs[id].ID)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("STATUS OK"))
}

func EnviarContextoCPU(id string, pcb *pcb.PCB) {
	globals.CPUMu.Lock()
	cpuData, ok := globals.CPUs[id]
	for !ok {
		globals.CPUCond.Wait()
		cpuData, ok = globals.CPUs[id]
	}
	globals.CPUMu.Unlock()

	url := fmt.Sprintf("http://%s:%d/cpu/kernel", cpuData.Ip, cpuData.Puerto)

	mensaje := MensajeACPU{
		Pid: pcb.PID, // PID real
		Pc:  pcb.PC,  // PC real
	}

	err := data.EnviarDatos(url, mensaje)
	if err != nil {
		logger.Info("Error enviando PID y PC a CPU: %s", err.Error())
		return
	}
	logger.Info("PID: %d y PC: %d enviados exitosamente a CPU", mensaje.Pid, mensaje.Pc)
}
