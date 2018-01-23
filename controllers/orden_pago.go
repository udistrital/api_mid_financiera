package controllers

import (
	// "fmt"
	// "strconv"
	// "strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/ss_crud_api/models"
	// "github.com/udistrital/api_mid_financiera/utilidades"
)

// OrdenPagoController operations for Orden_pago
type OrdenPagoController struct {
	beego.Controller
}

// URLMapping ...
func (c *OrdenPagoController) URLMapping() {
	c.Mapping("GetOrdenPagoByFuenteFinanciamiento", c.GetOrdenPagoByFuenteFinanciamiento)
}

// GetOrdenPagoByFuenteFinanciamiento ...
// @Title GetOrdenPagoByFuenteFinanciamiento
// @Description lista Ordenes de Pago por fuente de financiamineto.
// @Param	fuente	query	string	true	"Id Fuente de Financiamiento"
// @Param	vigencia	query	string	true	"Vigencia de registro Presupuestal"
// @Param	unidadEjecutora	query	string	true	"Unidad Ejecutora"
// @Success 201 {object} models.Alert
// @Failure 403 body is empty
// @router /GetOrdenPagoByFuenteFinanciamiento [get]
func (c *OrdenPagoController) GetOrdenPagoByFuenteFinanciamiento() {
	println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAa")
	fuente := c.GetString("fuente")
	vigencia := c.GetString("vigencia")
	unidadEjecutora := c.GetString("unidadEjecutora")
	if fuente != "" && vigencia != "" && unidadEjecutora != "" {
		println("AAAAAAAAAAAAAAA")
		// if mensaje.Type == "success" {
		// 	c.Data["json"] = Ordenes
		// } else {
		// 	c.Ctx.Output.SetStatus(201)
		// 	c.Data["json"] = mensaje
		// }
	} else {
		c.Ctx.Output.SetStatus(201)
		c.Data["json"] = models.Alert{Code: "E_0458", Body: "Not enough parameter in GetOrdenPagoByFuenteFinanciamiento", Type: "error"}
	}
	c.ServeJSON()

}
