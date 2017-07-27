package controllers

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/api_mid_financiera/models"
	"github.com/udistrital/api_mid_financiera/utilidades"
)

type RubroController struct {
	beego.Controller
}

// GenerarPac ...
// @Title GenerarPac
// @Description Get PAC By Rubro
// @Param	vigencia		path 	string	true		"vigencia del pac para los rubros"
// @Param	finicio		query 	string	true		"fecha inicial del rango del PAC"
// @Param	ffin		query 	string	true		"fecha final del rango del PAC"
// @Success 200 {object} interface{}
// @Failure 403
// @router /GenerarPac/:vigencia [get]
func (c *RubroController) GenerarPac() {
	vigStr := c.Ctx.Input.Param(":vigencia")
	vigencia, err := strconv.Atoi(vigStr)
	finicioStr := c.GetString("finicio")
	ffinStr := c.GetString("ffin")
	if len(finicioStr) == 0 || len(ffinStr) == 0 {

		c.Data["json"] = models.Alert{Code: "E_XXX", Body: nil, Type: "error"}
		c.ServeJSON()
	} else {
		var idapr string
		var idfuente string
		var rubrosHijos []map[string]interface{}

		if err != nil {
			fmt.Println(err.Error())
			c.Data["json"] = models.Alert{Code: "E_XXX", Body: err.Error(), Type: "error"}
			c.ServeJSON()
		} else {
			if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/apropiacion/GetApropiacionesHijo/"+strconv.Itoa(vigencia)+"?tipo=2", &rubrosHijos); err == nil {
				for _, rubro := range rubrosHijos {
					var ingresoRubro []map[string]interface{}
					err = utilidades.FillStruct(rubro["id"], &idapr)
					if err != nil {
						fmt.Println(err.Error())
						c.Data["json"] = models.Alert{Code: "E_XXX", Body: err.Error(), Type: "error"}
						c.ServeJSON()
					}
					if rubro["idfuente"] != nil {
						err = utilidades.FillStruct(rubro["idfuente"], &idfuente)
						if err != nil {
							fmt.Println(err.Error())
							c.Data["json"] = models.Alert{Code: "E_XXX", Body: err.Error(), Type: "error"}
							c.ServeJSON()
						}
						//fmt.Println("rubro ", idrubro)
						fmt.Println("fuente ", rubro["idfuente"])
					} else {
						idfuente = "0"
					}
					if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/rubro/GetRubroIngreso?apropiacion="+idapr+"&fuente="+idfuente+"&finicio="+finicioStr+"&ffin="+ffinStr, &ingresoRubro); err == nil {
						if ingresoRubro != nil {
							fmt.Println(ingresoRubro)
						}

					}
					idapr = ""
					idfuente = ""

				}

				c.Data["json"] = rubrosHijos
			} else {
				fmt.Println(err.Error())
				c.Data["json"] = models.Alert{Code: "E_XXX", Body: err.Error(), Type: "error"}
			}
		}

		c.ServeJSON()
	}

}
