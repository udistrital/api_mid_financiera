package models



type OrdenadorGastoTemporal struct {
	Id     int    `orm:"column(id);pk"`
	Nombre string `orm:"column(nombre)"`
}
