package models

// Organizacion ...
type Organizacion struct {
	Id               int               `orm:"column(id);pk;auto"`
	Nombre           string            `orm:"column(nombre)"`
	Ente             int               `orm:"column(ente)"`
	TipoOrganizacion *TipoOrganizacion `orm:"column(tipo_organizacion);rel(fk)"`

}
