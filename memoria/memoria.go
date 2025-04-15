package main

import (
	"fmt"
	"github.com/sisoputnfrba/tp-golang/memoria/globals"
	"github.com/sisoputnfrba/tp-golang/memoria/utils"
	"log"
	"net/http"
)

func main() {

	globals.MemoryConfig = utils.Config("config.json")

	if globals.MemoryConfig == nil {
		log.Fatal("No se pudo cargar el archivo de configuración")
	}
	var portMemory = globals.MemoryConfig.PortMemory
	log.Println("Comenzó ejecucion de la memoria")

	mux := http.NewServeMux()
	// ESTÁ ESPERANDO LOS MENSAJES DE LOS OTROS MODULOS
	mux.HandleFunc("/memoria/cpu", utils.RecibirMensajeDeCPU)
	mux.HandleFunc("/memoria/kernel", utils.RecibirMensajeDeKernel)
	//mux.HandleFunc("/memoria/cpu", utils.CreacionProceso)

	fmt.Printf("Servidor escuchando en http://localhost:%d/memoria\n", portMemory)

	direccion := fmt.Sprintf(":%d", portMemory)
	err := http.ListenAndServe(direccion, mux)
	if err != nil {
		panic(err)
	}

}
