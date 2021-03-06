package controllers

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/api_mid_financiera/models"
	"github.com/udistrital/utils_oas/optimize"
	"github.com/udistrital/utils_oas/request"
)

// GestionChequesController operations for Gestion_cheques
type GestionChequesController struct {
	beego.Controller
}

// URLMapping ...
func (c *GestionChequesController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Gestion_cheques
// @Param	body		body 	models.Gestion_cheques	true		"body for Gestion_cheques content"
// @Success 201 {object} models.Gestion_cheques
// @Failure 403 body is empty
// @router / [post]
func (c *GestionChequesController) Post() {

}

// GetOne ...
// @Title GetOne
// @Description get Gestion_cheques by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Gestion_cheques
// @Failure 403 :id is empty
// @router /:id [get]
func (c *GestionChequesController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get Gestion_cheques
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Gestion_cheques
// @Failure 403
// @router / [get]
func (c *GestionChequesController) GetAll() {

}

// Put ...
// @Title Put
// @Description update the Gestion_cheques
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Gestion_cheques	true		"body for Gestion_cheques content"
// @Success 200 {object} models.Gestion_cheques
// @Failure 403 :id is not int
// @router /:id [put]
func (c *GestionChequesController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Gestion_cheques
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *GestionChequesController) Delete() {

}

// CreateChequera ...
// @Title CreateChequera
// @Description create homologate category for an organization
// @Param	body		body 	interface	true		"body for Homologacion_rubro content"
// @Success 201 {object} interface{}
// @Failure 403 body is empty
// @router /CreateChequera [post]
func (c *GestionChequesController) CreateChequera() {
	defer c.ServeJSON()
	var chequera interface{}
	var response map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &chequera); err == nil {
		if err = request.SendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/chequera/CreateChequeraEstado", "POST", &response, chequera); err == nil {
			if strings.Compare(response["Type"].(string), "success") == 0 {
				c.Data["json"] = models.Alert{Type: "success", Code: "S_543", Body: response["Body"]}
				c.Ctx.Output.SetStatus(201)
			} else {
				beego.Error("Error", response)
				c.Data["json"] = models.Alert{Type: response["Type"].(string), Code: response["Code"].(string), Body: response["Body"]}
			}
		} else {
			beego.Error("Error", err)
			c.Data["json"] = models.Alert{Type: "error", Code: "E_0458", Body: err}
		}
	} else {
		beego.Error("Error", err)
		c.Data["json"] = models.Alert{Type: "error", Code: "E_0458", Body: err}
	}
}

// CreateCheque ...
// @Title CreateCheque
// @Description create cheque asociated to first state
// @Param	body		body 	interface	true		"body for Homologacion_rubro content"
// @Success 201 {object} interface{}
// @Failure 403 body is empty
// @router /CreateCheque [post]
func (c *GestionChequesController) CreateCheque() {
	defer c.ServeJSON()
	var cheque interface{}
	var response map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &cheque); err == nil {
		beego.Error("valor cheque ", cheque)
		if err = request.SendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/cheque/CreateChequeEstado", "POST", &response, cheque); err == nil {
			if strings.Compare(response["Type"].(string), "success") == 0 {
				c.Data["json"] = models.Alert{Type: "success", Code: "S_543", Body: response["Body"]}
				c.Ctx.Output.SetStatus(201)
			} else {
				beego.Error("Error", response)
				c.Data["json"] = models.Alert{Type: response["Type"].(string), Code: response["Code"].(string), Body: response["Body"]}
			}
		} else {
			beego.Error("Error", err)
			c.Data["json"] = models.Alert{Type: "error", Code: "E_0458", Body: err}
		}
	} else {
		beego.Error("Error", err)
		c.Data["json"] = models.Alert{Type: "error", Code: "E_0458", Body: err}
	}
}

// GetAllCheque ...
// @Title GetAllCheque
// @Description get all cheque and its information
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Gestion_cheques
// @Failure 403
// @router /GetAllCheque/ [get]
func (c *GestionChequesController) GetAllCheque() {
	defer c.ServeJSON()
	var cheques []interface{}
	var limit int64 = 10
	var offset int64
	var query string
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
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/cheque/?limit="+strconv.FormatInt(limit, 10)+"&offset="+strconv.FormatInt(offset, 10)+"&query="+query, &cheques); err == nil {
		if cheques != nil {
			respuesta := optimize.ProccDigest(cheques, getValuesCheques, nil, 3)
			c.Data["json"] = respuesta
		}
	} else {
		beego.Error("Error ", err)
		c.Data["json"] = models.Alert{Type: "error", Code: "E_0458", Body: err}
	}
}

func getValuesCheques(rpintfc interface{}, params ...interface{}) (res interface{}) {
	var resEstado []map[string]interface{}
	var resBeneficiario []map[string]interface{}
	chequeID := strconv.FormatFloat(rpintfc.(map[string]interface{})["Id"].(float64), 'f', -1, 64)
	chequera := rpintfc.(map[string]interface{})["Chequera"]
	chequeraID := strconv.FormatFloat(chequera.(map[string]interface{})["Id"].(float64), 'f', -1, 64)
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/chequera/"+chequeraID, &chequera); err == nil {
		if chequera != nil {
			rpintfc.(map[string]interface{})["Chequera"] = getValuesChequera(chequera)
		}
	} else {
		beego.Error("Error ", err)
	}
	ordenPago := rpintfc.(map[string]interface{})["OrdenPago"]
	ordenPagoID := strconv.FormatFloat(ordenPago.(map[string]interface{})["Id"].(float64), 'f', -1, 64)
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/orden_pago/"+ordenPagoID, &ordenPago); err == nil {
		if ordenPago != nil {
			rpintfc.(map[string]interface{})["OrdenPago"] = ordenPago
		}
	} else {
		beego.Error("Error ", err)
	}

	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/cheque_estado_cheque/?query=Activo:true"+",cheque.Id:"+chequeID, &resEstado); err == nil {
		if resEstado[0] != nil {
			rpintfc.(map[string]interface{})["Estado"] = resEstado[0]["Estado"]
		}
	} else {
		beego.Error("Error", err.Error())
	}
	beneficiario := strconv.FormatFloat(rpintfc.(map[string]interface{})["Beneficiario"].(float64), 'f', -1, 64)
	if err := request.GetJson("http://"+beego.AppConfig.String("AdministrativaAmazonService")+"/informacion_proveedor/?query=NumDocumento:"+beneficiario+"&limit=1", &resBeneficiario); err == nil {
		if resBeneficiario != nil {
			idBen,err := strconv.Atoi(beneficiario)
			if (err!=nil){
				beego.Error("Error", err.Error())
			}
			rpintfc.(map[string]interface{})["Beneficiario"] = models.BeneficiarioCheque{Id:idBen,Nombre:resBeneficiario[0]["NomProveedor"].(string)}
		}
	} else {
		beego.Error("Error", err.Error())
	}
	return rpintfc
}

// GetAllChequera ...
// @Title GetAllChequera
// @Description get all chequeras
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Gestion_cheques
// @Failure 403
// @router /GetAllChequera/ [get]
func (c *GestionChequesController) GetAllChequera() {
	defer c.ServeJSON()
	var chequeras []interface{}
	var limit int64 = 10
	var offset int64
	var query string
	var buscarDisponibles bool
	var complementation string
	var params []interface{}
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
	// buscarDisponibles:false (default is false)
	if v, err := c.GetBool("bDisponibles"); err == nil {
		buscarDisponibles = v
	}

	if buscarDisponibles {
		querybase := "Activo:true"
		if query != "" {
			query = query + "," + querybase
		} else {
			query = querybase
		}
		complementation = "/chequera_estado_chequera/?limit=" + strconv.FormatInt(limit, 10) + "&offset=" + strconv.FormatInt(offset, 10) + "&query=" + query
	} else {
		complementation = "/chequera/?limit=" + strconv.FormatInt(limit, 10) + "&offset=" + strconv.FormatInt(offset, 10) + "&query=" + query
	}
	params = append(params, buscarDisponibles)
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+complementation, &chequeras); err == nil {
		if chequeras != nil {
			respuesta := optimize.ProccDigest(chequeras, chequeraInfoDistpacher, params, 3)
			c.Data["json"] = respuesta
		}
	} else {
		beego.Error("Error ", err)
		c.Data["json"] = models.Alert{Type: "error", Code: "E_0458", Body: err}
	}
}

func chequeraInfoDistpacher(rpintfc interface{}, params ...interface{}) (res interface{}) {
	if len(params) > 0 {
		buscarDisponibles := params[0].(bool)
		if buscarDisponibles {
			chequera := rpintfc.(map[string]interface{})["Chequera"]
			chequeraID := strconv.FormatFloat(chequera.(map[string]interface{})["Id"].(float64), 'f', -1, 64)
			if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/chequera/"+chequeraID, &chequera); err == nil {
				if chequera != nil {
					rpintfc.(map[string]interface{})["Chequera"] = getValuesChequera(chequera)
					res = rpintfc
				}
			} else {
				beego.Error("Error ", err)
			}
		} else {
			res = getValuesChequera(rpintfc)
		}
	}

	return
}
func getValuesChequera(rpintfc interface{}) (res interface{}) {
	var resEstado []map[string]interface{}
	var unidadEjecutoraResp []map[string]interface{}
	var responsableResp []map[string]interface{}
	var resSucursal []map[string]interface{}

	chequeraID := strconv.FormatFloat(rpintfc.(map[string]interface{})["Id"].(float64), 'f', -1, 64)
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/chequera_estado_chequera/?query=Activo:true"+",chequera.Id:"+chequeraID, &resEstado); err == nil {
		if resEstado[0] != nil {
			rpintfc.(map[string]interface{})["Estado"] = resEstado[0]["Estado"]
		}
	} else {
		beego.Error("Error", err.Error())
	}
	UnidadEjecutora := strconv.FormatFloat(rpintfc.(map[string]interface{})["UnidadEjecutora"].(float64), 'f', -1, 64)
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/unidad_ejecutora?limit=-1&query=Id:"+UnidadEjecutora, &unidadEjecutoraResp); err == nil {
		if unidadEjecutoraResp != nil {
			rpintfc.(map[string]interface{})["UnidadEjecutora"] = unidadEjecutoraResp[0]
		}
	} else {
		beego.Error("Error", err.Error())
	}
	sucursalID := strconv.FormatFloat(rpintfc.(map[string]interface{})["CuentaBancaria"].(map[string]interface{})["Sucursal"].(float64), 'f', -1, 64)
	if err := request.GetJson(beego.AppConfig.String("coreOrganizacionService")+"organizacion/?query=Id:"+sucursalID, &resSucursal); err == nil {
		if resSucursal[0] != nil {
			rpintfc.(map[string]interface{})["Sucursal"] = resSucursal[0]
		}
	} else {
		beego.Error("Error", err.Error())
	}
	resBanco, err := GetBancoSucursal(sucursalID)
	if err == nil {
		rpintfc.(map[string]interface{})["Banco"] = resBanco.([]interface{})[0].(map[string]interface{})
	} else {
		beego.Error("Error", err.Error())
	}
	Responsable := strconv.FormatFloat(rpintfc.(map[string]interface{})["Responsable"].(float64), 'f', -1, 64)
	if err := request.GetJson("http://"+beego.AppConfig.String("AdministrativaAmazonService")+"/supervisor_contrato?limit=-1&query=Documento:"+Responsable, &responsableResp); err == nil {
		if responsableResp != nil {
			rpintfc.(map[string]interface{})["Responsable"] = responsableResp[0]
		}
	} else {
		beego.Error("Error", err.Error())
	}
	return rpintfc
}
