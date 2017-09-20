package models

import (
	"time"
)

type RegistroPresupuestal struct {
	Id                         int                         `orm:"auto;column(id);pk"`
	Vigencia                   float64                     `orm:"column(vigencia)"`
	FechaRegistro              time.Time                   `orm:"column(fecha_registro);type(date);null"`
	Responsable                int                         `orm:"column(responsable);null"`
	Estado                     *EstadoRegistroPresupuestal `orm:"column(estado);rel(fk)"`
	NumeroRegistroPresupuestal int                         `orm:"column(numero_registro_presupuestal)"`
	Beneficiario               int                         `orm:"column(beneficiario);rel(fk)"`
	Compromiso                 *Compromiso                 `orm:"column(compromiso);rel(fk)"`
	Solicitud                  int
	DatosSolicitud             *SolicitudRp
}
type DatosRubroRegistroPresupuestal struct {
	Id                 int
	Disponibilidad     *Disponibilidad
	Apropiacion        *Apropiacion
	FuenteFinanciacion *FuenteFinanciacion
	Valor              float64
	ValorAsignado      float64
	Saldo              float64
}
type DatosRegistroPresupuestal struct {
	Rp     *RegistroPresupuestal
	Rubros []DatosRubroRegistroPresupuestal
}

type InfoSolRp struct {
	Solicitud *SolicitudRp
	Rubros    []DisponibilidadApropiacionSolicitud_rp
}
