package models

import (
	"time"
)

type Disponibilidad struct {
	Id                   int                   `orm:"column(id);pk"`
	Vigencia             float64               `orm:"column(vigencia)"`
	NumeroDisponibilidad float64               `orm:"column(numero_disponibilidad);null"`
	Responsable          int                   `orm:"column(responsable);null"`
	FechaRegistro        time.Time             `orm:"column(fecha_registro);type(date);null"`
	Estado               *EstadoDisponibilidad `orm:"column(estado);rel(fk)"`
	Solicitud            int                   `orm:"column(solicitud)"`
	DatosNecesidad       *Necesidad
}

type InfoSolDisp struct {
	SolicitudDisponibilidad SolicitudDisponibilidad
	DependenciaSolicitante  Dependencia
	DependenciaDestino      Dependencia
}
