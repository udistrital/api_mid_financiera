package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/api_mid_financiera/models"
	"github.com/udistrital/utils_oas/optimize"
)

// OrdenPagoNominaController operations for Orden_pago_planta
type OrdenPagoNominaController struct {
	beego.Controller
}

// URLMapping ...
func (c *OrdenPagoNominaController) URLMapping() {
	c.Mapping("Post", c.Post)
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
					resch := optimize.GenChanInterface(listaLiquidacion...)
					chlistaLiquidacion := optimize.Digest(done, formatoListaLiquidacion, resch, nil)
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
				resch := optimize.GenChanInterface(listaDetalles...)
				f := homologacionFunctionDispatcher(listaDetalles[0].(map[string]interface{})["Preliquidacion"].(map[string]interface{})["Nomina"].(map[string]interface{})["TipoNomina"].(map[string]interface{})["Nombre"].(string))
				var params []interface{}
				params = append(params, "persona")
				params = append(params, nContrato)
				params = append(params, vigenciaContrato)
				chConcHomologados := optimize.Digest(done, f, resch, params)
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

// PreviewCargueMasivoOp ...
// @Title PreviewCargueMasivoOp
// @Description lista liquidaciones para ordenes de pago masivas.
// @Param	idNomina	query	string	false	"nomina a listar"
// @Param	mesLiquidacion	query	string	false	"mes de la liquidacion a listar"
// @Param	anioLiquidacion	query	string	false	"anio de la liquidacion a listar"
// @Param   OrdenPago       map[string]string	true		"body for OrdenPago content"
// @Success 201 {object} models.Alert
// @Failure 403 body is empty
// @router /PreviewCargueMasivoOp [get]
func (c *OrdenPagoNominaController) PreviewCargueMasivoOp() {
	idNomina, err1 := c.GetInt("idNomina")
	mesLiquidacion, err2 := c.GetInt("mesLiquidacion")
	anioLiquidacion, err3 := c.GetInt("anioLiquidacion")
	var params []interface{}
	if err1 == nil && err2 == nil && err3 == nil {
		var liquidacion interface{}
		if err := getJson("http://"+beego.AppConfig.String("titanService")+"preliquidacion/contratos_x_preliquidacion?idNomina="+strconv.Itoa(idNomina)+"&mesLiquidacion="+strconv.Itoa(mesLiquidacion)+"&anioLiquidacion="+strconv.Itoa(anioLiquidacion), &liquidacion); err == nil {
			if liquidacion != nil {
				f := formatoPreviewOpFunctionDispatcher(liquidacion.(map[string]interface{})["Nombre_tipo_nomina"].(string))
				var res interface{}
				if f != nil {
					res = f(liquidacion, params...)
				}
				c.Data["json"] = res
			}
		} else {
			//error consumo de servicio titan. Lista contratos por liqu
			c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
		}

	} else {
		//no se enviaron los parametros necesarios
		c.Data["json"] = models.Alert{Code: "E_0458", Body: "Not enough parameter", Type: "error"}
	}
	c.ServeJSON()
}

// RegistroCargueMasivoOp ...
// @Title RegistroCargueMasivoOp
// @Description lista liquidaciones para ordenes de pago masivas.
// @Param	body		body 	models.RegistroPresupuestal	true		"body for RegistroPresupuestal content"
// @Success 201 {object} models.Alert
// @Failure 403 body is empty
// @router /RegistroCargueMasivoOp [post]
func (c *OrdenPagoNominaController) RegistroCargueMasivoOp() {
	var v map[string]interface{}
	var alert interface{}
	var param []interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if v != nil {
			param = append(param, v["InfoGeneralOp"])
			done := make(chan interface{})
			defer close(done)
			if tipo, e := v["TipoLiquidacion"].(string); e {
				f := RegistroOpFunctionDispatcher(tipo)
				if f != nil {
					alert = f(v, param)
				}
			} else {
				var aux []interface{}
				alert = append(aux, models.Alert{Code: "E_0458", Body: nil, Type: "error"})

			}

		} else {
			var aux []interface{}
			alert = append(aux, models.Alert{Code: "E_0458", Body: nil, Type: "error"})

		}
	} else {
		var aux []interface{}
		alert = append(aux, models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"})

	}
	c.Data["json"] = alert
	c.ServeJSON()
}

func RegistroOpProveedor(datain map[string]interface{}, params ...interface{}) (res interface{}) {
	//"http://"+beego.AppConfig.String("kronosService")+
	dataconv, _ := datain["DetalleCargueOp"].([]interface{})
	alerts := []models.Alert{}
	alert := models.Alert{}
	for _, data := range dataconv {
		if auxmap, e := data.(map[string]interface{}); e {
			if auxbool, e := auxmap["Aprobado"].(bool); e {
				if auxbool {
					valorBase, e2 := auxmap["ValorBase"].(float64)
					if Opmap, e := auxmap["OrdenPago"].(map[string]interface{}); e && e2 {
						Opmap["UnidadEjecutora"], e = params[0].([]interface{})[0].(map[string]interface{})["UnidadEjecutora"]
						Opmap["SubTipoOrdenPago"], e = params[0].([]interface{})[0].(map[string]interface{})["SubTipoOrdenPago"]
						Opmap["FormaPago"], e = params[0].([]interface{})[0].(map[string]interface{})["FormaPago"]
						Opmap["Vigencia"], e = params[0].([]interface{})[0].(map[string]interface{})["Vigencia"]
						Opmap["Documento"], e = params[0].([]interface{})[0].(map[string]interface{})["Documento"]
						Opmap["ValorBase"] = valorBase
						auxmap["OrdenPago"] = Opmap
						if err := sendJson("http://"+beego.AppConfig.String("kronosService")+"orden_pago/RegistrarOpProveedor", "POST", &alert, &auxmap); err == nil {
							alerts = append(alerts, alert)
						} else {
							alerts = append(alerts, models.Alert{Code: "E_0458", Body: data, Type: "error"})
						}
					}
				} else {
					alerts = append(alerts, models.Alert{Code: "E_0458", Body: data, Type: "error"})
				}
			} else {
				//si no se aprobo la op para su registro. (notificar a quien corresponda)
				alerts = append(alerts, models.Alert{Code: "OPM_E005", Body: data, Type: "error"})
			}
		} else {
			alerts = append(alerts, models.Alert{Code: "OPM_E005", Body: data, Type: "error"})
		}
	}

	return alerts
}

func formatoResumenCargueOp(infoDetalleCargueint interface{}, params ...interface{}) (resumen interface{}) {
	resRubr := make(map[float64]map[string]interface{})
	resMov := make(map[float64]map[string]interface{})
	infoDetalleCargue, _ := infoDetalleCargueint.([]interface{})
	for _, detalle := range infoDetalleCargue {
		if detallemap, e := detalle.(map[string]interface{}); e {
			if auxbool, e := detallemap["Aprobado"].(bool); e {
				if auxbool {
					//construccion del resumen de la afectacion presupuestal...
					if copmap, e := detallemap["ConceptoOrdenPago"].([]map[string]interface{}); e {
						for _, conceptoOp := range copmap {

							if auxmap, e := conceptoOp["RegistroPresupuestalDisponibilidadApropiacion"].(map[string]interface{}); e {
								if auxmap, e := auxmap["DisponibilidadApropiacion"].(map[string]interface{}); e {
									if auxmap, e := auxmap["Apropiacion"].(map[string]interface{}); e {
										if rubromap, e := auxmap["Rubro"].(map[string]interface{}); e {
											idrbr := rubromap["Id"].(float64)
											if resRubr[idrbr] != nil {
												resRubr[idrbr] = map[string]interface{}{"Rubro": rubromap, "Afectacion": resRubr[idrbr]["Afectacion"].(float64) + detallemap["ValorBase"].(float64)}
											} else {
												resRubr[idrbr] = map[string]interface{}{"Rubro": rubromap, "Afectacion": detallemap["ValorBase"].(float64)}
											}
										} else {
											fmt.Println("1")
											return
										}
									} else {
										fmt.Println("11")
										return
									}
								} else {
									fmt.Println("12")
									return
								}
							} else {
								fmt.Println("13")
								return
							}

						}
					} else {

					}
					//construccion del resumen de la afectacion Contable...

					if movsCont, e := detallemap["MovimientoContable"].([]interface{}); e {
						for _, movint := range movsCont {
							if mov, e := movint.(map[string]interface{}); e {

								if auxmap, e := mov["CuentaContable"].(map[string]interface{}); e {
									idCuenta := auxmap["Id"].(float64)
									fmt.Println("cuenta id ", idCuenta)
									if resMov[idCuenta] != nil {
										resMov[idCuenta] = map[string]interface{}{"CuentaContable": auxmap, "Debito": resMov[idCuenta]["Debito"].(float64) + mov["Debito"].(float64), "Credito": resMov[idCuenta]["Credito"].(float64) + mov["Credito"].(float64)}
									} else {
										resMov[idCuenta] = map[string]interface{}{"CuentaContable": auxmap, "Debito": mov["Debito"].(float64), "Credito": mov["Credito"].(float64)}
									}
								} else {
									fmt.Println("err mov 3")
								}
							} else {
								fmt.Println("err movs 2")
							}

						}
					} else {
						fmt.Println("err movs 1")

					}

				}
			} else {
				fmt.Println("4")
				return
			}

		} else {
			fmt.Println("5")
			return
		}
	}
	fmt.Println("3")
	var resRubrArr []map[string]interface{}
	for _, aux := range resRubr {
		resRubrArr = append(resRubrArr, aux)
	}
	var resMovArr []map[string]interface{}
	for _, aux := range resMov {
		resMovArr = append(resMovArr, aux)
	}
	return map[string]interface{}{"ResumenPresupuestal": resRubrArr, "ResumenContable": resMovArr}

}

func formatoConceptoOrdenPago(desgrRp []map[string]interface{}, conceptos []map[string]interface{}) (res []map[string]interface{}, comp bool, code string) {
	comp = false
	code = "OPM_S001"
	acumConceptos := make(map[float64]map[string]interface{})
	for _, concepto := range conceptos {
		if auxconcp, e := concepto["Concepto"].(map[string]interface{}); e {
			value := concepto["Valor"].(float64)
			//idConcepto := auxconcp["Id"].(float64)
			if auxconcp, e = auxconcp["Rubro"].(map[string]interface{}); e {
				key := auxconcp["Id"].(float64)
				if acumConceptos[key] != nil {
					var auxconceptos []interface{}
					auxconceptos = append(auxconceptos, acumConceptos[key]["Concepto"].([]interface{})...)
					auxconceptos = append(auxconceptos, concepto)
					acumConceptos[key] = map[string]interface{}{"Valor": acumConceptos[key]["Valor"].(float64) + value, "Concepto": auxconceptos}
				} else {
					var auxconceptos []interface{}
					auxconceptos = append(auxconceptos, concepto)
					acumConceptos[key] = map[string]interface{}{"Valor": value, "Concepto": auxconceptos}
				}
			}

		} else {
			comp = false
			code = "E_0458"
			return
		}

	}
	for _, apRp := range desgrRp {
		if auxmap, e := apRp["Apropiacion"].(map[string]interface{}); e {
			if auxmap, e = auxmap["Rubro"].(map[string]interface{}); e {
				if idrbRp, e := auxmap["Id"].(float64); e {
					//fmt.Println(35645)
					fmt.Println(idrbRp)
					if acumConceptos[idrbRp] != nil {
						saldorp := apRp["Saldo"].(float64)
						beego.Info("acum. ", idrbRp)
						if valor, e := acumConceptos[idrbRp]["Valor"].(float64); e && saldorp >= valor {
							comp = true
							if concetosmap, e := acumConceptos[idrbRp]["Concepto"].([]interface{}); e {
								for _, cpt := range concetosmap {
									if mapcpt, e := cpt.(map[string]interface{}); e {
										row := make(map[string]interface{})
										row["RegistroPresupuestalDisponibilidadApropiacion"] = apRp["RegistroPresupuestalDisponibilidadApropiacion"]
										row["Apropiacion"] = apRp["Apropiacion"]
										row["Concepto"] = mapcpt["Concepto"]
										row["Valor"] = mapcpt["Valor"]
										res = append(res, row)
									}

								}
							}

						} else {
							comp = false
							code = "OPM_E002"
						}
					} else {
						comp = false
						code = "OPM_E001"
					}
				} else {
					code = "E_0458"
				}
			} else {
				code = "E_0458"
			}
		} else {
			code = "E_0458"
		}
	}
	if conceptos == nil {
		code = "OPM_E004"
	}
	if desgrRp == nil {
		code = "OPM_E003"
	}
	return
}

func formatoMovimientosContablesOp(concepto interface{}) (res []map[string]interface{}) {
	var out []map[string]interface{}
	cuentaContable, e := concepto.(map[string]interface{})["Concepto"].(map[string]interface{})["ConceptoCuentaContable"].([]interface{})
	if !e {
		fmt.Println(concepto)
		fmt.Println("1 concepto")
		return nil
	}
	if len(cuentaContable) == 2 {
		for _, cuentaComp := range cuentaContable {
			//fmt.Println(cuentaComp)
			if cuentaComp.(map[string]interface{})["CuentaContable"].(map[string]interface{})["Naturaleza"].(string) == "debito" {
				out = append(out, map[string]interface{}{"Debito": concepto.(map[string]interface{})["Valor"].(float64), "Credito": float64(0),
					"Concepto":       concepto.(map[string]interface{})["Concepto"],
					"CuentaContable": cuentaComp.(map[string]interface{})["CuentaContable"]})
			} else {
				out = append(out, map[string]interface{}{"Debito": float64(0), "Credito": concepto.(map[string]interface{})["Valor"].(float64),
					"Concepto":       concepto.(map[string]interface{})["Concepto"],
					"CuentaContable": cuentaComp.(map[string]interface{})["CuentaContable"]})
			}

		}
	} else {
		return nil
	}

	return out
}

func findMovimientoCredito(movimientos []interface{}) (movimiento interface{}) {
	for _, movimiento := range movimientos {
		if auxmap, e := movimiento.(map[string]interface{}); e {
			if auxmap, e := auxmap["CuentaContable"].(map[string]interface{}); e {

				if naturaleza, e := auxmap["Naturaleza"].(string); e && naturaleza == "credito" {
					return movimiento
				}

			}
		}
	}
	return
}

func formatoMovimientosContablesDescuentosOp(descuento interface{}, movimiento interface{}) (res []map[string]interface{}, resmovimiento map[string]interface{}) {
	var out []map[string]interface{}
	cuentaComp, e := descuento.(map[string]interface{})["Descuento"]
	if !e {
		//fmt.Println(descuento)
		beego.Info(descuento)
		fmt.Println("1")
		return
	}
	if movmap, e := movimiento.(map[string]interface{}); e {
		//fmt.Println(cuentaComp)
		if cuentaComp.(map[string]interface{})["CuentaContable"].(map[string]interface{})["Naturaleza"].(string) == "debito" {
			out = append(out, map[string]interface{}{"Debito": descuento.(map[string]interface{})["Valor"].(float64), "Credito": float64(0),
				"CuentaEspecial": descuento.(map[string]interface{})["Descuento"],
				"CuentaContable": cuentaComp.(map[string]interface{})["CuentaContable"],
				"Concepto":       movmap["Concepto"]})
		} else {
			out = append(out, map[string]interface{}{"Debito": float64(0), "Credito": descuento.(map[string]interface{})["Valor"].(float64),
				"CuentaEspecial": descuento.(map[string]interface{})["Descuento"],
				"CuentaContable": cuentaComp.(map[string]interface{})["CuentaContable"],
				"Concepto":       movmap["Concepto"]})
		}
		movmap["Credito"] = movmap["Credito"].(float64) - out[0]["Credito"].(float64)
		movmap["Debito"] = movmap["Debito"].(float64) - out[0]["Debito"].(float64)
		return out, movmap
	}

	return out, resmovimiento
}

func formatoInfoRp(nContrato string, vigenciaContrato float64) (desagregacionrp []map[string]interface{}) {
	var rp []interface{}
	var saldoRp map[string]float64
	//DVE48
	fmt.Println("http://" + beego.AppConfig.String("argoService") + "solicitud_rp?limit=-1&query=Expedida:true,NumeroContrato:" + nContrato + ",VigenciaContrato:" + strconv.Itoa(int(vigenciaContrato)))
	if err := getJson("http://"+beego.AppConfig.String("argoService")+"solicitud_rp?limit=-1&query=Expedida:true,NumeroContrato:"+nContrato+",VigenciaContrato:"+strconv.Itoa(int(vigenciaContrato)), &rp); err == nil && rp != nil {
		//if err := getJson("http://"+beego.AppConfig.String("argoService")+"solicitud_rp?limit=-1&query=Expedida:false,NumeroContrato:"+"DVE48"+",VigenciaContrato:"+"2017", &rp); err == nil && rp != nil {
		if rpmap, e := rp[0].(map[string]interface{}); e {
			if solicitudrp, e := rpmap["Id"].(float64); e {
				fmt.Println("sol rp : ", solicitudrp)
				if err = getJson("http://"+beego.AppConfig.String("kronosService")+"registro_presupuestal?limit=-1&query=Solicitud:"+strconv.Itoa(int(solicitudrp)), &rp); err == nil && rp != nil {
					//if err = getJson("http://"+beego.AppConfig.String("kronosService")+"registro_presupuestal?limit=-1&query=Solicitud:"+"307", &rp); err == nil && rp != nil {
					rpmap = nil
					if rpmap, e = rp[0].(map[string]interface{}); e {
						if desagregacionpresrp, e := rpmap["RegistroPresupuestalDisponibilidadApropiacion"].([]interface{}); e {
							for _, infopresrp := range desagregacionpresrp {
								row := make(map[string]interface{})
								if info, e := infopresrp.(map[string]interface{}); e {
									if dispoap, e := info["DisponibilidadApropiacion"].(map[string]interface{}); e {
										row["RegistroPresupuestalDisponibilidadApropiacion"] = info
										row["Rp"] = rp[0]
										row["Apropiacion"] = dispoap["Apropiacion"]
										row["FuenteFinanciacion"] = dispoap["FuenteFinanciamiento"]
										if err = sendJson("http://"+beego.AppConfig.String("kronosService")+"registro_presupuestal/SaldoRp", "POST", &saldoRp, row); err == nil && rp != nil {
											row["Saldo"] = saldoRp["saldo"]
										}

										desagregacionrp = append(desagregacionrp, row)
									}

								}
							}
							return
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
				return nil
			}
		} else {
			return nil
		}
	} else {
		return nil
	}
}

func formatoInfoRpById(idRp float64) (desagregacionrp []map[string]interface{}) {
	var rp []interface{}
	var saldoRp map[string]float64
	if err := getJson("http://"+beego.AppConfig.String("kronosService")+"registro_presupuestal?limit=-1&query=Id:"+strconv.Itoa(int(idRp)), &rp); err == nil && rp != nil {
		//if err = getJson("http://"+beego.AppConfig.String("kronosService")+"registro_presupuestal?limit=-1&query=Solicitud:"+"307", &rp); err == nil && rp != nil {
		if rpmap, e := rp[0].(map[string]interface{}); e {
			if desagregacionpresrp, e := rpmap["RegistroPresupuestalDisponibilidadApropiacion"].([]interface{}); e {
				for _, infopresrp := range desagregacionpresrp {
					row := make(map[string]interface{})
					if info, e := infopresrp.(map[string]interface{}); e {
						if dispoap, e := info["DisponibilidadApropiacion"].(map[string]interface{}); e {
							row["RegistroPresupuestalDisponibilidadApropiacion"] = info
							row["Rp"] = rp[0]
							row["Apropiacion"] = dispoap["Apropiacion"]
							row["FuenteFinanciacion"] = dispoap["FuenteFinanciamiento"]
							if err = sendJson("http://"+beego.AppConfig.String("kronosService")+"registro_presupuestal/SaldoRp", "POST", &saldoRp, row); err == nil && rp != nil {
								row["Saldo"] = saldoRp["saldo"]
							}

							desagregacionrp = append(desagregacionrp, row)
						}

					}
				}
				return
			} else {
				return nil
			}
		} else {
			return nil
		}
	} else {
		return nil
	}
}

func GetRpDesdeNecesidadProcesoExternoGeneral(idLiquidacion float64, CodigoAbreviacion string) (rpDisponibilidadApropiacion []map[string]interface{}, outputError map[string]interface{}) {
	//var outputError []map[string]interface{}
	if idLiquidacion != 0 {
		if idNecesidad, e := getNecesidadByProcesoExterno(int(idLiquidacion), CodigoAbreviacion); e == nil {
			if solicitudCDP, e := getSolicitudDisponibilidad(int(idNecesidad)); e == nil {
				if disponibilidad, e := getDisponibilidad(int(solicitudCDP)); e == nil {
					if rpDisponibilidadApropiacion, outputError = getRegistroPresupuestalDisponibilidadApropiacion(int(disponibilidad)); outputError == nil {
						fmt.Println("rp", rpDisponibilidadApropiacion[0]["Rp"].(map[string]interface{})["Id"])
						return rpDisponibilidadApropiacion, nil
					} else {
						return nil, outputError
					}
				} else {
					return nil, e
				}
			} else {
				return nil, e
			}
		} else {
			return nil, e
		}

	} else {
		outputError = map[string]interface{}{"Code": "E_0458", "Body": "Not enough parameter in GetRpDesdeNecesidadProcesoExterno", "Type": "error"}
		return nil, outputError
	}
}

func getNecesidadByProcesoExterno(idLiquidacion int, CodigoAbreviacion string) (necesidad float64, outputError map[string]interface{}) {
	var necesidadProcesoExterno []interface{}
	if idLiquidacion != 0 {
		//TipoNecesidad.CodigoAbreviacion:S  seguridad social
		// Necesidad.EstadoNecesidad.CodigoAbreviacion:C  => Solicitud de CDP creada
		fmt.Println("http://" + beego.AppConfig.String("argoService") + "necesidad_proceso_externo?query=TipoNecesidad.CodigoAbreviacion:" + CodigoAbreviacion + ",ProcesoExterno:" + strconv.Itoa(idLiquidacion) + ",Necesidad.EstadoNecesidad.CodigoAbreviacion:C&limit=1")
		if err := getJson("http://"+beego.AppConfig.String("argoService")+"necesidad_proceso_externo?query=TipoNecesidad.CodigoAbreviacion:"+CodigoAbreviacion+",ProcesoExterno:"+strconv.Itoa(idLiquidacion)+",Necesidad.EstadoNecesidad.CodigoAbreviacion:C&limit=1", &necesidadProcesoExterno); err == nil && necesidadProcesoExterno != nil && necesidadProcesoExterno[0].(map[string]interface{})["Necesidad"].(map[string]interface{})["Id"] != nil {
			necesidad = necesidadProcesoExterno[0].(map[string]interface{})["Necesidad"].(map[string]interface{})["Id"].(float64)
			return necesidad, nil
		} else {
			outputError = map[string]interface{}{"Code": "E_0458", "Body": "No existe necesidad de proceso externo para liquidacion de Seguridad Social en el periodo", "Type": "error"}
			return 0, outputError
		}
	} else {
		outputError = map[string]interface{}{"Code": "E_0458", "Body": "Not enough parameter in getNecesidadByProcesoExterno", "Type": "error"}
		return 0, outputError
	}
}

func resumenOpFunctionDispatcher(tipo string) (f func(data interface{}, params ...interface{}) interface{}) {
	switch os := tipo; os {
	case "HCS":
		return formatoResumenCargueOp
	case "HCH":
		return formatoResumenCargueOp
	case "FP":
		return formatoResumenCargueOp
	default:
		return nil
	}
}

func homologacionFunctionDispatcher(tipo string) (f func(data interface{}, params ...interface{}) interface{}) {
	switch os := tipo; os {
	case "HCS":
		return homologacionConceptosHC
	case "HCH":
		return homologacionConceptosHC
	case "FP":
		return homologacionConceptosDocentesPlanta
	default:
		return nil
	}
}
func formatoRegistroOpFunctionDispatcher(tipo string) (f func(data interface{}, params ...interface{}) interface{}) {
	switch os := tipo; os {
	case "HCS":
		return formatoRegistroOpHC
	case "HCH":
		return formatoRegistroOpHC
	case "FP":
		return formatoResumenOpPlanta
	default:
		return nil
	}
}

func formatoPreviewOpFunctionDispatcher(tipo string) (f func(data interface{}, params ...interface{}) interface{}) {
	switch os := tipo; os {
	case "HCS":
		return formatoPreViewCargueMasivoOpHc
	case "HCH":
		return formatoPreViewCargueMasivoOpHc
	case "FP":
		return formatoPreViewCargueMasivoOpPlanta
	default:
		return nil
	}
}

func RegistroOpFunctionDispatcher(tipo string) (f func(data map[string]interface{}, params ...interface{}) interface{}) {
	switch os := tipo; os {
	case "HCS":
		return RegistroOpProveedor
	case "HCH":
		return RegistroOpProveedor
	case "FP":
		return RegistroOpPlanta
	default:
		return nil
	}
}

func ConsultarDevengosNominaPorContrato(idLiquidacion float64, nContrato string, vigenciaContrato float64) (detalle []interface{}, err error) {
	//fmt.Println("http://" + beego.AppConfig.String("titanService") + "detalle_preliquidacion?limit=-1&query=Concepto.NaturalezaConcepto.Nombre:devengo,Preliquidacion.Id:" + strconv.Itoa(int(idLiquidacion)) + ",NumeroContrato:" + nContrato + ",VigenciaContrato:" + strconv.Itoa(int(vigenciaContrato)))
	if err = getJson("http://"+beego.AppConfig.String("titanService")+"detalle_preliquidacion?limit=-1&query=Concepto.NaturalezaConcepto.Nombre:devengo,Preliquidacion.Id:"+strconv.Itoa(int(idLiquidacion))+",NumeroContrato:"+nContrato+",VigenciaContrato:"+strconv.Itoa(int(vigenciaContrato)), &detalle); err == nil {
		return
	} else {
		return nil, err
	}
}

func ConsultarDescuentosNominaPorContrato(idLiquidacion float64, nContrato string, vigenciaContrato float64) (detalle []interface{}, err error) {
	if err = getJson("http://"+beego.AppConfig.String("titanService")+"detalle_preliquidacion?limit=-1&query=Concepto.NaturalezaConcepto.Nombre:descuento,Preliquidacion.Id:"+strconv.Itoa(int(idLiquidacion))+",NumeroContrato:"+nContrato+",VigenciaContrato:"+strconv.Itoa(int(vigenciaContrato)), &detalle); err == nil {
		return
	} else {
		return nil, err
	}
}
