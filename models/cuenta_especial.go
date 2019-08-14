package models

type CuentaEspecial struct {
	Id                         int                 `orm:"column(id);pk;auto"`
	Descripcion                string              `orm:"column(descripcion);null"`
	Porcentaje                 float64             `orm:"column(porcentaje);null"`
	TarifaUvt                  float64             `orm:"column(tarifa_uvt);null"`
	Deducible                  bool                `orm:"column(deducible);null"`
	CuentaContable             *CuentaContable     `orm:"column(cuenta_contable);rel(fk)"`
	InformacionPersonaJuridica int64               `orm:"column(informacion_persona_juridica)"`
}
