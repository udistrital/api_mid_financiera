package controllers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/udistrital/api_mid_financiera/models"
	"github.com/udistrital/api_mid_financiera/utilidades"
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
	var finicio time.Time
	var ffin time.Time
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &pacData); err == nil {
		utilidades.FillStruct(pacData["inicio"], &finicio)
		utilidades.FillStruct(pacData["fin"], &ffin)
		if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/rubro/RubroReporte", "POST", &reporteData, &pacData); err == nil {

			var ingresos []map[string]interface{}
			err := utilidades.FillStruct(reporteData["ingresos"], &ingresos)
			if err == nil {
				for _, ingresosRow := range ingresos { //recorrer los datos del reporte de ingresos para el rango actual
					var reporteIngresoData []map[string]interface{}
					err := utilidades.FillStruct(ingresosRow["reporte"], &reporteIngresoData)
					if err == nil {
						for _, reporteRow := range reporteIngresoData {
							var valor string
							var mes int
							err := utilidades.FillStructDeep(reporteRow["valores"].(map[string]interface{}), "valor", &valor)
							err = utilidades.FillStruct(reporteRow["mes"], &mes)
							if err == nil {
								fechaInicio := time.Date(finicio.Year(), time.Month(mes), finicio.Day(), 0, 0, 0, 0, time.Local)
								fechaFin := time.Date(finicio.Year(), time.Month(mes+1), finicio.Day(), 0, 0, 0, 0, time.Local)
								idFuente := 0
								if ingresosRow["idfuente"] != nil {
									err := utilidades.FillStruct(ingresosRow["idfuente"], &idFuente)
									if err != nil {
										fmt.Println("err 2")
										c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
									}
								}
								if fechaFin.After(ffin) {
									fechaFin = ffin
								}
								fmt.Println("Apropiacion: ", ingresosRow["id"])
								fmt.Println("Fuente: ", idFuente)
								fmt.Println("finicio: ", fechaInicio)
								fmt.Println("ffin: ", fechaFin)
							} else {
								fmt.Println("err 2")
								c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
							}
						}
					} else {
						fmt.Println("err 2")
						c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
					}
				}
			} else {
				fmt.Println("err 2")
				c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
			}

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
