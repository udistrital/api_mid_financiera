package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/api_financiera/models"
	"github.com/udistrital/utils_oas/request"
)

// InversionController operations for Inversion
type InversionController struct {
	beego.Controller
}


// URLMapping ...
func (c *InversionController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Inversion
// @Param	body		body 	models.Inversion	true		"body for Inversion content"
// @Success 201 {object} models.Inversion
// @Failure 403 body is empty
// @router / [post]
func (c *InversionController) Post() {

}


// CreacionInversion ...
// @Title Creacion Inversion
// @Description create Inversion
// @Param	body		body 	interface{}	true		"body for Inversion content"
// @Success 201 {object} interface{}
// @Failure 403 body is empty
// @router CreateInversion/ [post]
func (c *InversionController) CreateInversion() {
	defer c.ServeJSON()
	var v interface{}
	var respuesta interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {

		if err := request.SendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/cancelacion_inversion/CreateCancelacion", "POST", &respuesta, v); err == nil {
			c.Data["json"] = respuesta
			} else {
					beego.Error(err.Error())
					c.Data["json"] = models.Alert{Type:"error",Code:"E_0458",Body:err}
			}

	}else{
			beego.Error(err.Error())
			c.Data["json"] = models.Alert{Type:"error",Code:"E_0458",Body:err}
	}
}
// GetAllCancelaciones ...
// @Title GetAllCancelaciones
// @Description get cancelacion Inversion
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Homologacion_rubro
// @Failure 403
// @router GetAllCancelaciones/ [get]
func (c *InversionController) GetAllCancelaciones() {
	defer c.ServeJSON()
	var cancelaciones []interface{}
	var limit int64 = 10
	var offset int64
	var query string
	query = "Activo:true"
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	if r := c.GetString("query"); r != "" {
		query = r + "," + query
	}
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/cancelacion_inversion_estado_cancelacion?limit="+strconv.FormatInt(limit, 10)+"&offset="+strconv.FormatInt(offset, 10)+"&query="+query, &cancelaciones); err == nil {
		if cancelaciones != nil{
			c.Data["json"]=cancelaciones
		}
	}
}

// GetOne ...
// @Title GetOne
// @Description get Inversion by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Inversion
// @Failure 403 :id is empty
// @router /:id [get]
func (c *InversionController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get Inversion
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Inversion
// @Failure 403
// @router / [get]
func (c *InversionController) GetAll() {

}

// Put ...
// @Title Put
// @Description update the Inversion
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Inversion	true		"body for Inversion content"
// @Success 200 {object} models.Inversion
// @Failure 403 :id is not int
// @router /:id [put]
func (c *InversionController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Inversion
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *InversionController) Delete() {

}
