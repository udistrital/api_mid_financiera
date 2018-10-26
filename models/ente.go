package models

type Ente struct {
	Id       int      `orm:"column(id);pk;auto"`
	TipoEnte TipoEnte `orm:"column(tipo_ente);rel(fk)"`
}
