package controllers

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"

	"github.com/udistrital/api_mid_financiera/models"
	"github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/utils_oas/optimize"
	"github.com/udistrital/utils_oas/request"
)

// DevolucionesController operations for Devoluciones
type DevolucionesController struct {
	beego.Controller
}

type pagosAcademica struct {
	InformacionEstudiante infoEstudiante
	InformacionCarrera    []*infoCarrera
}

type infoRecibo struct {
	Total                float64
	Numero_Recibo        string
	Fecha_Extraordinario time.Time
	Fecha_Ordinario      time.Time
	Periodo              string
	Pago                 string
	DesagregaRecibos     []*infoPago
}

type infoEstudiante struct {
	Tipo_Docu string
	Documento string
	Tipo      string
	Nombre    string
}

type infoCarrera struct {
	Carrera            string
	Facultad           string
	Cod_Carrera        string
	Codigo             string
	InformacionRecibos []*infoRecibo
}

type carreraIntrinc struct {
	Carrera     string
	Facultad    string
	Cod_Carrera string
	Codigo      string
}

type infoPago struct {
	Descripcion string
	Valor       float64
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

// GetTransformRequest ...
// @Title GetTransformRequest
// @Description obtiene json de transformacion de respuesta de servicio de academica
// @Param	query	query	string	false	"Objecto con valores de pagos en academica"
// @Success 200 {object} interface{}
// @Failure 403
// @router /GetTransformRequest/ [post]
func (c *DevolucionesController) GetTransformRequest() {
	defer c.ServeJSON()
	var pagos pagosAcademica
	var ingresoData map[string]interface{}
	var ingresoData2 []interface{}
	//var estudiante infoEstudiante
	//var carreras []*infoCarrera

	var data map[string]interface{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &ingresoData); err == nil {

		ingresoData2 = ingresoData["pagos"].([]interface{})

		//_ = formatdata.FillStruct(ingresoData2[0], &carrera)

		if err = formatdata.FillStruct(ingresoData2[0], &pagos.InformacionEstudiante); err != nil {
			beego.Error(err)
			c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
		}

		for _, value := range ingresoData2 {
			carrera := new(infoCarrera)
			err = formatdata.FillStruct(value, &carrera)
			err = formatdata.FillStruct(value, &data)
			if len(pagos.InformacionCarrera) > 0 {
				agregar, carreraReturn := searchCarrera(pagos.InformacionCarrera, carrera)
				if !agregar {
					pagos.InformacionCarrera = append(pagos.InformacionCarrera, carrera)

				} else {
					informacionRecibo := new(infoRecibo)
					informacionRecibo = getPayInfo(data)
					beego.Error("informacion Recibo", informacionRecibo)
					carreraReturn.InformacionRecibos = append(carreraReturn.InformacionRecibos, informacionRecibo)
				}
			} else {
				informacionRecibo := new(infoRecibo)
				informacionRecibo = getPayInfo(data)
				beego.Error("informacion Recibo", informacionRecibo)
				carrera.InformacionRecibos = append(carrera.InformacionRecibos, informacionRecibo)
				pagos.InformacionCarrera = append(pagos.InformacionCarrera, carrera)
				beego.Info("agrega carrera 2 ", carrera)
			}
		}
		//pagos.informacionCarrera = carreras
		beego.Info("Pagos", pagos)
		c.Data["json"] = pagos
	} else {
		beego.Error(err)
		c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
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

func searchCarrera(array []*infoCarrera, varBusqueda *infoCarrera) (encontrado bool, carreraReturn *infoCarrera) {

	for _, a := range array {

		if reflect.DeepEqual(a.Codigo, varBusqueda.Codigo) {
			carreraReturn = a
			encontrado = true
			return
		}
	}
	encontrado = false
	return
}

func getPayInfo(data map[string]interface{}) (informacionRecibo *infoRecibo) {
	var totalRecibo float64
	var err error
	err = formatdata.FillStruct(data, &informacionRecibo)

	varIntrinsecas := []string{"fecha_ordinario", "pago", "cod_facultad", "periodo", "facultad",
		"tipo_docu", "fecha_extraordinario", "carrera", "cod_carrera", "numero_recibo", "documento", "tipo", "nombre", "codigo"}

	for key, _ := range data {
		if !contains(varIntrinsecas, key) && data[key] != nil {
			informacionPago := new(infoPago)
			informacionPago.Descripcion = key
			if informacionPago.Valor, err = strconv.ParseFloat(data[key].(string), 64); err != nil {
				beego.Error(err)
			} else {
				totalRecibo = totalRecibo + informacionPago.Valor
			}
			informacionRecibo.DesagregaRecibos = append(informacionRecibo.DesagregaRecibos, informacionPago)
		}
	}
	informacionRecibo.Total = totalRecibo
	for _, valor := range informacionRecibo.DesagregaRecibos {
		beego.Error(valor)
	}
	return
}

// GetAllDevolucionesTributarias ...
// @Title GetAllDevolucionesTributarias
// @Description get all devoluciones tributarias
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Devoluciones
// @Failure 403
// @router /GetAllDevolucionesTributarias [get]
func (c *DevolucionesController) GetAllDevolucionesTributarias() {
	defer c.ServeJSON()
	var devoluciones []interface{}
	var limit int64 = 10
	var offset int64
	var query string
	var regCuantity map[string]interface{}
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
	respuesta := make(map[string]interface{})
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/devolucion_tributaria/?limit="+strconv.FormatInt(limit, 10)+"&offset="+strconv.FormatInt(offset, 10)+"&query="+query, &devoluciones); err == nil {
		if devoluciones != nil {
			respuesta["Devolutions"] = optimize.ProccDigest(devoluciones, getValuesDevolTributaria, nil, 3)
			if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/devolucion_tributaria/GetDevolucionRecordsNumber/?query="+query, &regCuantity); err == nil {
				if strings.Compare(regCuantity["Type"].(string), "success") == 0 {
					respuesta["RegCuantity"] = regCuantity["Body"]
					c.Ctx.Output.SetStatus(201)
				}
			}
			c.Data["json"] = respuesta
		}
	} else {
		beego.Error("Error ", err)
		c.Data["json"] = models.Alert{Type: "error", Code: "E_0458", Body: err}
	}
}

func getValuesDevolTributaria(rpintfc interface{}, params ...interface{}) (res interface{}) {
	var resEstado []map[string]interface{}
	var resSolicitante []map[string]interface{}
	var resValorDevol []map[string]interface{}
	var resActa []map[string]interface{}
	var resDocumentoGen []map[string]interface{}
	var tipoCuenta map[string]interface{}
	var valorDevol float64
	devolID := strconv.FormatFloat(rpintfc.(map[string]interface{})["Id"].(float64), 'f', -1, 64)
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/devolucion_tributaria_estado_devolucion/?query=Activo:true"+",Devolucion.Id:"+devolID, &resEstado); err == nil {
		if resEstado[0] != nil {
			rpintfc.(map[string]interface{})["Estado"] = resEstado[0]["EstadoDevolucion"]
		}
	} else {
		beego.Error("Error", err.Error())
	}
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/devolucion_tributaria_concepto/?query=DevolucionTributaria.Id:"+devolID+"&fields=ValorDevolucion&limit=1", &resValorDevol); err == nil {
		if resValorDevol != nil {
			for _, v := range resValorDevol {
				valorDevol = valorDevol + v["ValorDevolucion"].(float64)
			}
			rpintfc.(map[string]interface{})["ValorDevolucion"] = valorDevol
		}
	} else {
		beego.Error("Error", err.Error())
	}
	tipoCuentaID := strconv.FormatFloat(rpintfc.(map[string]interface{})["CuentaBancariaEnte"].(map[string]interface{})["TipoCuenta"].(float64), 'f', -1, 64)
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/tipo_cuenta_bancaria/"+tipoCuentaID, &tipoCuenta); err == nil {
		if resValorDevol != nil {
			rpintfc.(map[string]interface{})["CuentaBancariaEnte"].(map[string]interface{})["TipoCuenta"] = tipoCuenta
		}
	} else {
		beego.Error("Error", err.Error())
	}

	actaID := strconv.FormatFloat(rpintfc.(map[string]interface{})["Acta"].(float64), 'f', -1, 64)
	if err := request.GetJson("http://"+beego.AppConfig.String("coreService")+"documento/?query=Id:"+actaID+"&limit=1", &resActa); err == nil {
		if resActa != nil {
			rpintfc.(map[string]interface{})["Acta"] = resActa[0]
		}
	} else {
		beego.Error("Error", err.Error())
	}
	documentoGenerador := rpintfc.(map[string]interface{})["DocumentoGenerador"]
	documentoID := strconv.FormatFloat(documentoGenerador.(map[string]interface{})["TipoDocumento"].(float64), 'f', -1, 64)
	if err := request.GetJson("http://"+beego.AppConfig.String("coreService")+"documento/?query=Id:"+documentoID+"&limit=1", &resDocumentoGen); err == nil {
		if resDocumentoGen != nil {
			rpintfc.(map[string]interface{})["DocumentoGenerador"].(map[string]interface{})["TipoDocumento"] = resDocumentoGen[0]
		}
	} else {
		beego.Error("Error", err.Error())
	}
	solicitante := strconv.FormatFloat(rpintfc.(map[string]interface{})["Solicitante"].(float64), 'f', -1, 64)
	if err := request.GetJson("http://"+beego.AppConfig.String("AdministrativaAmazonService")+"/informacion_proveedor/?query=NumDocumento:"+solicitante+"&limit=1", &resSolicitante); err == nil {
		if resSolicitante != nil {
			rpintfc.(map[string]interface{})["Solicitante"] = resSolicitante[0]
		}
	} else {
		beego.Error("Error", err.Error())
	}
	return rpintfc
}

// GetTributaDevolutionAccountantInf ...
// @Title GetTributaDevolutionAccountantInf
// @Description get accountant information related to devolution
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Devoluciones
// @Failure 403 :id is empty
// @router /GetTributaDevolutionAccountantInf/:id [get]
func (c *DevolucionesController) GetTributaDevolutionAccountantInf() {
	defer c.ServeJSON()

	idStr := c.Ctx.Input.Param(":id")
	var respMovimientosDevolucion []interface{}
	var params []interface{}
	params = append(params, idStr)
	beego.Error(params)
	var respConceptos []interface{}
	respuesta := make(map[string]interface{})
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/devolucion_tributaria_movimiento/?query=devolucion.id:"+idStr, &respMovimientosDevolucion); err == nil {
		if respMovimientosDevolucion != nil {
			respuesta["MovimientosAsociados"] = optimize.ProccDigest(respMovimientosDevolucion, getValuesMovimientosDevolucion, nil, 3)
		}
		if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/devolucion_tributaria_concepto/?query=DevolucionTributaria.Id:"+idStr+"&fields=ValorDevolucion,Id,Concepto", &respConceptos); err == nil {
			if respConceptos != nil {
				respuesta["Conceptos"] = optimize.ProccDigest(respConceptos, getConceptosDevolucion, params, 3)
			}
		} else {
			beego.Error("Error", err.Error())
			c.Data["json"] = models.Alert{Type: "error", Code: "E_0458", Body: err}
		}
		c.Data["json"] = respuesta
		c.Ctx.Output.SetStatus(201)
	} else {
		beego.Error("Error ", err)
		c.Data["json"] = models.Alert{Type: "error", Code: "E_0458", Body: err}
	}
}

func getValuesMovimientosDevolucion(rpintfc interface{}, params ...interface{}) (res interface{}) {
	var resValorBase []map[string]interface{}
	var resOp map[string]interface{}
	cuentaEspecialID := strconv.FormatFloat(rpintfc.(map[string]interface{})["MovimientoContable"].(map[string]interface{})["CuentaEspecial"].(map[string]interface{})["Id"].(float64), 'f', -1, 64)
	OrdenPagoID := strconv.FormatFloat(rpintfc.(map[string]interface{})["MovimientoContable"].(map[string]interface{})["CodigoDocumentoAfectante"].(float64), 'f', -1, 64)
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/orden_pago_cuenta_especial/?query=CuentaEspecial.Id:"+cuentaEspecialID+",OrdenPago.Id:"+OrdenPagoID+"&fields=ValorBase", &resValorBase); err == nil {
		if resValorBase != nil {
			rpintfc.(map[string]interface{})["ValorBase"] = resValorBase[0]["ValorBase"]
		}
	} else {
		beego.Error("Error", err.Error())
	}

	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/orden_pago/"+OrdenPagoID, &resOp); err == nil {
		if resOp != nil {
			rpintfc.(map[string]interface{})["OrdenPago"] = resOp
		}
	} else {
		beego.Error("Error", err.Error())
	}
	return rpintfc
}

func getConceptosDevolucion(rpintfc interface{}, params ...interface{}) (res interface{}) {
	var resMovimientoContable []map[string]interface{}
	rpintfcCp := rpintfc.(map[string]interface{})
	concepto := rpintfcCp["Concepto"].(map[string]interface{})
	devolucionID := params[0].(string)
	delete(rpintfcCp, "Id")
	for k, v := range concepto {
		rpintfcCp[k] = v
	}
	delete(rpintfcCp, "Concepto")

	conceptoID := strconv.FormatFloat(rpintfcCp["Id"].(float64), 'f', -1, 64)
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/movimiento_contable/?query=Concepto.Id:"+conceptoID+",CodigoDocumentoAfectante:"+devolucionID+",TipoDocumentoAfectante.NumeroOrden:6&fields=Id,Credito,Debito,CuentaContable", &resMovimientoContable); err == nil {
		if resMovimientoContable != nil {
			rpintfcCp["MovimientoContable"] = resMovimientoContable
		}
	} else {
		beego.Error("Error", err.Error())
	}

	return rpintfcCp
}
