package controllers

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/api_mid_financiera/utilidades"
)

// Orden_pago_plantaController operations for Orden_pago_planta
type Orden_pago_plantaController struct {
	beego.Controller
}

// URLMapping ...
func (c *Orden_pago_plantaController) URLMapping() {

}

//funcion para generar el formato de la op planta para su vista previa.
func formatoRegistroOpDocentesPlanta(dataLiquidacion interface{}, params ...interface{}) (res interface{}) {
	var idLiquidacion float64
	var idRp []map[string]interface{}
	e1 := true
	e2 := true
	if params != nil {
		//beego.Info(params)
		idLiquidacion, e1 = params[0].(float64)
		idRp, e2 = params[1].([]map[string]interface{})
		//beego.Info(idRp)
	} else {
		return nil
	}
	if e1 && e2 {
		var homologacionConceptos []map[string]interface{}
		//var homologacionDescuentos []map[string]interface{}
		var valorTotal float64
		var params []interface{}
		var err error
		valorTotal = 0
		nContrato, e := dataLiquidacion.(map[string]interface{})["NumeroContrato"].(string)
		if !e {
			return nil
		}
		vigenciaContrato, e := dataLiquidacion.(map[string]interface{})["VigenciaContrato"].(float64)
		if !e {
			return nil
		}
		var infoContrato interface{}
		//consulta de los devengos calculados por persona
		var devengosNomina []interface{}
		if devengosNomina, err = ConsultarDevengosNominaPorContrato(idLiquidacion, nContrato, vigenciaContrato); err != nil {
			//si no se consiguen datos de devengos para el contrato liquidado (nomina planta)
			return nil
		}
		//beego.Info(devengosNomina)
		//homlogacion de los devengos de la nomina de planta a conceptos de kronos...
		if devengosNomina != nil {
			done := make(chan interface{})
			defer close(done)
			resch := utilidades.GenChanInterface(devengosNomina...)
			f := homologacionFunctionDispatcher(devengosNomina[0].(map[string]interface{})["Preliquidacion"].(map[string]interface{})["Nomina"].(map[string]interface{})["TipoNomina"].(map[string]interface{})["Nombre"].(string))
			if f != nil {
				infoContrato = formatoListaLiquidacion(dataLiquidacion, nil)
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

			}
		}
		//---------------------------------------------------------------------------
		//consulta de los descuentos calculados por persona
		/*var descuentosNomina []interface{}
		if descuentosNomina, err := ConsultarDescuentosNominaPorContrato(idLiquidacion, nContrato, vigenciaContrato); err != nil {
			//si no se consiguen datos de devengos para el contrato liquidado (nomina planta)
			return nil
		}
		//homlogacion de los descuentos de la nomina de planta a conceptos de kronos...
		done = make(chan interface{})
		defer close(done)
		resch = utilidades.GenChanInterface(listaDetalles...)
		chDescHomologados := utilidades.Digest(done, homologacionDescuentosHC, resch, nil)
		for descuentoHomologado := range chDescHomologados {
			homologado, e := descuentoHomologado.(map[string]interface{})
			if e {
				existe := false
				for _, comp := range homologacionDescuentos {
					if comp["Descuento"] != nil && homologado["Descuento"] != nil {
						if comp["Descuento"].(map[string]interface{})["Id"].(float64) == homologado["Descuento"].(map[string]interface{})["Id"].(float64) {
							comp["Valor"] = comp["Valor"].(float64) + homologado["Valor"].(float64)
							existe = true

						}
					}

				}
				if !existe {
					if homologado["Descuento"] != nil {

						homologacionDescuentos = append(homologacionDescuentos, homologado)

					}

				}
			}
		}*/
		//-----------------------------------------------------------------------------
		var movimientosContables []interface{}
		for _, concepto := range homologacionConceptos {
			movimientoContable := formatoMovimientosContablesOp(concepto)
			for _, aux := range movimientoContable {
				movimientosContables = append(movimientosContables, aux)
			}

		}
		/*		movcredito := findMovimientoCredito(movimientosContables)
				for _, descuento := range homologacionDescuentos {
					if movimientosContables != nil {

						movimientoContable, mov := formatoMovimientosContablesDescuentosOp(descuento, movcredito)
						movcredito = mov
						beego.Info(movcredito.(map[string]interface{})["Credito"])
						for _, aux := range movimientoContable {

							movimientosContables = append(movimientosContables, aux)
						}

					}

				}*/
		res := make(map[string]interface{})
		res["ValorBase"] = valorTotal
		if auxmap, e := infoContrato.(map[string]interface{}); e {
			res["infoPersona"], e = auxmap["infoPersona"]
		} else {
			res["infoPersona"] = make(map[string]interface{})
		}
		res["Contrato"] = nContrato
		res["VigenciaContrato"] = vigenciaContrato
		res["MovimientoContable"] = movimientosContables
		res["Aprobado"] = true
		res["ConceptoOrdenPago"], res["Aprobado"], res["Code"] = formatoConceptoOrdenPago(idRp, homologacionConceptos)
		return res
	} else {
		res := make(map[string]interface{})
		nContrato, e := dataLiquidacion.(map[string]interface{})["NumeroContrato"].(string)
		if !e {
			return nil
		}
		vigenciaContrato, e := dataLiquidacion.(map[string]interface{})["VigenciaContrato"].(float64)
		if !e {
			return nil
		}
		infoContrato := formatoListaLiquidacion(dataLiquidacion, nil)
		if auxmap, e := infoContrato.(map[string]interface{}); e {
			res["infoPersona"], e = auxmap["infoPersona"]
		}
		res["Contrato"] = nContrato
		res["VigenciaContrato"] = vigenciaContrato
		res["Aprobado"] = false
		res["Code"] = "OPM_E003"
		return res
	}
	return
}

func homologacionConceptosDocentesPlanta(dataConcepto interface{}, params ...interface{}) (res interface{}) {
	dataConceptoAhomologar, e := dataConcepto.(map[string]interface{})
	var vigContrato float64
	//fmt.Println(dataConcepto)
	if !e {
		fmt.Println("e1")
		return nil
	}
	//fmt.Println(dataConceptoAhomologar)
	out := make(map[string]interface{})
	var homologacion []interface{}
	//aqui va la consulta sobre facultad y proyecto para HC (modificar para hacerla de forma genral)
	if params != nil {
		vigContrato, e = params[2].(float64)
		//fmt.Println(vigContrato)
		if !e {
			fmt.Println("e3")
			return nil
		}
		//fmt.Println("http://" + beego.AppConfig.String("Urlcrud") + ":" + beego.AppConfig.String("Portcrud") + "/" + beego.AppConfig.String("Nscrud") + "/homologacion_concepto?query=ConceptoTitan:" + strconv.Itoa(int(dataConceptoAhomologar["Concepto"].(map[string]interface{})["Id"].(float64))) + ",ConceptoKronos.ConceptoTesoralFacultadProyecto.Facultad:" + strconv.Itoa(int(idFacultad)) + ",ConceptoKronos.ConceptoTesoralFacultadProyecto.ProyectoCurricular:" + strconv.Itoa(int(idProyecto)))
		if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/homologacion_concepto?query=ConceptoTitan:"+strconv.Itoa(int(dataConceptoAhomologar["Concepto"].(map[string]interface{})["Id"].(float64)))+",Vigencia:"+strconv.FormatFloat(vigContrato, 'f', -1, 64), &homologacion); err == nil {
			//fmt.Println("Hom ", homologacion)
			if homologacion != nil {
				//cuando hay homologacion de un concepto para concepto kronos.
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
				//cuando no encuentra la homologacion del concepto (buscar en descuentos).

			}
		} else {
			fmt.Println(err.Error())
			return nil
		}
	} else {
		return nil
	}

	return out

}

func formatoResumenCargueOpPlanta(infoDetalleCargueint interface{}, params ...interface{}) (resumen interface{}) {
	resRubr := make(map[float64]map[string]interface{})
	resMov := make(map[float64]map[string]interface{})
	infoDetalleCargue, _ := infoDetalleCargueint.([]interface{})
	if params != nil {
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

	} else {
		return nil
	}
}

func formatoPreViewCargueMasivoOpPlanta(liquidacion interface{}, params ...interface{}) (res interface{}) {
	var respuesta []interface{}
	var rp []map[string]interface{}
	if params != nil {
		done := make(chan interface{})
		defer close(done)
		_, e := liquidacion.(map[string]interface{})
		if e {
			if liquidacion.(map[string]interface{})["Contratos_por_preliq"] != nil {
				listaLiquidacion := liquidacion.(map[string]interface{})["Contratos_por_preliq"].([]interface{})
				resch := utilidades.GenChanInterface(listaLiquidacion...)
				var params2 []interface{}

				params2 = append(params2, liquidacion.(map[string]interface{})["Id_Preliq"].(interface{}))
				rp = formatoInfoRpById(params[0].(float64))
				if rp != nil {
					params2 = append(params2, rp)
				}
				f := formatoRegistroOpFunctionDispatcher(liquidacion.(map[string]interface{})["Nombre_tipo_nomina"].(string))
				//beego.Info(liquidacion.(map[string]interface{})["Nombre_tipo_nomina"].(string))
				if f != nil {

					chlistaLiquidacion := utilidades.Digest(done, f, resch, params2)
					for dataLiquidacion := range chlistaLiquidacion {
						if dataLiquidacion != nil {
							respuesta = append(respuesta, dataLiquidacion)
						}
					}
					fresumen := resumenOpFunctionDispatcher(liquidacion.(map[string]interface{})["Nombre_tipo_nomina"].(string))
					var resultado interface{}
					if fresumen != nil {
						resultado = formatoResumenCargueOp(respuesta)
					}
					res = map[string]interface{}{"DetalleCargueOp": respuesta, "ResumenCargueOp": resultado, "TipoLiquidacion": liquidacion.(map[string]interface{})["Nombre_tipo_nomina"].(string)}
				} else {
					beego.Info("1")
					res = map[string]interface{}{"Code": "E_0458", "Body": nil, "Type": "error"}
				}
			} else {
				beego.Info("2")
				res = map[string]interface{}{"Code": "E_0458", "Body": nil, "Type": "error"}
			}
		} else {
			beego.Info("3")
			res = map[string]interface{}{"Code": "E_0458", "Body": nil, "Type": "error"}
		}
	}
	return
}
