package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/manucorporat/try"
	"github.com/udistrital/api_mid_financiera/models"
	"github.com/udistrital/utils_oas/optimize"
	"github.com/udistrital/utils_oas/request"
)

// GiroController operations for Giro
type GiroController struct {
	beego.Controller
}

// URLMapping ...
func (c *GiroController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Giro
// @Param	body		body 	models.Giro	true		"body for Giro content"
// @Success 201 {object} models.Giro
// @Failure 403 body is empty
// @router / [post]
func (c *GiroController) Post() {

}

// GetOne ...
// @Title GetOne
// @Description get Giro by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Giro
// @Failure 403 :id is empty
// @router /:id [get]
func (c *GiroController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get Giro
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Giro
// @Failure 403
// @router / [get]
func (c *GiroController) GetAll() {

}

// Put ...
// @Title Put
// @Description update the Giro
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Giro	true		"body for Giro content"
// @Success 200 {object} models.Giro
// @Failure 403 :id is not int
// @router /:id [put]
func (c *GiroController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Giro
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *GiroController) Delete() {

}

// CreateGiro ...
// @Title CreateGiro
// @Description Guarda objeto de Giro
// @Param	query	query	string	false	"Objecto del Giro"
// @Success 201 {object} models.Alert
// @Failure 403 body is empty
// @router /CreateGiro [post]
func (c *GiroController) CreateGiro() {
	defer c.ServeJSON()
	var giro map[string]interface{}
	var resProveedor []map[string]interface{}
	var idCuentasEspeciales []map[string]interface{}
	var response map[string]interface{}
	var responseDescuentos map[string]interface{}
	giroDescuentos := make(map[string]interface{})
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &giro); err == nil {
		urlcrud := "http://" + beego.AppConfig.String("Urlcrud") + ":" + beego.AppConfig.String("Portcrud") + "/" + beego.AppConfig.String("Nscrud")
		if err = request.SendJson(urlcrud+"/giro/RegistrarGiro", "POST", &response, giro); err == nil {
			if strings.Compare(response["Type"].(string), "success") == 0 {
				for _, element := range response["OrdenesPago"].([]interface{}) {
					if err = request.GetJson(urlcrud+"/giro/GetCuentasEspeciales?idordenpago="+strconv.FormatFloat(element.(map[string]interface{})["Id"].(float64), 'f', -1, 64), &idCuentasEspeciales); err == nil {
						for _, cuenta := range idCuentasEspeciales {
							if err = request.GetJson("http://"+beego.AppConfig.String("AdministrativaAmazonService")+"informacion_proveedor/?query=Id:"+cuenta["informacion_persona_juridica"].(string)+"&limit=1", &resProveedor); err == nil {
								giroDescuentos["idCuenta"] = cuenta["cuenta_especial"].(string)
								giroDescuentos["idOrdenPago"] = strconv.FormatFloat(element.(map[string]interface{})["Id"].(float64), 'f', -1, 64)
								giroDescuentos["idGiro"] = strconv.FormatFloat(response["IdGiro"].(float64), 'f', -1, 64)
								giroDescuentos["resProveedor"] = resProveedor
								err = request.SendJson(urlcrud+"/giro/RegistrarGiroDescuentos", "POST", &responseDescuentos, giroDescuentos)
							}
						}
					}
				}
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

// ListarGiros ...
// @Title ListarGiros
// @Description Listar todos los giros por vigencia
// @Param	vigencia	query	string	false	"vigencia de la lista"
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Param	rangoinicio	query	string	false	"rango inicial del periodo a consultar"
// @Param	rangofin	query	string	false	"rango final del periodo a consultar"
// @Param	query	query	string	false	"query de filtrado "
// @Success 200 {object} models.Giro
// @Failure 403
// @router ListarGiros/:vigencia [get]
func (c *GiroController) ListarGiros() {
	vigenciaStr := c.Ctx.Input.Param(":vigencia")
	vigencia, err1 := strconv.Atoi(vigenciaStr)
	var giro []interface{}
	//var respuesta []map[string]interface{}
	var limit int64 = 10
	var offset int64
	var startrange string
	var endrange string
	var query string
	var querybase string
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	if r := c.GetString("rangoinicio"); r != "" {
		startrange = r

	}

	if r := c.GetString("rangofin"); r != "" {
		endrange = r

	}
	if r := c.GetString("query"); r != "" {
		querybase = r

	}
	if startrange != "" && endrange != "" {
		query = querybase + ",FechaRegistro__gte:" + startrange + ",FechaRegistro__lte:" + endrange

	} else if querybase != "" {
		query = "," + querybase
	}
	if err1 == nil {
		if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/giro?limit="+strconv.FormatInt(limit, 10)+"&offset="+strconv.FormatInt(offset, 10)+"&query=Vigencia:"+strconv.Itoa(vigencia)+query, &giro); err == nil {
			if giro != nil {
				c.Data["json"] = giro
			} else {
				c.Data["json"] = models.Alert{Code: "E_0458", Body: nil, Type: "error"}
			}
		} else {
			c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
		}
	} else {
		c.Data["json"] = models.Alert{Code: "E_0458", Body: "Not enough parameter", Type: "error"}
	}

	c.ServeJSON()
}

// GetGirosById ...
// @Title GetGirosById
// @Description get Giro by Id
// @Param	Id	query	string	false	"Id del giro"
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Param	rangoinicio	query	string	false	"rango inicial del periodo a consultar"
// @Param	rangofin	query	string	false	"rango final del periodo a consultar"
// @Param	query	query	string	false	"query de filtrado "
// @Success 200 {object} models.Giro
// @Failure 403
// @router GetGirosById/:Id [get]
func (c *GiroController) GetGirosById() {
	giroIdStr := c.Ctx.Input.Param(":Id")
	giroId, err1 := strconv.Atoi(giroIdStr)
	var giro []interface{}
	var respuesta []map[string]interface{}
	var limit int64 = 10
	var offset int64
	var startrange string
	var endrange string
	var query string
	var querybase string
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	if r := c.GetString("rangoinicio"); r != "" {
		startrange = r

	}

	if r := c.GetString("rangofin"); r != "" {
		endrange = r

	}
	if r := c.GetString("query"); r != "" {
		querybase = r

	}
	if startrange != "" && endrange != "" {
		query = querybase + ",FechaRegistro__gte:" + startrange + ",FechaRegistro__lte:" + endrange

	} else if querybase != "" {
		query = "," + querybase
	}
	if err1 == nil {
		urlcrud := "http://" + beego.AppConfig.String("Urlcrud") + ":" + beego.AppConfig.String("Portcrud") + "/" + beego.AppConfig.String("Nscrud")
		if err := request.GetJson(urlcrud+"/giro_detalle?limit="+strconv.FormatInt(limit, 10)+"&offset="+strconv.FormatInt(offset, 10)+"&query=Giro:"+strconv.Itoa(giroId)+query, &giro); err == nil {
			if giro != nil {
				done := make(chan interface{})
				defer close(done)
				resch := optimize.GenChanInterface(giro...)
				chrgiroDetalle := optimize.Digest(done, FormatoGiro, resch, nil)
				for gd := range chrgiroDetalle {
					if gd != nil {
						respuesta = append(respuesta, gd.(map[string]interface{}))
					}

				}
				c.Data["json"] = respuesta
			} else {
				c.Data["json"] = models.Alert{Code: "E_0458", Body: nil, Type: "error"}
			}
		} else {
			c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
		}
	} else {
		c.Data["json"] = models.Alert{Code: "E_0458", Body: "Not enough parameter", Type: "error"}
	}

	c.ServeJSON()
}

// GetSumGiro ...
// @Title GetSumGiro
// @Description get sum values by Id for models.GiroGiro
// @Param	Id	query	string	false	"vigencia de la lista"
// @Success 200 {object} models.Giro
// @Failure 403
// @router GetSumGiro/:Id [get]
func (c *GiroController) GetSumGiro() {
	giroIdStr := c.Ctx.Input.Param(":Id")
	giroId, err1 := strconv.Atoi(giroIdStr)
	var totalGiro []interface{}
	var infoGiro interface{}
	// var respuesta []map[string]interface{}

	if err1 == nil {
		urlcrud := "http://" + beego.AppConfig.String("Urlcrud") + ":" + beego.AppConfig.String("Portcrud") + "/" + beego.AppConfig.String("Nscrud")
		if err := request.GetJson(urlcrud+"/giro/GetSumGiro/?IdGiro="+strconv.Itoa(giroId), &totalGiro); err == nil {
			if totalGiro != nil {
				if err := request.GetJson(urlcrud+"/giro/"+strconv.Itoa(giroId), &infoGiro); err == nil {
					totalGiro[0].(map[string]interface{})["total_op"] = strconv.FormatFloat(infoGiro.(map[string]interface{})["ValorTotal"].(float64), 'f', -1, 64)
					if totalGiro[0].(map[string]interface{})["total_cuentas_especiales"] == nil {
						totalGiro[0].(map[string]interface{})["total_cuentas_especiales"] = 0
					}
				} else {
					c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
				}

				c.Data["json"] = totalGiro
			} else {
				c.Data["json"] = models.Alert{Code: "E_0458", Body: nil, Type: "error"}
			}
		} else {
			c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
		}
	} else {
		c.Data["json"] = models.Alert{Code: "E_0458", Body: "Not enough parameter", Type: "error"}
	}

	c.ServeJSON()
}

//funcion para recopilar valores de la cuenta especial
func getValueGiroCuentaEspecial(element interface{}) (res interface{}) {

	if element.(map[string]interface{})["CuentaEspecial"].(map[string]interface{})["TipoCuentaEspecial"].(map[string]interface{})["Nombre"] == "Endoso" {
		urlcrud := "http://" + beego.AppConfig.String("Urlcrud") + ":" + beego.AppConfig.String("Portcrud") + "/" + beego.AppConfig.String("Nscrud")
		if err := request.GetJson(urlcrud+"/giro/GetValueEndoso/?IdOrdenPago="+strconv.FormatFloat(element.(map[string]interface{})["OrdenPago"].(map[string]interface{})["Id"].(float64), 'f', -1, 64)+"&IdCuentaEspecial="+strconv.FormatFloat(element.(map[string]interface{})["CuentaEspecial"].(map[string]interface{})["Id"].(float64), 'f', -1, 64), &res); err == nil {
			res = res.([]interface{})[0].(map[string]interface{})["valor_endoso"]
		}
	} else {
		res = element.(map[string]interface{})["ValorBase"].(float64) * element.(map[string]interface{})["CuentaEspecial"].(map[string]interface{})["Porcentaje"].(float64)
	}
	return
}

//funcion para recopilar datos detallados del giro
func FormatoGiro(girointfc interface{}, params ...interface{}) (res interface{}) {
	giroDetalle := girointfc.(map[string]interface{})
	var resProveedor []map[string]interface{}
	try.This(func() {
		idCuentaEspecial := giroDetalle["CuentaEspecial"].(map[string]interface{})["Id"].(float64)
		urladministrativa := "http://" + beego.AppConfig.String("AdministrativaAmazonService") + "informacion_proveedor/?query=Id:"
		if idCuentaEspecial == 0 {
			if err := request.GetJson(urladministrativa+strconv.FormatFloat(giroDetalle["OrdenPago"].(map[string]interface{})["OrdenPagoRegistroPresupuestal"].([]interface{})[0].(map[string]interface{})["RegistroPresupuestal"].(map[string]interface{})["Beneficiario"].(float64), 'f', -1, 64), &resProveedor); err == nil {
				giroDetalle["InfoProveedor"] = resProveedor
				giroDetalle["ValorBasePago"] = giroDetalle["OrdenPago"].(map[string]interface{})["ValorBase"].(float64)
				giroDetalle["TipoMov"] = giroDetalle["OrdenPago"].(map[string]interface{})["SubTipoOrdenPago"].(map[string]interface{})["TipoOrdenPago"].(map[string]interface{})["CodigoAbreviacion"].(string)
			}
		} else {
			if err := request.GetJson(urladministrativa+strconv.FormatFloat(giroDetalle["CuentaEspecial"].(map[string]interface{})["InformacionPersonaJuridica"].(float64), 'f', -1, 64), &resProveedor); err == nil {
				giroDetalle["InfoProveedor"] = resProveedor
				for _, element := range giroDetalle["OrdenPago"].(map[string]interface{})["OrdenPagoCuentaEspecial"].([]interface{}) {
					if giroDetalle["CuentaEspecial"].(map[string]interface{})["Id"].(float64) == element.(map[string]interface{})["CuentaEspecial"].(map[string]interface{})["Id"].(float64) {
						giroDetalle["ValorBasePago"] = getValueGiroCuentaEspecial(element)
						giroDetalle["TipoMov"] = element.(map[string]interface{})["CuentaEspecial"].(map[string]interface{})["TipoCuentaEspecial"].(map[string]interface{})["Nombre"].(string)
					}
				}
			}
		}
	}).Catch(func(e try.E) {
		// Print crash
		fmt.Println("expc ", e)
	})
	return giroDetalle
}
