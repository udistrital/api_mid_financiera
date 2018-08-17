package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/request"
	"encoding/json"
	"strings"
	"github.com/udistrital/api_mid_financiera/models"
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


// Post ...
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
			if err = request.SendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/chequera/CreateChequeraEstado", "POST", &response, chequera);err == nil{
	 		 if (strings.Compare(response["Type"].(string),"success")==0){
	 			 c.Data["json"]= models.Alert{Type:"success",Code:"S_543",Body:response["Body"]}
	 			 c.Ctx.Output.SetStatus(201)
	 		 }else{
				 beego.Error("Error",response)
	 			 c.Data["json"]= models.Alert{Type:response["Type"].(string),Code:response["Code"].(string),Body:response["Body"]}
	 		 }
			}else{
				beego.Error("Error",err)
			 	c.Data["json"]= models.Alert{Type:"error",Code:"E_0458",Body:err}
			}
	}else{
		beego.Error("Error",err)
		c.Data["json"]= models.Alert{Type:"error",Code:"E_0458",Body:err}
	}
}
