package models



type DisponibilidadRubro struct {
	Id             int             `orm:"column(id);pk"`
	Vigencia       float64         `orm:"column(vigencia)"`
	Rubro          *Rubro          `orm:"column(rubro);rel(fk)"`
	Disponibilidad *Disponibilidad `orm:"column(disponibilidad);rel(fk)"`
	Valor          float64         `orm:"column(valor);null"`
}
