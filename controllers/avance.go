package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/request"
)

// AvanceController operations for Avance
type AvanceController struct {
	beego.Controller
}

// URLMapping ...
func (c *AvanceController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Avance
// @Param	body		body 	models.Avance	true		"body for Avance content"
// @Success 201 {object} models.Avance
// @Failure 403 body is empty
// @router / [post]
func (c *AvanceController) Post() {

}

// GetOne ...
// @Title GetOne
// @Description get Avance by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Avance
// @Failure 403 :id is empty
// @router /:id [get]
func (c *AvanceController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get Avance
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Avance
// @Failure 403
// @router / [get]
func (c *AvanceController) GetAll() {

}

// Put ...
// @Title Put
// @Description update the Avance
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Avance	true		"body for Avance content"
// @Success 200 {object} models.Avance
// @Failure 403 :id is not int
// @router /:id [put]
func (c *AvanceController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Avance
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *AvanceController) Delete() {

}

// GetAvanceById ...
// @Title GetAvanceById
// @Description get All information of an advance payment by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Legalizacion_avance
// @Failure 403 :id is empty
// @router /:id [get]
func (c *AvanceController) GetAvanceById() {
	idStr := c.Ctx.Input.Param(":id")
	var resEstado []map[string]interface{}
	var rpintfc interface{}
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/solicitud_tipo_avance/?query=query=SolicitudAvance.Id:"+",AvanceLegalizacion.Id:"+idStr+"&sortby=Id&limit=-1&order=asc", &resEstado); err == nil {
		if resEstado[0] != nil {
			rpintfc.(map[string]interface{})["Estado"] = resEstado[0]["Estado"]
		}
	} else {
		beego.Error("Error", err.Error())
	}
}
