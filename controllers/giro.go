package controllers

import (
	"encoding/json"
	"strconv"
	"strings"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/manucorporat/try"
	"github.com/udistrital/api_mid_financiera/models"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/optimize"
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
		//beego.Error("valor giro ", giro)
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
// @Description get RP by vigencia
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
				// done := make(chan interface{})
				// defer close(done)
				// resch := optimize.GenChanInterface(rpresupuestal...)
				// chrpresupuestal := optimize.Digest(done, FormatoListaRP, resch, nil)
				// for rp := range chrpresupuestal {
				// 	if rp != nil {
				// 		respuesta = append(respuesta, rp.(map[string]interface{}))
				// 	}

				// }
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
// @Description get RP by vigencia
// @Param	Id	query	string	false	"vigencia de la lista"
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Param	rangoinicio	query	string	false	"rango inicial del periodo a consultar"
// @Param	rangofin	query	string	false	"rango final del periodo a consultar"
// @Param	query	query	string	false	"query de filtrado "
// @Success 200 {object} models.Giro
// @Failure 403
// @router GetGirosById/:Id [get]
func (c *GiroController) GetGirosById() {
	vigenciaStr := c.Ctx.Input.Param(":Id")
	vigencia, err1 := strconv.Atoi(vigenciaStr)
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
		if err := request.GetJson(urlcrud + "/giro?limit="+strconv.FormatInt(limit, 10)+"&offset="+strconv.FormatInt(offset, 10)+"&query=Id:"+strconv.Itoa(vigencia)+query, &giro); err == nil {
			if giro != nil {
				done := make(chan interface{})
				defer close(done)
				resch := optimize.GenChanInterface(giro[0].(map[string]interface{})["GiroDetalle"].([]interface{})...)
				chrgiroDetalle := optimize.Digest(done, FormatoGiro, resch, nil)
				for gd := range chrgiroDetalle {
					if gd != nil {
						respuesta = append(respuesta, gd.(map[string]interface{}))
					}

				}
				giro[0].(map[string]interface{})["GiroDetalle"] = respuesta
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

//funcion para recopilar datos detallados del giro
func FormatoGiro(girointfc interface{}, params ...interface{}) (res interface{}) {
	giroDetalle := girointfc.(map[string]interface{})
	try.This(func() {
		idCuentaEspecial := giroDetalle["CuentaEspecial"].(map[string]interface{})["Id"].(float64)
		// beego.Info("idCuentaEspecial:",idCuentaEspecial)
		//solicituddisp, err := DetalleSolicitudDisponibilidadById(strconv.Itoa(idSolicitudDisponibilidad))

		if idCuentaEspecial == 0 {
			giroDetalle["InfoCuentaTercero"] = idCuentaEspecial
		}
	}).Catch(func(e try.E) {
		// Print crash
		fmt.Println("expc ",e)
	})
	return giroDetalle
}