package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/api_mid_financiera/tools"
)

// PartidadobleController operations for Partidadoble
type PartidadobleController struct {
	beego.Controller
}

// PDMovimientoContable movimienos
type PDMovimientoContable struct {
	Debito               int64
	Credito              int64
	CodigoCuentaContable string
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
// @Description valida en base a los movimientos obtenidos si se cumple o no el principio de partida doble
// @Param	body		body 	[]PDMovimientoContable	true		"cuerpo para el post del servicio, debe ser en un array"
// @Success 201 {object} bool
// @Failure 403 body is empty
// @router / [post]
func (c *PartidadobleController) Post() {
	println("entro")
	tool := tools.EntornoReglas{}
	//tool.Agregar_dominio("presupuesto")
	var v []PDMovimientoContable
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		predicados := ``
		tool.Agregar_dominio("Contabilidad_partida")

		for _, mov := range v {
			predicados += `mov_cuenta('` + mov.CodigoCuentaContable + `', ` + strconv.FormatInt(mov.Debito, 10) + `, ` + strconv.FormatInt(mov.Credito, 10) + `).`
		}
		tool.Agregar_dominio("Contabilidad")
		tool.Agregar_predicado(predicados)

		if resp := tool.Ejecutar_result("aplica_partida_doble(SD).", "SD"); resp == "1" {
			c.Data["json"] = true
		} else {
			c.Data["json"] = false
		}
	} else {
		c.Data["json"] = err.Error()
	}
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

	err := tool.Agregar_predicado_dinamico("ValorNecesidad:argoService.fuente_financiacion_rubro_necesidad|SolicitudNecesidad.Id|MontoParcial")
	fmt.Println("err: ", err)
	//fmt.Println(tool)
	c.Data["json"] = tool.Obtener_predicados()
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
