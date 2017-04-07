package models



type FuenteFinanciacion struct {
	Id          int    `orm:"column(id);pk"`
	Descripcion string `orm:"column(descripcion);null"`
	Sigla       string `orm:"column(sigla)"`
}
