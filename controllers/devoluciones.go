package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"time"
	"strconv"
	"reflect"


	"github.com/udistrital/utils_oas/formatdata"

)

// DevolucionesController operations for Devoluciones
type DevolucionesController struct {
	beego.Controller
}

type PagosAcademica struct{
	informacionEstudiante infoEstudiante
	informacionCarrera    []*infoCarrera
}

type infoRecibo struct {
	Total        int64
	Numero_Recibo int64
	FechaExtraordinario time.Time
	FechaOrdinario time.Time
	Periodo			string
	Pago				string
	DesagregaRecibos []*infoPago
}

type infoEstudiante struct {
	Tipo_Docu  	string
	Documento   string
	Tipo				string
	Nombre      string
}

type infoCarrera struct {
	Carrera			string
	Facultad		string
	Cod_Carrera string
	Codigo			string
	informacionRecibos  []*infoRecibo
}

type infoPago struct {
	descripcion string
	valor				float64
}


// URLMapping ...
func (c *DevolucionesController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create Devoluciones
// @Param	body		body 	models.Devoluciones	true		"body for Devoluciones content"
// @Success 201 {object} models.Devoluciones
// @Failure 403 body is empty
// @router / [post]
func (c *DevolucionesController) Post() {

}

// GetOne ...
// @Title GetOne
// @Description get Devoluciones by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Devoluciones
// @Failure 403 :id is empty
// @router /:id [get]
func (c *DevolucionesController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get Devoluciones
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Devoluciones
// @Failure 403
// @router / [get]
func (c *DevolucionesController) GetAll() {

}

// Put ...
// @Title Put
// @Description update the Devoluciones
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Devoluciones	true		"body for Devoluciones content"
// @Success 200 {object} models.Devoluciones
// @Failure 403 :id is not int
// @router /:id [put]
func (c *DevolucionesController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Devoluciones
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *DevolucionesController) Delete() {

}
// GetTransformRequest...
// @Title GetTransformRequest
// @Description obtiene json de transformacion de respuesta de servicio de academica
// @Param	query	query	string	false	"Objecto con valores de pagos en academica"
// @Success 200 {object} interface{}
// @Failure 403
// @router /GetTransformRequest/ [post]
func (c *DevolucionesController) GetTransformRequest() {
	defer c.ServeJSON()

	var ingresoData map[string]interface{}
	var ingresoData2 []interface{}
	var estudiante infoEstudiante
	var carrera infoCarrera
	var carreras []*infoCarrera
	var informacionPago infoPago
	var totalRecibo float64
	var informacionRecibo infoRecibo


	var data map[string]interface{}


	if err:= json.Unmarshal(c.Ctx.Input.RequestBody,&ingresoData);err == nil {


		varIntrinsecas:= []string{"fecha_ordinario","pago","cod_facultad","periodo","facultad",
															"tipo_docu","fecha_extraordinario","carrera","cod_carrera","numero_recibo","documento","tipo","nombre","codigo"}


		ingresoData2 = ingresoData["pagos"].([]interface{})

		_ = formatdata.FillStruct(ingresoData2[0], &carrera)

		if err = formatdata.FillStruct(ingresoData2[0], &estudiante);err!=nil{
			beego.Error(err)
		}


		for  _, value := range ingresoData2 {
				err = formatdata.FillStruct(value, &carrera)
				err = formatdata.FillStruct(value, &data)


				if len(carreras) > 0{
					for _,car:= range carreras {
						if !reflect.DeepEqual(car, &carrera) {
								carreras = append(carreras,&carrera)
							 	beego.Info("agrega carrera ", carrera)
						}
					}
				}else {
					carreras = append(carreras,&carrera)
					beego.Info("agrega carrera ",carrera)

					for key , _ := range data {
						if !contains(varIntrinsecas,key) && data[key] != nil {
							informacionPago.descripcion = key
							if informacionPago.valor,err = strconv.ParseFloat(data[key].(string),64);err != nil {
								beego.Error(err)
							}else{
								totalRecibo = totalRecibo + informacionPago.valor
							}
							err = formatdata.FillStruct(value, &informacionRecibo)
							beego.Info(informacionPago.valor)
						}
					}


				}



		}


	}else{
		beego.Error(err)
	}
}


func contains(array []string, varBusqueda string) bool {
    for _, a := range array {
        if a == varBusqueda {
            return true
        }
    }
    return false
}
