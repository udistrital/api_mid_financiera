package controllers

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/optimize"
)

// OrganizacionesCoreNewController operations for OrganizacionesCoreNew
type OrganizacionesCoreNewController struct {
	beego.Controller
}

// URLMapping ...
func (c *OrganizacionesCoreNewController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create OrganizacionesCoreNew
// @Param	body		body 	models.OrganizacionesCoreNew	true		"body for OrganizacionesCoreNew content"
// @Success 201 {object} models.OrganizacionesCoreNew
// @Failure 403 body is empty
// @router / [post]
func (c *OrganizacionesCoreNewController) Post() {

}

// GetOne ...
// @Title GetOne
// @Description get OrganizacionesCoreNew by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.OrganizacionesCoreNew
// @Failure 403 :id is empty
// @router /:id [get]
func (c *OrganizacionesCoreNewController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get OrganizacionesCoreNew
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.OrganizacionesCoreNew
// @Failure 403
// @router / [get]
func (c *OrganizacionesCoreNewController) GetAll() {

}


// GetOrganizacion ...
// @Title GetOrganizacion
// @Description get OrganizacionesCoreNew
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.OrganizacionesCoreNew
// @Failure 403
// @router getOrganizacion/
func (c *OrganizacionesCoreNewController) GetOrganizacion()() {

	var ente  []interface{}
	var tipoEnte []map[string]interface{}
	var limit int64 = 10
	var offset int64
	var query string
	var respuesta []map[string]interface{}
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


	if err := request.GetJson(beego.AppConfig.String("coreEnteService")+"tipo_ente?limit="+strconv.FormatInt(limit,10)+"&offset="+strconv.FormatInt(offset, 10)+"&query="+query, &tipoEnte); err == nil {
				beego.Info(tipoEnte);
				idEnte:=int(tipoEnte[0]["Id"].(float64));
				beego.Error(beego.AppConfig.String("coreEnteService")+"ente?limit=-1&query=TipoEnte.Id:"+strconv.Itoa(idEnte))
				if request.GetJson(beego.AppConfig.String("coreEnteService")+"ente?limit=-1&query=TipoEnte.Id:"+strconv.Itoa(idEnte), &ente);err==nil{
					beego.Info(ente)
				if ente!=nil {
					done := make(chan interface{})
					defer close(done)
					resch := optimize.GenChanInterface(ente...)
					chentes := optimize.Digest(done, getOrganizacion, resch, nil)

					for organizacion := range chentes {
						beego.Info(organizacion)
						//respuesta = append(respuesta, organizacion.(map[string]interface{}))
					}
					c.Data["json"] = respuesta
				}

			}else{
  			beego.Error(err);
			}

}else{
	beego.Error(err);
}
c.ServeJSON()
}

func getTipoOrganizacion(tipoEnte int) (f func(data map[string]interface{}, params ...interface{}) interface{}){
	switch tipoEnte {
		case 3:
			return getOrganizacionEnte
		default:
			return nil
	}
}

func getOrganizacionEnte(ente map[string]interface{}, params ...interface{}) interface{} {
	var idEnte int64
	var organizacion interface{}
	idEnte = int64(ente["Id"].(float64))
	beego.Error(beego.AppConfig.String("coreOrganizacionService")+"organizacion?limit=-1&query=Ente:"+strconv.FormatInt(idEnte,10))
	if err := request.GetJson(beego.AppConfig.String("coreOrganizacionService")+"organizacion?limit=-1&query=Ente:"+strconv.FormatInt(idEnte,10), &organizacion); err == nil {
	if organizacion != nil {
		return organizacion
	}else{
		return map[string]interface{}{"Code": "E_0458", "Body": "Core New Service", "Type": "error"}
	}
	}else{
		return map[string]interface{}{"Code": "E_0458", "Body": "Core Organizacion Service", "Type": "error"}
	}
}

func getOrganizacion(ente interface{}, params ...interface{}) (res interface{}) {
 enteMap := ente.(map[string]interface{})
	if funcion := getTipoOrganizacion(int(enteMap["TipoEnte"].(map[string]interface{})["Id"].(float64))); funcion != nil {
			res = funcion(enteMap, params)
		} else {
			beego.Error("Error")
			return enteMap
		}

	beego.Error("respuesta",res)
	return
}

// Put ...
// @Title Put
// @Description update the OrganizacionesCoreNew
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.OrganizacionesCoreNew	true		"body for OrganizacionesCoreNew content"
// @Success 200 {object} models.OrganizacionesCoreNew
// @Failure 403 :id is not int
// @router /:id [put]
func (c *OrganizacionesCoreNewController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the OrganizacionesCoreNew
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *OrganizacionesCoreNewController) Delete() {

}
