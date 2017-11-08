package controllers

import (
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
	if err1 == nil && err2 == nil && err3 == nil {
		var respuesta []interface{}
		var liquidacion interface{}
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
							resultado := formatoResumenCargueOp(respuesta)
							c.Data["json"] = map[string]interface{}{"DetalleCargueOp": respuesta, "ResumenCargueOp": resultado}
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
		//no se enviaron los parametros necesarios
		c.Data["json"] = models.Alert{Code: "E_0458", Body: "Not enough parameter", Type: "error"}
	}
	c.ServeJSON()
}

func formatoResumenCargueOp(infoDetalleCargue []interface{}) (resumen interface{}) {
	resRubr := make(map[float64]map[string]interface{})
	resMov := make(map[float64]map[string]interface{})

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
						fmt.Println("2")
						return
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
						return
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

func formatoRegistroOpHC(dataLiquidacion interface{}, params ...interface{}) (res interface{}) {
	idLiquidacion, e := params[0].(float64)
	if e {
		var homologacionConceptos []map[string]interface{}
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
		//consulta del rp asociado al contrato de la persona... strconv.Itoa(int(vigenciaContrato)) ... strconv.Itoa(int(solicitudrp))
		desagregacionrp := formatoInfoRp(nContrato, vigenciaContrato)
		//fin consulta del rp ...
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
						//return nil
						params = append(params, "0")
					}
					params = append(params, nContrato)
					params = append(params, vigenciaContrato)
					chConcHomologados := utilidades.Digest(done, f, resch, params)
					for conceptoHomologadoint := range chConcHomologados {
						conceptoHomologado, e := conceptoHomologadoint.(map[string]interface{})
						if e {
							existe := false
							for _, comp := range homologacionConceptos {
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
									homologacionConceptos = append(homologacionConceptos, conceptoHomologado)
								}

							}
						}
					}
					var movimientosContables []interface{}
					for _, concepto := range homologacionConceptos {
						movimientoContable := formatoMovimientosContablesOp(concepto)
						for _, aux := range movimientoContable {
							movimientosContables = append(movimientosContables, aux)
						}

					}
					res := make(map[string]interface{})
					res["ValorBase"] = valorTotal
					if rpint, e := desagregacionrp[0]["Rp"].(interface{}); e {
						ordenPago := make(map[string]interface{})
						ordenPago["RegistroPresupuestal"] = rpint
						res["OrdenPago"] = ordenPago
					} else {
						ordenPago := make(map[string]interface{})
						ordenPago["RegistroPresupuestal"] = nil
						res["OrdenPago"] = ordenPago
					}
					if auxmap, e := infoContrato.(map[string]interface{}); e {
						res["infoPersona"], e = auxmap["infoPersona"]
					}
					res["ConceptoOrdenPago"] = homologacionConceptos
					res["Contrato"] = nContrato
					res["VigenciaContrato"] = vigenciaContrato
					res["MovimientoContable"] = movimientosContables
					res["ConceptoOrdenPago"], res["Aprobado"], res["Code"] = formatoConceptoOrdenPago(desagregacionrp, homologacionConceptos)
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

func formatoConceptoOrdenPago(desgrRp []map[string]interface{}, conceptos []map[string]interface{}) (res []map[string]interface{}, comp bool, code string) {
	comp = false
	code = "OPM_S001"
	acumConceptos := make(map[float64]map[string]interface{})
	for _, concepto := range conceptos {
		if auxconcp, e := concepto["Concepto"].(map[string]interface{}); e {
			value := concepto["Valor"].(float64)
			idConcepto := auxconcp["Id"].(float64)
			if auxconcp, e = auxconcp["Rubro"].(map[string]interface{}); e {
				key := auxconcp["Id"].(float64)
				if acumConceptos[key] != nil {
					acumConceptos[key] = map[string]interface{}{"Valor": acumConceptos[key]["Valor"].(float64) + value, "Concepto": map[string]interface{}{"Id": idConcepto}}
				} else {
					acumConceptos[key] = map[string]interface{}{"Valor": value, "Concepto": map[string]interface{}{"Id": idConcepto}}
				}
			}

		} else {
			comp = false
			code = "OPM_E001"
			return
		}

	}
	for _, apRp := range desgrRp {
		if auxmap, e := apRp["Apropiacion"].(map[string]interface{}); e {
			if auxmap, e = auxmap["Rubro"].(map[string]interface{}); e {
				if idrbRp, e := auxmap["Id"].(float64); e {
					fmt.Println(35644)
					if acumConceptos[idrbRp] != nil {
						saldorp := apRp["Saldo"].(float64)
						fmt.Println("acum. ", idrbRp)
						if valor := acumConceptos[idrbRp]["Valor"].(float64); true && saldorp >= valor {
							comp = true
							acumConceptos[idrbRp]["RegistroPresupuestalDisponibilidadApropiacion"] = apRp["RegistroPresupuestalDisponibilidadApropiacion"]
							acumConceptos[idrbRp]["Apropiacion"] = apRp["Apropiacion"]
							res = append(res, acumConceptos[idrbRp])
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
	return
}

func formatoMovimientosContablesOp(concepto interface{}) (res []map[string]interface{}) {
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
				out = append(out, map[string]interface{}{"Debito": concepto.(map[string]interface{})["Valor"].(float64), "Credito": float64(0),
					"Concepto":       concepto.(map[string]interface{})["Concepto"],
					"CuentaContable": cuentaComp})
			} else {
				out = append(out, map[string]interface{}{"Debito": float64(0), "Credito": concepto.(map[string]interface{})["Valor"].(float64),
					"Concepto":       concepto.(map[string]interface{})["Concepto"],
					"CuentaContable": cuentaComp})
			}

		}
	} else {
		return nil
	}

	return out
}

func formatoInfoRp(nContrato string, vigenciaContrato float64) (desagregacionrp []map[string]interface{}) {
	var rp []interface{}
	var saldoRp map[string]float64
	if err := getJson("http://"+beego.AppConfig.String("argoService")+"solicitud_rp?limit=-1&query=Expedida:false,NumeroContrato:"+"DVE48"+",VigenciaContrato:"+"2017", &rp); err == nil && rp != nil {
		if rpmap, e := rp[0].(map[string]interface{}); e {
			if solicitudrp, e := rpmap["Id"].(float64); e {
				fmt.Println("sol rp : ", solicitudrp)
				if err = getJson("http://"+beego.AppConfig.String("kronosService")+"registro_presupuestal?limit=-1&query=Solicitud:"+"307", &rp); err == nil && rp != nil {
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
