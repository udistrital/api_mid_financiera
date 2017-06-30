package controllers

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/api_mid_financiera/tools"
)

// PartidadobleController operations for Partidadoble
type PartidadobleController struct {
	beego.Controller
}

// URLMapping ...
func (c *PartidadobleController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Partidadoble
// @Param	body		body 	models.Partidadoble	true		"body for Partidadoble content"
// @Success 201 {object} models.Partidadoble
// @Failure 403 body is empty
// @router / [post]
func (c *PartidadobleController) Post() {
	tool := tools.EntornoReglas{}
	tool.Agregar_dominio("presupuesto")
	c.Data["json"] = tool
	c.ServeJSON()
}

// GetOne ...
// @Title GetOne
// @Description get Partidadoble by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Partidadoble
// @Failure 403 :id is empty
// @router /:id [get]
func (c *PartidadobleController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get Partidadoble
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Partidadoble
// @Failure 403
// @router / [get]
func (c *PartidadobleController) GetAll() {
	tool := new(tools.EntornoReglas)
	tool.Agregar_dominio("Presupuesto")
	res, err := tool.Agregar_predicado_dinamico("nombre_regl:argoService.necesidad|Id|Valor")
	fmt.Println("err: ", err)
	//fmt.Println(tool)
	c.Data["json"] = res
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Partidadoble
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Partidadoble	true		"body for Partidadoble content"
// @Success 200 {object} models.Partidadoble
// @Failure 403 :id is not int
// @router /:id [put]
func (c *PartidadobleController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Partidadoble
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *PartidadobleController) Delete() {

}
