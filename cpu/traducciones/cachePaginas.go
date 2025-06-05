package traducciones

import "log"

func Leer(dirFisica int, tamanio int) (string, error) {

	log.Printf("Direccion leida: %d - Tamanio leido: %d", dirFisica, tamanio)

	return "", nil
}

func Escribir(dirFisica int, datos string) error {

	log.Printf("Se escribio %s en la direccion fisica %d", datos, dirFisica)

	return nil
}

/*LOGS MINIMOS RESTANTES:
CACHE HIT --> Si con la dirFisica encuentra una Pagina
CACHE MISS --> Si con la dirFisica no encuentra una Pagina
CACHE ADD --> Despues de no haber encontrado la pagina, la agrega
*/
