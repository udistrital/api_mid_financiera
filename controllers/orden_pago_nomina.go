package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/api_mid_financiera/models"
	"github.com/udistrital/api_mid_financiera/utilidades"
)

// OrdenPagoNominaController operations for Orden_pago_planta
type OrdenPagoNominaController struct {
	beego.Controller
}

// URLMapping ...
func (c *OrdenPagoNominaController) URLMapping() {
	c.Mapping("Post", c.Post)

}

// Post ...
// @Title Create
// @Description create Orden_pago_planta
// @Param	body		body 	models.Orden_pago_planta	true		"body for Orden_pago_planta content"
// @Success 201 {object} models.Orden_pago_planta
// @Failure 403 body is empty
// @router / [post]
func (c *OrdenPagoNominaController) Post() {
	var alerta models.Alert
	var v interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		m := v.(map[string]interface{})
		var detalle []interface{}
		var idLiquidacion int
		//
		err1 := utilidades.FillStruct(m["Liquidacion"], &idLiquidacion)
		if err1 != nil {
			alerta.Type = "error"
			alerta.Code = "E_OPN_01_2"
			alerta.Body = err1.Error()
			c.Data["json"] = alerta
			c.ServeJSON()
		}
		// get data titan
		if err := getJson("http://"+beego.AppConfig.String("titanService")+"detalle_liquidacion?query=Liquidacion:"+strconv.Itoa(idLiquidacion)+"&sortby=Concepto&order=desc", &detalle); err == nil {
		} else {
			alerta.Type = "error"
			alerta.Code = "E_OPN_01_3"
			alerta.Body = err.Error()
			c.Data["json"] = alerta
			c.ServeJSON()
		}
		// Control si no existe detalle de liquidacion
		if len(detalle) == 0 {
			alerta.Type = "error"
			alerta.Code = "E_OPN_01_4"
			alerta.Body = ""
			c.Data["json"] = alerta
			c.ServeJSON()
		}
		// estructura para enviar data a kronos
		type Send struct {
			OrdenPago          interface{}
			DetalleLiquidacion []interface{}
		}
		total := Send{OrdenPago: m, DetalleLiquidacion: detalle}
		var outputData interface{}
		//Envia data to kronos
		if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"orden_pago/RegistrarOpPlanta", "POST", &outputData, &total); err == nil {
		} else {
			alerta.Type = "error"
			alerta.Code = "E_OPN_01_5"
			alerta.Body = ""
			c.Data["json"] = alerta
			c.ServeJSON()
		}
		c.Data["json"] = outputData
		c.ServeJSON()
	} else {
		alerta.Type = "error"
		alerta.Code = "E_OPN_01_1"
		alerta.Body = err.Error()
		c.Data["json"] = alerta
		c.ServeJSON()
	}
}
