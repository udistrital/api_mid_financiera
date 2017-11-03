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
			fmt.Println(infoPersona)
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
		//fmt.Println("http://" + beego.AppConfig.String("titanService") + "preliquidacion/contratos_x_preliquidacion?idNomina=" + strconv.Itoa(idNomina) + "&mesLiquidacion=" + strconv.Itoa(mesLiquidacion) + "&anioLiquidacion=" + strconv.Itoa(anioLiquidacion))
		if err := getJson("http://"+beego.AppConfig.String("titanService")+"preliquidacion/contratos_x_preliquidacion?idNomina="+strconv.Itoa(idNomina)+"&mesLiquidacion="+strconv.Itoa(mesLiquidacion)+"&anioLiquidacion="+strconv.Itoa(anioLiquidacion), &liquidacion); err == nil {
			if liquidacion != nil {
				fmt.Println(liquidacion)
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
	vigenciaContrato, err2 := c.GetFloat("vigenciaContrato")
	idLiquidacion, err3 := c.GetInt("idLiquidacion")
	if nContrato != "" && err2 == nil && err3 == nil {
		var respuesta []map[string]interface{}
		var listaDetalles []interface{}
		if err := getJson("http://"+beego.AppConfig.String("titanService")+"detalle_preliquidacion?limit=-1&query=Preliquidacion.Id:"+strconv.Itoa(idLiquidacion)+",NumeroContrato:"+nContrato+",VigenciaContrato:"+strconv.Itoa(int(vigenciaContrato)), &listaDetalles); err == nil {
			if listaDetalles != nil {
				done := make(chan interface{})
				defer close(done)
				resch := utilidades.GenChanInterface(listaDetalles...)
				f := homologacionFunctionDispatcher(listaDetalles[0].(map[string]interface{})["Preliquidacion"].(map[string]interface{})["Nomina"].(map[string]interface{})["TipoNomina"].(map[string]interface{})["Nombre"].(string))
				var params []interface{}
				params = append(params, "persona")
				params = append(params, nContrato)
				params = append(params, vigenciaContrato)
				chConcHomologados := utilidades.Digest(done, f, resch, params)
				for conceptoHomologadoint := range chConcHomologados {
					conceptoHomologado, e := conceptoHomologadoint.(map[string]interface{})
					if e {
						existe := false
						for _, comp := range respuesta {
							if comp["Concepto"] != nil && conceptoHomologado["Concepto"] != nil {
								if comp["Concepto"].(map[string]interface{})["Id"].(float64) == conceptoHomologado["Concepto"].(map[string]interface{})["Id"].(float64) {
									comp["Valor"] = comp["Valor"].(float64) + conceptoHomologado["Valor"].(float64)
									existe = true
									//valorTotal = valorTotal + comp["Valor"].(float64)
								}
							}

						}
						if !existe {
							if conceptoHomologado["Concepto"] != nil {
								//valorTotal = valorTotal + conceptoHomologado["Valor"].(float64)
								movcont := formatoMovimientosContablesOp(conceptoHomologado)
								conceptoHomologado["MovimientoContable"] = movcont
								respuesta = append(respuesta, conceptoHomologado)
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

func homologacionConceptosHC(dataConcepto interface{}, params ...interface{}) (res interface{}) {
	dataConceptoAhomologar, e := dataConcepto.(map[string]interface{})
	//fmt.Println(dataConcepto)
	if !e {
		fmt.Println("e1")
		return nil
	}
	//fmt.Println(dataConceptoAhomologar)
	out := make(map[string]interface{})
	var numContrato string
	var vigContrato float64
	//var idPersona float64
	if params != nil {
		idPersona, e := params[0].(string)
		fmt.Println(idPersona)
		if !e {
			fmt.Println("e1")
			return nil
		}
		numContrato, e = params[1].(string)
		//fmt.Println(numContrato)
		if !e {
			fmt.Println("e2")
			return nil
		}
		vigContrato, e = params[2].(float64)
		//fmt.Println(vigContrato)
		if !e {
			fmt.Println("e3")
			return nil
		}
		var homologacion []interface{}
		//aqui va la consulta sobre facultad y proyecto para HC (modificar para hacerla de forma genral)
		var infoVinculacion []interface{}
		//fmt.Println("http://" + beego.AppConfig.String("argoService") + "vinculacion_docente?query=NumeroContrato:" + numContrato + ",Vigencia:" + strconv.FormatFloat(vigContrato, 'f', -1, 64))
		if err := getJson("http://"+beego.AppConfig.String("argoService")+"vinculacion_docente?query=NumeroContrato:"+numContrato+",Vigencia:"+strconv.FormatFloat(vigContrato, 'f', -1, 64), &infoVinculacion); err == nil {
			if infoVinculacion != nil {
				//fmt.Println("Facultad: ", infoVinculacion[0].(map[string]interface{})["IdResolucion"].(map[string]interface{})["IdFacultad"], "Proyecto: ", infoVinculacion[0].(map[string]interface{})["IdProyectoCurricular"])
				idFacultad, e := infoVinculacion[0].(map[string]interface{})["IdResolucion"].(map[string]interface{})["IdFacultad"].(float64)
				if !e {
					fmt.Println("err idres")
					return nil
				}
				idProyecto, e := infoVinculacion[0].(map[string]interface{})["IdProyectoCurricular"].(float64)
				if !e {
					fmt.Println("err idPro")
					return nil
				}
				//fmt.Println("http://" + beego.AppConfig.String("Urlcrud") + ":" + beego.AppConfig.String("Portcrud") + "/" + beego.AppConfig.String("Nscrud") + "/homologacion_concepto?query=ConceptoTitan:" + strconv.Itoa(int(dataConceptoAhomologar["Concepto"].(map[string]interface{})["Id"].(float64))) + ",ConceptoKronos.ConceptoTesoralFacultadProyecto.Facultad:" + strconv.Itoa(int(idFacultad)) + ",ConceptoKronos.ConceptoTesoralFacultadProyecto.ProyectoCurricular:" + strconv.Itoa(int(idProyecto)))
				if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/homologacion_concepto?query=ConceptoTitan:"+strconv.Itoa(int(dataConceptoAhomologar["Concepto"].(map[string]interface{})["Id"].(float64)))+",ConceptoKronos.ConceptoTesoralFacultadProyecto.Facultad:"+strconv.Itoa(int(idFacultad))+",ConceptoKronos.ConceptoTesoralFacultadProyecto.ProyectoCurricular:"+strconv.Itoa(int(idProyecto)), &homologacion); err == nil {
					//fmt.Println("Hom ", homologacion)
					for _, conceptoKronos := range homologacion {
						row, e := conceptoKronos.(map[string]interface{})
						//fmt.Println(row)
						if e {
							out["Concepto"] = row["ConceptoKronos"]
							out["Valor"] = dataConceptoAhomologar["ValorCalculado"]
						} else {
							fmt.Println("err  concKron")
							return nil
						}

					}
				} else {
					fmt.Println(err.Error())
					return nil
				}
			} else {
				fmt.Println("no vinculacion data")
				return nil
			}

		} else {
			fmt.Println(err.Error())
			return nil
		}

		return out
	} else {
		fmt.Println("e2")
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
					_, e := liquidacion.(map[string]interface{})
					if e {
						if liquidacion.(map[string]interface{})["Contratos_por_preliq"] != nil {
							listaLiquidacion := liquidacion.(map[string]interface{})["Contratos_por_preliq"].([]interface{})
							resch := utilidades.GenChanInterface(listaLiquidacion...)
							var params []interface{}
							params = append(params, liquidacion.(map[string]interface{})["Id_Preliq"].(interface{}))
							f := formatoRegistroOpFunctionDispatcher(liquidacion.(map[string]interface{})["Nombre_tipo_nomina"].(string))

							if f != nil {

								chlistaLiquidacion := utilidades.Digest(done, f, resch, params)
								for dataLiquidacion := range chlistaLiquidacion {
									if dataLiquidacion != nil {
										respuesta = append(respuesta, dataLiquidacion)
									}
								}
								c.Data["json"] = respuesta
							} else {
								c.Data["json"] = models.Alert{Code: "E_0458", Body: nil, Type: "error"}
							}
						} else {
							c.Data["json"] = models.Alert{Code: "E_0458", Body: nil, Type: "error"}
						}
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

	} else {
		//no se enviaron los parametros necesarios
		c.Data["json"] = models.Alert{Code: "E_0458", Body: "Not enough parameter", Type: "error"}
	}
	c.ServeJSON()
}

// ListaPagoSsPorPersona ...
// @Title ListaPagoSsPorPersonaAllDetalle
// @Description lista pagos de seguridad socila por persona.
// @Param	idNomina	query	string	false	"nomina a listar"
// @Param	mesLiquidacion	query	string	false	"mes de la liquidacion a listar"
// @Param	anioLiquidacion	query	string	false	"anio de la liquidacion a listar"
// @Success 201 {object} models.Alert
// @Failure 403 body is empty
// @router /ListaPagoSsPorPersona [get]
func (c *OrdenPagoNominaController) ListaPagoSsPorPersona() {
	idNomina, err1 := c.GetInt("idNomina")
	mesLiquidacion, err2 := c.GetInt("mesLiquidacion")
	anioLiquidacion, err3 := c.GetInt("anioLiquidacion")
	if err1 == nil && err2 == nil && err3 == nil {
		pagosAgrupados := pagoSsPorPersonaF(idNomina, mesLiquidacion, anioLiquidacion)
		if pagosAgrupados != nil {
			fmt.Println(pagosAgrupados)
			done := make(chan interface{})
			defer close(done)

			c.Data["json"] = "hola"
		} else {
			c.Data["json"] = nil
		}
	} else {
		c.Data["json"] = models.Alert{Code: "E_0458", Body: "Not enough parameter", Type: "error"}
	}
	//c.Data["json"] = allData
	c.ServeJSON()
}

func pagoSsPorPersonaF(idNomina, mesLiquidacion, anioLiquidacion int) (dataOutp interface{}) {
	var liquidacion interface{}
	var periodoPago []interface{}
	var pagosPorDetalle []interface{}
	allData := make(map[string]interface{})
	if idNomina != 0 && mesLiquidacion != 0 && anioLiquidacion != 0 {
		// get id preliquidacion
		if err := getJson("http://"+beego.AppConfig.String("titanService")+"preliquidacion/contratos_x_preliquidacion?idNomina="+strconv.Itoa(idNomina)+"&mesLiquidacion="+strconv.Itoa(mesLiquidacion)+"&anioLiquidacion="+strconv.Itoa(anioLiquidacion), &liquidacion); err == nil {
			if liquidacion != nil {
				if liquidacion.(map[string]interface{})["Id_Preliq"] != nil {
					allData["IdPreliquidacion"] = liquidacion.(map[string]interface{})["Id_Preliq"]
				} else {
					allData["IdPreliquidacion"] = nil
				}
			} else {
				allData["IdPreliquidacion"] = nil
			}
		} else {
			return nil
		}
		// get periodo pago
		if allData["IdPreliquidacion"] != nil {
			idPeriodoLiquidacion := allData["IdPreliquidacion"].(float64)
			if err := getJson("http://"+beego.AppConfig.String("SsService")+"periodo_pago/?query=Mes:"+strconv.Itoa(mesLiquidacion)+"&Anio:"+strconv.Itoa(anioLiquidacion)+"&Liquidacion:"+strconv.FormatFloat(idPeriodoLiquidacion, 'f', -1, 64)+"&limit:1", &periodoPago); err == nil {
				if periodoPago != nil {
					if periodoPago[0].(map[string]interface{})["Id"] != nil {
						allData["PeriodoPago"] = periodoPago[0].(map[string]interface{})["Id"]
					} else {
						allData["PeriodoPago"] = nil
					}
				} else {
					allData["PeriodoPago"] = nil
				}
			} else {
				return nil
			}
		}
		//get pagos por persona
		if allData["PeriodoPago"] != nil {
			if err := getJson("http://"+beego.AppConfig.String("SsService")+"pago/PagosPorPeriodoPago?idPeriodoPago="+strconv.FormatFloat(allData["PeriodoPago"].(float64), 'f', -1, 64), &pagosPorDetalle); err == nil {
				if pagosPorDetalle != nil {
					allData["Pagos"] = pagosPorDetalle
				} else {
					allData["Pagos"] = nil
				}
			} else {
				return nil
			}
		}
		dataOutp = allData
	}
	return
}

func formatoRegistroOpHC(dataLiquidacion interface{}, params ...interface{}) (res interface{}) {
	idLiquidacion, e := params[0].(float64)
	if e {
		var respuesta []map[string]interface{}
		var listaDetalles []interface{}
		var valorTotal float64
		var params []interface{}
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
				f := homologacionFunctionDispatcher(listaDetalles[0].(map[string]interface{})["Preliquidacion"].(map[string]interface{})["Nomina"].(map[string]interface{})["TipoNomina"].(map[string]interface{})["Nombre"].(string))
				if f != nil {
					infoContrato := formatoListaLiquidacion(dataLiquidacion, nil)
					idProveedor, e := infoContrato.(map[string]interface{})["infoPersona"].(map[string]interface{})["id_persona"]
					if e {
						params = append(params, idProveedor)
					} else {
						return nil
					}
					numContrato, e := infoContrato.(map[string]interface{})["NumeroContrato"]
					if e {
						params = append(params, numContrato)
					} else {
						return nil
					}
					vigContrato, e := infoContrato.(map[string]interface{})["VigenciaContrato"]
					if e {
						params = append(params, vigContrato)
					} else {
						return nil
					}
					chConcHomologados := utilidades.Digest(done, f, resch, params)
					for conceptoHomologadoint := range chConcHomologados {
						conceptoHomologado, e := conceptoHomologadoint.(map[string]interface{})
						if e {
							existe := false
							for _, comp := range respuesta {
								if comp["Concepto"] != nil && conceptoHomologado["Concepto"] != nil {
									if comp["Concepto"].(map[string]interface{})["Id"].(float64) == conceptoHomologado["Concepto"].(map[string]interface{})["Id"].(float64) {
										comp["Valor"] = comp["Valor"].(float64) + conceptoHomologado["Valor"].(float64)
										existe = true
										valorTotal = valorTotal + comp["Valor"].(float64)
									}
								}

							}
							if !existe {
								if conceptoHomologado["Concepto"] != nil {
									valorTotal = valorTotal + conceptoHomologado["Valor"].(float64)
									respuesta = append(respuesta, conceptoHomologado)
								}

							}
						}
					}
					var movimientosContables []interface{}
					for _, concepto := range respuesta {
						movimientoContable := formatoMovimientosContablesOp(concepto)
						movimientosContables = append(movimientosContables, movimientoContable)
					}
					res := make(map[string]interface{})
					res["ValorBase"] = valorTotal
					res["id_proveedor"], err = strconv.Atoi(infoContrato.(map[string]interface{})["infoPersona"].(map[string]interface{})["id_persona"].(string))
					res["Conceptos"] = respuesta
					res["Contrato"] = nContrato
					res["MovimientoContable"] = movimientosContables
					return res
				} else {
					return nil
				}
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

func formatoMovimientosContablesOp(concepto interface{}) (res interface{}) {
	var out []map[string]interface{}
	cuentaContable, e := concepto.(map[string]interface{})["Concepto"].(map[string]interface{})["ConceptoCuentaContable"].([]interface{})
	if !e {
		//fmt.Println(concepto)
		fmt.Println("1")
		return nil
	}
	if len(cuentaContable) == 2 {
		for _, cuentaComp := range cuentaContable {
			fmt.Println(cuentaComp)
			if cuentaComp.(map[string]interface{})["CuentaContable"].(map[string]interface{})["Naturaleza"].(string) == "debito" {
				out = append(out, map[string]interface{}{"Debito": concepto.(map[string]interface{})["Valor"], "Credito": 0,
					"Concepto":       concepto.(map[string]interface{})["Concepto"],
					"CuentaContable": cuentaComp})
			} else {
				out = append(out, map[string]interface{}{"Debito": 0, "Credito": concepto.(map[string]interface{})["Valor"],
					"Concepto":       concepto.(map[string]interface{})["Concepto"],
					"CuentaContable": cuentaComp})
			}

		}
	} else {
		return nil
	}

	return out
}

func homologacionFunctionDispatcher(tipo string) (f func(data interface{}, params ...interface{}) interface{}) {
	switch os := tipo; os {
	case "HCS":
		return homologacionConceptosHC
	case "HCH":
		return nil
	default:
		return nil
	}
}
func formatoRegistroOpFunctionDispatcher(tipo string) (f func(data interface{}, params ...interface{}) interface{}) {
	switch os := tipo; os {
	case "HCS":
		return formatoRegistroOpHC
	case "HCH":
		return nil
	default:
		return nil
	}
}
