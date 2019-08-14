

package models

type CuentaContable struct {
	Id                 int                 `orm:"column(id);pk;auto"`
	Saldo              int64               `orm:"column(saldo)"`
	Nombre             string              `orm:"column(nombre)"`
	Naturaleza         string              `orm:"column(naturaleza)"`
	Descripcion        string              `orm:"column(descripcion);null"`
	Codigo             string              `orm:"column(codigo)"`
	
}
