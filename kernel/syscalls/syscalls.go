package syscalls

import (
	"log"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/kernel/globals"
	"github.com/sisoputnfrba/tp-golang/kernel/utils"
	"github.com/sisoputnfrba/tp-golang/utils/data"
)

// Body JSON a recibir
type MensajeInit struct {
	Filename string `json:"filename"`
	Tamanio  int `json:"tamanio_memoria"`
}

type MensajeIo struct {
	Nombre string `json:"nombre"`
	Duracion  int `json:"duracion"`
}

func ContextoInterrumpido(w http.ResponseWriter, r *http.Request) {
    // tu código...
}

func InitProc(w http.ResponseWriter, r *http.Request) {
    log.Printf("Syscall recibida: “## (<PID>) - Solicitó syscall: <INIT_PROC>”")
  
    var mensajeRecibido MensajeInit
	if err := data.LeerJson(w, r, &mensajeRecibido); err != nil {
		return //hubo error
	}

	//Cargar en
	/*globals.Init = globals.InitProc{
		Filename:     mensajeRecibido.Filename,
		Tamanio: mensajeRecibido.Tamanio,
	}*/

	log.Printf("Se ha recibido: Filename: %s Tamaño Memoria: %d", 
        mensajeRecibido.Filename, mensajeRecibido.Tamanio)

	//Crear Proceso
	//CrearProceso(mensajeRecibido.Filename,mensajeRecibido.Tamanio)
}

func Exit(w http.ResponseWriter, r *http.Request) {
    // tu código...
}

func DumpMemory(w http.ResponseWriter, r *http.Request) {
  
}

func Io(w http.ResponseWriter, r *http.Request) {
    log.Printf("Syscall recibida: “## (<PID>) - Solicitó syscall: <IO>”")
  
    var mensajeRecibido MensajeIo
	if err := data.LeerJson(w, r, &mensajeRecibido); err != nil {
		return //hubo error
	}

	log.Printf("Se ha recibido: Nombre: %s Duracion: %d", 
        mensajeRecibido.Nombre, mensajeRecibido.Duracion)

	if(globals.IO.Nombre == mensajeRecibido.Nombre){
        utils.EnviarContextoIO(globals.IO.Ip,globals.IO.Puerto)
    }
}