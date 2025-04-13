package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/kernel/utils"
)

func main() {
	log.Println("Comenzó ejecucion del Kernel")

	mux := http.NewServeMux()

	//FUNCIONES SERVER DE LOS OTROS MODULOS: Escuchar sus mensajes
	mux.HandleFunc("/kernel/mensaje", utils.RecibirMensajeDeCpu)

	fmt.Println("Servidor escuchando en http://localhost:8081/kernel/mensaje")

	err := http.ListenAndServe(":8081", mux)
	if err != nil {
		panic(err)
	}

	//FUNCIONES DE CLIENTE. CONEXION CON OTROS MODULOS:
	//enviar mensaje
	//generar y enviar paquete

	//TODO ver guía de golang en la sección de "Protocolo HTTP".
}
