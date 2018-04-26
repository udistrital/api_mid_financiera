package models

import (
	"time"
)

type PagosAcademica struct{
	informacionEstudiante infoEstudiante
	informacionCarrera    infoCarrera
}

type infoRecibo struct {
	Id           int64
	Total        int64
	NumeroRecibo int64
	fechaExtraordinario time.Time
	fechaOrdinario time.Time
	periodo			string
	pago				string
	desagregaRecibos []*infoPago
}

type infoEstudiante struct {
	tipoDocu  	string
	documento   string
	tipo				string
	nombre      string
}

type infoCarrera struct {
	carrera			string
	facultad		string
	codCarrera 	string
	codigoEst			int64
	informacionRecibos  []*infoRecibo
}

type infoPago struct {
	descripcion string
	valor				float64
}
