package controllers

import (
	// "fmt"
	// "strconv"
	// "strings"

	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/api_financiera/utilidades"
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
	fuente := c.GetString("fuente")
	vigencia := c.GetString("vigencia")
	unidadEjecutora := c.GetString("unidadEjecutora")
	if fuente != "" && vigencia != "" && unidadEjecutora != "" {
		var registro_presupuestales []interface{}
		respuestas := make(map[string]interface{})
		// search registro_presupuestal
		if err := getJson("http://"+beego.AppConfig.String("kronosService")+"registro_presupuestal/?query=RegistroPresupuestalDisponibilidadApropiacion.DisponibilidadApropiacion.FuenteFinanciamiento.Id:"+fuente+",Vigencia:"+vigencia+"&limit:-1", &registro_presupuestales); err == nil && registro_presupuestales != nil {
			println("1111111")
			var outputOrdenPago []interface{}
			done := make(chan interface{})
			defer close(done)
			resch := utilidades.GenChanInterface(registro_presupuestales...)
			chlistaOrdenes := utilidades.Digest(done, searchOrdenPagoByRpId, resch, nil)
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

func searchOrdenPagoByRpId(registro_resupuestal interface{}, params ...interface{}) (res interface{}) {
	rowRp, e := registro_resupuestal.(map[string]interface{})
	if e {
		var ordenes_pagos []interface{}
		// if err := getJson("http://"+beego.AppConfig.String("kronosService")+"orden_pago/?query=RegistroPresupuestal.Id:200", &ordenes_pagos); err == nil && ordenes_pagos != nil {
		// 	return ordenes_pagos
		// }

		if err := getJson("http://"+beego.AppConfig.String("kronosService")+"orden_pago/?query=RegistroPresupuestal.Id:"+strconv.Itoa(int(rowRp["Id"].(float64))), &ordenes_pagos); err == nil && ordenes_pagos != nil {
			for _, orden := range ordenes_pagos {
				row := orden.(map[string]interface{})
				row["RegistroPresupuestal"] = rowRp
			}
			return ordenes_pagos
		}
	}
	return
}
