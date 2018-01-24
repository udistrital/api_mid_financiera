package controllers

import (
	// "fmt"
	// "strconv"
	// "strings"

	"strconv"

	"github.com/astaxie/beego"
	// "github.com/udistrital/api_mid_financiera/utilidades"
)

// NecesidadController operations for Necesidad
type NecesidadController struct {
	beego.Controller
}

// URLMapping ...
// func (c *NecesidadController) URLMapping() {
// 	c.Mapping("GetNecesidadDesdeRp", c.GetNecesidadDesdeRp)
// }
func getNecesidadDesdeRp(registroPresupuestal interface{}) (outputNecesidad interface{}, outputError interface{}) {
	beego.Info("getNecesidadDesdeRp")
	if rowRp, e := registroPresupuestal.(map[string]interface{}); e {
		var solicitudRp []interface{}
		if err := getJson("http://"+beego.AppConfig.String("argoService")+"solicitud_rp/?query=Id:"+strconv.Itoa(int(rowRp["Solicitud"].(float64)))+"&limit:1", &solicitudRp); err == nil && solicitudRp != nil {
			beego.Info("solicitudRp: ", solicitudRp[0].(map[string]interface{})["Id"])
			rowSolicitudRp := solicitudRp[0].(map[string]interface{})
			// disponibilidad
			var disponibilidad []interface{}
			if err := getJson("http://"+beego.AppConfig.String("kronosService")+"disponibilidad/?query=Id:"+strconv.Itoa(int(rowSolicitudRp["Cdp"].(float64)))+"&limit:1", &disponibilidad); err == nil && disponibilidad != nil {
				beego.Info("Disponibilidad: ", disponibilidad[0].(map[string]interface{})["Id"])
				rowDisponibilidad := disponibilidad[0].(map[string]interface{})
				// Solicitud de disponibilidad
				var solicitudDisponibilidad []interface{}
				if err := getJson("http://"+beego.AppConfig.String("argoService")+"solicitud_disponibilidad/?query=Id:"+strconv.Itoa(int(rowDisponibilidad["Solicitud"].(float64)))+"&limit:1", &solicitudDisponibilidad); err == nil && solicitudDisponibilidad != nil {
					beego.Info("solicitudDisponibilidad: ", solicitudDisponibilidad[0].(map[string]interface{})["Id"])
					rowSolicitudDisponibilidad := solicitudDisponibilidad[0].(map[string]interface{})
					//dependencia_necesidad
					var dependenciaNecesidad []interface{}
					if err := getJson("http://"+beego.AppConfig.String("argoService")+"dependencia_necesidad/?query=Necesidad.Id:"+strconv.Itoa(int(rowSolicitudDisponibilidad["Necesidad"].(map[string]interface{})["Id"].(float64)))+"&limit:1", &dependenciaNecesidad); err == nil && dependenciaNecesidad != nil {
						beego.Info("DependenciaNecesidad: ", dependenciaNecesidad[0].(map[string]interface{})["Id"])
						rowDependenciaNecesidad := dependenciaNecesidad[0].(map[string]interface{})
						//dependencia
						var dependencia []interface{}
						if err := getJson("http://"+beego.AppConfig.String("oikosService")+"dependencia/?query=Id:"+strconv.Itoa(int(rowDependenciaNecesidad["JefeDependenciaSolicitante"].(float64)))+"&limit:1", &dependencia); err == nil && dependencia != nil {
							beego.Info("Dependencia: ", dependencia[0])
							//rowDependenciaNecesidad := dependencia[0].(map[string]interface{})
						}
					}
				}
			}
		}
	}
	return nil, nil

}
