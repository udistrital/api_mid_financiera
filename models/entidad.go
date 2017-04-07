package models

type Entidad struct {
	Id            int    `orm:"column(id);pk"`
	Nombre        string `orm:"column(nombre)"`
	CodigoEntidad string `orm:"column(codigo_entidad)"`
}
