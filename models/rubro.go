package models

type Rubro struct {
	Id             int64    `pk;orm:"column(id);serial"`
	Entidad        *Entidad `orm:"column(entidad);rel(fk)"`
	Codigo         string   `orm:"column(codigo);null"`
	Vigencia       float64  `orm:"column(vigencia)"`
	Descripcion    string   `orm:"column(descripcion);null"`
	TipoPlan       int16    `orm:"column(tipo_plan);null"`
	Administracion string   `orm:"column(administracion);null"`
	Estado         int16    `orm:"column(estado);null"`
}
