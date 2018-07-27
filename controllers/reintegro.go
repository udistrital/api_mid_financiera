package controllers

import (
	"github.com/astaxie/beego"
	"strings"
	"encoding/json"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/api_mid_financiera/models"
)

// ReintegroController operations for Reintegro
type ReintegroController struct {
	beego.Controller
}

// URLMapping ...
func (c *ReintegroController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Reintegro
// @Param	body		body 	models.Reintegro	true		"body for Reintegro content"
// @Success 201 {object} models.Reintegro
// @Failure 403 body is empty
// @router / [post]
func (c *ReintegroController) Post() {

}

// GetOne ...
// @Title GetOne
// @Description get Reintegro by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Reintegro
// @Failure 403 :id is empty
// @router /:id [get]
func (c *ReintegroController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get Reintegro
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Reintegro
// @Failure 403
// @router / [get]
func (c *ReintegroController) GetAll() {

}

// Put ...
// @Title Put
// @Description update the Reintegro
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Reintegro	true		"body for Reintegro content"
// @Success 200 {object} models.Reintegro
// @Failure 403 :id is not int
// @router /:id [put]
func (c *ReintegroController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Reintegro
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *ReintegroController) Delete() {

}

// Create ...
// @Title Create
// @Description create Reingreso
// @Param	body		body interface{} true		"body for Reingreso content"
// @Success 201 {object} interface{}
// @Failure 403 body is empty
// @router /Create [post]
func (c *ReintegroController) Create() {
	var reintegro map[string]interface{}
	var resReintegro map[string]interface{}


  defer c.ServeJSON()

  if err := json.Unmarshal(c.Ctx.Input.RequestBody, &reintegro); err == nil {
				if err = request.SendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/reintegro/Create", "POST", &resReintegro, reintegro);err == nil{
				 if (strings.Compare(resReintegro["Type"].(string),"success")==0){
					 c.Data["json"]= models.Alert{Type:"success",Code:"S_543",Body:resReintegro["Body"]}
					 c.Ctx.Output.SetStatus(201)
				 }else{
					 c.Data["json"]= models.Alert{Type:resReintegro["Type"].(string),Code:resReintegro["Code"].(string),Body:resReintegro["Body"]}
				 }
			 	}
			}else{
				c.Data["json"]= models.Alert{Type:"error",Code:"E_0458",Body:err}
			}
		}
