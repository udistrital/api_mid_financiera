package models

import (
	"time"
)

type SolicitudNecesidad struct {
	Id                     int                       `orm:"column(id);pk"`
	Numero                 int                       `orm:"column(numero)"`
	Vigencia               float64                   `orm:"column(vigencia)"`
	Dependencia            *DependenciaTemporal                       `orm:"column(dependencia);rel(fk)"`
	JefeDependencia        int  										`orm:"column(jefe_dependencia)"`
	DependenciaDestino     *DependenciaTemporal      `orm:"column(dependencia_destino);rel(fk)"`
	ObjetoContractual      string                    `orm:"column(objeto_contractual)"`
	Estado                 *EstadoSolicitudNecesidad `orm:"column(estado);rel(fk)"`
	FechaSolicitud         time.Time                 `orm:"column(fecha_solicitud);type(date)"`
	ValorContratacion      float64                   `orm:"column(valor_contratacion)"`
	Justificacion          string                    `orm:"column(justificacion)"`
	UnidadEjecutora      *UnidadEjecutora                  `orm:"column(unidad_ejecutora);rel(fk);null"`
	DiasDuracion           float64                   `orm:"column(dias_duracion)"`
	UnicoPago              bool                      `orm:"column(unico_pago)"`
	AgotarPresupuesto      bool                      `orm:"column(agotar_presupuesto)"`
	NovedadOtroSi          bool                      `orm:"column(novedad_otro_si)"`
	CodigoContratoOtroSi   int                       `orm:"column(codigo_contrato_otro_si);null"`
	ModalidadSeleccion     *ModalidadSeleccion       `orm:"column(modalidad_seleccion);rel(fk)"`
	Entidad                int                       `orm:"column(entidad)"`
	Servicio               int                       `orm:"column(servicio)"`
	PlanAnualAdquisiciones int                       `orm:"column(plan_anual_adquisiciones)"`
	EstudioMercado         string                    `orm:"column(estudio_mercado);null"`
	FechaEvaluacion        time.Time                 `orm:"column(fecha_evaluacion);type(date);null"`
	OrdenadorGasto         int   										 `orm:"column(ordenador_gasto)"`
	JustificacionRechazo   string                    `orm:"column(justificacion_rechazo);null"`
	TipoContratacion       *TipoContratacion         `orm:"column(tipo_contratacion);rel(fk)"`
	AnalisisRiesgo         string                    `orm:"column(analisis_riesgo);null"`
	TecnicasUniformes      bool                    `orm:"column(tecnicas_uniformes);null"`
	VigenciaContratoOtroSi float64									 `orm:"column(vigencia_contrato_otro_si)"`
}
