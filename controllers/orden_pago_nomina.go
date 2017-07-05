package controllers

import (
	"encoding/json"
	"fmt"
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
		// en local no se a logrado comunicar el mid con el api de financiera
		//if err := sendJson("http://127.0.0.1:8084/v1/orden_pago/RegistrarOpNomina", "POST", &outputData, &total); err == nil {
		if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/orden_pago/RegistrarOpNomina", "POST", &outputData, &total); err == nil {
		} else {
			fmt.Println("Error ----------- ", "http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/orden_pago/RegistrarOpNomina")
			fmt.Print(err.Error())
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

// CrearOPSeguridadSocial ...
// @Title Create
// @Description create Orden Pago Seguridad Social
// @Param	body		body 	models.Orden_pago_planta	true		"body for Orden_pago_planta content"
// @Success 201 {object} models.Orden_pago_planta
// @Failure 403 body is empty
// @router CrearOPSeguridadSocial [post]
func (c *OrdenPagoNominaController) CrearOPSeguridadSocial() {
	var alerta models.Alert
	var v interface{}
	if err1 := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err1 == nil {
		m := v.(map[string]interface{})
		var DataSeguridadSocial interface{}
		var PagosSeguridadSocial []interface{}
		//
		err2 := utilidades.FillStruct(m["SeguridadSocial"], &DataSeguridadSocial)
		if err2 != nil {
			alerta.Type = "error"
			alerta.Code = "E_OPN_01_2"
			alerta.Body = err2.Error()
			c.Data["json"] = alerta
			c.ServeJSON()
		}
		fmt.Print(DataSeguridadSocial) //
		// get data administarativa seguridad social
		// debe ser por mes y año el filtro, en el momento el api no cuenta con esos datos.
		fmt.Print("\n http://" + beego.AppConfig.String("SsService") + "pago?query=PeriodoPago.Id:1")
		if err3 := getJson("http://"+beego.AppConfig.String("SsService")+"pago?query=PeriodoPago.Id:1", &PagosSeguridadSocial); err3 == nil {
		} else {
			alerta.Type = "error"
			alerta.Code = "E_OPN_01_3"
			alerta.Body = err3.Error()
			c.Data["json"] = alerta
			c.ServeJSON()
		}
		fmt.Print("\nPAGOS SS:", PagosSeguridadSocial)

	} else {
		alerta.Type = "error"
		alerta.Code = "E_OPN_01_1"
		alerta.Body = err1.Error()
		c.Data["json"] = alerta
		c.ServeJSON()
	}
}
