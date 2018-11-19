package controllers

import (
	"strconv"
	"encoding/json"

	"time"
	"github.com/astaxie/beego"
		"github.com/udistrital/api_mid_financiera/models"
	"github.com/udistrital/utils_oas/optimize"
	"github.com/udistrital/utils_oas/request"
	"fmt"
)

// OrdenPagoController operations for Orden_pago
type OrdenPagoController struct {
	beego.Controller
}

// URLMapping ...
func (c *OrdenPagoController) URLMapping() {
	c.Mapping("GetOrdenPagoByFuenteFinanciamiento", c.GetOrdenPagoByFuenteFinanciamiento)
	c.Mapping("AnularOrdenPago", c.AnularOrdenPago)
}


// AnularOrdenPago ...
// @Title AnularOrdenPago
// @Description Cambia el estado y registra el histórico de la Orden de Pago
// @Param	body		body 	var v models.OrdenPago	true		"body for OrdenPago content"
// @Success 201 {object}  models.Alert
// @Failure 403 body is empty
// @router /anulacion_orden_pago/ [post]
func (c *OrdenPagoController) AnularOrdenPago() {

	var v models.OrdenPago
	var res map[string]interface{}
	var respuestaMovimientos models.Alert
	if errUnmarshal := json.Unmarshal(c.Ctx.Input.RequestBody, &v); errUnmarshal == nil {
		var NuevoOPEOP models.OrdenPagoEstadoOrdenPago
		NuevoOPEOP.OrdenPago  = &models.OrdenPago {Id: v.Id}
		NuevoOPEOP.FechaRegistro = time.Now()
		NuevoOPEOP.EstadoOrdenPago = &models.EstadoOrdenPago {Id: 11}
		NuevoOPEOP.Usuario = 1;
		urlcrud := "http://" + beego.AppConfig.String("Urlcrud") + ":" + beego.AppConfig.String("Portcrud") + "/" + beego.AppConfig.String("Nscrud") + "/orden_pago_estado_orden_pago/"

		respuestaMovimientos = anularMovimientosContables(v.Id)

		if respuestaMovimientos.Type == "success" {
			if errPost := request.SendJson(urlcrud, "POST", &res, &NuevoOPEOP); errPost == nil {
				c.Data["json"] = map[string]interface{}{"Code": "E_2", "Body": res, "Type": "success"}

			}else{
					fmt.Println(errPost)
					c.Data["json"] = map[string]interface{}{"Code": "E_2", "Body": errPost, "Type": "error"}
			}

		}else{
			c.Data["json"] = respuestaMovimientos
		}





	}else{
		c.Data["json"] = map[string]interface{}{"Code": "E_3", "Body": errUnmarshal, "Type": "error"}

	}

		c.ServeJSON()
}

// GetOrdenPagoByFuenteFinanciamiento ...
// @Title GetOrdenPagoByFuenteFinanciamiento
// @Description lista Ordenes de Pago por fuente de financiamineto.
// @Param	fuente	query	string	true	"Id Fuente de Financiamiento"
// @Param	vigencia	query	string	true	"Vigencia de registro Presupuestal"
// @Param	unidadEjecutora	query	string	true	"Unidad Ejecutora"
// @Success 201 {object} models.Alert
// @Failure 403 body is empty
// @router /GetOrdenPagoByFuenteFinanciamiento [get]
func (c *OrdenPagoController) GetOrdenPagoByFuenteFinanciamiento() {
	fuente := c.GetString("fuente")
	vigencia := c.GetString("vigencia")
	unidadEjecutora := c.GetString("unidadEjecutora")
	if fuente != "" && vigencia != "" && unidadEjecutora != "" {
		var registroPresupuestales []interface{}
		respuestas := make(map[string]interface{})
		var parametro []interface{}
		parametro = append(parametro, unidadEjecutora)
		// search registro_presupuestal
		beego.Info("http://"+beego.AppConfig.String("kronosService")+"registro_presupuestal/?query=RegistroPresupuestalDisponibilidadApropiacion.DisponibilidadApropiacion.FuenteFinanciamiento.Id:"+fuente+",Vigencia:"+vigencia+"&limit:-1")
		if err := request.GetJson("http://"+beego.AppConfig.String("kronosService")+"registro_presupuestal/?query=RegistroPresupuestalDisponibilidadApropiacion.DisponibilidadApropiacion.FuenteFinanciamiento.Id:"+fuente+",Vigencia:"+vigencia+"&limit:-1", &registroPresupuestales); err == nil && registroPresupuestales != nil {
			var outputOrdenPago []interface{}
			done := make(chan interface{})
			defer close(done)
			resch := optimize.GenChanInterface(registroPresupuestales...)
			chlistaOrdenes := optimize.Digest(done, searchOrdenPagoByRpID, resch, parametro)
			for arrayOrdenPago := range chlistaOrdenes {
				if dataOrden, e := arrayOrdenPago.([]interface{}); e {
					for _, ordenPago := range dataOrden {
						outputOrdenPago = append(outputOrdenPago, ordenPago.(interface{}))
					}
				}
			}
			if outputOrdenPago != nil {
				respuestas["OrdenPago"] = outputOrdenPago
				c.Data["json"] = respuestas
			} else {
				c.Data["json"] = nil
			}
		} else {
			c.Data["json"] = nil
		}
	} else {
		c.Ctx.Output.SetStatus(201)
		c.Data["json"] = models.Alert{Code: "E_0458", Body: "Not enough parameter in GetOrdenPagoByFuenteFinanciamiento", Type: "error"}
	}
	c.ServeJSON()

}

func searchOrdenPagoByRpID(inputRegistroPresupuestal interface{}, params ...interface{}) (res interface{}) {
	unidadEjecutoraID, e1 := params[0].(string)
	rowRp, e2 := inputRegistroPresupuestal.(map[string]interface{})
	// beego.Info("unidadEjecutoraID: ",unidadEjecutoraID)
	if e1 && e2 {
		var ordenesPagos []interface{}
		// seach dependencia de necesidad
		if necesidad := getNecesidadDesdeRp(inputRegistroPresupuestal, unidadEjecutoraID); necesidad != nil {
			// beego.Info("necesidad: ",necesidad)
			// beego.Info("necesidad?")
			if areaNecesidad := getAreaDeNecesidad(necesidad); areaNecesidad != nil {
				//op
				if err := request.GetJson("http://"+beego.AppConfig.String("kronosService")+"orden_pago/?query=RegistroPresupuestal.Id:"+strconv.Itoa(int(rowRp["Id"].(float64))), &ordenesPagos); err == nil && ordenesPagos != nil {
					for _, orden := range ordenesPagos {
						row := orden.(map[string]interface{})
						row["RegistroPresupuestal"] = rowRp
						row["Necesidad"] = areaNecesidad
					}
					return ordenesPagos
				}
			}
		} else {
			//beego.Info("No hay Necesidad desde Rp")

		}
	}
	return
}


func anularMovimientosContables(iDOp int) (res models.Alert){
	var movimientoContables []models.MovimientoContable
	var retorno models.Alert
	var cuentaEsp *models.CuentaEspecial
	var respuesta interface{}
	if errMov := request.GetJson("http://" + beego.AppConfig.String("Urlcrud") + ":" + beego.AppConfig.String("Portcrud") + "/" + beego.AppConfig.String("Nscrud") + "/movimiento_contable?limit=-1&query=CodigoDocumentoAfectante:"+strconv.Itoa(iDOp), &movimientoContables); errMov == nil && movimientoContables != nil {
		for x, dato := range movimientoContables {

			cuentaEsp = nil;

			if (dato.Debito == 0){

				movimientoContables[x].Debito = movimientoContables[x].Credito
				movimientoContables[x].Credito = 0;

			}
			if (dato.Credito == 0){

				movimientoContables[x].Credito = movimientoContables[x].Debito
				movimientoContables[x].Debito = 0
			}

			if  (movimientoContables[x].CuentaEspecial != nil) {
					cuentaEsp = &models.CuentaEspecial {Id:movimientoContables[x].CuentaEspecial.Id }
			}




			nuevoObjetoMovimiento:= &models.MovimientoContable {
			  Debito:  movimientoContables[x].Debito      ,
			  Credito:  movimientoContables[x].Credito     ,
			  Fecha:    time.Now()   ,
			  Concepto: &models.Concepto {Id:  movimientoContables[x].Concepto.Id},
			  CuentaContable : &models.CuentaContable {Id: movimientoContables[x].CuentaContable.Id },
			  TipoDocumentoAfectante: &models.TipoDocumentoAfectante {Id:  movimientoContables[x].TipoDocumentoAfectante.Id },
			  CodigoDocumentoAfectante: movimientoContables[x].CodigoDocumentoAfectante,
			  EstadoMovimientoContable : &models.EstadoMovimientoContable {Id: movimientoContables[x].EstadoMovimientoContable.Id },
			  CuentaEspecial : cuentaEsp,
			}

			Urlcrud := "http://" + beego.AppConfig.String("Urlcrud") + ":" + beego.AppConfig.String("Portcrud") + "/" + beego.AppConfig.String("Nscrud") + "/movimiento_contable"
			if errNuevosMovi := request.SendJson(Urlcrud, "POST", &respuesta, nuevoObjetoMovimiento); errNuevosMovi == nil && respuesta != nil {
					retorno = models.Alert{Code: "E_", Body: res, Type: "success"}
					
			}else{
					retorno = models.Alert{Code: "E_", Body: errNuevosMovi, Type: "error"}
				}
		}




	}else{
		fmt.Println(errMov)
		retorno = models.Alert{Code: "E_", Body: errMov, Type: "error"}
	}

	return retorno
}
