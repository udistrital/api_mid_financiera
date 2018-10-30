package models

//ConceptoOrdenPago ...
type ConceptoOrdenPago struct {
	Id          int        `orm:"column(id);pk;auto"`
	Valor       float64    `orm:"column(valor)"`
	Concepto    *Concepto  `orm:"column(concepto);rel(fk)"`
	OrdenDePago *OrdenPago `orm:"column(orden_de_pago);rel(fk)"`
}
