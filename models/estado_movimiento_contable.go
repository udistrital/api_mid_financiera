package models 

type EstadoMovimientoContable struct {
	Id                int     `orm:"column(id);pk;auto"`
	Nombre            string  `orm:"column(nombre)"`
	Descripcion       string  `orm:"column(descripcion);null"`
	CodigoAbreviacion string  `orm:"column(codigo_abreviacion);null"`
	Estado            bool    `orm:"column(estado)"`
	NumeroOrden       float64 `orm:"column(numero_orden);null"`
}
