package controllers

import (
	"encoding/json"
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/api_financiera/models"
	"github.com/udistrital/utils_oas/request"
)

// InversionController operations for Inversion
type InversionController struct {
	beego.Controller
}

type estadoCancelacion struct {
	CancelacionInversion       interface{}
	EstadoCancelacionInversion interface{}
	Activo                     bool
	Usuario                    int
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
	var v map[string]interface{}
	var respuesta map[string]interface{}
	var estadoResp []interface{}
	var respuestaCreaEst map[string]interface{}
	var inversionCanc map[string]interface{}
	var usuario int

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		cancelacionInversion := v["cancelacionInversion"]
		inversionCanc = cancelacionInversion.(map[string]interface{})
		if err := request.SendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/ingreso_sin_situacion_fondos", "POST", &respuesta, cancelacionInversion); err == nil {
			if strings.Compare(respuesta["Type"].(string), "success") == 0 {
				if err = request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/estado_cancelacion_inversion?query=numeroOrden:1", &estadoResp); err == nil {
					usuario = int(inversionCanc["UsuarioEjecucion"].(float64))
					estadoCancInv := &estadoCancelacion{CancelacionInversion: respuesta["Body"], EstadoCancelacionInversion: estadoResp[0], Activo: true, Usuario: usuario}
					if err := request.SendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/cancelacion_inversion_estado_cancelacion", "POST", &respuestaCreaEst, estadoCancInv); err == nil {
						if strings.Compare(respuestaCreaEst["Type"].(string), "success") == 0 {
							alert := models.Alert{Type: "success", Code: "S_543", Body: respuesta["Body"]}
							c.Data["json"] = alert
							beego.Error(c.Data["json"])
							c.Ctx.Output.SetStatus(201)
						}
					} else {
						alert := models.Alert{Type: "error", Code: "E_0458", Body: err}
						c.Data["json"] = alert
					}
				} else {
					alert := models.Alert{Type: "error", Code: "E_0458", Body: err}
					c.Data["json"] = alert
				}
			}
		} else {
			alert := models.Alert{Type: "error", Code: "E_0458", Body: err}
			beego.Error(err.Error())
			c.Data["json"] = alert

		}
	} else {
		beego.Error(err.Error())
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
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
