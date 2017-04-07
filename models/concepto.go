package models

import (
	"time"
)

type Concepto struct {
	Id              int       `orm:"column(id);pk;auto"`
	Codigo          string    `orm:"column(codigo)"`
	Nombre          string    `orm:"column(nombre)"`
	FechaCreacion   time.Time `orm:"column(fecha_creacion);type(date)"`
	Cabeza          bool      `orm:"column(cabeza)"`
	FechaExpiracion time.Time `orm:"column(fecha_expiracion);type(date);null"`
	Descripcion     string    `orm:"column(descripcion);null"`
	Rubro           *Rubro    `orm:"column(rubro);rel(fk);null"`
}
