package utils

import (
	"encoding/json"
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"log"
	"net/http"
	"os"
)

func Config(filepath string) *globals.Config {
	var config *globals.Config
	configFile, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}

func LeerJson(w http.ResponseWriter, r *http.Request, mensaje any) {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&mensaje)

	if err != nil {
		log.Printf("Error al decodificar el mensaje: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error al decodificar mensaje"))
		return
	}

	log.Println("Me llego un mensaje JSON:")
	log.Printf("%+v\n", mensaje)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("STATUS OK"))
}

// FUNCION PARA RECIBIR LOS MENSAJES PROVENIENTES DE LA CPU
func RecibirMensajeDeCPU(w http.ResponseWriter, r *http.Request) {
	var mensaje globals.DatosDeCPU
	LeerJson(w, r, &mensaje)

	globals.CPU = globals.DatosDeCPU{
		PID: mensaje.PID,
		PC:  mensaje.PC,
	}

	log.Printf("PID Pedido: %d\n", mensaje.PID)
	log.Printf("PC Pedido: %d\n", mensaje.PC)
}

// FUNCION PARA RECIBIR LOS MENSAJES PROVENIENTES DEL KERNEL
func RecibirMensajeDeKernel(w http.ResponseWriter, r *http.Request) {
	var mensaje globals.DatosDeKernel
	LeerJson(w, r, &mensaje)

	globals.Kernel = globals.DatosDeKernel{
		TamanioMemoria: mensaje.TamanioMemoria,
	}

	log.Printf("Tamaño de Memoria Pedido: %d\n", mensaje.TamanioMemoria)
}
