package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/manucorporat/try"
	"github.com/udistrital/api_mid_financiera/models"
	"github.com/udistrital/utils_oas/request"
)

type FuenteFinanciamientoController struct {
	beego.Controller
}

// URLMapping ...
func (c *FuenteFinanciamientoController) URLMapping() {
	c.Mapping("Post", c.Post)
}

// Post ...
// @Title Post
// @Description create FuenteFinanciamiento
// @Param	body		body 	string	true		"body for Rubro content"
// @Success 201 {int} models.Rubro
// @Failure 403 body is empty
// @router /RegistrarFuente/:UnidadEjecutora [post]
func (c *FuenteFinanciamientoController) Post() {
	var v map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		ue, err := c.GetInt(":UnidadEjecutora")
		if err != nil {
			panic("Parameter Error")
		}
		var resCrud map[string]interface{}
		if err := request.SendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/fuente_financiamiento/RegistrarFuenteFinanciamientoTr", "POST", &resCrud, &v); err == nil {
			resCrud["FuenteFinanciamiento"].(map[string]interface{})["UnidadEjecutora"] = ue
			alert := models.Alert{Type: "success", Code: "S_F0001", Body: resCrud}
			c.Data["json"] = alert
		} else {
			alert := models.Alert{Type: "error", Code: "E_0458", Body: err}
			c.Data["json"] = alert
		}
	} else {
		alert := models.Alert{Type: "error", Code: "E_0458", Body: err}
		c.Data["json"] = alert
	}
	c.ServeJSON()
}

// RegistrarModificacionFuente ...
// @Title RegistrarModificacionFuente
// @Description create ModificacionFuenteFinanciamiento
// @Param	body		body 	string	true		"body for Rubro content"
// @Success 201 {int} models.Rubro
// @Failure 403 body is empty
// @router /RegistrarModificacionFuente [post]
func (c *FuenteFinanciamientoController) RegistrarModificacionFuente() {
	var v []map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		var resCrud interface{}
		if err := request.SendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/fuente_financiamiento/MovimientoFuenteFinanciamientoTr", "POST", &resCrud, &v); err == nil {
			alert := models.Alert{Type: "success", Code: "S_F0001", Body: resCrud}
			c.Data["json"] = alert
		} else {
			alert := models.Alert{Type: "error", Code: "E_0458", Body: err}
			c.Data["json"] = alert
		}
	} else {
		alert := models.Alert{Type: "error", Code: "E_0458", Body: err}
		c.Data["json"] = alert
	}
	c.ServeJSON()
}

// AddFuenteFinanciamientoMongo... agrega la informacion de una fuente a mongo.
func AddFuenteFinanciamientoMongo(parameter ...interface{}) (err interface{}) {
	try.This(func() {
		//Convertir Datos retornados para registrarlos en mongo.
		dataMongo := parameter[0].(map[string]interface{})
		resM := make(map[string]interface{})
		//infoFuente := parameter[0].(map[string]interface{})["FuenteFinanciamiento"].(map[string]interface{})
		//infoAfectacion := parameter[0].(map[string]interface{})["AfectacionFuente"].(map[string]interface{})
		beego.Info("Data from job ", parameter[0])
		//beego.Info(infoAfectacion)
		Urlmongo := "http://" + beego.AppConfig.String("financieraMongoCurdApiService") + "/fuente_financiamiento"
		if err1 := request.SendJson(Urlmongo, "POST", &resM, &dataMongo); err1 == nil {
			if resM["Type"].(string) == "success" {
				err = err1
			} else {
				panic("Mongo api error")
			}
		} else {
			panic("Mongo Not Found")
		}
		//panic("No sé que más hacer !!!")
	}).Catch(func(e try.E) {
		beego.Error("Retroceder Tr ")
		beego.Error("Error ", e)
		idFuente := parameter[0].(map[string]interface{})["FuenteFinanciamiento"].(map[string]interface{})["Id"].(float64)
		beego.Error("IdFuente ", idFuente)
		var resCrud interface{}
		request.SendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/fuente_financiamiento/DeleteMovimientoFuenteFinanciamientoTr/"+strconv.Itoa(int(idFuente)), "DELETE", &resCrud, nil)
		beego.Error("res crud ", resCrud)
	})

	return
}

// AddModificacionFuenteFinanciamientoMongo... agrega la informacion de una Modificacion de fuente a mongo.
func AddModificacionFuenteFinanciamientoMongo(parameter ...interface{}) (err interface{}) {
	try.This(func() {
		//Convertir Datos retornados para registrarlos en mongo.
		//dataMongo := make(map[string]interface{})
		infoFuente := parameter[0].(map[string]interface{})["Body"].([]interface{})
		beego.Info(infoFuente)
		panic("No sé que más hacer !!!")
	}).Catch(func(e try.E) {

		beego.Error("Retroceder Tr ")
		beego.Error("Data ", parameter[0])
		infoFuente := parameter[0].(map[string]interface{})["Body"].([]interface{})
		var resCrud interface{}
		request.SendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/fuente_financiamiento/DeleteModificacionFuenteFinanciamientoTr", "POST", &resCrud, infoFuente)
		beego.Error("res crud ", resCrud)
	})

	return
}
