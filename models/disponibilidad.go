package models

import (
	"time"
)

type Disponibilidad struct {
	Id                   int                   `orm:"column(id);pk"`
	UnidadEjecutora      *UnidadEjecutora      `orm:"column(unidad_ejecutora);rel(fk)"`
	Vigencia             float64               `orm:"column(vigencia)"`
	NumeroDisponibilidad float64               `orm:"column(numero_disponibilidad);null"`
	Responsable          int                   `orm:"column(responsable);null"`
	Solicitante          int                   `orm:"column(solicitante);null"`
	FechaRegistro        time.Time             `orm:"column(fecha_registro);type(date);null"`
	Estado               *EstadoDisponibilidad `orm:"column(estado);rel(fk)"`
	NumeroOficio         string                `orm:"column(numero_oficio);null"`
	Destino              int                   `orm:"column(destino);null"`
	Solicitud            int                   `orm:"column(solicitud)"`
	DatosNecesidad       *Necesidad
}

type InfoSolDisp struct {
	SolicitudDisponibilidad SolicitudDisponibilidad
	DependenciaSolicitante  Dependencia
	DependenciaDestino      Dependencia
}
