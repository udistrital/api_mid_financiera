package controllers

import (
	"errors"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/manucorporat/try"
	"github.com/udistrital/api_mid_financiera/models"
	"github.com/udistrital/utils_oas/formatdata"
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

// MigrarMovimiento ...
// @Title MigrarMovimiento
// @Description Get Arbol Rubros By UE
// @Param	vigencia		path 	int64	true		"vigencia a migrar"
// @Success 200 {object} models.Rubro
// @Failure 403
// @router /MigrarMovimiento/:tipo/:vigencia [get]
func (c *MigracionMongoController) MigrarMovimiento() {
	response := make(map[string]interface{})
	var f func(int) (interface{}, error)
	try.This(func() {
		vigencia, err := c.GetInt(":vigencia")
		if err != nil {
			panic("E_0458")
		}
		tipo := c.GetString(":tipo")
		switch tipo {
		case "Cdp":
			f = migrarCdp
		case "Rp":
			f = migrarRp
		default:
			panic("Migration Error")
		}
		if data, err := f(vigencia); err == nil {
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

func migrarCdp(vigencia int) (res interface{}, err error) {
	try.This(func() {
		var (
			cdpData          []map[string]interface{}
			mongoData        map[string]interface{}
			serviceResponse  []interface{}
			externalResponse map[string]interface{}
			afectacionArray  []interface{}
			afetcacionData   map[string]interface{}
		)

		// uri := urlCrud + "/disponibilidad?limit=-1&query=Vigencia:" + strconv.Itoa(vigencia)
		uri := urlCrud + "/disponibilidad?limit=-1&query=Vigencia:" + strconv.Itoa(vigencia)
		if err := request.GetJson(uri, &cdpData); err != nil {
			beego.Info(err.Error())
			panic("Crud Service Error")
		}
		uri = urlMongo + "/arbol_rubro_apropiaciones/RegistrarMovimiento/Cdp"
		for _, data := range cdpData {
			mongoData = make(map[string]interface{})
			mongoData = data
			dateStr := mongoData["FechaRegistro"].(string)
			mongoData["Vigencia"] = strconv.Itoa(int(mongoData["Vigencia"].(float64)))
			t, err := time.Parse(time.RFC3339, dateStr)
			if err != nil {
				panic(err.Error())
			}
			dataAfectacion := mongoData["DisponibilidadApropiacion"].([]interface{})
			for _, data := range dataAfectacion {
				afetcacionData = make(map[string]interface{})
				// dispAfectInt := make(map[string]interface{})
				dispAfect := models.DisponibilidadApropiacion{}
				err = formatdata.FillStruct(data, &dispAfect)
				if err != nil {
					panic(err.Error())
				}
				afetcacionData["Apropiacion"] = dispAfect.Apropiacion.Id
				afetcacionData["FuenteCodigo"] = dispAfect.FuenteFinanciamiento.Codigo
				afetcacionData["FuenteNombre"] = dispAfect.FuenteFinanciamiento.Nombre
				afetcacionData["Rubro"] = dispAfect.Apropiacion.Rubro.Codigo
				afetcacionData["UnidadEjecutora"] = strconv.Itoa(int(dispAfect.Apropiacion.Rubro.UnidadEjecutora))
				afetcacionData["Valor"] = dispAfect.Valor
				afectacionArray = append(afectacionArray, afetcacionData)
			}
			mongoData["Afectacion"] = afectacionArray
			mongoData["MesRegistro"] = strconv.Itoa(int(t.Month()))
			if err = request.SendJson(uri, "POST", &externalResponse, &mongoData); err == nil {
				if externalResponse["Type"].(string) == "success" {
					//serviceResponse = append(serviceResponse, "Success Migrate "+strconv.Itoa(int(data["Id"].(float64))))
				} else {
					serviceResponse = append(serviceResponse, "Error Migrate "+strconv.Itoa(int(data["Id"].(float64))))
				}
			} else {
				serviceResponse = append(serviceResponse, "Error Migrate "+strconv.Itoa(int(data["Id"].(float64))))
			}
			afectacionArray = nil
		}
		res = serviceResponse
	}).Catch(func(e try.E) {
		beego.Error("error ", e)
		err = errors.New(e.(string))
	})

	return
}

func migrarRp(vigencia int) (res interface{}, err error) {
	try.This(func() {
		var (
			rpData           []map[string]interface{}
			mongoData        map[string]interface{}
			serviceResponse  []interface{}
			externalResponse map[string]interface{}
			afectacionArray  []interface{}
			afetcacionData   map[string]interface{}
		)

		// uri := urlCrud + "/disponibilidad?limit=-1&query=Vigencia:" + strconv.Itoa(vigencia)
		uri := urlCrud + "/registro_presupuestal?limit=-1&query=Vigencia:" + strconv.Itoa(vigencia)
		if err := request.GetJson(uri, &rpData); err != nil {
			beego.Info(err.Error())
			panic("Crud Service Error")
		}
		uri = urlMongo + "/arbol_rubro_apropiaciones/RegistrarMovimiento/Rp"
		for _, data := range rpData {
			mongoData = make(map[string]interface{})
			mongoData = data
			dateStr := mongoData["FechaRegistro"].(string)
			mongoData["Vigencia"] = strconv.Itoa(int(mongoData["Vigencia"].(float64)))
			t, err := time.Parse(time.RFC3339, dateStr)
			if err != nil {
				panic(err.Error())
			}
			dataAfectacion := mongoData["RegistroPresupuestalDisponibilidadApropiacion"].([]interface{})
			for _, data := range dataAfectacion {
				afetcacionData = make(map[string]interface{})
				// dispAfectInt := make(map[string]interface{})
				rpAfect := models.RegistroPresupuestalDisponibilidadApropiacion{}
				err = formatdata.FillStruct(data, &rpAfect)
				if err != nil {
					panic(err.Error())
				}
				afetcacionData["Apropiacion"] = rpAfect.DisponibilidadApropiacion.Apropiacion.Id
				afetcacionData["FuenteCodigo"] = rpAfect.DisponibilidadApropiacion.FuenteFinanciamiento.Codigo
				afetcacionData["FuenteNombre"] = rpAfect.DisponibilidadApropiacion.FuenteFinanciamiento.Nombre
				afetcacionData["Rubro"] = rpAfect.DisponibilidadApropiacion.Apropiacion.Rubro.Codigo
				afetcacionData["UnidadEjecutora"] = strconv.Itoa(int(rpAfect.DisponibilidadApropiacion.Apropiacion.Rubro.UnidadEjecutora))
				afetcacionData["Valor"] = rpAfect.Valor
				mongoData["Disponibilidad"] = rpAfect.DisponibilidadApropiacion.Disponibilidad.Id

				afectacionArray = append(afectacionArray, afetcacionData)
			}
			mongoData["Afectacion"] = afectacionArray
			mongoData["MesRegistro"] = strconv.Itoa(int(t.Month()))
			if err = request.SendJson(uri, "POST", &externalResponse, &mongoData); err == nil {
				if externalResponse["Type"].(string) == "success" {
					serviceResponse = append(serviceResponse, "Success Migrate "+strconv.Itoa(int(data["Id"].(float64))))
				} else {
					serviceResponse = append(serviceResponse, "Error Migrate "+strconv.Itoa(int(data["Id"].(float64))))
				}
			} else {
				serviceResponse = append(serviceResponse, "Error Migrate "+strconv.Itoa(int(data["Id"].(float64))))
			}
			afectacionArray = nil
		}
		res = serviceResponse
	}).Catch(func(e try.E) {
		beego.Error("error ", e)
		err = errors.New(e.(string))
	})

	return
}
