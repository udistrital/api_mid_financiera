package controllers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/api_mid_financiera/models"
	"github.com/udistrital/utils_oas/optimize"
	"github.com/udistrital/utils_oas/request"
)

// OrdenPagoSsController operations for Orden_pago_seguridad_social
type OrdenPagoSsController struct {
	beego.Controller
}

// URLMapping ...
func (c *OrdenPagoSsController) URLMapping() {
	c.Mapping("GetConceptosMovimeintosContablesSs", c.GetConceptosMovimeintosContablesSs)
}

// TestJota01 ...
// @Title TestJota01
// @Description lista pagos de seguridad socila por persona.
// @Success 201 {object} models.Alert
// @Failure 403 body is empty
// @router /TestJota01 [get]
func (c *OrdenPagoSsController) TestJota01() {
	if data, e := pagoSsPorPersona(23); e == nil {
		c.Data["json"] = data
	} else {
		c.Data["json"] = e
	}
	c.ServeJSON()
}

// agrupa los pagos por detalle de liquidacion donde se referencia el contrato y usuario
func pagoSsPorPersona(idPeriodoPago float64) (dataOutp interface{}, outputError map[string]interface{}) {
	if idPeriodoPago != 0 {
		var pagosPorDetalle []interface{}
		var pagosVigenciaContrato []interface{}
		if err := request.GetJson("http://"+beego.AppConfig.String("SsService")+"pago/PagosPorPeriodoPago?idPeriodoPago="+strconv.FormatFloat(idPeriodoPago, 'f', -1, 64), &pagosPorDetalle); err == nil && pagosPorDetalle != nil {
			done := make(chan interface{})
			defer close(done)
			resch := optimize.GenChanInterface(pagosPorDetalle...)
			// add vigencia y contrato a la agrupacion de pagos
			chlistaPagos := optimize.Digest(done, getContratoVigenciaDetalleLiquidacion, resch, nil)
			for datalistaPagos := range chlistaPagos {
				if datalistaPagos != nil {
					pagosVigenciaContrato = append(pagosVigenciaContrato, datalistaPagos.(interface{}))
				}
			}
			for _, pagoAgrupago := range pagosVigenciaContrato {
				if rowPagoAgrupado, e := pagoAgrupago.(map[string]interface{}); e {
					if detallePagos, e := listaDetalleDePagoSsPorPersona(int(idPeriodoPago), int(rowPagoAgrupado["DetalleLiquidacion"].(float64))); e == nil {
						rowPagoAgrupado["DetallePagos"] = detallePagos
					} else {
						rowPagoAgrupado["DetallePagos"] = e
					}
				}
			}
			return pagosVigenciaContrato, nil
		} else {
			outputError = map[string]interface{}{"Code": "E_0458", "Body": "No se encontraron pagos para el periodo pago establecido", "Type": "error"}
			return nil, outputError
		}
	} else {
		outputError = map[string]interface{}{"Code": "E_0458", "Body": "Not enough parameter in pagoSsPorPersona", "Type": "error"}
		return nil, outputError
	}
}

func getContratoVigenciaDetalleLiquidacion(idsLiquidacionDesdePagos interface{}, params ...interface{}) (res interface{}) {
	row, e := idsLiquidacionDesdePagos.(map[string]interface{})
	var infoDetallePreliquidacion []interface{}
	var infoPersona interface{}
	if e {
		if err := request.GetJson("http://"+beego.AppConfig.String("titanService")+"detalle_preliquidacion/?query=Id:"+strconv.FormatFloat(row["DetalleLiquidacion"].(float64), 'f', 0, 64)+"&limit=1", &infoDetallePreliquidacion); err == nil && infoDetallePreliquidacion != nil {
			if data, e1 := infoDetallePreliquidacion[0].(map[string]interface{}); e1 {
				row["NumeroContrato"] = data["NumeroContrato"]
				row["VigenciaContrato"] = data["VigenciaContrato"]
				if err := request.GetJsonWSO2("http://jbpm.udistritaloas.edu.co:8280/services/contrato_suscrito_DataService.HTTPEndpoint/informacion_contrato_elaborado_contratista/"+row["NumeroContrato"].(string)+"/"+strconv.Itoa(int(row["VigenciaContrato"].(float64))), &infoPersona); err == nil {
					if row["infoPersona"], e = infoPersona.(map[string]interface{})["informacion_contratista"]; e {
						return row
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

func listaDetalleDePagoSsPorPersona(idPeriodoPago, idDetalleLiquidacion int) (DetallePagos []interface{}, outputError map[string]interface{}) {
	if idPeriodoPago != 0 && idDetalleLiquidacion != 0 {
		var listaPagos []interface{}
		var allData []interface{}
		if err := request.GetJson("http://"+beego.AppConfig.String("SsService")+"pago/?query=DetalleLiquidacion:"+strconv.Itoa(idDetalleLiquidacion)+"&PeriodoPago.Id:"+strconv.Itoa(idPeriodoPago), &listaPagos); err == nil && listaPagos != nil {
			done := make(chan interface{})
			defer close(done)
			resch := optimize.GenChanInterface(listaPagos...)
			chlistaPagos := optimize.Digest(done, getTipoPagoTitan, resch, nil)
			for dataChlistaPagos := range chlistaPagos {
				if dataChlistaPagos != nil {
					allData = append(allData, dataChlistaPagos.(interface{}))
				}
			}
			return allData, nil
		} else {
			outputError = map[string]interface{}{"Code": "E_0458", "Body": "No se encontraron pagos para el periodo pago y detalle de liquidacion especificados", "Type": "error"}
			return nil, outputError
		}
	} else {
		outputError = map[string]interface{}{"Code": "E_0458", "Body": "Not enough parameter in ListaDetalleDePagoSsPorPersona", "Type": "error"}
		return nil, outputError
	}
}

func getTipoPagoTitan(listaPagos interface{}, params ...interface{}) (res interface{}) {
	row, e := listaPagos.(map[string]interface{})
	var tipoPago []interface{}
	if e {
		if err := request.GetJson("http://"+beego.AppConfig.String("titanService")+"concepto_nomina/?query=Id:"+strconv.FormatFloat(row["TipoPago"].(float64), 'f', 0, 64)+"&limit=1", &tipoPago); err == nil {
			row["TipoPago"] = tipoPago[0].(map[string]interface{})
			return row
		} else {
			return nil
		}
	} else {
		return nil
	}
}

// GetConceptosMovimeintosContablesSs ...
// @Title GetConceptosMovimeintosContablesSs
// @Description lista pagos de seguridad socila por persona.
// @Param	idNomina	query	string	false	"nomina a listar"
// @Param	mesLiquidacion	query	string	false	"mes de la liquidacion a listar"
// @Param	anioLiquidacion	query	string	false	"anio de la liquidacion a listar"
// @Success 201 {object} models.Alert
// @Failure 403 body is empty
// @router /GetConceptosMovimeintosContablesSs [get]
func (c *OrdenPagoSsController) GetConceptosMovimeintosContablesSs() {
	idNomina, err1 := c.GetInt("idNomina")
	mesLiquidacion, err2 := c.GetInt("mesLiquidacion")
	anioLiquidacion, err3 := c.GetInt("anioLiquidacion")
	if err1 == nil && err2 == nil && err3 == nil {
		var homologacionConceptos []map[string]interface{}
		var idLiquidacion float64
		if rpCorrespondiente, e := getRpDesdeNecesidadProcesoExterno(idNomina, mesLiquidacion, anioLiquidacion); e == nil {
			//if idLiquidacion, outputError := getIdliquidacionForSs(idNomina, mesLiquidacion, anioLiquidacion); outputError == nil {
			if liquidacionTipoNomina, outputError := getIdliquidacionTipoNominaForSs(idNomina, mesLiquidacion, anioLiquidacion); outputError == nil {
				idLiquidacion = liquidacionTipoNomina["Id_Preliq"].(float64)
				if idPeriodoPago, outputError := getIDPeriodoPagoForSs(int(idLiquidacion), mesLiquidacion, anioLiquidacion); outputError == nil {
					allPago := getPagosConDetalleLiquidacion(int(idPeriodoPago))
					if allPago != nil {
						done := make(chan interface{})
						defer close(done)
						resch := optimize.GenChanInterface(allPago...)
						chConcHomologados := optimize.Digest(done, homologacionConceptosSS, resch, nil)
						for conceptoHomologadoint := range chConcHomologados {
							conceptoHomologado, e := conceptoHomologadoint.(map[string]interface{})
							if e {
								existe := false
								for _, comp := range homologacionConceptos {
									if comp["Concepto"] != nil && conceptoHomologado["Concepto"] != nil {
										if comp["Concepto"].(map[string]interface{})["Id"].(float64) == conceptoHomologado["Concepto"].(map[string]interface{})["Id"].(float64) {
											comp["Valor"] = comp["Valor"].(float64) + conceptoHomologado["Valor"].(float64)
											existe = true
										}
									}
								}
								if !existe {
									if conceptoHomologado["Concepto"] != nil {
										homologacionConceptos = append(homologacionConceptos, conceptoHomologado)
									}
								}
							}
						}
						//c.Data["json"] = homologacionConceptos
						// -- movimeintos
						if homologacionConceptos != nil {
							var movimientosContables []map[string]interface{}
							for _, concepto := range homologacionConceptos {
								movimientoContable := formatoMovimientosContablesOp(concepto)
								for _, aux := range movimientoContable {
									movimientosContables = append(movimientosContables, aux)
								}
							}
							//estructura out fin
							ordenData := make(map[string]interface{})
							ordenData["Liquidacion"] = idPeriodoPago
							ordenData["RegistroPresupuestal"] = rpCorrespondiente[0]["Rp"].(interface{})
							allDataOuput := make(map[string]interface{})
							allDataOuput["OrdenPago"] = ordenData
							//totalizar los movimientos
							if movimientosDeOP, e := getMovimientosDescuentoDeLiquidacion(int(idLiquidacion), idNomina); e == nil {
								if allMovimientos, allAfectacion, e := afectarDescuentosToMovimientos(movimientosDeOP, movimientosContables); e == nil {
									allDataOuput["MovimientoContable"] = allMovimientos
									if AllConceptos, e := afectarDescuentosToConceptos(allAfectacion, homologacionConceptos); e == nil {
										allDataOuput["ConceptoOrdenPago"], allDataOuput["Aprobado"], allDataOuput["Code"] = formatoConceptoOrdenPago(rpCorrespondiente, AllConceptos)
										if valorBase, e := getTotalAfectacionOfConceptos(AllConceptos); e == nil {
											allDataOuput["ValorBase"] = valorBase
										} else {
											allDataOuput["Aprobado"] = false
											allDataOuput["Code"] = e["Code"]
											allDataOuput["ValorBase"] = e
										}
									} else {
										allDataOuput["Aprobado"] = false
										allDataOuput["Code"] = e["Code"]
										allDataOuput["ConceptoOrdenPago"] = e
									}
								} else {
									allDataOuput["Aprobado"] = false
									allDataOuput["Code"] = e["Code"]
									allDataOuput["MovimientoContable"] = e
								}
							} else {
								allDataOuput["Aprobado"] = false
								allDataOuput["Code"] = e["Code"]
								allDataOuput["MovimientoContable"] = e
							}
							// pintar pagos por personas
							if data, e := pagoSsPorPersona(idPeriodoPago); e == nil {
								allDataOuput["ViewPagosPorPersona"] = data
							} else {
								allDataOuput["ViewPagosPorPersona"] = e
							}
							allDataOuput["RegistroPresupuestal"] = rpCorrespondiente[0]["Rp"].(interface{})
							// for test
							if movimientosSoloDescuentos, e := getMovimientosDescuentoDeLiquidacion(int(idLiquidacion), idNomina); e == nil {
								allDataOuput["MovimientosDeDescuento"] = movimientosSoloDescuentos
							} else {
								allDataOuput["MovimientosDeDescuento"] = e
							}
							// for test
							var dataformato []map[string]interface{}
							dataformato = append(dataformato, allDataOuput)
							output := make(map[string]interface{})
							output["TipoLiquidacion"] = liquidacionTipoNomina["Nombre_tipo_nomina"].(string)
							output["DetalleCargueOp"] = dataformato

							c.Data["json"] = output
						} else {
							c.Data["json"] = models.Alert{Code: "E_0458", Body: "Erro en la homologacion de los conceptos", Type: "error"}
						}
					} else {
						c.Data["json"] = models.Alert{Code: "E_0458", Body: "no se logro asocial informacion del detalle de liquidacion a los pagos de Seguridad Social para el periodo", Type: "error"}
					}
				} else {
					c.Data["json"] = outputError
				}
			} else {
				c.Data["json"] = outputError
			}

		} else {
			c.Data["json"] = e
		}
	} else {
		c.Data["json"] = models.Alert{Code: "E_0458", Body: "Not enough parameter", Type: "error"}
	}
	c.ServeJSON()
}

func getRpDesdeNecesidadProcesoExterno(idNomina, mesLiquidacion, anioLiquidacion int) (rpDisponibilidadApropiacion []map[string]interface{}, outputError map[string]interface{}) {
	//var outputError []map[string]interface{}
	if idNomina != 0 && mesLiquidacion != 0 && anioLiquidacion != 0 {
		if idLiquidacion, e := getIdliquidacionForSs(idNomina, mesLiquidacion, anioLiquidacion); e == nil {
			fmt.Println("Liquidacion: ", idLiquidacion)
			if idPeriodoPago, e := getIDPeriodoPagoForSs(int(idLiquidacion), mesLiquidacion, anioLiquidacion); e == nil {
				fmt.Println("periodo pago ", idPeriodoPago)
				if idNecesidad, e := getNecesidadByProcesoExternoSS(int(idPeriodoPago)); e == nil {
					fmt.Println("Nececidad id", idNecesidad)
					if solicitudCDP, e := getSolicitudDisponibilidad(int(idNecesidad)); e == nil {
						fmt.Println("solicitud disponibilidad id, ", solicitudCDP)
						if disponibilidad, e := getDisponibilidad(int(solicitudCDP)); e == nil {
							fmt.Println("disponibilidad", disponibilidad)
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

// se consulta servicio que retorna las liquidacions en un mes, año y titpo nomina que ya esten en estado EnOrdenPago
func getIdliquidacionForSs(idNomina, mesLiquidacion, anioLiquidacion int) (IdLiquidacion float64, outputError map[string]interface{}) {
	var liquidacion interface{}
	if idNomina != 0 && mesLiquidacion != 0 && anioLiquidacion != 0 {
		if err := request.GetJson("http://"+beego.AppConfig.String("titanService")+"preliquidacion/contratos_x_preliquidacion?idNomina="+strconv.Itoa(idNomina)+"&mesLiquidacion="+strconv.Itoa(mesLiquidacion)+"&anioLiquidacion="+strconv.Itoa(anioLiquidacion), &liquidacion); err == nil {
			if liquidacion != nil && liquidacion.(map[string]interface{})["Id_Preliq"].(float64) != 0 {
				IdLiquidacion = liquidacion.(map[string]interface{})["Id_Preliq"].(float64)
				return IdLiquidacion, nil
			} else {
				outputError = map[string]interface{}{"Code": "E_0458", "Body": "No existe liquidacion en estado EnOrdenPago para el periodo", "Type": "error"}
				return 0, outputError
			}
		} else {
			outputError = map[string]interface{}{"Code": "E_0458", "Body": "No existe liquidacion en estado EnOrdenPago para el periodo", "Type": "error"}
			return 0, outputError
		}
	} else {
		outputError = map[string]interface{}{"Code": "E_0458", "Body": "Not enough parameter in getIdliquidacionForSs", "Type": "error"}
		return 0, outputError
	}
}

// se consulta servicio que retorna las liquidacions en un mes, año y titpo nomina que ya esten en estado EnOrdenPago
func getIdliquidacionTipoNominaForSs(idNomina, mesLiquidacion, anioLiquidacion int) (IdLiquidacion map[string]interface{}, outputError map[string]interface{}) {
	var liquidacion map[string]interface{}
	if idNomina != 0 && mesLiquidacion != 0 && anioLiquidacion != 0 {
		if err := request.GetJson("http://"+beego.AppConfig.String("titanService")+"preliquidacion/contratos_x_preliquidacion?idNomina="+strconv.Itoa(idNomina)+"&mesLiquidacion="+strconv.Itoa(mesLiquidacion)+"&anioLiquidacion="+strconv.Itoa(anioLiquidacion), &liquidacion); err == nil {
			if liquidacion != nil && liquidacion["Id_Preliq"].(float64) != 0 {
				return liquidacion, nil
			} else {
				outputError = map[string]interface{}{"Code": "E_0458", "Body": "No existe liquidacion en estado EnOrdenPago para el periodo in getIdliquidacionTipoNominaForSs", "Type": "error"}
				return nil, outputError
			}
		} else {
			outputError = map[string]interface{}{"Code": "E_0458", "Body": "No existe liquidacion en estado EnOrdenPago para el periodo in getIdliquidacionTipoNominaForSs", "Type": "error"}
			return nil, outputError
		}
	} else {
		outputError = map[string]interface{}{"Code": "E_0458", "Body": "Not enough parameter in getIdliquidacionTipoNominaForSs", "Type": "error"}
		return nil, outputError
	}
}

// se consulta servicio de periodo_pago en un mes, año y con id liquidacion
func getIDPeriodoPagoForSs(idLiquidacion, mesLiquidacion, anioLiquidacion int) (idPeriodoPago float64, outputError map[string]interface{}) {
	var periodoPago []interface{}
	if idLiquidacion != 0 && mesLiquidacion != 0 && anioLiquidacion != 0 {
		if err := request.GetJson("http://"+beego.AppConfig.String("SsService")+"periodo_pago/?query=Mes:"+strconv.Itoa(mesLiquidacion)+"&Anio:"+strconv.Itoa(anioLiquidacion)+"&Liquidacion:"+strconv.Itoa(idLiquidacion)+"&limit:1", &periodoPago); err == nil {
			if periodoPago != nil && periodoPago[0].(map[string]interface{})["Id"] != nil {
				idPeriodoPago = periodoPago[0].(map[string]interface{})["Id"].(float64)
				return idPeriodoPago, nil
			} else {
				outputError = map[string]interface{}{"Code": "E_0458", "Body": "No existe periodo pago de Seguridad Social para el periodo", "Type": "error"}
				return 0, outputError
			}
		} else {
			outputError = map[string]interface{}{"Code": "E_0458", "Body": "No existe periodo pago de Seguridad Social para el periodo", "Type": "error"}
			return 0, outputError
		}
	} else {
		outputError = map[string]interface{}{"Code": "E_0458", "Body": "Not enough parameter in getIDPeriodoPagoForSs", "Type": "error"}
		return 0, outputError
	}
}

func getNecesidadByProcesoExternoSS(idPeriodoPagoSs int) (necesidad float64, outputError map[string]interface{}) {
	var necesidadProcesoExterno []interface{}
	if idPeriodoPagoSs != 0 {
		//TipoNecesidad.CodigoAbreviacion:S  seguridad social
		// Necesidad.EstadoNecesidad.CodigoAbreviacion:C  => Solicitud de CDP creada
		if err := request.GetJson("http://"+beego.AppConfig.String("argoService")+"necesidad_proceso_externo?query=TipoNecesidad.CodigoAbreviacion:S,ProcesoExterno:"+strconv.Itoa(idPeriodoPagoSs)+",Necesidad.EstadoNecesidad.CodigoAbreviacion:C&limit:1", &necesidadProcesoExterno); err == nil && necesidadProcesoExterno != nil && necesidadProcesoExterno[0].(map[string]interface{})["Necesidad"].(map[string]interface{})["Id"] != nil {
			necesidad = necesidadProcesoExterno[0].(map[string]interface{})["Necesidad"].(map[string]interface{})["Id"].(float64)
			return necesidad, nil
		} else {
			outputError = map[string]interface{}{"Code": "E_0458", "Body": "No existe necesidad de proceso externo para liquidacion de Seguridad Social en el periodo", "Type": "error"}
			return 0, outputError
		}
	} else {
		outputError = map[string]interface{}{"Code": "E_0458", "Body": "Not enough parameter in getNecesidadByProcesoExternoSS", "Type": "error"}
		return 0, outputError
	}
}

func getSolicitudDisponibilidad(idNecesidad int) (solicitudDisponibilidad float64, outputError map[string]interface{}) {
	var solicitudDisponibilidadData []interface{}
	if idNecesidad != 0 {
		if err := request.GetJson("http://"+beego.AppConfig.String("argoService")+"solicitud_disponibilidad?query=Expedida:true,Necesidad.Id:"+strconv.Itoa(idNecesidad), &solicitudDisponibilidadData); err == nil && solicitudDisponibilidadData != nil && solicitudDisponibilidadData[0].(map[string]interface{})["Id"] != nil {
			solicitudDisponibilidad = solicitudDisponibilidadData[0].(map[string]interface{})["Id"].(float64)
			return solicitudDisponibilidad, nil
		} else {
			outputError = map[string]interface{}{"Code": "E_0458", "Body": "No existe Solicitud de Disponibilidad para La Necesidad", "Type": "error"}
			return 0, outputError
		}
	} else {
		outputError = map[string]interface{}{"Code": "E_0458", "Body": "Not enough parameter in getSolicitudDisponibilidad", "Type": "error"}
		return 0, outputError
	}
}

func getDisponibilidad(idSolicitudDisponibilidad int) (idDisponibilidad float64, outputError map[string]interface{}) {
	var solicitudDisponibilidadData []interface{}
	if idSolicitudDisponibilidad != 0 {
		if err := request.GetJson("http://"+beego.AppConfig.String("kronosService")+"disponibilidad?query=Solicitud:"+strconv.Itoa(idSolicitudDisponibilidad)+"&limit:1", &solicitudDisponibilidadData); err == nil && solicitudDisponibilidadData != nil && solicitudDisponibilidadData[0].(map[string]interface{})["Id"] != nil {
			idDisponibilidad = solicitudDisponibilidadData[0].(map[string]interface{})["Id"].(float64)
			return idDisponibilidad, nil
		} else {
			outputError = map[string]interface{}{"Code": "E_0458", "Body": "No existe Disponibilidad para La Necesidad", "Type": "error"}
			return 0, outputError
		}
	} else {
		outputError = map[string]interface{}{"Code": "E_0458", "Body": "Not enough parameter in getDisponibilidad", "Type": "error"}
		return 0, outputError
	}
}

func getRegistroPresupuestalDisponibilidadApropiacion(idDisponibilidad int) (rpDisponibilidadApropiacion []map[string]interface{}, outputError map[string]interface{}) {
	var dataSolicitudRp []interface{}
	var dataRp []interface{}
	var saldoRp map[string]float64
	if idDisponibilidad != 0 {
		if err := request.GetJson("http://"+beego.AppConfig.String("argoService")+"solicitud_rp?limit=-1&query=Expedida:true,Cdp:"+strconv.Itoa(int(idDisponibilidad)), &dataSolicitudRp); err == nil && dataSolicitudRp != nil {
			if idSolicitudRp, e := dataSolicitudRp[0].(map[string]interface{})["Id"].(float64); e {
				if err = request.GetJson("http://"+beego.AppConfig.String("kronosService")+"registro_presupuestal?limit=-1&query=Solicitud:"+strconv.Itoa(int(idSolicitudRp)), &dataRp); err == nil && dataRp != nil {
					if rpmap, e := dataRp[0].(map[string]interface{}); e {
						if desagregacionpresrp, e := rpmap["RegistroPresupuestalDisponibilidadApropiacion"].([]interface{}); e {
							for _, infopresrp := range desagregacionpresrp {
								row := make(map[string]interface{})
								if info, e := infopresrp.(map[string]interface{}); e {
									if dispoap, e := info["DisponibilidadApropiacion"].(map[string]interface{}); e {
										row["RegistroPresupuestalDisponibilidadApropiacion"] = info
										row["Rp"] = dataRp[0]
										row["Apropiacion"] = dispoap["Apropiacion"]
										row["FuenteFinanciacion"] = dispoap["FuenteFinanciamiento"]
										if err = request.SendJson("http://"+beego.AppConfig.String("kronosService")+"registro_presupuestal/SaldoRp", "POST", &saldoRp, row); err == nil && dataRp != nil {
											row["Saldo"] = saldoRp["saldo"]
										}
										rpDisponibilidadApropiacion = append(rpDisponibilidadApropiacion, row)
									}
								}
							}
							return rpDisponibilidadApropiacion, nil
						} else {
							//get data rp
							outputError = map[string]interface{}{"Code": "E_0458", "Body": "no existe Resistro Presupuestal para La Necesidad", "Type": "error"}
							return nil, outputError
						}
					} else {
						// conversion data del rp
						outputError = map[string]interface{}{"Code": "E_0458", "Body": "no existe Resistro Presupuestal para La Necesidad", "Type": "error"}
						return nil, outputError
					}
				} else {
					//get data registro presupuestal
					outputError = map[string]interface{}{"Code": "E_0458", "Body": "no existe Resistro Presupuestal para La Necesidad", "Type": "error"}
					return nil, outputError
				}
			} else {
				outputError = map[string]interface{}{"Code": "E_0458", "Body": "no existe Solicitud de Rp para La Necesidad", "Type": "error"}
				return nil, outputError
			}
		} else {
			outputError = map[string]interface{}{"Code": "E_0458", "Body": "no existe Solicitud de Rp para La Necesidad", "Type": "error"}
			return nil, outputError
		}
	} else {
		outputError = map[string]interface{}{"Code": "E_0458", "Body": "no existe Dispinibildida para La Necesidad", "Type": "error"}
		return nil, outputError
	}
}

func getPagosConDetalleLiquidacion(idPeriodoPago int) (respuestaCV []interface{}) {
	var dataPagos []interface{}
	if err := request.GetJson("http://"+beego.AppConfig.String("SsService")+"pago/?query=PeriodoPago.Id:"+strconv.Itoa(idPeriodoPago)+"&limit=-1", &dataPagos); err == nil && dataPagos != nil {
		done := make(chan interface{})
		defer close(done)
		resch := optimize.GenChanInterface(dataPagos...)
		chlistaPagos := optimize.Digest(done, getContratoVigenciaDetalleLiquidacion, resch, nil)
		for datalistaPagos := range chlistaPagos {
			if datalistaPagos != nil {
				respuestaCV = append(respuestaCV, datalistaPagos.(interface{}))
			}
		}
		return respuestaCV
	} else {
		return nil
	}
}

func homologacionConceptosSS(dataPagos interface{}, params ...interface{}) (res interface{}) {
	if dataPago, e := dataPagos.(map[string]interface{}); e {
		var infoVinculacion []interface{}
		var homologacion []interface{}
		outputConceptoHomologado := make(map[string]interface{})
		if err := request.GetJson("http://"+beego.AppConfig.String("AdministrativaAmazonService")+"vinculacion_docente?query=NumeroContrato:"+dataPago["NumeroContrato"].(string)+",Vigencia:"+strconv.FormatFloat(dataPago["VigenciaContrato"].(float64), 'f', -1, 64), &infoVinculacion); err == nil && infoVinculacion != nil {
			idFacultad, e := infoVinculacion[0].(map[string]interface{})["IdResolucion"].(map[string]interface{})["IdFacultad"].(float64)
			if !e {
				return nil
			}
			if err := request.GetJson("http://"+beego.AppConfig.String("kronosService")+"/homologacion_concepto?query=ConceptoTitan:"+strconv.Itoa(int(dataPago["TipoPago"].(float64)))+",ConceptoKronos.ConceptoTesoralFacultadProyecto.Facultad:"+strconv.Itoa(int(idFacultad))+",ConceptoKronos.ConceptoTesoralFacultadProyecto.ProyectoCurricular:0", &homologacion); err == nil && homologacion != nil {
				outputConceptoHomologado["Concepto"] = homologacion[0].(map[string]interface{})["ConceptoKronos"]
				outputConceptoHomologado["Valor"] = dataPago["Valor"]
				return outputConceptoHomologado
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

func getConceptosEnRpDisponibilidadApropiacion(listaRpDispoApropi interface{}, params ...interface{}) (res interface{}) {
	row, e := listaRpDispoApropi.(map[string]interface{})
	rubro := row["DisponibilidadApropiacion"].(map[string]interface{})["Apropiacion"].(map[string]interface{})["Rubro"].(map[string]interface{})
	var conceptos []interface{}
	if e {
		if err := request.GetJson("http://"+beego.AppConfig.String("kronosService")+"concepto/?query=Rubro:"+strconv.FormatFloat(rubro["Id"].(float64), 'f', 0, 64), &conceptos); err == nil {
			if conceptos != nil {
				row["DisponibilidadApropiacion"].(map[string]interface{})["Apropiacion"].(map[string]interface{})["Rubro"].(map[string]interface{})["Concepto"] = conceptos
				return row
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
func getMovimientosDescuentoDeLiquidacion(idLiquidacion, idNomina int) (DataMovimientoDescuento []map[string]interface{}, outputError map[string]interface{}) {
	if idLiquidacion != 0 && idNomina != 0 {
		var ordenespago []interface{}
		var allMovimientos []map[string]interface{}
		var params []interface{}
		if err := request.GetJson("http://"+beego.AppConfig.String("kronosService")+"orden_pago/?query=Liquidacion:"+strconv.Itoa(idLiquidacion)+"&limit:-1", &ordenespago); err == nil && ordenespago != nil {
			done := make(chan interface{})
			defer close(done)
			params = append(params, idNomina)
			resch := optimize.GenChanInterface(ordenespago...)
			chlistaMovimientos := optimize.Digest(done, getMovimeintosContables, resch, params)
			for dataChListaMovimientos := range chlistaMovimientos {
				if movimientosPorOrdenP, e := dataChListaMovimientos.([]interface{}); e {
					for _, movimientoOp := range movimientosPorOrdenP {
						if rowMovimiento, e := movimientoOp.(map[string]interface{}); e {
							existe := false
							for _, allM := range allMovimientos {
								if allM["CuentaContable"] != nil && rowMovimiento["CuentaContable"] != nil && allM["CuentaContable"].(map[string]interface{})["Id"].(float64) == rowMovimiento["CuentaContable"].(map[string]interface{})["Id"].(float64) {
									allM["Debito"] = allM["Debito"].(float64) + rowMovimiento["Debito"].(float64)
									allM["Credito"] = allM["Credito"].(float64) + rowMovimiento["Credito"].(float64)
									existe = true
								}
							}
							if !existe {
								allMovimientos = append(allMovimientos, rowMovimiento)
							}
						}
					}
				}
			}
			return allMovimientos, nil
		} else {
			outputError = map[string]interface{}{"Code": "E_OPM_0001", "Body": "No se encontraron ordenes de pago relacionadas a la liquidacion", "Type": "error"}
			return nil, outputError
		}
	} else {
		outputError = map[string]interface{}{"Code": "E_0458", "Body": "Not enough parameter in getMovimientosDescuentoDeLiquidacion", "Type": "error"}
		return nil, outputError
	}
}

func getMovimeintosContables(listaOrdenesPago interface{}, params ...interface{}) (res interface{}) {
	if ordenPago, e := listaOrdenesPago.(map[string]interface{}); e {
		var movimientosContables []interface{}
		var outputMovimientosContables []interface{}
		// aray de regla
		descuentosPermitidos := reglaGetDescuentosDeLiquidacion(params[0].(int))

		if err := request.GetJson("http://"+beego.AppConfig.String("kronosService")+"movimiento_contable/?query=TipoDocumentoAfectante.CodigoAbreviacion:DA-OP,CuentaEspecial__isnull:false,CodigoDocumentoAfectante:"+strconv.Itoa(int(ordenPago["Id"].(float64)))+"&limit:-1", &movimientosContables); err == nil && movimientosContables != nil {
			for _, movimientoContable := range movimientosContables {
				if rowMovimientoC, e := movimientoContable.(map[string]interface{}); e {
					for _, descuento := range descuentosPermitidos {
						if rowDescuento, e := descuento.(map[string]interface{}); e {
							if rowDescuento["Id"].(float64) == rowMovimientoC["CuentaEspecial"].(map[string]interface{})["Id"].(float64) {
								outputMovimientosContables = append(outputMovimientosContables, rowMovimientoC)
							}
						}
					}
				}
			}
			return outputMovimientosContables
		} else {
			return nil
		}
	} else {
		return nil
	}
}

func reglaGetDescuentosDeLiquidacion(idNomina int) (DataDescuentos []interface{}) {
	var nomina []interface{}
	var idDescuentos [3]int
	var descuento []interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("titanService")+"nomina/?query=Id:"+strconv.Itoa(idNomina), &nomina); err == nil && nomina != nil {
		if nominaName, e := nomina[0].(map[string]interface{})["TipoNomina"].(map[string]interface{})["Nombre"]; e {
			if nominaName == "HCH" {
				idDescuentos = [3]int{40, 42, 41} //titan descuentos 239, 290, 291
			} else {
				idDescuentos = [3]int{43, 45, 44}
			}
			for _, id := range idDescuentos {
				if err := request.GetJson("http://"+beego.AppConfig.String("kronosService")+"cuenta_especial/?query=Id:"+strconv.Itoa(id)+"&limit:1", &descuento); err == nil && descuento != nil {
					DataDescuentos = append(DataDescuentos, descuento[0])
				}
			}
		} else {
			return nil
		}
	} else {
		return nil
	}
	return
}

func reglaGetCuentaAfectarPorDescuento(codigoCuentaEspecial string) (codigoCuentaAfectar string, outputError map[string]interface{}) {
	if codigoCuentaEspecial != "" {
		var joinCuentaEspecial string
		splitCuentaEspecial := strings.Split(codigoCuentaEspecial, "-")

		if splitCuentaEspecial[len(splitCuentaEspecial)-1] == "02" {
			splitCuentaEspecial[len(splitCuentaEspecial)-1] = "01"
			joinCuentaEspecial = strings.Join(splitCuentaEspecial, "-")
		} else if splitCuentaEspecial[len(splitCuentaEspecial)-1] == "01" {
			splitCuentaEspecial[len(splitCuentaEspecial)-1] = "02"
			joinCuentaEspecial = strings.Join(splitCuentaEspecial, "-")
		} else {
			outputError = map[string]interface{}{"Code": "E_0458", "Body": "No se especificaron mas casos para afectar la cuenta del descuento", "Type": "error"}
			return "", outputError
		}
		return joinCuentaEspecial, nil
	}
	outputError = map[string]interface{}{"Code": "E_0458", "Body": "Not enough parameter in reglaGetCuentaAfectarPorDescuento", "Type": "error"}
	return "", outputError
}

func afectarDescuentosToMovimientos(cuentasDescuentoOP, allMovimientos []map[string]interface{}) (totalCuentas interface{}, allAfectacionConceto []map[string]interface{}, outputError map[string]interface{}) {
	if cuentasDescuentoOP != nil && allMovimientos != nil {
		for _, cuentaOP := range cuentasDescuentoOP {
			if codigoDescuentoOP, e := cuentaOP["CuentaContable"].(map[string]interface{})["Codigo"].(string); e {
				if CodigoDescuentoBuscar, e := reglaGetCuentaAfectarPorDescuento(codigoDescuentoOP); e == nil {
					rowDescuento := make(map[string]interface{})
					//rowDescuento["Concepto"] = cuentaOP["Concepto"]
					rowDescuento["Debito"] = cuentaOP["Credito"] //como son descuentos solo titnen credito
					rowDescuento["Credito"] = 0
					rowDescuento["CuentaContable"] = cuentaOP["CuentaContable"]
					rowDescuento["CuentaEspecial"] = cuentaOP["CuentaEspecial"]
					for _, movimineto := range allMovimientos {
						if codeCuentaC, e := movimineto["CuentaContable"].(map[string]interface{}); e {
							if codeCuentaC["Codigo"].(string) == CodigoDescuentoBuscar {
								rowDescuento["Concepto"] = movimineto["Concepto"].(map[string]interface{})
								movimineto["Credito"] = movimineto["Credito"].(float64) + cuentaOP["Credito"].(float64)
								allMovimientos = append(allMovimientos, rowDescuento)
								allAfectacionConceto = append(allAfectacionConceto, rowDescuento)
							}
						} else {
							outputError = map[string]interface{}{"Code": "E_0458", "Body": "Error En la estructura de datos del parametro allMovimientos en  afectarDescuentosToMovimientos", "Type": "error"}
							return nil, nil, outputError
						}
					}
				} else {
					return nil, nil, e
				}
			} else {
				outputError = map[string]interface{}{"Code": "E_0458", "Body": "Error En la estructura de datos en el  parametro cuentasDescuentoOP de  afectarDescuentosToMovimientos", "Type": "error"}
				return nil, nil, outputError
			}
		}
		return allMovimientos, allAfectacionConceto, nil
	}
	outputError = map[string]interface{}{"Code": "E_0458", "Body": "Not enough parameter in afectarDescuentosToMovimientos", "Type": "error"}
	return nil, nil, outputError
}

func afectarDescuentosToConceptos(ConceptosDescuentos, allConceptos []map[string]interface{}) (conceptos []map[string]interface{}, outputError map[string]interface{}) {
	if ConceptosDescuentos != nil && allConceptos != nil {
		for _, conceptoDesc := range ConceptosDescuentos {
			if codigoAfectado, e := conceptoDesc["Concepto"].(map[string]interface{})["Codigo"].(string); e {
				for _, concepto := range allConceptos {
					if codigoConcepto, e := concepto["Concepto"].(map[string]interface{})["Codigo"].(string); e {
						if codigoAfectado == codigoConcepto {
							concepto["Valor"] = concepto["Valor"].(float64) + conceptoDesc["Debito"].(float64)
						}
					} else {
						outputError = map[string]interface{}{"Code": "E_0458", "Body": "Error En la estructura de datos del parametro allConceptos en afectarDescuentosToConceptos", "Type": "error"}
						return nil, outputError
					}
				}
			} else {
				outputError = map[string]interface{}{"Code": "E_0458", "Body": "Error En la estructura de datos del parametro ConceptosDescuentos en afectarDescuentosToConceptos", "Type": "error"}
				return nil, outputError
			}
		}
		return allConceptos, nil
	}
	outputError = map[string]interface{}{"Code": "E_0458", "Body": "Not enough parameter in afectarDescuentosToConceptos", "Type": "error"}
	return nil, outputError
}

func getTotalAfectacionOfConceptos(allConceptos []map[string]interface{}) (total float64, outputError map[string]interface{}) {
	if allConceptos != nil {
		var acumulador float64
		for _, concepto := range allConceptos {
			acumulador = acumulador + concepto["Valor"].(float64)
		}
		return acumulador, nil
	}
	outputError = map[string]interface{}{"Code": "E_0458", "Body": "Not enough parameter in getTotalAfectacionOfConceptos", "Type": "error"}
	return 0, outputError
}
