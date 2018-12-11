package controllers

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/manucorporat/try"
	"github.com/udistrital/api_mid_financiera/models"
	"github.com/udistrital/utils_oas/optimize"
	"github.com/udistrital/utils_oas/request"
)

// AvanceController operations for Avance
type AvanceController struct {
	beego.Controller
}

// URLMapping ...
func (c *AvanceController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Avance
// @Param	body		body 	models.Avance	true		"body for Avance content"
// @Success 201 {object} models.Avance
// @Failure 403 body is empty
// @router / [post]
func (c *AvanceController) Post() {

}

// GetOne ...
// @Title GetOne
// @Description get Avance by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Avance
// @Failure 403 :id is empty
// @router /:id [get]
func (c *AvanceController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get Avance
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Avance
// @Failure 403
// @router / [get]
func (c *AvanceController) GetAll() {

}

// Put ...
// @Title Put
// @Description update the Avance
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Avance	true		"body for Avance content"
// @Success 200 {object} models.Avance
// @Failure 403 :id is not int
// @router /:id [put]
func (c *AvanceController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Avance
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *AvanceController) Delete() {

}

// GetAvanceByID ...
// @Title GetAvanceByID
// @Description get All information of an advance payment by id
// @Param	idAvance	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	idAvceLeg	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Success 200 {object} models.Legalizacion_avance
// @Failure 403 :id is empty
// @router /GetAvanceById [get]
func (c *AvanceController) GetAvanceById() {
	var idStr string
	var idAvceLeg string
	if v := c.GetString("idAvance"); v != "" {
		idStr = v
	}
	if v := c.GetString("idAvceLeg"); v != "" {
		idAvceLeg = v
	}
	defer c.ServeJSON()
	var solicitudAvance map[string]interface{}
	var infoBeneficiario map[string]interface{}
	var total float64
	var resTipo []interface{}
	var params []interface{}
	var valorLegalizado float64
	params = append(params, total)
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/solicitud_avance/"+idStr, &solicitudAvance); err == nil {
		if solicitudAvance != nil {
			if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/solicitud_tipo_avance/?query=SolicitudAvance.Id:"+idStr+"&sortby=Id&limit=-1&order=asc", &resTipo); err == nil {
				if resTipo != nil {
					solicitudAvance["Tipos"] = optimize.ProccDigest(resTipo, getTiposInfo, params, 3)
					solicitudAvance["Total"] = params[0]
				}
			} else {
				beego.Error("Error", err.Error())
			}

			if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/avance_legalizacion_tipo/GetLegalizationValue/"+idAvceLeg, &valorLegalizado); err == nil {
				solicitudAvance["valorLegalizado"] = valorLegalizado
			} else {
				beego.Error("Error", err.Error())
			}
			idBenStr := strconv.FormatFloat(solicitudAvance["Beneficiario"].(float64), 'f', -1, 64)
			if err := request.GetJson("http://10.20.0.127/urano/index.php?data=B-7djBQWvIdLAEEycbH1n6e-3dACi5eLUOb63vMYhGq0kPBs7NGLYWFCL0RSTCu1yTlE5hH854MOgmjuVfPWyvdpaJDUOyByX-ksEPFIrrQQ7t1p4BkZcBuGD2cgJXeD&documento="+idBenStr, &infoBeneficiario); err == nil {
				if infoBeneficiario != nil {
					solicitudAvance["Tercero"] = infoBeneficiario
				}
			} else {
				beego.Error("Error", err.Error())
			}
		}
	} else {
		c.Data["json"] = models.Alert{Type: "error", Code: "E_0458", Body: err}
		beego.Error("Error", err.Error())
		return
	}
	c.Data["json"] = solicitudAvance
}

func getTiposInfo(rpintfc interface{}, params ...interface{}) (res interface{}) {
	var requisitos interface{}
	tipoId := strconv.FormatFloat(rpintfc.(map[string]interface{})["TipoAvance"].(map[string]interface{})["Id"].(float64), 'f', -1, 64)
	params[0] = params[0].(float64) + rpintfc.(map[string]interface{})["Valor"].(float64)
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/requisito_tipo_avance/?query=TipoAvance:"+tipoId+",Activo:1&limit=-1&fields:RequisitoAvance,TipoAvance,Id&sortby=TipoAvance&order=asc", &requisitos); err == nil {
		if requisitos != nil {
			rpintfc.(map[string]interface{})["Requisitos"] = requisitos
		}
	} else {
		beego.Error("Error", err.Error())
	}
	return rpintfc
}

// GetSolicitudes ...
// @Title GetSolicitudes
// @Description get All information of an advance payment by id
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Legalizacion_avance
// @Failure 403 body is empty
// @router /GetSolicitudes [get]
func (c *AvanceController) GetSolicitudes() {

	var solicitudes []interface{}
	var limit int64 = 10
	var offset int64
	var query string
	var sortby string
	var order string
	defer c.ServeJSON()

	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	if r := c.GetString("query"); r != "" {
		query = r
	}
	if v := c.GetString("sortby"); v != "" {
		sortby = "&sortby=" + v
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = "&order=" + v
	}

	try.This(func() {
		if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/solicitud_avance/?limit="+strconv.FormatInt(limit, 10)+"&offset="+strconv.FormatInt(offset, 10)+"&query="+query+sortby+order, &solicitudes); err == nil {
			if solicitudes != nil {
				response := optimize.ProccDigest(solicitudes, getSolicitudInfo, nil, 3)
				c.Data["json"] = response
			}
		}
	}).Catch(func(e try.E) {
		beego.Error("expc ", e)
		c.Data["json"] = models.Alert{Type: "error", Code: "E_0458", Body: e}
	})
}

func getSolicitudInfo(rpintfc interface{}, params ...interface{}) (res interface{}) {
	var estadosSolicitud []interface{}
	var solTipoAvance []interface{}
	var resDocente map[string]interface{}
	var requisitosTipo []interface{}
	solicitudId := strconv.FormatFloat(rpintfc.(map[string]interface{})["Id"].(float64), 'f', -1, 64)
	beneficiarioID := strconv.FormatFloat(rpintfc.(map[string]interface{})["Beneficiario"].(float64), 'f', -1, 64)
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/avance_estado_avance/?query=SolicitudAvance.Id:"+solicitudId+"&limit=-1&sortby=FechaRegistro&order=desc", &estadosSolicitud); err == nil {
		if estadosSolicitud != nil {
			rpintfc.(map[string]interface{})["Estado"] = estadosSolicitud
		}
	} else {
		beego.Error("Error", err.Error())
	}

	if err := request.GetJsonWSO2(beego.AppConfig.String("Wso2Service")+"servicios_academicos/consulta_datos_docente_planta/"+beneficiarioID, &resDocente); err == nil {
		if resDocente["datosCollection"].(map[string]interface{})["datos"] != nil {
			rpintfc.(map[string]interface{})["Tercero"] = resDocente["datosCollection"].(map[string]interface{})["datos"].([]interface{})[0]
		}
	} else {
		beego.Error("Error", err.Error())
	}
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/avance_estado_avance/?query=SolicitudAvance.Id:"+solicitudId+"&limit=-1&sortby=FechaRegistro&order=desc", &estadosSolicitud); err == nil {
		if estadosSolicitud != nil {
			rpintfc.(map[string]interface{})["Estado"] = estadosSolicitud
		}
	} else {
		beego.Error("Error", err.Error())
	}

	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/solicitud_tipo_avance/?query=SolicitudAvance.Id:"+solicitudId+"&limit=-1&sortby=Id&order=asc", &solTipoAvance); err == nil {
		if solTipoAvance != nil {
			if rpintfc.(map[string]interface{})["Total"] == nil {
				rpintfc.(map[string]interface{})["Total"] = float64(0)
			}
			var Sol float64
			var Leg float64

			rpintfc.(map[string]interface{})["Tipos"] = solTipoAvance
			for _, v := range solTipoAvance {
				rpintfc.(map[string]interface{})["Total"] = rpintfc.(map[string]interface{})["Total"].(float64) + v.(map[string]interface{})["Valor"].(float64)
				strIdTipoAvc := strconv.FormatFloat(v.(map[string]interface{})["TipoAvance"].(map[string]interface{})["Id"].(float64), 'f', -1, 64)
				solicitudTipoAvcID := v.(map[string]interface{})["Id"]
				if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/requisito_tipo_avance/?query=TipoAvance:"+strIdTipoAvc+",Activo:1&fields=RequisitoAvance,TipoAvance,Id&limit=-1&sortby=TipoAvance&order=asc", &requisitosTipo); err == nil {
					v.(map[string]interface{})["Requisitos"] = requisitosTipo
					for _, vr := range requisitosTipo {
						requisitoTipoAvance := vr.(map[string]interface{})["RequisitoAvance"]
						idEtapaAvance := requisitoTipoAvance.(map[string]interface{})["EtapaAvance"].(map[string]interface{})["Id"].(float64)
						switch idEtapaAvance {
						case 1:
							Sol += 1
						case 2:
							Leg += 1
						}
						vr.(map[string]interface{})["SolicitudTipoAvance"] = map[string]interface{}{"Id": solicitudTipoAvcID}
						vr.(map[string]interface{})["RequisitoTipoAvance"] = map[string]interface{}{"Id": requisitoTipoAvance.(map[string]interface{})["Id"]}
					}
				} else {
					beego.Error("Error ", err.Error())
				}
				v.(map[string]interface{})["n_solicitar"] = Sol
				v.(map[string]interface{})["n_legalizar"] = Leg
			}
		}
	} else {
		beego.Error("Error", err.Error())
	}

	return rpintfc
}
