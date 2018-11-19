package models

//ActividadesCentroCostos ...
type ActividadesCentroCostos struct {
	Id           int           `orm:"column(id);pk"`
	Nombre       string        `orm:"column(nombre)"`
	Descripcion  string        `orm:"column(descripcion);null"`
	CentroCostos *CentroCostos `orm:"column(centro_costos);rel(fk)"`
}
