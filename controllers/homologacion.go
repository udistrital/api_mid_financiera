package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/api_mid_financiera/models"
	"github.com/udistrital/utils_oas/request"
)

// HomologacionController operations for homologation fo liquidation of titan
type HomologacionController struct {
	beego.Controller
}

// URLMapping ...
func (c *HomologacionController) URLMapping() {
	c.Mapping("MidHomologacionLiquidacion", c.MidHomologacionLiquidacion)
}

// MidHomologacionLiquidacion ...
// @Title MidHomologacionLiquidacion
// @Description homologa conceptos de titan con conceptos kronos
// @Param	idPreliquidacion	identificador de la liquidación de titan
// @Param	body		body 	models.IdLiquidacion, models.RegistroPresupuestal	"body for Homologacion content"
// @Success 201 {object} models.Conceptos
// @Failure 403 body is empty
// @router MidHomologacionLiquidacion [post]
func (c *HomologacionController) MidHomologacionLiquidacion() {
	var alerta models.Alert
	var v interface{}
	var DetallePreliquidacion []interface{}
	var outputData interface{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		m := v.(map[string]interface{})
		idPreliquidacion, e := m["IdLiquidacion"].(float64)
		if e != true {
			alerta.Type = "error"
			alerta.Code = "E_OPN_01_2"
			alerta.Body = err.Error()
			c.Data["json"] = alerta
			c.ServeJSON()
			return
		}
		registroPresupuestal, e := m["RegistroPresupuestal"].(map[string]interface{})
		if e != true {
			alerta.Type = "error"
			alerta.Code = "E_OPN_01_2"
			alerta.Body = err.Error()
			c.Data["json"] = alerta
			c.ServeJSON()
			return
		}

		// get data titan
		if err := request.GetJson("http://"+beego.AppConfig.String("titanService")+"detalle_preliquidacion?query=Preliquidacion:"+strconv.FormatFloat(idPreliquidacion, 'f', 0, 64)+"&sortby=Concepto&order=desc&limit=-1", &DetallePreliquidacion); err == nil {
		} else {
			alerta.Type = "error"
			alerta.Code = "E_OPN_01_3"
			alerta.Body = err.Error()
			c.Data["json"] = alerta
			c.ServeJSON()
			return
		}
		// Control si no existe detalle de liquidacion
		if len(DetallePreliquidacion) == 0 {
			alerta.Type = "error"
			alerta.Code = "E_OPN_01_4"
			alerta.Body = ""
			c.Data["json"] = alerta
			c.ServeJSON()
			return
		}

		// estructura para enviar data a kronos
		type Send struct {
			DetalleLiquidacion   []interface{}
			RegistroPresupuestal interface{}
		}
		sendData2Kronos := Send{DetalleLiquidacion: DetallePreliquidacion, RegistroPresupuestal: registroPresupuestal}
		//Envia data to kronos
		if err := request.SendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/homologacion_concepto/HomolgacionConceptosTitan/", "POST", &outputData, &sendData2Kronos); err == nil {
		} else {
			alerta.Type = "error"
			alerta.Code = "E_OPN_01_5"
			alerta.Body = "[POST] /homologacion_concepto/HomolgacionConceptosTitan/"
			c.Data["json"] = alerta
			c.ServeJSON()
			return
		}
		c.Data["json"] = outputData
		c.ServeJSON()
	} else {
		alerta.Type = "error"
		alerta.Code = "E_OPN_01_1"
		alerta.Body = err.Error()
		c.Data["json"] = alerta
		c.ServeJSON()
		return
	}
}
