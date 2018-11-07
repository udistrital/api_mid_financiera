package controllers

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/api_mid_financiera/models"
	"github.com/udistrital/utils_oas/optimize"
	"github.com/udistrital/utils_oas/request"
)

// CuentasBancariasController operations for Cuentas_bancarias
type CuentasBancariasController struct {
	beego.Controller
}

// URLMapping ...
func (c *CuentasBancariasController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Cuentas_bancarias
// @Param	body		body 	models.Cuentas_bancarias	true		"body for Cuentas_bancarias content"
// @Success 201 {object} models.Cuentas_bancarias
// @Failure 403 body is empty
// @router / [post]
func (c *CuentasBancariasController) Post() {

}

// GetOne ...
// @Title GetOne
// @Description get Cuentas_bancarias by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Cuentas_bancarias
// @Failure 403 :id is empty
// @router /:id [get]
func (c *CuentasBancariasController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get Cuentas_bancarias
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Cuentas_bancarias
// @Failure 403
// @router / [get]
func (c *CuentasBancariasController) GetAll() {
	var cuentasBancarias []interface{}
	var limit int64 = 10
	var offset int64
	var query string
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	if r := c.GetString("query"); r != "" {
		query = r
	}
	beego.Info("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/cuenta_bancaria/?limit="+strconv.FormatInt(limit, 10)+"&offset="+strconv.FormatInt(offset, 10)+"&query="+query)
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/cuenta_bancaria/?limit="+strconv.FormatInt(limit, 10)+"&offset="+strconv.FormatInt(offset, 10)+"&query="+query, &cuentasBancarias); err == nil {
		if cuentasBancarias != nil {
			beego.Info("Cuentas Bancarias",cuentasBancarias)
			respuesta := optimize.ProccDigest(cuentasBancarias, getValuesCuentas, nil, 3)
			beego.Info("respuesta ",respuesta)
			c.Data["json"] = respuesta
		}else{
			beego.Error("RESPUESTA VACIA")
		}
	} else {
		beego.Error("Error ", err)
		c.Data["json"] = models.Alert{Type: "error", Code: "E_0458", Body: err}
	}
	c.ServeJSON()
}

func getValuesCuentas(rpintfc interface{}, params ...interface{}) (res interface{}) {

	var sucursal []interface{}
	beego.Info("rpintfc inicio ",rpintfc)
	idSucursalStr := strconv.FormatFloat(rpintfc.(map[string]interface{})["Sucursal"].(float64), 'f', -1, 64)
	beego.Info("IdSucursal ",idSucursalStr)

	if err := request.GetJson(beego.AppConfig.String("coreOrganizacionService")+"organizacion/?query=TipoOrganizacion.CodigoAbreviacion:SU,Id:"+idSucursalStr, &sucursal); err == nil {
		rpintfc.(map[string]interface{})["Sucursal"] = sucursal[0]
	}else{
		beego.Error("Error",err.Error())
	}
	resBanco, err := GetBancoSucursal(idSucursalStr)
	if err == nil {
		rpintfc.(map[string]interface{})["Banco"] = resBanco.([]interface{})[0].(map[string]interface{})
	} else {
		beego.Error("Error", err.Error())
	}
		beego.Info("rpintfc fin ",rpintfc)
	return rpintfc
}

// Put ...
// @Title Put
// @Description update the Cuentas_bancarias
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Cuentas_bancarias	true		"body for Cuentas_bancarias content"
// @Success 200 {object} models.Cuentas_bancarias
// @Failure 403 :id is not int
// @router /:id [put]
func (c *CuentasBancariasController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Cuentas_bancarias
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *CuentasBancariasController) Delete() {

}
