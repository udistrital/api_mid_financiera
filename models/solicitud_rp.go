package models

import (
	"time"
)

type SolicitudRp struct {
	Id                  int
	Vigencia            int       `orm:"column(vigencia)"`
	FechaSolicitud      time.Time `orm:"column(fecha_solicitud);type(date);null"`
	Cdp                 int       `orm:"column(cdp)"`
	Expedida            bool      `orm:"column(expedida)"`
	NumeroContrato      string    `orm:"column(numero_contrato)"`
	VigenciaContrato    int       `orm:"column(vigencia_contrato)"`
	Compromiso          int       `orm:"column(compromiso)"`
	Masivo              bool
	DatosDisponibilidad *Disponibilidad
	DatosProveedor      *InformacionProveedor
	DatosCompromiso     interface{}
	Rubros              []interface{}
}
