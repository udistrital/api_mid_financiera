package models

import (
	"time"
)

type RegistroPresupuestal struct {
	Id                         int                         `orm:"auto;column(id);pk"`
	UnidadEjecutora            *UnidadEjecutora            `orm:"column(unidad_ejecutora);rel(fk)"`
	Vigencia                   float64                     `orm:"column(vigencia)"`
	FechaMovimiento            time.Time                   `orm:"column(fecha_movimiento);type(date);null"`
	Responsable                int                         `orm:"column(responsable);null"`
	Estado                     *EstadoRegistroPresupuestal `orm:"column(estado);rel(fk)"`
	NumeroRegistroPresupuestal int                         `orm:"column(numero_registro_presupuestal)"`
	Beneficiario               int                         `orm:"column(beneficiario);rel(fk)"`
	Compromiso                 *Compromiso                 `orm:"column(compromiso);rel(fk)"`
}
type DatosRubroRegistroPresupuestal struct {
	Id             int
	Disponibilidad *Disponibilidad
	Apropiacion    *Apropiacion
	Valor          float64
	ValorAsignado  float64
	Saldo          float64
}
type DatosRegistroPresupuestal struct { //estructura temporal para el registro con relacion a las apropiaciones
	Rp     *RegistroPresupuestal
	Rubros []DatosRubroRegistroPresupuestal
}

type InfoSolRp struct {
	Solicitud *SolicitudRp
	Rubros    []DisponibilidadApropiacionSolicitud_rp
}
