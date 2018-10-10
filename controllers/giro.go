package controllers

import (
	"encoding/json"
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/api_mid_financiera/models"
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
	var response map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &giro); err == nil {
		//beego.Error("valor giro ", giro)
		if err = request.SendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/giro/RegistrarGiro", "POST", &response, giro); err == nil {
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
