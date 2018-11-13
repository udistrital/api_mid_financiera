package models

import (
	"time"
)

type MovimientoContable struct {
	Id                       int                       `orm:"column(id);pk;auto"`
	Debito                   int64                     `orm:"column(debito)"`
	Credito                  int64                     `orm:"column(credito)"`
	Fecha                    time.Time                 `orm:"column(fecha);type(timestamp without time zone)"`
	Concepto                 *Concepto                 `orm:"column(concepto_tesoral);rel(fk)"`
	CuentaContable           *CuentaContable           `orm:"column(cuenta_contable);rel(fk)"`
	TipoDocumentoAfectante   *TipoDocumentoAfectante   `orm:"column(tipo_documento_afectante);rel(fk)"`
	CodigoDocumentoAfectante int                       `orm:"column(codigo_documento_afectante)"`
	EstadoMovimientoContable *EstadoMovimientoContable `orm:"column(estado_movimiento_contable);rel(fk);null"`
	CuentaEspecial           *CuentaEspecial           `orm:"column(cuenta_especial);rel(fk);null"`
}
