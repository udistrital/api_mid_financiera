package controllers

import (

	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/udistrital/api_mid_financiera/utilidades"
	"strconv"
)

// Orden_pago_plantaController operations for Orden_pago_planta
type Orden_pago_plantaController struct {
	beego.Controller
}

// URLMapping ...
func (c *Orden_pago_plantaController) URLMapping() {
	c.Mapping("Post", c.Post)

}

// Post ...
// @Title Create
// @Description create Orden_pago_planta
// @Param	body		body 	models.Orden_pago_planta	true		"body for Orden_pago_planta content"
// @Success 201 {object} models.Orden_pago_planta
// @Failure 403 body is empty
// @router / [post]
func (c *Orden_pago_plantaController) Post() {
	var v interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil{
		m := v.(map[string]interface{})
		var detalle []interface{}
		var id_liquidacion int
		utilidades.FillStruct(m["Liquidacion"], &id_liquidacion)
		if err := getJson("http://"+beego.AppConfig.String("titanService")+"detalle_liquidacion?query=Liquidacion:"+strconv.Itoa(id_liquidacion)+"&sortby=Concepto&order=desc", &detalle); err == nil {
		}else{
			c.Data["json"] = err.Error()
		}
		type Send struct{
			OrdenPago map[string]interface{}
			DetalleLiquidacion []interface{}
		}
		total := Send{OrdenPago: m, DetalleLiquidacion:detalle,}
		//kronos
		if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/disponibilidad/SaldoCdp", "POST", &saldoCDP, &datos); err == nil {
		}else{

		}

		c.Data["json"] = total
		c.ServeJSON()
	}
}
