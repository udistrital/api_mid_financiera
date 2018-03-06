package controllers

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/request"
)

// NecesidadController operations for Necesidad
type NecesidadController struct {
	beego.Controller
}

func getNecesidadDesdeRp(registroPresupuestal interface{}, unidadEjecutora string) (outputNecesidad interface{}) {
	if rowRp, e := registroPresupuestal.(map[string]interface{}); e {
		var solicitudRp []interface{}
		if err := request.GetJson("http://"+beego.AppConfig.String("argoService")+"solicitud_rp/?query=Id:"+strconv.Itoa(int(rowRp["Solicitud"].(float64)))+"&limit:1", &solicitudRp); err == nil && solicitudRp != nil {
			//beego.Info("solicitudRp: ", solicitudRp[0].(map[string]interface{})["Id"])
			rowSolicitudRp := solicitudRp[0].(map[string]interface{})
			// disponibilidad
			var disponibilidad []interface{}
			//beego.Info("http://"+beego.AppConfig.String("kronosService")+"disponibilidad/?query=Id:"+strconv.Itoa(int(rowSolicitudRp["Cdp"].(float64)))+"&limit:1")
			if err := request.GetJson("http://"+beego.AppConfig.String("kronosService")+"disponibilidad/?query=Id:"+strconv.Itoa(int(rowSolicitudRp["Cdp"].(float64)))+"&limit:1", &disponibilidad); err == nil && disponibilidad != nil {
				//beego.Info("Disponibilidad: ", disponibilidad[0].(map[string]interface{})["Id"])
				rowDisponibilidad := disponibilidad[0].(map[string]interface{})
				// Solicitud de disponibilidad
				var solicitudDisponibilidad []interface{}
				numSolicitud := rowDisponibilidad["DisponibilidadProcesoExterno"].([]interface{})[0].(map[string]interface{})["ProcesoExterno"]
				//beego.Info("solicitud", numSolicitud)
				//beego.Info("http://"+beego.AppConfig.String("argoService")+"solicitud_disponibilidad/?query=Id:"+strconv.Itoa(int(numSolicitud.(float64)))+",Necesidad.UnidadEjecutora:"+unidadEjecutora+"&limit:1")
				if err := request.GetJson("http://"+beego.AppConfig.String("argoService")+"solicitud_disponibilidad/?query=Id:"+strconv.Itoa(int(numSolicitud.(float64)))+",Necesidad.UnidadEjecutora:"+unidadEjecutora+"&limit:1", &solicitudDisponibilidad); err == nil && solicitudDisponibilidad != nil {
					beego.Info("solicitudDisponibilidad: ", solicitudDisponibilidad[0].(map[string]interface{})["Id"])
					outputNecesidad := solicitudDisponibilidad[0].(map[string]interface{})["Necesidad"].(map[string]interface{})
					//beego.Info("Necesidad return : ", outputNecesidad)
					return outputNecesidad
				}
			}
		}
	}
	return
}

func getAreaDeNecesidad(necesidad interface{}) (outputAreaNecesidad interface{}) {
	//dependencia_necesidad
	rowNecesidad := necesidad.(map[string]interface{})
	var dependenciaNecesidad []interface{}
	if err := request.GetJson("http://"+beego.AppConfig.String("argoService")+"dependencia_necesidad/?query=Necesidad.Id:"+strconv.Itoa(int(rowNecesidad["Id"].(float64)))+"&limit:1", &dependenciaNecesidad); err == nil && dependenciaNecesidad != nil {
		//beego.Info("DependenciaNecesidad: ", dependenciaNecesidad[0].(map[string]interface{})["Id"])
		rowDependenciaNecesidad := dependenciaNecesidad[0].(map[string]interface{})
		// core jefe dependencia
		var jefeDependencia []interface{}
		if err := request.GetJson("http://"+beego.AppConfig.String("coreService")+"jefe_dependencia/?query=Id:"+strconv.Itoa(int(rowDependenciaNecesidad["JefeDependenciaDestino"].(float64)))+"&limit:1", &jefeDependencia); err == nil && jefeDependencia != nil {
			//beego.Info("JefeDependencia: ", jefeDependencia[0].(map[string]interface{})["Id"])
			rowJefeDependencia := jefeDependencia[0].(map[string]interface{})
			//dependencia
			var dependencia []interface{}
			if err := request.GetJson("http://"+beego.AppConfig.String("oikosService")+"dependencia/?query=Id:"+strconv.Itoa(int(rowJefeDependencia["DependenciaId"].(float64)))+"&limit:1", &dependencia); err == nil && dependencia != nil {
				//beego.Info("Dependencia: ", dependencia[0].(map[string]interface{})["Id"])
				rowNecesidad["Dependencia"] = dependencia[0]
				return rowNecesidad
			}
		}
	}
	return
}
