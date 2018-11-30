package controllers

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/api_financiera/models"
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
	var v map[string]interface{}
	var version map[string]interface{}
	var tipoTransaccionVersion map[string]interface{}
	var detalleTransaccion map[string]interface{}
	var response map[string]interface{}
	var respDetalle map[string]interface{}
	var resDelete map[string]interface{}
	urlCrud := "http://" + beego.AppConfig.String("Urlcrud") + ":" + beego.AppConfig.String("Portcrud") + "/" + beego.AppConfig.String("Nscrud")
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		version = v["version"].(map[string]interface{})
		detalleTransaccion = v["detalleTransaccion"].(map[string]interface{})
		if err = request.SendJson(urlCrud+"/version_tipo_transaccion", "POST", &response, version); err == nil {
			if strings.Compare(response["Type"].(string), "success") == 0 {
				if err = request.SendJson(urlCrud+"/tipo_transaccion_version", "POST", &tipoTransaccionVersion, response["Body"]); err == nil {
					detalleTransaccion["TipoTransaccionVersion"] = tipoTransaccionVersion
					if err = request.SendJson(urlCrud+"/detalle_tipo_transaccion_version", "POST", &respDetalle, detalleTransaccion); err == nil {
						c.Data["json"] = models.Alert{Type: "success", Code: "S_543", Body: respDetalle["Body"]}
						c.Ctx.Output.SetStatus(201)
					} else {

					}
				} else {
					idVersion := strconv.Itoa(int(response["Body"].(map[string]interface{})["Id"].(float64)))
					if errorDelete := request.SendJson(urlCrud+"/version_tipo_transaccion/"+idVersion, "DELETE", &resDelete, nil); errorDelete == nil {
						beego.Info("Data ", resDelete)
						panic("Mongo API Error")
					} else {
						beego.Info("Error delete ", errorDelete)
						panic("Delete API Error")
					}
				}

			} else {
				beego.Error("Error", response)
				c.Data["json"] = response
			}
		}
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
