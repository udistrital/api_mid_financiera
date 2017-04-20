package models

type FuenteFinanciacionRubroNecesidad struct {
	Id                 int     `orm:"column(id);pk"`
	Apropiacion        int     `orm:"column(apropiacion)"`
	MontoParcial       float64 `orm:"column(monto_parcial)"`
	FuenteFinanciacion int     `orm:"column(fuente_financiacion);null"`
}
