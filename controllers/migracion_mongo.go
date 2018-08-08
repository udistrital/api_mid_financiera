package controllers

import (
	"errors"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/manucorporat/try"
	"github.com/udistrital/utils_oas/request"
)

var (
	urlMongo = "http://" + beego.AppConfig.String("financieraMongoCurdApiService")
	urlCrud  = "http://" + beego.AppConfig.String("Urlcrud") + ":" + beego.AppConfig.String("Portcrud") + "/" + beego.AppConfig.String("Nscrud")
)

// MigracionApropiacionController operations for Inversion
type MigracionMongoController struct {
	beego.Controller
}

// MigrarApropiacion ...
// @Title MigrarApropiacion
// @Description Get Arbol Rubros By UE
// @Param	vigencia		path 	int64	true		"vigencia a migrar"
// @Success 200 {object} models.Rubro
// @Failure 403
// @router /MigrarApropiacion/:vigencia [get]
func (c *MigracionMongoController) MigrarApropiacion() {
	response := make(map[string]interface{})

	try.This(func() {
		vigencia, err := c.GetInt(":vigencia")
		if err != nil {
			panic("E_0458")
		}
		if data, err := migrarApropiacion(vigencia); err == nil {
			c.Data["json"] = data
		} else {
			panic("Migration Error")
		}

	}).Catch(func(e try.E) {
		response["Type"] = "Error"
		response["Code"] = e
		c.Data["json"] = response
	})

	c.ServeJSON()
}

//Migrar Apropiaciones a Mongo desde api Crud
func migrarApropiacion(vigencia int) (res interface{}, err error) {

	try.This(func() {
		var (
			apropiacionesData []map[string]interface{}
			mongoData         map[string]interface{}
			serviceResponse   []interface{}
			externalResponse  map[string]interface{}
		)

		uri := urlCrud + "/apropiacion?limit=-1&query=Vigencia:" + strconv.Itoa(vigencia)
		if err := request.GetJson(uri, &apropiacionesData); err != nil {
			beego.Info(err.Error())
			panic("Crud Service Error")
		}
		uri = urlMongo + "/arbol_rubro_apropiaciones/RegistrarApropiacionInicial/" + strconv.Itoa(vigencia)
		for _, data := range apropiacionesData {
			mongoData = make(map[string]interface{})
			mongoData = data["Rubro"].(map[string]interface{})
			mongoData["UnidadEjecutora"] = strconv.Itoa(int(mongoData["UnidadEjecutora"].(float64)))
			mongoData["Id"] = data["Id"]
			mongoData["ApropiacionInicial"] = data["Valor"]
			if err = request.SendJson(uri, "POST", &externalResponse, &mongoData); err == nil {
				if externalResponse["Type"].(string) == "success" {
					serviceResponse = append(serviceResponse, "Success Migrate "+strconv.Itoa(int(data["Id"].(float64))))
				} else {
					serviceResponse = append(serviceResponse, "Error Migrate "+strconv.Itoa(int(data["Id"].(float64))))
				}
			} else {
				serviceResponse = append(serviceResponse, "Error Migrate "+strconv.Itoa(int(data["Id"].(float64))))
			}
		}
		res = serviceResponse
	}).Catch(func(e try.E) {
		beego.Error("error ", e)
		err = errors.New(e.(string))
	})

	return
}
