package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/manucorporat/try"
	"github.com/udistrital/api_mid_financiera/models"
	"github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/utils_oas/request"
)

type MovimientoApropiacionController struct {
	beego.Controller
}

func (c *MovimientoApropiacionController) URLMapping() {

}

// AprobarMovimietnoApropiacion ...
// @Title AprobarMovimietnoApropiacion
// @Description create MovimientoApropiacion
// @Param	body		body 	models.MovimientoApropiacion	true		"body for MovimientoApropiacion content"
// @Success 201 {int} models.MovimientoApropiacion
// @Failure 403 body is empty
// @router /AprobarMovimietnoApropiacion [post]
func (c *MovimientoApropiacionController) AprobarMovimietnoApropiacion() {
	try.This(func() {
		var v map[string]interface{}
		var res interface{}
		if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
			beego.Info("Data Recibed ", v)
			Urlcrud := "http://" + beego.AppConfig.String("Urlcrud") + ":" + beego.AppConfig.String("Portcrud") + "/" + beego.AppConfig.String("Nscrud") + "/movimiento_apropiacion/AprobarMovimietnoApropiacion"
			if err := request.SendJson(Urlcrud, "POST", &res, &v); err == nil {
				beego.Info("Data to Send ", res)
				c.Data["json"] = res
			} else {
				panic(err.Error())
			}
		} else {
			panic(err.Error())
		}
	}).Catch(func(e try.E) {
		beego.Error("catch error registrar valores: ", e)
		var alert []models.Alert
		alt := models.Alert{}
		alt.Code = "E_0458"
		alt.Body = e
		alt.Type = "error"
		alert = append(alert, alt)
		c.Data["json"] = alert
	})
	c.ServeJSON()
}

func AddMovimientoApropiacionMongo(parameter ...interface{}) (err interface{}) {
	idMov := 0.0
	try.This(func() {
		dataMongo := make(map[string]interface{})
		infoMovimiento := parameter[0].(map[string]interface{})["Movimiento"].(map[string]interface{})
		idMov = infoMovimiento["Id"].(float64)
		var afectacionArr []map[string]interface{}

		dataMongo["FechaMovimiento"] = infoMovimiento["FechaMovimiento"]
		dataMongo["Vigencia"] = infoMovimiento["Vigencia"]
		dataMongo["UnidadEjecutora"] = infoMovimiento["UnidadEjecutora"]
		dataMongo["Id"] = infoMovimiento["Id"]
		var movimientos []map[string]interface{}
		err1 := formatdata.FillStruct(infoMovimiento["MovimientoApropiacionDisponibilidadApropiacion"], &movimientos)
		if err1 != nil {
			panic(err1.Error())
		}
		for _, data := range movimientos {
			var CuentaContraCredito string
			var CuentaCredito string
			if CuentaContraCreditoInt, e := data["CuentaContraCredito"].(map[string]interface{}); e {
				CuentaContraCredito = CuentaContraCreditoInt["Rubro"].(map[string]interface{})["Codigo"].(string)
			}
			if CuentaCreditoInt, e := data["CuentaCredito"].(map[string]interface{}); e {
				CuentaCredito = CuentaCreditoInt["Rubro"].(map[string]interface{})["Codigo"].(string)
			}
			Valor := data["Valor"]
			TipoMovimiento := data["TipoMovimientoApropiacion"].(map[string]interface{})["Nombre"]
			afectacion := map[string]interface{}{
				"CuentaContraCredito": CuentaContraCredito,
				"CuentaCredito":       CuentaCredito,
				"Valor":               Valor,
				"TipoMovimiento":      TipoMovimiento,
			}
			afectacionArr = append(afectacionArr, afectacion)
		}
		dataMongo["Afectacion"] = afectacionArr
		beego.Info("Data Send ", dataMongo)
	}).Catch(func(e try.E) {
		var resC interface{}
		infoMovimiento := parameter[0].(map[string]interface{})["Movimiento"].(map[string]interface{})
		Urlcrud := "http://" + beego.AppConfig.String("Urlcrud") + ":" + beego.AppConfig.String("Portcrud") + "/" + beego.AppConfig.String("Nscrud") + "/movimiento_apropiacion/" + strconv.Itoa(int(idMov))
		estadoMov := infoMovimiento["EstadoMovimientoApropiacion"].(map[string]interface{})
		estadoMov["Id"] = 1
		infoMovimiento["EstadoMovimientoApropiacion"] = estadoMov
		request.SendJson(Urlcrud, "PUT", &resC, &infoMovimiento)
		beego.Error("error job ", e)
	})
	return
}
