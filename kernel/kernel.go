package main

import (
	"fmt"
	"log"
	"net/http"
	
)

func init() {

}

func main() {
	log.Println("Comenzó ejecucion del Kernel")

	//FUNCIONES DE CLIENTE. CONEXION CON OTROS MODULOS:
	//enviar mensaje
	//generar y enviar paquete
	

	fmt.Println("Hello, World!")

	//TODO ver guía de golang en la sección de "Protocolo HTTP".

	

	//FUNCIONES SERVER DE LOS OTROS MODULOS:
	//mux.HandleFunc("/paquetes", utils.RecibirPaquetes)
	//mux.HandleFunc("/mensaje", utils.RecibirMensaje)

}
