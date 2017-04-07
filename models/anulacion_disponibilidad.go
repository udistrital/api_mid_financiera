package models

import (
	"time"
)

type AnulacionDisponibilidad struct {
	Id            int       `orm:"column(id);pk"`
	Motivo        string    `orm:"column(motivo)"`
	FechaRegistro time.Time `orm:"column(fecha_registro);type(date)"`
	TipoAnulacion string    `orm:"column(tipo_anulacion)"`
}
