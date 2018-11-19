package models

type ValorAtributoUbicacion struct {
	Id                int                `orm:"column(id);pk;auto"`
	UbicacionEnte     *UbicacionEnte     `orm:"column(ubicacion_ente);rel(fk)"`
	AtributoUbicacion interface{}       `orm:"column(atributo_ubicacion);rel(fk)"`
	Valor             string             `orm:"column(valor)"`
}
