package controllers

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/ss_crud_api/models"
	"github.com/udistrital/utils_oas/optimize"
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
	fuente := c.GetString("fuente")
	vigencia := c.GetString("vigencia")
	unidadEjecutora := c.GetString("unidadEjecutora")
	if fuente != "" && vigencia != "" && unidadEjecutora != "" {
		var registroPresupuestales []interface{}
		respuestas := make(map[string]interface{})
		var parametro []interface{}
		parametro = append(parametro, unidadEjecutora)
		// search registro_presupuestal
		if err := getJson("http://"+beego.AppConfig.String("kronosService")+"registro_presupuestal/?query=RegistroPresupuestalDisponibilidadApropiacion.DisponibilidadApropiacion.FuenteFinanciamiento.Id:"+fuente+",Vigencia:"+vigencia+"&limit:-1", &registroPresupuestales); err == nil && registroPresupuestales != nil {
			var outputOrdenPago []interface{}
			done := make(chan interface{})
			defer close(done)
			resch := optimize.GenChanInterface(registroPresupuestales...)
			chlistaOrdenes := optimize.Digest(done, searchOrdenPagoByRpId, resch, parametro)
			for arrayOrdenPago := range chlistaOrdenes {
				if dataOrden, e := arrayOrdenPago.([]interface{}); e {
					for _, ordenPago := range dataOrden {
						outputOrdenPago = append(outputOrdenPago, ordenPago.(interface{}))
					}
				}
			}
			if outputOrdenPago != nil {
				respuestas["OrdenPago"] = outputOrdenPago
				c.Data["json"] = respuestas
			} else {
				c.Data["json"] = nil
			}
		} else {
			c.Data["json"] = nil
		}
	} else {
		c.Ctx.Output.SetStatus(201)
		c.Data["json"] = models.Alert{Code: "E_0458", Body: "Not enough parameter in GetOrdenPagoByFuenteFinanciamiento", Type: "error"}
	}
	c.ServeJSON()

}

func searchOrdenPagoByRpId(inputRegistroPresupuestal interface{}, params ...interface{}) (res interface{}) {
	unidadEjecutoraId, e1 := params[0].(string)
	rowRp, e2 := inputRegistroPresupuestal.(map[string]interface{})
	if e1 && e2 {
		var ordenesPagos []interface{}
		// seach dependencia de necesidad
		if necesidad := getNecesidadDesdeRp(inputRegistroPresupuestal, unidadEjecutoraId); necesidad != nil {
			if areaNecesidad := getAreaDeNecesidad(necesidad); areaNecesidad != nil {
				//op
				if err := getJson("http://"+beego.AppConfig.String("kronosService")+"orden_pago/?query=RegistroPresupuestal.Id:"+strconv.Itoa(int(rowRp["Id"].(float64))), &ordenesPagos); err == nil && ordenesPagos != nil {
					for _, orden := range ordenesPagos {
						row := orden.(map[string]interface{})
						row["RegistroPresupuestal"] = rowRp
						row["Necesidad"] = areaNecesidad
					}
					return ordenesPagos
				}
			}
		}
	}
	return
}
