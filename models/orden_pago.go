package models

import (
	"time"
)

type OrdenPago struct {
	Id                   int                   `orm:"column(id);pk;auto"`
	Vigencia             float64               `orm:"column(vigencia)"`
	FechaCreacion        time.Time             `orm:"column(fecha_creacion);type(date)"`
	RegistroPresupuestal *RegistroPresupuestal `orm:"column(registro_presupuestal);rel(fk)"`
	ValorBase            float64               `orm:"column(valor_base)"`
	PersonaElaboro       int                   `orm:"column(persona_elaboro)"`
	Convenio             int                   `orm:"column(convenio);null"`
}
