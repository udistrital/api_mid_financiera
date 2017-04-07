package models


type DependenciaTemporal struct {
	Id             int                     `orm:"column(id);pk"`
	Nombre         string                  `orm:"column(nombre)"`
	OrdenadorGasto *OrdenadorGastoTemporal `orm:"column(ordenador_gasto);rel(fk)"`
}
