package models

//ContactoEnte ...
type ContactoEnte struct {
	Id           int           `orm:"column(id);pk;auto"`
	TipoContacto *TipoContacto `orm:"column(tipo_contacto);rel(fk)"`
	Ente         *Ente         `orm:"column(ente);rel(fk)"`
	Valor        string        `orm:"column(valor)"`
}
