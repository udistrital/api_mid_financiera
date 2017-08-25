package models

type ContratoGeneral struct {
	Id          string `orm:"column(numero_contrato);pk"`
	Contratista int    `orm:"column(contratista)"`
}
