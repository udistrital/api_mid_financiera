package models



type JefeDependenciaTemporal struct {
	Id                  int                  `orm:"column(id);pk"`
	Nombre              string               `orm:"column(nombre)"`
	DependenciaTemporal *DependenciaTemporal `orm:"column(dependencia_temporal);rel(fk)"`
}
