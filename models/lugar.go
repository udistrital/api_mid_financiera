package models

type Lugar struct {
	Id        int        `orm:"column(id);pk;auto"`
	Nombre    string     `orm:"column(nombre)"`
	TipoLugar *TipoLugar `orm:"column(tipo_lugar);rel(fk)"`
	Activo    bool       `orm:"column(activo)"`
}
