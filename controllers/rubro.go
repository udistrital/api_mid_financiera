package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/api_mid_financiera/models"
)

type RubroController struct {
	beego.Controller
}

// GenerarPac ...
// @Title GenerarPac
// @Description Get PAC By Rubro
// @Param	pacData		query 	interface{}	true		"objeto con fechas del rango del PAC"
// @Success 200 {object} interface{}
// @Failure 403
// @router /GenerarPac/ [post]
func (c *RubroController) GenerarPac() {
	var pacData map[string]interface{} //definicion de la interface que recibe los datos del reporte y proyecciones
	var reporteData map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &pacData); err == nil {
		if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/rubro/RubroReporte", "POST", &reporteData, &pacData); err == nil {
			c.Data["json"] = reporteData
		} else {
			fmt.Println("err 2")
			c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
		}
	} else {
		fmt.Println("err 1")
		c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
	}
	c.ServeJSON()
}
