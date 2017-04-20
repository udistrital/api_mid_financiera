package models

import (
	"time"
)

type JefeDependencia struct {
	Id             int       `orm:"column(id);pk;auto"`
	FechaInicio    time.Time `orm:"column(fecha_inicio)"`
	FechaFin       time.Time `orm:"column(fecha_fin)"`
	TerceroId      int       `orm:"column(tercero_id)"`
	DependenciaId  int       `orm:"column(dependencia_id)"`
	ActaAprobacion string    `orm:"column(acta_aprobacion)"`
}
