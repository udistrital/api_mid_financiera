package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/manucorporat/try"
	"github.com/udistrital/utils_oas/request"
)

//  ApropiacionController operations for Apropiacion
type ApropiacionController struct {
	beego.Controller
}

// URLMapping ...
func (c *ApropiacionController) URLMapping() {
	c.Mapping("Post", c.Post)
}

// Post ...
// @Title Post
// @Description create Apropiacion
// @Param	body		body 	models.Apropiacion	true		"body for Apropiacion content"
// @Success 201 {int} models.Apropiacion
// @Failure 403 body is empty
// @router / [post]
func (c *ApropiacionController) Post() {
	var v map[string]interface{}
	var res map[string]interface{}
	var resM map[string]interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		try.This(func() {
			mongoData := v["Rubro"].(map[string]interface{})
			v["Rubro"] = map[string]interface{}{"Id": v["Rubro"].(map[string]interface{})["Id"]}
			beego.Info("Data Send By Client: ", v)
			urlcrud := "http://" + beego.AppConfig.String("Urlcrud") + ":" + beego.AppConfig.String("Portcrud") + "/" + beego.AppConfig.String("Nscrud") + "/apropiacion/"
			if err = request.SendJson(urlcrud, "POST", &res, &v); err == nil {
				if res["Type"] != nil && res["Type"].(string) == "success" {
					mongoData["Id"] = res["Body"].(map[string]interface{})["Id"]
					mongoData["ApropiacionInicial"] = v["Valor"]
					beego.Info("Post data: ", mongoData)
					beego.Info("Crud Recibed data: ", res["Body"])
					urlmongo := "http://" + beego.AppConfig.String("financieraMongoCurdApiService") + "/arbol_rubro_apropiaciones/RegistrarApropiacionInicial/" + strconv.Itoa(int(res["Body"].(map[string]interface{})["Vigencia"].(float64)))
					if err = request.SendJson(urlmongo, "POST", &resM, &mongoData); err == nil {
						if resM["Type"].(string) == "success" {
							c.Data["json"] = res
						} else {
							resul := res["Body"].(map[string]interface{})
							urlcrud = urlcrud + "/" + strconv.Itoa(int(resul["Id"].(float64)))
							request.SendJson(urlcrud, "DELETE", &resM, nil)
							beego.Info("Data ", resM)
							panic("Mongo API error")
						}
					} else {
						resul := res["Body"].(map[string]interface{})
						urlcrud = urlcrud + "/" + strconv.Itoa(int(resul["Id"].(float64)))
						request.SendJson(urlcrud, "DELETE", &resM, nil)
						beego.Info("Data ", resM)
						panic("Mongo API not Found")
					}
				} else {
					beego.Error(res)
					panic("Financiera CRUD not Found")
				}
			} else {
				panic("Financiera CRUD Service Error")
			}
		}).Catch(func(e try.E) {
			beego.Error("expc ", e)
			c.Data["json"] = map[string]interface{}{"Code": "E_0458", "Body": e, "Type": "error"}
		})
	} else {
		c.Data["json"] = map[string]interface{}{"Code": "E_0458", "Body": err, "Type": "error"}
	}
	c.ServeJSON()
}
