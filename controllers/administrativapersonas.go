package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/request"
)

// AdministrativaPersonasController operations for AdministrativaPersonas
type AdministrativaPersonasController struct {
	beego.Controller
}

// URLMapping ...
func (c *AdministrativaPersonasController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create AdministrativaPersonas
// @Param	body		body 	models.AdministrativaPersonas	true		"body for AdministrativaPersonas content"
// @Success 201 {object} models.AdministrativaPersonas
// @Failure 403 body is empty
// @router / [post]
func (c *AdministrativaPersonasController) Post() {

}

// GetOne ...
// @Title GetOne
// @Description get AdministrativaPersonas by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.AdministrativaPersonas
// @Failure 403 :id is empty
// @router /:id [get]
func (c *AdministrativaPersonasController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get AdministrativaPersonas
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.AdministrativaPersonas
// @Failure 403
// @router / [get]
func (c *AdministrativaPersonasController) GetAll() {

}

// Put ...
// @Title Put
// @Description update the AdministrativaPersonas
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.AdministrativaPersonas	true		"body for AdministrativaPersonas content"
// @Success 200 {object} models.AdministrativaPersonas
// @Failure 403 :id is not int
// @router /:id [put]
func (c *AdministrativaPersonasController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the AdministrativaPersonas
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *AdministrativaPersonasController) Delete() {

}

// GetPersona ...
// @Title GetPersona
// @Description get search person by id type and number
// @Param	numberId	query	string	false	"identification number"
// @Param	typeId	query	string	false	"type id"
// @Success 200 {object} models.AdministrativaPersonas
// @Failure 403
// @router /GetPersona/ [get]
func (c *AdministrativaPersonasController) GetPersona() {
	defer c.ServeJSON()
	var numberIDStr string
	var typeIDStr string
	var resProveedor []map[string]interface{}
	var resPersonaNat []map[string]interface{}
	var resPersonaJur []map[string]interface{}
	beego.Error("going on get persona")
	if v := c.GetString("numberId"); v != "" {
		numberIDStr = v
	}

	if v := c.GetString("typeId"); v != "" {
		typeIDStr = v
	}
	if err := request.GetJson("http://"+beego.AppConfig.String("AdministrativaAmazonService")+"informacion_proveedor/?query=NumDocumento:"+numberIDStr+"&limit=1", &resProveedor); err == nil {
		if resProveedor != nil {
			for _, v := range resProveedor {
				if v["Tipopersona"].(string) == "NATURAL" {
					if err := request.GetJson("http://"+beego.AppConfig.String("AdministrativaAmazonService")+"informacion_persona_natural/?query=Id:"+numberIDStr+",TipoDocumento.Id:"+typeIDStr+"&limit=1", &resPersonaNat); err == nil {
						if resPersonaNat != nil {
							c.Data["json"] = v
							return
						}
					} else {
						beego.Error("Error" + err.Error())
					}
				} else {
					if v["Tipopersona"].(string) == "JURIDICA" {
						if err := request.GetJson("http://"+beego.AppConfig.String("AdministrativaAmazonService")+"/informacion_persona_juridica/"+numberIDStr, &resPersonaJur); err != nil {
							if resPersonaJur != nil && typeIDStr == "11" {
								c.Data["json"] = v
								return
							}
						} else {
							beego.Error("Error" + err.Error())
						}
					}
				}
			}
		}
	} else {
		beego.Error("Error", err.Error())
		return
	}
}
