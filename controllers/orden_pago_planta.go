package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/api_mid_financiera/models"
	"github.com/udistrital/api_mid_financiera/utilidades"
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
			alerta.Code = "E_OPP_01_3"
			alerta.Body = err1.Error()
			c.Data["json"] = alerta
			c.ServeJSON()
		}
		// get data titan
		if err := getJson("http://"+beego.AppConfig.String("titanService")+"detalle_liquidacion?query=Liquidacion:"+strconv.Itoa(idLiquidacion)+"&sortby=Concepto&order=desc", &detalle); err == nil {
		} else {
			alerta.Type = "error"
			alerta.Code = "E_OPP_01_4"
			alerta.Body = err.Error()
			c.Data["json"] = alerta
			c.ServeJSON()
		}
		fmt.Print("detalle:", detalle)
		if len(detalle) == 0 {
			alerta.Type = "error"
			alerta.Code = "E_OPP_01_5"
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
		//fmt.Println("Enviar Data a Kronos ", detalle)
		if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"orden_pago/RegistrarOpPlanta", "POST", &outputData, &total); err == nil {
			fmt.Println("**111111111***")
		} else {
			fmt.Println("Error ----------- ", "http://"+beego.AppConfig.String("titanService")+"detalle_liquidacion?query=Liquidacion:"+strconv.Itoa(idLiquidacion)+"&sortby=Concepto&order=desc")
			fmt.Println(err.Error())
		}
		c.Data["json"] = outputData
		c.ServeJSON()
	}
}
