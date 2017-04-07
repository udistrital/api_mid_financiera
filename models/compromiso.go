package models

import (
	"time"
)

type Compromiso struct {
	Id                    int                    `orm:"column(id);pk;auto"`
	Objeto                string                 `orm:"column(objeto)"`
	Vigencia              float64                `orm:"column(vigencia)"`
	FechaInicio           time.Time              `orm:"column(fecha_inicio);type(date)"`
	FechaFin              time.Time              `orm:"column(fecha_fin);type(date)"`
	FechaModificacion     time.Time              `orm:"column(fecha_modificacion);type(date)"`
	EstadoCompromiso      *EstadoCompromiso      `orm:"column(estado_compromiso);rel(fk)"`
	TipoCompromisoTesoral *TipoCompromisoTesoral `orm:"column(tipo_compromiso_tesoral);rel(fk)"`
	UnidadEjecutora       *UnidadEjecutora       `orm:"column(unidad_ejecutora);rel(fk)"`
}
