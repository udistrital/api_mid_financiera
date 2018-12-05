package controllers

import (
	"github.com/astaxie/beego"

	"github.com/udistrital/api_mid_financiera/models"
	"github.com/udistrital/utils_oas/request"
)

// AcademicaPersonasController operations for Academica_personas
type AcademicaPersonasController struct {
	beego.Controller
}

// URLMapping ...
func (c *AcademicaPersonasController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Academica_personas
// @Param	body		body 	models.Academica_personas	true		"body for Academica_personas content"
// @Success 201 {object} models.Academica_personas
// @Failure 403 body is empty
// @router / [post]
func (c *AcademicaPersonasController) Post() {

}

// GetOne ...
// @Title GetOne
// @Description get Academica_personas by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Academica_personas
// @Failure 403 :id is empty
// @router /:id [get]
func (c *AcademicaPersonasController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get Academica_personas
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Academica_personas
// @Failure 403
// @router / [get]
func (c *AcademicaPersonasController) GetAll() {

}

// Put ...
// @Title Put
// @Description update the Academica_personas
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Academica_personas	true		"body for Academica_personas content"
// @Success 200 {object} models.Academica_personas
// @Failure 403 :id is not int
// @router /:id [put]
func (c *AcademicaPersonasController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Academica_personas
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *AcademicaPersonasController) Delete() {

}

// GetDocentebyId ...
// @Title GetDocentebyId
// @Description get search tercero(profesor) by id number
// @Param	Id	query	string	false	"id tercero"
// @Success 200 {object} models.Academica_personas
// @Failure 403
// @router /GetDocentebyId/ [get]
func (c *AcademicaPersonasController) GetDocentebyId() {
	defer c.ServeJSON()
	var Id string
	var resDocente map[string]interface{}
	if v := c.GetString("Id"); v != "" {
		Id = v
	}
	if err := request.GetJsonWSO2(beego.AppConfig.String("Wso2Service")+"servicios_academicos/consulta_datos_docente_planta/"+Id, &resDocente); err == nil {
		if resDocente["datosCollection"].(map[string]interface{})["datos"] != nil {
			c.Data["json"] = resDocente["datosCollection"].(map[string]interface{})["datos"].([]interface{})[0]
		}
	} else {
		c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
		beego.Error("Error", err.Error())
		return
	}
}
