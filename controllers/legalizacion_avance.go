package controllers

import (
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/fatih/structs"
	"github.com/udistrital/api_mid_financiera/models"
	"github.com/udistrital/utils_oas/optimize"
	"github.com/udistrital/utils_oas/request"
)

// LegalizacionAvanceController operations for Legalizacion_avance
type LegalizacionAvanceController struct {
	beego.Controller
}

// URLMapping ...
func (c *LegalizacionAvanceController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Legalizacion_avance
// @Param	body		body 	models.Legalizacion_avance	true		"body for Legalizacion_avance content"
// @Success 201 {object} models.Legalizacion_avance
// @Failure 403 body is empty
// @router / [post]
func (c *LegalizacionAvanceController) Post() {

}

// GetOne ...
// @Title GetOne
// @Description get Legalizacion_avance by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Legalizacion_avance
// @Failure 403 :id is empty
// @router /:id [get]
func (c *LegalizacionAvanceController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get Legalizacion_avance
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Legalizacion_avance
// @Failure 403
// @router / [get]
func (c *LegalizacionAvanceController) GetAll() {

}

// Delete ...
// @Title Delete
// @Description delete the Legalizacion_avance
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *LegalizacionAvanceController) Delete() {

}

// GetAllLegalizacionTipo ...
// @Title GetAllLegalizacionTipo
// @Description get Legalizacion_avance
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Legalizacion_avance
// @Failure 403
// @router /GetAllLegalizacionTipo [get]
func (c *LegalizacionAvanceController) GetAllLegalizacionTipo() {

	defer c.ServeJSON()
	var legalizaciones []interface{}
	var limit int64 = 10
	var offset int64
	var query string
	var respuesta interface{}
	// limit: 10 (default is 10)
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

	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/avance_legalizacion_tipo?limit="+strconv.FormatInt(limit, 10)+"&offset="+strconv.FormatInt(offset, 10)+"&query="+query, &legalizaciones); err == nil {
		if legalizaciones != nil {
			respuesta = optimize.ProccDigest(legalizaciones, formatoLegalizacion, nil, 3)
			c.Ctx.Output.SetStatus(201)
		} else {
			respuesta = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
		}
	}
	c.Data["json"] = respuesta
}

func formatoLegalizacion(legalizacion interface{}, params ...interface{}) (res interface{}) {
	legalMap := legalizacion.(map[string]interface{})
	tipoLegAvance := int(legalMap["TipoAvanceLegalizacion"].(map[string]interface{})["Id"].(float64))
	var conceptoAvanceLegalizacion []map[string]interface{}
	if f := formatoLegalizacionDispatcher(tipoLegAvance); f != nil {
		res = f(legalMap, nil)
	} else {
		res = legalMap
	}
	idAvaLeg := strconv.FormatFloat(res.(map[string]interface{})["Id"].(float64), 'f', -1, 64)
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/concepto_avance_legalizacion_tipo/?query=AvanceLegalizacion.Id:"+idAvaLeg, &conceptoAvanceLegalizacion); err == nil {
		res.(map[string]interface{})["Valor"] = conceptoAvanceLegalizacion[0]["Valor"]
	} else {
		beego.Error("Error", err.Error())
	}
	return
}
func formatoLegalizacionDispatcher(tipo int) (f func(data map[string]interface{}, params ...interface{}) interface{}) {
	switch tipo {
	case 2:
		return getLegalizacionCompra
	case 1:
		return getLegalizacionPracticaAcadem
	default:
		return nil
	}

}
func getLegalizacionCompra(data map[string]interface{}, params ...interface{}) (res interface{}) {
	var resProveedor interface{}
	var resPersonaNat interface{}
	var tipoDocumento interface{}

	tercero := data["Tercero"].(string)
	if err := request.GetJson("http://"+beego.AppConfig.String("AdministrativaAmazonService")+"informacion_proveedor/"+tercero, &resProveedor); err == nil {
		if resProveedor != nil {
			beego.Error(resProveedor)
			numberIdStr := resProveedor.(map[string]interface{})["NumDocumento"].(string)
			resProveedor.(map[string]interface{})["numero_documento"] = resProveedor.(map[string]interface{})["NumDocumento"].(string)
			if resProveedor.(map[string]interface{})["Tipopersona"].(string) == "NATURAL" {
				if err := request.GetJson("http://"+beego.AppConfig.String("AdministrativaAmazonService")+"informacion_persona_natural/"+numberIdStr, &resPersonaNat); err == nil {
					idTipoDoc := strconv.FormatFloat(resPersonaNat.(map[string]interface{})["TipoDocumento"].(map[string]interface{})["Id"].(float64), 'f', -1, 64)
					if resPersonaNat != nil {
						if err := request.GetJson("http://"+beego.AppConfig.String("AdministrativaAmazonService")+"parametro_estandar/"+idTipoDoc, &tipoDocumento); err == nil {
							resProveedor.(map[string]interface{})["tipoDocTercero"] = tipoDocumento
							resProveedor.(map[string]interface{})["tipo_documento"] = tipoDocumento.(map[string]interface{})["ValorParametro"]
						}
					}
				} else {
					beego.Error("Error" + err.Error())
				}
			} else {
				if resProveedor.(map[string]interface{})["Tipopersona"].(string) == "JURIDICA" {
					if err := request.GetJson("http://"+beego.AppConfig.String("AdministrativaAmazonService")+"parametro_estandar/11", &tipoDocumento); err == nil {
						resProveedor.(map[string]interface{})["tipoDocTercero"] = tipoDocumento
						resProveedor.(map[string]interface{})["tipo_documento"] = tipoDocumento.(map[string]interface{})["ValorParametro"]
					} else {
						beego.Error("Error" + err.Error())
					}
				}
			}
		}
		data["InformacionProveedor"] = resProveedor
	} else {
		beego.Error("Error", err.Error())
	}
	return data
}

func getLegalizacionPracticaAcadem(data map[string]interface{}, params ...interface{}) (res interface{}) {
	var infoEstudiante map[string]interface{}
	tercero := data["Tercero"].(string)
	if err := request.GetJsonWSO2(beego.AppConfig.String("Wso2Service")+"bienestarProxy/info_basica/"+tercero, &infoEstudiante); err == nil {
		data["Estudiante"] = infoEstudiante["datosCollection"].(map[string]interface{})["datos"].([]interface{})[0]
	} else {
		beego.Error("Error", err.Error())
	}
	return data
}

// GetLegalizacionInformation ...
// @Title GetLegalizacionInformation
// @Description get legalization information by avance id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Legalizacion_avance
// @Failure 403 :id is empty
// @router /GetLegalizacionInformation/:idAvance [get]
func (c *LegalizacionAvanceController) GetLegalizacionInformation() {
	idAvcStr := c.Ctx.Input.Param(":idAvance")
	defer c.ServeJSON()
	var avanceLegalizacion []map[string]interface{}
	var valorLegalizado float64

	respuesta := make(map[string]interface{})
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/avance_legalizacion//?query=Avance.Id:"+idAvcStr, &avanceLegalizacion); err == nil {
		if avanceLegalizacion != nil {
			respuesta["avanceLegalizacion"] = avanceLegalizacion
			idAvceLeg := strconv.FormatFloat(avanceLegalizacion[0]["Id"].(float64), 'f', -1, 64)
			if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/avance_legalizacion_tipo/GetLegalizationValue/"+idAvceLeg, &valorLegalizado); err == nil {
				respuesta["Total"] = valorLegalizado
			}
			c.Ctx.Output.SetStatus(201)
		}
	} else {
		res := models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
		respuesta = structs.Map(res)
	}
	c.Data["json"] = respuesta
}

// GetAllLegalizacionAvance ...
// @Title GetAllLegalizacionAvance
// @Description get Legalizacion_avance
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Legalizacion_avance
// @Failure 403
// @router /GetAllLegalizacionAvance [get]
func (c *LegalizacionAvanceController) GetAllLegalizacionAvance() {
	defer c.ServeJSON()
	var legalizaciones []interface{}
	var limit int64 = 10
	var offset int64
	var query string
	var regCuantity map[string]interface{}
	// limit: 10 (default is 10)
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
	respuesta := make(map[string]interface{})
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/avance_legalizacion/?limit="+strconv.FormatInt(limit, 10)+"&offset="+strconv.FormatInt(offset, 10)+"&query="+query, &legalizaciones); err == nil {
		if legalizaciones != nil {
			respuesta["Legalizaciones"] = optimize.ProccDigest(legalizaciones, getValuesLegalizacion, nil, 3)
			if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/avance_legalizacion/GetLegalizacionRecordsNumber/?query="+query, &regCuantity); err == nil {
				if strings.Compare(regCuantity["Type"].(string), "success") == 0 {
					respuesta["RegCuantity"] = regCuantity["Body"]
					c.Ctx.Output.SetStatus(201)
				}
			}
			c.Data["json"] = respuesta
		}
	} else {
		beego.Error("Error ", err)
		c.Data["json"] = models.Alert{Type: "error", Code: "E_0458", Body: err}
	}
}

func getValuesLegalizacion(rpintfc interface{}, params ...interface{}) (res interface{}) {
	var resEstado []map[string]interface{}
	var resValLegalizacion float64
	legalID := strconv.FormatFloat(rpintfc.(map[string]interface{})["Id"].(float64), 'f', -1, 64)
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/estado_legalizacion_avance_legalizacion/?query=Activo:true"+",AvanceLegalizacion.Id:"+legalID, &resEstado); err == nil {
		if resEstado != nil {
			rpintfc.(map[string]interface{})["Estado"] = resEstado[0]["Estado"]
		}
	} else {
		beego.Error("Error", err.Error())
	}
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/avance_legalizacion_tipo/GetLegalizationValue/"+legalID, &resValLegalizacion); err == nil {
		rpintfc.(map[string]interface{})["Valor"] = resValLegalizacion
	} else {
		beego.Error("Error", err.Error())
	}
	return rpintfc
}

// GetLegalizacionAccountantInformation ...
// @Title GetLegalizacionAccountantInformation
// @Description get accountant information to a legalization
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Legalizacion_avance
// @Failure 403 :id is empty
// @router /GetLegalizacionAccountantInformation/:idAvcLegalizacion [get]
func (c *LegalizacionAvanceController) GetLegalizacionAccountantInformation() {
	idAvceLeg := c.Ctx.Input.Param(":idAvcLegalizacion")
	defer c.ServeJSON()
	var avanceLegalizacionTipo []interface{}
	var conceptos []interface{}

	respuesta := make(map[string]interface{})

	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/avance_legalizacion_tipo/?query=AvanceLegalizacion.Id:"+idAvceLeg+",EstadoAvanceLegalizacionTipo.Id:1"+"&limit=-1", &avanceLegalizacionTipo); err == nil {
		respuesta["InformacionContable"] = optimize.ProccDigest(avanceLegalizacionTipo, getAccountantInfoLeg, nil, 3)
	} else {
		res := models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
		respuesta = structs.Map(res)
	}
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/concepto_avance_legalizacion_tipo/GetConceptoAvanceLegalizacionId/"+idAvceLeg, &conceptos); err == nil {
		respuesta["Conceptos"] = conceptos
	} else {
		res := models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
		respuesta = structs.Map(res)
	}

	c.Data["json"] = respuesta
}

func getAccountantInfoLeg(rpintfc interface{}, params ...interface{}) (res interface{}) {
	var conceptoLegalizacionAvance []map[string]interface{}
	var resMovimientoContable []map[string]interface{}
	var rpintfcCp map[string]interface{}
	var infoLegalizacion interface{}
	idLegTipo := strconv.FormatFloat(rpintfc.(map[string]interface{})["Id"].(float64), 'f', -1, 64)
	idTipoDocAfectante := strconv.FormatFloat(rpintfc.(map[string]interface{})["TipoDocumentoAfectante"].(map[string]interface{})["Id"].(float64), 'f', -1, 64)
	infoLegalizacion = formatoLegalizacion(rpintfc, nil)
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/concepto_avance_legalizacion_tipo/?query=AvanceLegalizacion.Id:"+idLegTipo, &conceptoLegalizacionAvance); err == nil {
		conceptoID := strconv.FormatFloat(conceptoLegalizacionAvance[0]["Concepto"].(map[string]interface{})["Id"].(float64), 'f', -1, 64)
		rpintfcCp = conceptoLegalizacionAvance[0]["Concepto"].(map[string]interface{})
		if infoLegalizacion.(map[string]interface{})["Estudiante"] != nil {
			rpintfcCp["Tercero"] = infoLegalizacion.(map[string]interface{})["Estudiante"]
		} else {
			rpintfcCp["Tercero"] = infoLegalizacion.(map[string]interface{})["InformacionProveedor"]
		}
		if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/movimiento_contable/?query=Concepto.Id:"+conceptoID+",CodigoDocumentoAfectante:"+idLegTipo+",TipoDocumentoAfectante.Id:"+idTipoDocAfectante+"&fields=Id,Credito,Debito,CuentaContable,Concepto", &resMovimientoContable); err == nil {
			if resMovimientoContable != nil {
				rpintfcCp["MovimientoContable"] = resMovimientoContable
			}
		} else {
			beego.Error("Error", err.Error())
		}
	} else {
		beego.Error("Error", err.Error())
	}
	return rpintfcCp
}
