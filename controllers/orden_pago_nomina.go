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

// MidCrearOPNomina ...
// @Title Create
// @Description create Orden_pago_planta
// @Param	body		body 	models.Orden_pago_planta	true		"body for Orden_pago_planta content"
// @Success 201 {object} models.Orden_pago_planta
// @Failure 403 body is empty
// @router MidCrearOPNomina [post]
func (c *OrdenPagoNominaController) MidCrearOPNomina() {
	var alerta models.Alert
	var v interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		m := v.(map[string]interface{})
		var DetallePreliquidacion []interface{}

		OrdenPago, e := m["OrdenPago"].(map[string]interface{})
		if e != true {
			alerta.Type = "error"
			alerta.Code = "E_OPN_01_2"
			alerta.Body = err.Error()
			c.Data["json"] = alerta
			c.ServeJSON()
			return
		}
		Preliquidacion := OrdenPago["Liquidacion"].(float64)
		Usuario, e := m["Usuario"].(map[string]interface{})
		if e != true {
			alerta.Type = "error"
			alerta.Code = "E_OPN_01_2"
			alerta.Body = err.Error()
			c.Data["json"] = alerta
			c.ServeJSON()
			return
		}
		// get data titan
		if err := getJson("http://"+beego.AppConfig.String("titanService")+"detalle_preliquidacion?query=Preliquidacion:"+strconv.FormatFloat(Preliquidacion, 'f', 0, 64)+"&sortby=Concepto&order=desc", &DetallePreliquidacion); err == nil {
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
			OrdenPago          interface{}
			DetalleLiquidacion []interface{}
			Usuario            interface{}
		}
		sendData2Kronos := Send{OrdenPago: OrdenPago, DetalleLiquidacion: DetallePreliquidacion, Usuario: Usuario}
		var outputData interface{}
		//Envia data to kronos
		if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/orden_pago/RegistrarOpNomina", "POST", &outputData, &sendData2Kronos); err == nil {
		} else {
			alerta.Type = "error"
			alerta.Code = "E_OPN_01_5"
			alerta.Body = ""
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

// CrearOPSeguridadSocial ...
// @Title CrearOPSeguridadSocial
// @Description create Orden Pago Seguridad Social
// @Param	body		body 	models.Orden_pago_planta	true		"body for Orden_pago_planta content"
// @Success 201 {object} models.Orden_pago_planta
// @Failure 403 body is empty
// @router CrearOPSeguridadSocial [post]
func (c *OrdenPagoNominaController) CrearOPSeguridadSocial() {
	var alerta models.Alert
	var v interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		m := v.(map[string]interface{})
		var DataSeguridadSocial map[string]interface{}
		var DataOrdenPago map[string]interface{}
		var PagosSeguridadSocial []interface{}
		var PeriodoPago []interface{}
		//
		err = utilidades.FillStruct(m["SeguridadSocial"], &DataSeguridadSocial)
		Mes := fmt.Sprintf("%v", DataSeguridadSocial["Mes"])
		Anio := fmt.Sprintf("%v", DataSeguridadSocial["Vigencia"])
		err = utilidades.FillStruct(m["OrdenPago"], &DataOrdenPago)

		fmt.Print(Mes)
		fmt.Print("-")
		fmt.Print(Anio)
		// get id periodo pago
		//fmt.Println("\n", "http://"+beego.AppConfig.String("SsService")+"periodo_pago?query=Mes:"+Mes+",Anio:"+Anio)
		if err = getJson("http://"+beego.AppConfig.String("SsService")+"periodo_pago?query=Mes:"+Mes+",Anio:"+Anio, &PeriodoPago); err == nil {
		} else {
			alerta.Type = "error"
			alerta.Code = "E_OPN_01_3"
			alerta.Body = err.Error()
			c.Data["json"] = alerta
			c.ServeJSON()
			return
		}
		fmt.Println("\nAAAAAAAAAAA \nPeriodoPago")
		fmt.Print(PeriodoPago)
		fmt.Println("\nAAAAAAAAAAA")

		// // get data administarativa seguridad social
		// debe ser por mes y año el filtro, en el momento el api no cuenta con esos datos.
		if err = getJson("http://"+beego.AppConfig.String("SsService")+"pago?query=PeriodoPago.Id:1", &PagosSeguridadSocial); err == nil {
		} else {
			alerta.Type = "error"
			alerta.Code = "E_OPN_01_3"
			alerta.Body = err.Error()
			c.Data["json"] = alerta
			c.ServeJSON()
			return
		}

		fmt.Println("\nAAAAAAAAAAA \nPagosSeguridadSocial")
		fmt.Print(PagosSeguridadSocial)
		fmt.Println("\nAAAAAAAAAAA")

		// Control si no existe detalle de liquidacion
		if len(PagosSeguridadSocial) == 0 {
			alerta.Type = "error"
			alerta.Code = "E_OPN_01_4"
			alerta.Body = ""
			c.Data["json"] = alerta
			c.ServeJSON()
			return
		}
		// estructura para enviar data a kronos
		type SendData struct {
			OrdenPago            interface{}
			PagosSeguridadSocial []interface{}
		}
		SendData2Kronos := SendData{OrdenPago: DataOrdenPago, PagosSeguridadSocial: PagosSeguridadSocial}
		var outputData interface{}
		//Envia data to kronos
		if err = sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/orden_pago/RegistrarOpSeguridadSocial", "POST", &outputData, &SendData2Kronos); err == nil {
		} else {
			alerta.Type = "error"
			alerta.Code = "E_OPN_01_5"
			alerta.Body = ""
			c.Data["json"] = alerta
			c.ServeJSON()
			return
		}
		c.Data["json"] = outputData
		c.ServeJSON()
		//fin
	} else {
		alerta.Type = "error"
		alerta.Code = "E_OPN_01_1"
		alerta.Body = err.Error()
		c.Data["json"] = alerta
		c.ServeJSON()
		return
	}
}
