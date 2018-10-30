package models


// NecesidadOtroSi ...
type NecesidadOtroSi struct {
	Id       int    `orm:"column(id);pk"`
	Contrato string `orm:"column(contrato)"`
	Vigencia int    `orm:"column(vigencia)"`
}
