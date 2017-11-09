package controllers

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/api_mid_financiera/models"
	"github.com/udistrital/api_mid_financiera/utilidades"
)

// OrdenPagoSsController operations for Orden_pago_seguridad_social
type OrdenPagoSsController struct {
	beego.Controller
}

// URLMapping ...
func (c *OrdenPagoSsController) URLMapping() {
	c.Mapping("ListaPagoSsPorPersona", c.ListaPagoSsPorPersona)

}

// DetalleListaPagoSsPorPersona ...
// @Title DetalleListaPagoSsPorPersona
// @Description lista pagos de seguridad socila por persona.
// @Param	idPeriodoPago	query	string	false	"nomina a listar"
// @Param	idDetalleLiquidacion	query	string	false	"mes de la liquidacion a listar"
// @Success 201 {object} models.Alert
// @Failure 403 body is empty
// @router /DetalleListaPagoSsPorPersona [get]
func (c *OrdenPagoSsController) DetalleListaPagoSsPorPersona() {
	idPeriodoPago, err1 := c.GetInt("idPeriodoPago")
	idDetalleLiquidacion, err2 := c.GetInt("idDetalleLiquidacion")
	if err1 == nil && err2 == nil {
		var listaPagos []interface{}
		var allData []interface{}
		if err := getJson("http://"+beego.AppConfig.String("SsService")+"pago/?query=DetalleLiquidacion:"+strconv.Itoa(idDetalleLiquidacion)+"&PeriodoPago.Id:"+strconv.Itoa(idPeriodoPago), &listaPagos); err == nil {
			if listaPagos != nil {
				done := make(chan interface{})
				defer close(done)
				resch := utilidades.GenChanInterface(listaPagos...)
				chlistaPagos := utilidades.Digest(done, getTipoPagoTitan, resch, nil)
				for dataChlistaPagos := range chlistaPagos {
					if dataChlistaPagos != nil {
						allData = append(allData, dataChlistaPagos.(interface{}))
					}
				}
				c.Data["json"] = allData
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

func getTipoPagoTitan(listaPagos interface{}, params ...interface{}) (res interface{}) {
	row, e := listaPagos.(map[string]interface{})
	var tipoPago []interface{}
	if e {
		if err := getJson("http://"+beego.AppConfig.String("titanService")+"concepto_nomina/?query=Id:"+strconv.FormatFloat(row["TipoPago"].(float64), 'f', 0, 64)+"&limit=1", &tipoPago); err == nil {
			row["TipoPago"] = tipoPago[0].(map[string]interface{})
			return row
		} else {
			return nil
		}
	} else {
		return nil
	}
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
func (c *OrdenPagoSsController) ListaPagoSsPorPersona() {
	idNomina, err1 := c.GetInt("idNomina")
	mesLiquidacion, err2 := c.GetInt("mesLiquidacion")
	anioLiquidacion, err3 := c.GetInt("anioLiquidacion")
	if err1 == nil && err2 == nil && err3 == nil {
		var respuestaCV []interface{}
		pagosAgrupados := pagoSsPorPersonaF(idNomina, mesLiquidacion, anioLiquidacion)
		if pagosAgrupados != nil {
			done := make(chan interface{})
			defer close(done)
			if pagosAgrupados.(map[string]interface{})["Pagos"] != nil {
				listaPagos := pagosAgrupados.(map[string]interface{})["Pagos"].([]interface{})
				fmt.Println(listaPagos)
				resch := utilidades.GenChanInterface(listaPagos...)
				chlistaPagos := utilidades.Digest(done, getContratoVigenciaDetalleLiquidacion, resch, nil)
				for datalistaPagos := range chlistaPagos {
					if datalistaPagos != nil {
						respuestaCV = append(respuestaCV, datalistaPagos.(interface{}))
					}
				}
				res := pagosAgrupados.(map[string]interface{})
				res["Pagos"] = respuestaCV
				c.Data["json"] = res
			} else {
				c.Data["json"] = models.Alert{Code: "E_0458", Body: nil, Type: "error"}
			}
		} else {
			c.Data["json"] = nil
		}
	} else {
		c.Data["json"] = models.Alert{Code: "E_0458", Body: "Not enough parameter", Type: "error"}
	}
	c.ServeJSON()
}

func getContratoVigenciaDetalleLiquidacion(idsLiquidacionDesdePagos interface{}, params ...interface{}) (res interface{}) {
	row, e := idsLiquidacionDesdePagos.(map[string]interface{})
	var infoDetallePreliquidacion []interface{}
	var infoPersona interface{}
	if e {
		if err := getJson("http://"+beego.AppConfig.String("titanService")+"detalle_preliquidacion/?query=Id:"+strconv.FormatFloat(row["DetalleLiquidacion"].(float64), 'f', 0, 64)+"&limit=1", &infoDetallePreliquidacion); err == nil {
			row["NumeroContrato"], e = infoDetallePreliquidacion[0].(map[string]interface{})["NumeroContrato"]
			row["VigenciaContrato"], e = infoDetallePreliquidacion[0].(map[string]interface{})["VigenciaContrato"]
			if e {
				if err := getJsonWSO2("http://jbpm.udistritaloas.edu.co:8280/services/contrato_suscrito_DataService.HTTPEndpoint/informacion_contrato_elaborado_contratista/"+row["NumeroContrato"].(string)+"/"+strconv.Itoa(int(row["VigenciaContrato"].(float64))), &infoPersona); err == nil {
					row["infoPersona"], e = infoPersona.(map[string]interface{})["informacion_contratista"]
					if e {
						return row
					} else {
						return nil
					}
				} else {
					return
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

// Getjota ...
// @Title Getjota
// @Description lista pagos de seguridad socila por persona.
// @Param	idNomina	query	string	false	"nomina a listar"
// @Param	mesLiquidacion	query	string	false	"mes de la liquidacion a listar"
// @Param	anioLiquidacion	query	string	false	"anio de la liquidacion a listar"
// @Success 201 {object} models.Alert
// @Failure 403 body is empty
// @router /Getjota [get]
func (c *OrdenPagoSsController) Getjota() {
	idNomina, err1 := c.GetInt("idNomina")
	mesLiquidacion, err2 := c.GetInt("mesLiquidacion")
	anioLiquidacion, err3 := c.GetInt("anioLiquidacion")
	if err1 == nil && err2 == nil && err3 == nil {
		idLiquidacion := getIdliquidacionForSs(idNomina, mesLiquidacion, anioLiquidacion)
		fmt.Println("Liquidacion: ", idLiquidacion)
		if idLiquidacion != 0 {
			idPeriodoPago := getIdPeriodoPagoForSs(int(idLiquidacion), mesLiquidacion, anioLiquidacion)
			fmt.Println("periodo pago ", idPeriodoPago)
			if idPeriodoPago != 0 {
				idNecesidad := getNecesidadByProcesoExternoSS(int(idPeriodoPago))
				fmt.Println("Nececidad id", idNecesidad)
				if idNecesidad != 0 {
					solicitudCDP := getSolicitudDisponibilidad(int(idNecesidad))
					fmt.Println("solicitud id, ", solicitudCDP)
					if solicitudCDP != 0 {
						disponibilidad := getDisponibilidad(int(solicitudCDP))
						fmt.Println("disponibilidad", disponibilidad)
						if disponibilidad != 0 {
							rpDisponibilidadApropiacion := getRegistroPresupuestalDisponibilidadApropiacion(int(disponibilidad))
							if rpDisponibilidadApropiacion != nil {
								fmt.Println("RPDATA ", rpDisponibilidadApropiacion)
							} else {
								c.Data["json"] = models.Alert{Code: "E_0458", Body: "no existe Registro Presupuestal para La Necesidad", Type: "error"}
							}
						} else {
							c.Data["json"] = models.Alert{Code: "E_0458", Body: "no existe Disponibilidad para La Necesidad", Type: "error"}
						}
					} else {
						c.Data["json"] = models.Alert{Code: "E_0458", Body: "no existe Solicitud de Disponibilidad para La Necesidad", Type: "error"}
					}
				} else {
					c.Data["json"] = models.Alert{Code: "E_0458", Body: "no existe necesidad para liquidacion de Seguridad Social en el periodo", Type: "error"}
				}
			} else {
				c.Data["json"] = models.Alert{Code: "E_0458", Body: "no existe periodo pago de Seguridad Social para el periodo", Type: "error"}
			}
		} else {
			c.Data["json"] = models.Alert{Code: "E_0458", Body: "no existe liquidacion en estado EnOrdenPago para el periodo", Type: "error"}
		}
	} else {
		c.Data["json"] = models.Alert{Code: "E_0458", Body: "Not enough parameter", Type: "error"}
	}
	fmt.Println("fin")
	c.ServeJSON()
}

// se consulta servicio que retorna las liquidacions en un mes, año y titpo nomina que ya esten en estado EnOrdenPago
func getIdliquidacionForSs(idNomina, mesLiquidacion, anioLiquidacion int) (IdLiquidacion float64) {
	var liquidacion interface{}
	if idNomina != 0 && mesLiquidacion != 0 && anioLiquidacion != 0 {
		if err := getJson("http://"+beego.AppConfig.String("titanService")+"preliquidacion/contratos_x_preliquidacion?idNomina="+strconv.Itoa(idNomina)+"&mesLiquidacion="+strconv.Itoa(mesLiquidacion)+"&anioLiquidacion="+strconv.Itoa(anioLiquidacion), &liquidacion); err == nil {
			if liquidacion != nil && liquidacion.(map[string]interface{})["Id_Preliq"] != nil {
				IdLiquidacion = liquidacion.(map[string]interface{})["Id_Preliq"].(float64)
			} else {
				IdLiquidacion = 0
			}
		} else {
			return 0
		}
	} else {
		return 0
	}
	return
}

// se consulta servicio de periodo_pago en un mes, año y con id liquidacion
func getIdPeriodoPagoForSs(idLiquidacion, mesLiquidacion, anioLiquidacion int) (idPeriodoPago float64) {
	var periodoPago []interface{}
	if idLiquidacion != 0 && mesLiquidacion != 0 && anioLiquidacion != 0 {
		if err := getJson("http://"+beego.AppConfig.String("SsService")+"periodo_pago/?query=Mes:"+strconv.Itoa(mesLiquidacion)+"&Anio:"+strconv.Itoa(anioLiquidacion)+"&Liquidacion:"+strconv.Itoa(idLiquidacion)+"&limit:1", &periodoPago); err == nil {
			if periodoPago != nil && periodoPago[0].(map[string]interface{})["Id"] != nil {
				idPeriodoPago = periodoPago[0].(map[string]interface{})["Id"].(float64)
			} else {
				idPeriodoPago = 0
			}
		} else {
			idPeriodoPago = 0
		}
	} else {
		idPeriodoPago = 0
	}
	return
}

func getNecesidadByProcesoExternoSS(idPeriodoPagoSs int) (necesidad float64) {
	var necesidadProcesoExterno []interface{}
	if idPeriodoPagoSs != 0 {
		//TipoNecesidad.CodigoAbreviacion:S  seguridad social
		// Necesidad.EstadoNecesidad.CodigoAbreviacion:C  => Solicitud de CDP creada
		if err := getJson("http://"+beego.AppConfig.String("argoServiceFlayway")+"necesidad_proceso_externo?query=TipoNecesidad.CodigoAbreviacion:S,ProcesoExterno:"+strconv.Itoa(idPeriodoPagoSs)+",Necesidad.EstadoNecesidad.CodigoAbreviacion:C&limit:1", &necesidadProcesoExterno); err == nil {
			if necesidadProcesoExterno != nil && necesidadProcesoExterno[0].(map[string]interface{})["Necesidad"].(map[string]interface{})["Id"] != nil {
				necesidad = necesidadProcesoExterno[0].(map[string]interface{})["Necesidad"].(map[string]interface{})["Id"].(float64)
			} else {
				necesidad = 0
			}
		} else {
			necesidad = 0
		}
	} else {
		necesidad = 0
	}
	return
}

func getSolicitudDisponibilidad(idNecesidad int) (solicitudDisponibilidad float64) {
	var solicitudDisponibilidadData []interface{}
	if idNecesidad != 0 {
		if err := getJson("http://"+beego.AppConfig.String("argoServiceFlayway")+"solicitud_disponibilidad?query=Expedida:true,Necesidad.Id:"+strconv.Itoa(idNecesidad), &solicitudDisponibilidadData); err == nil {
			if solicitudDisponibilidadData != nil && solicitudDisponibilidadData[0].(map[string]interface{})["Id"] != nil {
				solicitudDisponibilidad = solicitudDisponibilidadData[0].(map[string]interface{})["Id"].(float64)
			} else {
				solicitudDisponibilidad = 0
			}
		} else {
			solicitudDisponibilidad = 0
		}
	} else {
		solicitudDisponibilidad = 0
	}
	return
}

func getDisponibilidad(idSolicitudDisponibilidad int) (idDisponibilidad float64) {
	var solicitudDisponibilidadData []interface{}
	if idSolicitudDisponibilidad != 0 {
		if err := getJson("http://"+beego.AppConfig.String("kronosService")+"disponibilidad?query=Solicitud:"+strconv.Itoa(idSolicitudDisponibilidad)+"&limit:1", &solicitudDisponibilidadData); err == nil {
			if solicitudDisponibilidadData != nil && solicitudDisponibilidadData[0].(map[string]interface{})["Id"] != nil {
				idDisponibilidad = solicitudDisponibilidadData[0].(map[string]interface{})["Id"].(float64)
			} else {
				idDisponibilidad = 0
			}
		} else {
			idDisponibilidad = 0
		}
	} else {
		idDisponibilidad = 0
	}
	return
}

func getRegistroPresupuestalDisponibilidadApropiacion(idDisponibilidad int) (RpDisponibilidadApropiacion interface{}) {
	var dataRpDisponibilidadApropiacio interface{}
	if idDisponibilidad != 0 {
		if err := getJson("http://"+beego.AppConfig.String("kronosService")+"registro_presupuestal_disponibilidad_apropiacion?query=DisponibilidadApropiacion.Disponibilidad.Id:"+strconv.Itoa(idDisponibilidad), &dataRpDisponibilidadApropiacio); err == nil {
			if dataRpDisponibilidadApropiacio != nil {
				RpDisponibilidadApropiacion = dataRpDisponibilidadApropiacio
			} else {
				fmt.Println("1------------")
				RpDisponibilidadApropiacion = 0
			}
		} else {
			fmt.Println("2------------")
			RpDisponibilidadApropiacion = 0
		}
	} else {
		RpDisponibilidadApropiacion = 0
	}
	return
}
