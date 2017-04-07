package models

import (
	"time"
)

type Apropiacion struct {
	Id              int                `orm:"column(id);pk"`
	Vigencia        float64            `orm:"column(vigencia);null"`
	Rubro           *Rubro             `orm:"column(rubro);rel(fk)"`
	UnidadEjecutora int                `orm:"column(unidad_ejecutora);null"`
	ValorRezago     float64            `orm:"column(valor_rezago);null"`
	Valor           float64            `orm:"column(valor);null"`
	TipoDocumento   string             `orm:"column(tipo_documento);null"`
	DocumentoNumero string             `orm:"column(documento_numero);null"`
	DocumentoFecha  time.Time          `orm:"column(documento_fecha);type(date);null"`
	Estado          *EstadoApropiacion `orm:"column(estado);rel(fk)"`
}
