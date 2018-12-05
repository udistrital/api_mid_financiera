package controllers

import (
	"encoding/json"
	"errors"
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
			beego.Info("detalle transaccion ", detalleTransaccion)
			err, responseDT := SaveForTipoTr(urlCrud+"/detalle_tipo_transaccion_version", detalleTransaccion)
			if err != nil {
				beego.Error(err)
				panic(err)
			}
			responseRoute["detalle_tipo_transaccion_version"] = responseDT
			panic("errorrrrrrr!!!")
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = models.Alert{Type: "success", Code: "S_543", Body: responseRoute}
		}).Rollback(func(response interface{}, error interface{}) {
			beego.Error("Error Rollback ", error)
			urlCrud := "http://" + beego.AppConfig.String("Urlcrud") + ":" + beego.AppConfig.String("Portcrud") + "/" + beego.AppConfig.String("Nscrud")
			respuestas := response.(map[string]interface{})
			for key, value := range respuestas {
				body := value.(map[string]interface{})["Body"].(map[string]interface{})
				beego.Error("value ", value, "key ", key)
				if body["Id"] != nil {
					id := strconv.Itoa(int(body["Id"].(float64)))
					err = request.SendJson(urlCrud+"/"+key+"/"+id, "DELETE", &response, nil)
					beego.Error("response delete ", response)
					if err != nil {
						beego.Error(err)
						panic(err)
					}
				}
			}
			c.Data["json"] = models.Alert{Type: "error", Code: "E_0458", Body: respuestas}
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
