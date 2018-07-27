package models


type UbicacionEnte struct {
	Id                        int                        `orm:"column(id);pk;auto"`
	Lugar                     int                        `orm:"column(lugar)"`
	Ente                      *Ente                      `orm:"column(ente);rel(fk)"`
	TipoRelacionUbicacionEnte *TipoRelacionUbicacionEnte `orm:"column(tipo_relacion_ubicacion_ente);rel(fk)"`
	Activo                    bool                       `orm:"column(activo)"`
}
