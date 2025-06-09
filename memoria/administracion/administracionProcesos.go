package administracion

import (
	"github.com/sisoputnfrba/tp-golang/utils/logger"
	"net/http"
)

// ------------------------------------------------------------------
// ----------- FORMA PARTE DE LA MODIFICACIÓN DE PROCESOS -----------
// ------------------------------------------------------------------

func InicializacionProceso(w http.ResponseWriter, r *http.Request) {
	// TODO: VERIFICAR EL TAMAÑO NECESARIO

	// TODO: CREAR ESTRUCTURAS ADMINISTRATIVAS NECESARIAS

	// TODO: RESPONDER CON EL NUMERO DE PAGINA DE 1ER NIVEL DEL PROCESO
	logger.Info("## PID: <%d>  - Proceso Creado - Tamaño: <%d>")
}

func FinalizacionProceso(w http.ResponseWriter, r *http.Request) {
	//toDO

	logger.Info("## PID: <PID>  - Proceso Destruido - Métricas - Acc.T.Pag: <ATP>; Inst.Sol.: <Inst.Sol>; SWAP: <SWAP>; Mem. Prin.: <Mem.Prin.>; Lec.Mem.: <Lec.Mem.>; Esc.Mem.: <Esc.Mem.>")
}

func SuspensionProceso(w http.ResponseWriter, r *http.Request) {
	// TODO: NO ES NECESARIO EL SWAPEO DE TABLAS DE PAGINAS

	// TODO: SE LIBERA EN MEMORIA
	// TODO: SE ESCRIBE EN SWAP LA INFO NECESARIA

}

func DesSuspensionProceso(w http.ResponseWriter, r *http.Request) {
	// TODO: VERIFICAR EL TAMAÑO NECESARIO

	// TODO: LEER EL CONTENIDO DEL SWAP, ESCRIBIERLO EN EL FRAME ASIGNADO
	// TODO: LIBERAR ESPACIO EN SWAP
	// TODO: ACTUALIZAR ESTRUCTURAS NECESARIAS

	// TODO: RETORNAR OK
}
