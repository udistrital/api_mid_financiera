package models

type ContratoGeneral struct {
	Id          string                `orm:"column(numero_contrato);pk"`
	Contratista *InformacionProveedor `orm:"column(contratista);rel(fk)"`
}
