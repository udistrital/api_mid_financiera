package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/api_mid_financiera/models"
	"github.com/udistrital/api_mid_financiera/utilidades"
)

// OrdenPagoNominaController operations for Orden_pago_planta
type OrdenPagoNominaController struct {
	beego.Controller
}

// URLMapping ...
func (c *OrdenPagoNominaController) URLMapping() {
	c.Mapping("Post", c.Post)

}

// MidCrearOPNomina ...
// @Title Create
// @Description create Orden_pago_planta
// @Param	body		body 	models.Orden_pago_planta	true		"body for Orden_pago_planta content"
// @Success 201 {object} models.Orden_pago_planta
// @Failure 403 body is empty
// @router MidCrearOPNomina [post]
func (c *OrdenPagoNominaController) MidCrearOPNomina() {
	var alerta models.Alert
	var v interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		m := v.(map[string]interface{})
		var DetallePreliquidacion []interface{}

		OrdenPago, e := m["OrdenPago"].(map[string]interface{})
		if e != true {
			alerta.Type = "error"
			alerta.Code = "E_OPN_01_2"
			alerta.Body = err.Error()
			c.Data["json"] = alerta
			c.ServeJSON()
			return
		}
		Preliquidacion := OrdenPago["Liquidacion"].(float64)
		Usuario, e := m["Usuario"].(map[string]interface{})
		if e != true {
			alerta.Type = "error"
			alerta.Code = "E_OPN_01_2"
			alerta.Body = err.Error()
			c.Data["json"] = alerta
			c.ServeJSON()
			return
		}
		// get data titan
		if err := getJson("http://"+beego.AppConfig.String("titanService")+"detalle_preliquidacion?query=Preliquidacion:"+strconv.FormatFloat(Preliquidacion, 'f', 0, 64)+"&sortby=Concepto&order=desc", &DetallePreliquidacion); err == nil {
		} else {
			alerta.Type = "error"
			alerta.Code = "E_OPN_01_3"
			alerta.Body = err.Error()
			c.Data["json"] = alerta
			c.ServeJSON()
			return
		}
		// Control si no existe detalle de liquidacion
		if len(DetallePreliquidacion) == 0 {
			alerta.Type = "error"
			alerta.Code = "E_OPN_01_4"
			alerta.Body = ""
			c.Data["json"] = alerta
			c.ServeJSON()
			return
		}
		// estructura para enviar data a kronos
		type Send struct {
			OrdenPago          interface{}
			DetalleLiquidacion []interface{}
			Usuario            interface{}
		}
		sendData2Kronos := Send{OrdenPago: OrdenPago, DetalleLiquidacion: DetallePreliquidacion, Usuario: Usuario}
		var outputData interface{}
		//Envia data to kronos
		if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/orden_pago/RegistrarOpNomina", "POST", &outputData, &sendData2Kronos); err == nil {
		} else {
			alerta.Type = "error"
			alerta.Code = "E_OPN_01_5"
			alerta.Body = ""
			c.Data["json"] = alerta
			c.ServeJSON()
			return
		}
		c.Data["json"] = outputData
		c.ServeJSON()
	} else {
		alerta.Type = "error"
		alerta.Code = "E_OPN_01_1"
		alerta.Body = err.Error()
		c.Data["json"] = alerta
		c.ServeJSON()
		return
	}
}

// CrearOPSeguridadSocial ...
// @Title CrearOPSeguridadSocial
// @Description create Orden Pago Seguridad Social
// @Param	body		body 	models.Orden_pago_planta	true		"body for Orden_pago_planta content"
// @Success 201 {object} models.Orden_pago_planta
// @Failure 403 body is empty
// @router CrearOPSeguridadSocial [post]
func (c *OrdenPagoNominaController) CrearOPSeguridadSocial() {
	var alerta models.Alert
	var v interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		m := v.(map[string]interface{})
		var DataSeguridadSocial map[string]interface{}
		var DataOrdenPago map[string]interface{}
		var PagosSeguridadSocial []interface{}
		var PeriodoPago []interface{}
		//
		err = utilidades.FillStruct(m["SeguridadSocial"], &DataSeguridadSocial)
		Mes := fmt.Sprintf("%v", DataSeguridadSocial["Mes"])
		Anio := fmt.Sprintf("%v", DataSeguridadSocial["Vigencia"])
		err = utilidades.FillStruct(m["OrdenPago"], &DataOrdenPago)

		fmt.Print(Mes)
		fmt.Print("-")
		fmt.Print(Anio)
		// get id periodo pago
		//fmt.Println("\n", "http://"+beego.AppConfig.String("SsService")+"periodo_pago?query=Mes:"+Mes+",Anio:"+Anio)
		if err = getJson("http://"+beego.AppConfig.String("SsService")+"periodo_pago?query=Mes:"+Mes+",Anio:"+Anio, &PeriodoPago); err == nil {
		} else {
			alerta.Type = "error"
			alerta.Code = "E_OPN_01_3"
			alerta.Body = err.Error()
			c.Data["json"] = alerta
			c.ServeJSON()
			return
		}
		fmt.Println("\nAAAAAAAAAAA \nPeriodoPago")
		fmt.Print(PeriodoPago)
		fmt.Println("\nAAAAAAAAAAA")

		// // get data administarativa seguridad social
		// debe ser por mes y año el filtro, en el momento el api no cuenta con esos datos.
		if err = getJson("http://"+beego.AppConfig.String("SsService")+"pago?query=PeriodoPago.Id:1", &PagosSeguridadSocial); err == nil {
		} else {
			alerta.Type = "error"
			alerta.Code = "E_OPN_01_3"
			alerta.Body = err.Error()
			c.Data["json"] = alerta
			c.ServeJSON()
			return
		}

		fmt.Println("\nAAAAAAAAAAA \nPagosSeguridadSocial")
		fmt.Print(PagosSeguridadSocial)
		fmt.Println("\nAAAAAAAAAAA")

		// Control si no existe detalle de liquidacion
		if len(PagosSeguridadSocial) == 0 {
			alerta.Type = "error"
			alerta.Code = "E_OPN_01_4"
			alerta.Body = ""
			c.Data["json"] = alerta
			c.ServeJSON()
			return
		}
		// estructura para enviar data a kronos
		type SendData struct {
			OrdenPago            interface{}
			PagosSeguridadSocial []interface{}
		}
		SendData2Kronos := SendData{OrdenPago: DataOrdenPago, PagosSeguridadSocial: PagosSeguridadSocial}
		var outputData interface{}
		//Envia data to kronos
		if err = sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/orden_pago/RegistrarOpSeguridadSocial", "POST", &outputData, &SendData2Kronos); err == nil {
		} else {
			alerta.Type = "error"
			alerta.Code = "E_OPN_01_5"
			alerta.Body = ""
			c.Data["json"] = alerta
			c.ServeJSON()
			return
		}
		c.Data["json"] = outputData
		c.ServeJSON()
		//fin
	} else {
		alerta.Type = "error"
		alerta.Code = "E_OPN_01_1"
		alerta.Body = err.Error()
		c.Data["json"] = alerta
		c.ServeJSON()
		return
	}
}

func formatoListaLiquidacion(dataLiquidacion interface{}, params ...interface{}) (res interface{}) {
	row, e := dataLiquidacion.(map[string]interface{})
	var infoPersona interface{}
	if e {
		if err := getJsonWSO2("http://jbpm.udistritaloas.edu.co:8280/services/contrato_suscrito_DataService.HTTPEndpoint/informacion_contrato_elaborado_contratista/"+row["NumeroContrato"].(string)+"/"+strconv.Itoa(int(row["VigenciaContrato"].(float64))), &infoPersona); err == nil {
			row["infoPersona"], e = infoPersona.(map[string]interface{})["informacion_contratista"]
			fmt.Println(row["infoPersona"])
			if e {
				return row
			} else {
				fmt.Println("e")
				return
			}

		} else {
			return
		}
	} else {
		return
	}

	return
}

// ListaLiquidacionNominaHomologada ...
// @Title ListaLiquidacionNominaHomologada
// @Description lista liquidaciones para ordenes de pago masivas.
// @Param	idNomina	query	string	false	"nomina a listar"
// @Param	mesLiquidacion	query	string	false	"mes de la liquidacion a listar"
// @Param	anioLiquidacion	query	string	false	"anio de la liquidacion a listar"
// @Success 201 {object} models.Alert
// @Failure 403 body is empty
// @router /ListaLiquidacionNominaHomologada [get]
func (c *OrdenPagoNominaController) ListaLiquidacionNominaHomologada() {
	idNomina, err1 := c.GetInt("idNomina")
	mesLiquidacion, err2 := c.GetInt("mesLiquidacion")
	anioLiquidacion, err3 := c.GetInt("anioLiquidacion")
	if err1 == nil && err2 == nil && err3 == nil {
		var respuesta []map[string]interface{}
		var liquidacion interface{}
		if err := getJson("http://"+beego.AppConfig.String("titanService")+"preliquidacion/contratos_x_preliquidacion?idNomina="+strconv.Itoa(idNomina)+"&mesLiquidacion="+strconv.Itoa(mesLiquidacion)+"&anioLiquidacion="+strconv.Itoa(anioLiquidacion), &liquidacion); err == nil {
			if liquidacion != nil {
				done := make(chan interface{})
				defer close(done)
				if liquidacion.(map[string]interface{})["Contratos_por_preliq"] != nil {
					listaLiquidacion := liquidacion.(map[string]interface{})["Contratos_por_preliq"].([]interface{})
					resch := utilidades.GenChanInterface(listaLiquidacion...)
					chlistaLiquidacion := utilidades.Digest(done, formatoListaLiquidacion, resch, nil)
					for dataLiquidacion := range chlistaLiquidacion {
						if dataLiquidacion != nil {
							respuesta = append(respuesta, dataLiquidacion.(map[string]interface{}))
						}
					}
					res := liquidacion.(map[string]interface{})
					res["Contratos_por_preliq"] = respuesta
					c.Data["json"] = res
				} else {
					c.Data["json"] = models.Alert{Code: "E_0458", Body: nil, Type: "error"}
				}

			} else {
				c.Data["json"] = liquidacion
			}

		} else {
			c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
		}
	} else {
		c.Data["json"] = models.Alert{Code: "E_0458", Body: "Not enough parameter", Type: "error"}
	}
	c.ServeJSON()
}

// ListaConceptosNominaHomologados ...
// @Title ListaConceptosNominaHomologados
// @Description lista liquidaciones para ordenes de pago masivas.
// @Param	nContrato	query	string	false	"nomina a listar"
// @Param	vigenciaContrato	query	string	false	"mes de la liquidacion a listar"
// @Param	idLiquidacion	query	string	false	"anio de la liquidacion a listar"
// @Success 201 {object} models.Alert
// @Failure 403 body is empty
// @router /ListaConceptosNominaHomologados [get]
func (c *OrdenPagoNominaController) ListaConceptosNominaHomologados() {
	nContrato := c.GetString("nContrato")
	vigenciaContrato, err2 := c.GetInt("vigenciaContrato")
	idLiquidacion, err3 := c.GetInt("idLiquidacion")
	if nContrato != "" && err2 == nil && err3 == nil {
		var respuesta []map[string]interface{}
		var listaDetalles []interface{}
		if err := getJson("http://"+beego.AppConfig.String("titanService")+"detalle_preliquidacion?limit=-1&query=Preliquidacion.Id:"+strconv.Itoa(idLiquidacion)+",NumeroContrato:"+nContrato+",VigenciaContrato:"+strconv.Itoa(vigenciaContrato), &listaDetalles); err == nil {
			if listaDetalles != nil {
				done := make(chan interface{})
				defer close(done)
				resch := utilidades.GenChanInterface(listaDetalles...)
				chlistaDetalles := utilidades.Digest(done, homologacionConceptos, resch, nil)
				for dataLiquidacion := range chlistaDetalles {
					if dataLiquidacion != nil {
						existe := false
						for _, comp := range respuesta {
							fmt.Println(dataLiquidacion)
							if comp["Concepto"] != nil {
								if comp["Concepto"].(map[string]interface{})["Id"].(float64) == dataLiquidacion.(map[string]interface{})["Concepto"].(map[string]interface{})["Id"].(float64) {
									comp["Valor"] = comp["Valor"].(float64) + dataLiquidacion.(map[string]interface{})["Valor"].(float64)
									existe = true
								}
							}

						}
						if !existe {
							if dataLiquidacion.(map[string]interface{})["Concepto"] != nil {

								respuesta = append(respuesta, dataLiquidacion.(map[string]interface{}))
							}

						}
					}
				}
				c.Data["json"] = respuesta
			} else {
				c.Data["json"] = nil
			}

		} else {
			c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
		}
	} else {
		c.Data["json"] = models.Alert{Code: "E_0458", Body: "Not enough parameter", Type: "error"}
	}
	c.ServeJSON()
}

func homologacionConceptos(dataConcepto interface{}, params ...interface{}) (res interface{}) {
	dataConceptoAhomologar, e := dataConcepto.(map[string]interface{})
	out := make(map[string]interface{})
	if e {
		var homologacion []interface{}
		if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/homologacion_concepto?query=ConceptoTitan:"+strconv.Itoa(int(dataConceptoAhomologar["Concepto"].(map[string]interface{})["Id"].(float64))), &homologacion); err == nil {

			for _, conceptoKronos := range homologacion {
				row, e := conceptoKronos.(map[string]interface{})
				if e {
					out["Concepto"] = row["ConceptoKronos"]
					out["Valor"] = dataConceptoAhomologar["ValorCalculado"]
				} else {
					return nil
				}

			}
		} else {
			return nil
		}
		return out
	} else {
		return nil
	}
}

// CargueMasivoOp ...
// @Title CargueMasivoOp
// @Description lista liquidaciones para ordenes de pago masivas.
// @Param	idNomina	query	string	false	"nomina a listar"
// @Param	mesLiquidacion	query	string	false	"mes de la liquidacion a listar"
// @Param	anioLiquidacion	query	string	false	"anio de la liquidacion a listar"
// @Param   OrdenPago       map[string]string	true		"body for OrdenPago content"
// @Success 201 {object} models.Alert
// @Failure 403 body is empty
// @router /CargueMasivoOp [post]
func (c *OrdenPagoNominaController) CargueMasivoOp() {
	idNomina, err1 := c.GetInt("idNomina")
	mesLiquidacion, err2 := c.GetInt("mesLiquidacion")
	anioLiquidacion, err3 := c.GetInt("anioLiquidacion")
	if err1 == nil && err2 == nil && err3 == nil {
		var respuesta []interface{}
		var liquidacion interface{}
		infoOpGeneral := make(map[string]interface{})
		//leer json con info general de la op
		if err := json.Unmarshal(c.Ctx.Input.RequestBody, &infoOpGeneral); err == nil {
			if err := getJson("http://"+beego.AppConfig.String("titanService")+"preliquidacion/contratos_x_preliquidacion?idNomina="+strconv.Itoa(idNomina)+"&mesLiquidacion="+strconv.Itoa(mesLiquidacion)+"&anioLiquidacion="+strconv.Itoa(anioLiquidacion), &liquidacion); err == nil {
				if liquidacion != nil {
					done := make(chan interface{})
					defer close(done)
					if liquidacion.(map[string]interface{})["Contratos_por_preliq"] != nil {
						listaLiquidacion := liquidacion.(map[string]interface{})["Contratos_por_preliq"].([]interface{})
						resch := utilidades.GenChanInterface(listaLiquidacion...)
						var params []interface{}
						params = append(params, liquidacion.(map[string]interface{})["Id_Preliq"].(interface{}))
						chlistaLiquidacion := utilidades.Digest(done, formatoRegistroOp, resch, params)
						for dataLiquidacion := range chlistaLiquidacion {
							if dataLiquidacion != nil {
								respuesta = append(respuesta, dataLiquidacion)
							}
						}
						c.Data["json"] = respuesta
					} else {
						c.Data["json"] = models.Alert{Code: "E_0458", Body: nil, Type: "error"}
					}

				}
			} else {
				//error consumo de servicio titan. Lista contratos por liqu
				c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
			}
		} else {
			//error al recibir datos genrales de la op
			c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
		}
		fmt.Println(idNomina, mesLiquidacion, anioLiquidacion)
	} else {
		//no se enviaron los parametros necesarios
		c.Data["json"] = models.Alert{Code: "E_0458", Body: "Not enough parameter", Type: "error"}
	}
	c.ServeJSON()
}

func formatoRegistroOp(dataLiquidacion interface{}, params ...interface{}) (res interface{}) {
	idLiquidacion, e := params[0].(float64)
	if e {
		var respuesta []map[string]interface{}
		var listaDetalles []interface{}
		var valorTotal float64
		valorTotal = 0
		nContrato, e := dataLiquidacion.(map[string]interface{})["NumeroContrato"].(string)
		if !e {
			return nil
		}
		vigenciaContrato, e := dataLiquidacion.(map[string]interface{})["VigenciaContrato"].(float64)
		if !e {
			return nil
		}
		if err := getJson("http://"+beego.AppConfig.String("titanService")+"detalle_preliquidacion?limit=-1&query=Preliquidacion.Id:"+strconv.Itoa(int(idLiquidacion))+",NumeroContrato:"+nContrato+",VigenciaContrato:"+strconv.Itoa(int(vigenciaContrato)), &listaDetalles); err == nil {
			if listaDetalles != nil {
				done := make(chan interface{})
				defer close(done)
				resch := utilidades.GenChanInterface(listaDetalles...)
				chlistaDetalles := utilidades.Digest(done, homologacionConceptos, resch, nil)
				for dataLiquidacion := range chlistaDetalles {
					if dataLiquidacion != nil {
						existe := false

						for _, comp := range respuesta {

							if comp["Concepto"] != nil {
								if comp["Concepto"].(map[string]interface{})["Id"].(float64) == dataLiquidacion.(map[string]interface{})["Concepto"].(map[string]interface{})["Id"].(float64) {
									comp["Valor"] = comp["Valor"].(float64) + dataLiquidacion.(map[string]interface{})["Valor"].(float64)
									existe = true
									valorTotal = valorTotal + comp["Valor"].(float64)
								}
							}

						}
						if !existe {
							if dataLiquidacion.(map[string]interface{})["Concepto"] != nil {
								valorTotal = valorTotal + dataLiquidacion.(map[string]interface{})["Valor"].(float64)
								respuesta = append(respuesta, dataLiquidacion.(map[string]interface{}))
							}

						}
					}
				}
				res := make(map[string]interface{})
				res["ValorBase"] = valorTotal
				filtrorp := formatoListaLiquidacion(dataLiquidacion, nil)
				if filtrorp != nil {
					res["FiltroRp"], e = filtrorp.(map[string]interface{})["infoPersona"].(map[string]interface{})["Documento"].(map[string]interface{})["numero"]
				}
				res["Conceptos"] = respuesta
				return res
			} else {
				return nil
			}

		} else {
			return nil
		}
	} else {
		fmt.Println("err")
		return nil
	}
	return
}
