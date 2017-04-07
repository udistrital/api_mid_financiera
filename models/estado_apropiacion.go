package models

type EstadoApropiacion struct {
	Id          int    `orm:"column(id);pk"`
	Nombre      string `orm:"column(nombre);null"`
	Descripcion string `orm:"column(descripcion);null"`
}
