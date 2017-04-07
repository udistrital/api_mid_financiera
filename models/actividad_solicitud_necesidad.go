package models


type ActividadSolicitudNecesidad struct {
	Id                      int                      `orm:"column(id);pk"`
	Actividad               *ActividadesCentroCostos `orm:"column(actividad);rel(fk)"`
	Necesidad *Necesidad 		`orm:"column(necesidad);rel(fk)"`
	MontoParcial            float64                  `orm:"column(monto_parcial)"`
}
