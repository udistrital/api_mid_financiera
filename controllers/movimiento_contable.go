package controllers

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/manucorporat/try"
	"github.com/udistrital/api_mid_financiera/models"
	"github.com/udistrital/utils_oas/request"
)

// MovimientoContableController operations for movimiento_contable
type MovimientoContableController struct {
	beego.Controller
}

// URLMapping ...
func (c *MovimientoContableController) URLMapping() {
	c.Mapping("ResumenMovimientos", c.ResumenMovimientos)
}

// ResumenMovimientos ...
// @Title ResumenMovimientos
// @Description get ResumenMovimientos
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Cuentas_bancarias
// @Failure 403
// @router /ResumenMovimientos [get]
func (c *MovimientoContableController) ResumenMovimientos() {
	var resumenMovimientos []map[string]interface{}
	var movimientosContables []map[string]interface{}
	var limit int64 = 10
	var offset int64
	var query string
	newMov := make([]interface{}, 0, len(movimientosContables))
	//limit: 10 (default is 10)
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
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/movimiento_contable/GetSumMovimientos?"+"query="+query, &resumenMovimientos); err == nil {
		// respuesta := optimize.ProccDigest(movimientosContables, sumValuesMov, nil, 3)
		beego.Info(resumenMovimientos)
	} else {
		beego.Error("Error ", err)
		c.Data["json"] = models.Alert{Type: "error", Code: "E_0458", Body: err}
	}
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/movimiento_contable/?limit="+strconv.FormatInt(limit, 10)+"&offset="+strconv.FormatInt(offset, 10)+"&query="+query, &movimientosContables); err == nil {
		// respuesta := optimize.ProccDigest(movimientosContables, sumValuesMov, nil, 3)
		for _, value := range movimientosContables {
			if value["CuentaEspecial"] == nil && value["Credito"].(float64) != 0 {
				beego.Info("entra", value["Credito"])
				value["Credito"], _ = strconv.Atoi(resumenMovimientos[0]["credito"].(string))
				newMov = append(newMov, value)
			} else if value["CuentaEspecial"] == nil && value["Debito"] != 0 {
				newMov = append(newMov, value)
			}
		}
		c.Data["json"] = newMov
	} else {
		beego.Error("Error ", err)
		c.Data["json"] = models.Alert{Type: "error", Code: "E_0458", Body: err}
	}
	c.ServeJSON()
}

func sumValuesMov(rpintfc interface{}, params ...interface{}) (res interface{}) {

	var sucursal []interface{}
	idSucursalStr := strconv.FormatFloat(rpintfc.(map[string]interface{})["Sucursal"].(float64), 'f', -1, 64)
	if err := request.GetJson(beego.AppConfig.String("coreOrganizacionService")+"organizacion/?query=TipoOrganizacion.CodigoAbreviacion:SU,Id:"+idSucursalStr, &sucursal); err == nil {
		rpintfc.(map[string]interface{})["Sucursal"] = sucursal[0]
	} else {
		beego.Error("Error", err.Error())
	}
	resBanco, err := GetBancoSucursal(idSucursalStr)
	if err == nil {
		try.This(func() {
			rpintfc.(map[string]interface{})["Banco"] = resBanco.([]interface{})[0].(map[string]interface{})
		}).Catch(func(e try.E) {
			beego.Error("expc ", e)
		})
	} else {
		beego.Error("Error", err.Error())
	}
	return rpintfc
}
