package controllers

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/manucorporat/try"
	"github.com/udistrital/api_financiera/models"
	"github.com/udistrital/utils_oas/optimize"
	"github.com/udistrital/utils_oas/request"
)

// TipoTransaccionController operations for TipoTransaccion
type TipoTransaccionController struct {
	beego.Controller
}

// URLMapping ...
func (c *TipoTransaccionController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create TipoTransaccion
// @Param	body		body 	models.TipoTransaccion	true		"body for TipoTransaccion content"
// @Success 201 {object} models.TipoTransaccion
// @Failure 403 body is empty
// @router / [post]
func (c *TipoTransaccionController) Post() {
	defer c.ServeJSON()
	var v map[string]interface{}
	var version map[string]interface{}
	var detalleTransaccion map[string]interface{}
	urlCrud := "http://" + beego.AppConfig.String("Urlcrud") + ":" + beego.AppConfig.String("Portcrud") + "/" + beego.AppConfig.String("Nscrud")
	responseRoute := make(map[string]interface{})
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		version = v["version"].(map[string]interface{})
		detalleTransaccion = v["detalleTransaccion"].(map[string]interface{})
		request.Commit(func() {
			err, respV := SaveForTipoTr(urlCrud+"/version_tipo_transaccion", version)
			if err != nil {
				beego.Error(err)
				panic(err)
			}
			responseRoute["version_tipo_transaccion"] = respV
			err, responseTTV := SaveForTipoTr(urlCrud+"/tipo_transaccion_version/CreateTipoVersion", responseRoute["version_tipo_transaccion"].(map[string]interface{})["Body"])
			if err != nil {
				beego.Error(err)
				panic(err)
			}
			responseRoute["tipo_transaccion_version"] = responseTTV
			detalleTransaccion["TipoTransaccionVersion"] = responseTTV.(map[string]interface{})["Body"]
			err, responseDT := SaveForTipoTr(urlCrud+"/detalle_tipo_transaccion_version", detalleTransaccion)
			if err != nil {
				beego.Error(err)
				panic(err)
			}
			responseRoute["detalle_tipo_transaccion_version"] = responseDT
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = models.Alert{Type: "success", Code: "S_543", Body: responseRoute}
		}).Rollback(func(response interface{}, error interface{}) {
			beego.Error("Error Rollback ", error)
			err, _ := removeElement("detalle_tipo_transaccion_version", response)
			if err != nil {
				beego.Error(err)
				panic(err)
			}
			err, _ = removeElement("tipo_transaccion_version", response)
			if err != nil {
				beego.Error(err)
				panic(err)
			}
			err, _ = removeElement("version_tipo_transaccion", response)
			if err != nil {
				beego.Error(err)
				panic(err)
			}
			c.Data["json"] = models.Alert{Type: "error", Code: "E_0458", Body: response}
		}, responseRoute)
	} else {
		beego.Error(err.Error())
	}
}

// GetOne ...
// @Title GetOne
// @Description get TipoTransaccion by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.TipoTransaccion
// @Failure 403 :id is empty
// @router /:id [get]
func (c *TipoTransaccionController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get TipoTransaccionTipoTransaccionController
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.TipoTransaccion
// @Failure 403
// @router / [get]
func (c *TipoTransaccionController) GetAll() {

}

// Put ...
// @Title Put
// @Description update the TipoTransaccion
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.TipoTransaccion	true		"body for TipoTransaccion content"
// @Success 200 {object} models.TipoTransaccion
// @Failure 403 :id is not int
// @router /:id [put]
func (c *TipoTransaccionController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the TipoTransaccion
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *TipoTransaccionController) Delete() {

}

func SaveForTipoTr(object ...interface{}) (err error, response interface{}) {
	route := object[0]
	sendData := object[1]
	err = request.SendJson(route.(string), "POST", &response, sendData)
	if strings.Compare(response.(map[string]interface{})["Type"].(string), "success") != 0 {
		beego.Error(response)
		err = errors.New(response.(map[string]interface{})["Code"].(string))
	}
	return
}

func removeElement(idMap string, object interface{}) (err error, response interface{}) {
	urlCrud := "http://" + beego.AppConfig.String("Urlcrud") + ":" + beego.AppConfig.String("Portcrud") + "/" + beego.AppConfig.String("Nscrud")
	value := object.(map[string]interface{})[idMap]
	if value != nil {
		body := value.(map[string]interface{})["Body"].(map[string]interface{})
		if body["Id"] != nil {
			id := strconv.Itoa(int(body["Id"].(float64)))
			err = request.SendJson(urlCrud+"/"+idMap+"/"+id, "DELETE", &response, nil)
		}
	}
	return
}

func getElement(route string) (err error, response interface{}) {
	request.GetJson(route, &response)
	if strings.Compare(response.(map[string]interface{})["Type"].(string), "success") != 0 {
		beego.Error(response)
		err = errors.New(response.(map[string]interface{})["Code"].(string))
	}
	return
}

// GetTipoTransaccionByVersion ...
// @Title GetTipoTransaccionByVersion
// @Description get TipoTransaccionTipoTransaccionController given initial version
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.TipoTransaccion
// @Failure 403
// @router /GetTipoTransaccionByVersion/ [get]
func (c *TipoTransaccionController) GetTipoTransaccionByVersion() {
	var versionesTipoT []interface{}
	var limit int64 = 10
	var offset int64
	var query string
	var sortby string
	var order string
	var regCuantity map[string]interface{}
	defer c.ServeJSON()
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
	if v := c.GetString("sortby"); v != "" {
		sortby = "&sortby=" + v
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = "&order=" + v
	}
	try.This(func() {
		respuesta := make(map[string]interface{})
		if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/version_tipo_transaccion/?limit="+strconv.FormatInt(limit, 10)+"&offset="+strconv.FormatInt(offset, 10)+"&query="+query+sortby+order, &versionesTipoT); err == nil {
			if versionesTipoT != nil {
				respuesta["TipoTransaccion"] = optimize.ProccDigest(versionesTipoT, getValuesTiposTr, nil, 3)
				if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/version_tipo_transaccion/GetVersionesTipoNumber/?query="+query, &regCuantity); err == nil {
					if strings.Compare(regCuantity["Type"].(string), "success") == 0 {
						respuesta["RegCuantity"] = regCuantity["Body"]
						c.Ctx.Output.SetStatus(201)
					}
				}
				c.Data["json"] = respuesta
			}
		} else {
			panic(err)
		}
	}).Catch(func(e try.E) {
		beego.Error("expc ", e)
		c.Data["json"] = models.Alert{Type: "error", Code: "E_0458", Body: e}
	})
}

func getValuesTiposTr(rpintfc interface{}, params ...interface{}) (res interface{}) {
	var tipoTransaccionVersion []interface{}
	var detalleTipoTr []interface{}
	idVersionStr := strconv.Itoa(int(rpintfc.(map[string]interface{})["Id"].(float64)))
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/tipo_transaccion_version/?&query=Version.Id:"+idVersionStr, &tipoTransaccionVersion); err == nil {
		idTipoTranVStr := strconv.Itoa(int(tipoTransaccionVersion[0].(map[string]interface{})["Id"].(float64)))
		if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_tipo_transaccion_version/?&query=TipoTransaccionVersion.Id:"+idTipoTranVStr, &detalleTipoTr); err == nil {
			rpintfc.(map[string]interface{})["DetalleTipoTransaccion"] = detalleTipoTr[0]
		} else {
			beego.Error("Error", err.Error())
		}
	} else {
		beego.Error("Error", err.Error())
	}
	return rpintfc
}

// GetTipoTransaccionByTipo ...
// @Title GetTipoTransaccionByTipo
// @Description get TipoTransaccionTipoTransaccionController given initial version
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.TipoTransaccion
// @Failure 403
// @router /GetTipoTransaccionByTipo/ [get]
func (c *TipoTransaccionController) GetTipoTransaccionByTipo() {
	var tipoTransaccion []interface{}
	var limit int64 = 10
	var offset int64
	var query string
	var sortby string
	var order string
	var regCuantity map[string]interface{}
	defer c.ServeJSON()
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
	if v := c.GetString("sortby"); v != "" {
		sortby = "&sortby=" + v
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = "&order=" + v
	}
	try.This(func() {
		respuesta := make(map[string]interface{})
		if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/tipo_transaccion_version/?limit="+strconv.FormatInt(limit, 10)+"&offset="+strconv.FormatInt(offset, 10)+"&query="+query+sortby+order, &tipoTransaccion); err == nil {
			if tipoTransaccion != nil {
				respuesta["TipoTransaccion"] = optimize.ProccDigest(tipoTransaccion, getValuesDetallesTiposTr, nil, 3)
				if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/tipo_transaccion_version/GetTipoTransaccionVersionesNumber/?query="+query, &regCuantity); err == nil {
					if strings.Compare(regCuantity["Type"].(string), "success") == 0 {
						respuesta["RegCuantity"] = regCuantity["Body"]
						c.Ctx.Output.SetStatus(201)
					}
				}
				c.Data["json"] = respuesta
			}
		} else {
			panic(err)
		}
	}).Catch(func(e try.E) {
		beego.Error("expc ", e)
		c.Data["json"] = models.Alert{Type: "error", Code: "E_0458", Body: e}
	})
}

func getValuesDetallesTiposTr(rpintfc interface{}, params ...interface{}) (res interface{}) {
	var detalleTipoTr []interface{}
	idTipoTrStr := strconv.Itoa(int(rpintfc.(map[string]interface{})["Id"].(float64)))
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_tipo_transaccion_version/?query=TipoTransaccionVersion.Id:"+idTipoTrStr, &detalleTipoTr); err == nil {
		if detalleTipoTr != nil {
			rpintfc.(map[string]interface{})["DetalleTipoTransaccion"] = detalleTipoTr[0]
		}
	} else {
		beego.Error("Error", err.Error())
	}
	return rpintfc
}

// NewTipoTransaccionVersion ...
// @Title CreateNewTipoTransaccionVersion
// @Description create NewTipoTransaccionVersion
// @Param	body		body 	models.TipoTransaccion	true		"body for TipoTransaccion content"
// @Success 201 {object} models.TipoTransaccion
// @Failure 403 body is empty
// @router /NewTipoTransaccionVersion/ [post]
func (c *TipoTransaccionController) NewTipoTransaccionVersion() {
	defer c.ServeJSON()
	var v map[string]interface{}
	var tipoTransaccionVersion map[string]interface{}
	var version map[string]interface{}
	var detalleTransaccion map[string]interface{}
	urlCrud := "http://" + beego.AppConfig.String("Urlcrud") + ":" + beego.AppConfig.String("Portcrud") + "/" + beego.AppConfig.String("Nscrud")
	responseRoute := make(map[string]interface{})
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		tipoTransaccionVersion = v["tipoTransaccionVersion"].(map[string]interface{})
		detalleTransaccion = v["detalleTransaccion"].(map[string]interface{})
		version = v["version"].(map[string]interface{})
		fechaInicioStr := version["FechaInicio"].(string)
		fechaFinStr := version["FechaFin"].(string)
		tipoStr := strconv.Itoa(int(tipoTransaccionVersion["TipoTransaccion"].(float64)))
		err, versionesFecha := getElement(urlCrud + "/version_tipo_transaccion/GetVersionInEspecifiedDate/?fechaInicio=" + fechaInicioStr + "&fechaFin=" + fechaFinStr + "&tipo=" + tipoStr)
		if err != nil || versionesFecha.(map[string]interface{})["Body"] != nil {
			beego.Error(err)
			c.Data["json"] = models.Alert{Type: "error", Code: "E_LT0001", Body: "version date out of range"}
			return
		}
		request.Commit(func() {
			idTipoStr := strconv.Itoa(int(tipoTransaccionVersion["Id"].(float64)))
			err, respVersionNumero := getElement(urlCrud + "/version_tipo_transaccion/GetVersionToType/" + idTipoStr)
			if err != nil {
				beego.Error(err)
				panic(err)
			}
			version["NumeroVersion"] = respVersionNumero.(map[string]interface{})["Body"]
			err, respV := SaveForTipoTr(urlCrud+"/version_tipo_transaccion", version)
			if err != nil {
				beego.Error(err)
				panic(err)
			}
			responseRoute["version_tipo_transaccion"] = respV
			tipoTransaccionCrear := make(map[string]interface{})
			tipoTransaccionCrear["TipoTransaccion"] = tipoTransaccionVersion["TipoTransaccion"]
			tipoTransaccionCrear["Version"] = respV.(map[string]interface{})["Body"]
			err, responseTTV := SaveForTipoTr(urlCrud+"/tipo_transaccion_version", tipoTransaccionCrear)
			if err != nil {
				beego.Error(err)
				panic(err)
			}
			responseRoute["tipo_transaccion_version"] = responseTTV
			detalleTransaccion["TipoTransaccionVersion"] = responseTTV.(map[string]interface{})["Body"]
			err, responseDT := SaveForTipoTr(urlCrud+"/detalle_tipo_transaccion_version", detalleTransaccion)
			if err != nil {
				beego.Error(err)
				panic(err)
			}
			responseRoute["detalle_tipo_transaccion_version"] = responseDT
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = models.Alert{Type: "success", Code: "S_543", Body: responseRoute}
		}).Rollback(func(response interface{}, error interface{}) {
			err, _ := removeElement("detalle_tipo_transaccion_version", response)
			if err != nil {
				beego.Error(err)
				panic(err)
			}
			err, _ = removeElement("tipo_transaccion_version", response)
			if err != nil {
				beego.Error(err)
				panic(err)
			}
			err, _ = removeElement("version_tipo_transaccion", response)
			if err != nil {
				beego.Error(err)
				panic(err)
			}
			c.Data["json"] = models.Alert{Type: "error", Code: "E_0458", Body: response}
		}, responseRoute)
	} else {
		beego.Error(err.Error())
	}
}

// UpdateTipoTransaccion ...
// @Title UpdateTipoTransaccion
// @Description update the TipoTransaccion
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.TipoTransaccion	true		"body for TipoTransaccion content"
// @Success 200 {object} models.TipoTransaccion
// @Failure 403 body us empty
// @router /UpdateTipoTransaccion/ [post]
func (c *TipoTransaccionController) UpdateTipoTransaccion() {
	var v interface{}
	var respuesta interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := request.SendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_tipo_transaccion_version/UpdateTipoTransaccion", "PUT", &respuesta, v); err == nil {
			c.Data["json"] = respuesta
		} else {
			beego.Error("error ", err)
			alert := models.Alert{Type: "error", Code: "E_0458", Body: err}
			c.Data["json"] = alert
		}
	} else {
		beego.Error("error ", err)
		alert := models.Alert{Type: "error", Code: "E_0458", Body: err}
		c.Data["json"] = alert
	}
	c.ServeJSON()
}

// GetTipoTransaccionByTipo ...
// @Title GetTipoTransaccionByTipo
// @Description get TipoTransaccionTipoTransaccionController given initial version
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.TipoTransaccion
// @Failure 403
// @router /GetTipoTransaccionByDefinitiveVersion/ [get]
func (c *TipoTransaccionController) GetTipoTransaccionByDefinitiveVersion() {
	var tipoTransaccion []interface{}
	defer c.ServeJSON()

	try.This(func() {
		beego.Error("http://" + beego.AppConfig.String("Urlcrud") + ":" + beego.AppConfig.String("Portcrud") + "/" + beego.AppConfig.String("Nscrud") + "/version_tipo_transaccion/GetDefinitiveVersion/?fecha=" + time.Now().Format("2006-01-02"))
		if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/version_tipo_transaccion/GetDefinitiveVersion/?fecha="+time.Now().Format("2006-01-02"), &tipoTransaccion); err == nil {
			if tipoTransaccion != nil {
				respuesta := optimize.ProccDigest(tipoTransaccion, getValuesDetallesTiposTr, nil, 3)
				c.Data["json"] = respuesta
			}
		} else {
			panic(err)
		}
	}).Catch(func(e try.E) {
		beego.Error("expc ", e)
		c.Data["json"] = models.Alert{Type: "error", Code: "E_0458", Body: e}
	})
}
