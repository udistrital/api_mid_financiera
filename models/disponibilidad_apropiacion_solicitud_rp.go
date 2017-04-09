package models

type DisponibilidadApropiacionSolicitud_rp struct {
	Id                        int
	DisponibilidadApropiacion int `orm:"column(disponibilidad_apropiacion)"`
	SolicitudRp               int `orm:"column(solicitud_rp)"`
	Monto                     int `orm:"column(monto)"`
}
