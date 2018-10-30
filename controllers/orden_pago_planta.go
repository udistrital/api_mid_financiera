package controllers

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/api_mid_financiera/models"
	"github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/utils_oas/optimize"
	"github.com/udistrital/utils_oas/request"
)

// OrdenPagoPlantaController operations for OrdenPagoPlanta
type OrdenPagoPlantaController struct {
	beego.Controller
}

// URLMapping ...
func (c *OrdenPagoPlantaController) URLMapping() {

}

//funcion para generar el formato de la op planta para su vista previa.
func formatoResumenOpPlanta(dataLiquidacion interface{}, params ...interface{}) (res interface{}) {
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
		var homologacionDescuentos []map[string]interface{}
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
			resch := optimize.GenChanInterface(devengosNomina...)
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
				chConcHomologados := optimize.Digest(done, f, resch, params)
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
		var descuentosNomina []interface{}
		if descuentosNomina, err = ConsultarDescuentosNominaPorContrato(idLiquidacion, nContrato, vigenciaContrato); err != nil {
			//si no se consiguen datos de devengos para el contrato liquidado (nomina planta)
			return nil
		}
		//homlogacion de los descuentos de la nomina de planta a conceptos de kronos...
		done := make(chan interface{})
		defer close(done)
		resch := optimize.GenChanInterface(descuentosNomina...)
		chDescHomologados := optimize.Digest(done, homologacionDescuentosHC, resch, nil)
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
		}
		//-----------------------------------------------------------------------------
		var movimientosContables []interface{}
		for _, concepto := range homologacionConceptos {
			movimientoContable := formatoMovimientosContablesOp(concepto)
			for _, aux := range movimientoContable {
				movimientosContables = append(movimientosContables, aux)
			}

		}
		movcredito := findMovimientoCredito(movimientosContables)
		for _, descuento := range homologacionDescuentos {
			if movimientosContables != nil {

				movimientoContable, mov := formatoMovimientosContablesDescuentosOp(descuento, movcredito)
				movcredito = mov
				beego.Info(movcredito.(map[string]interface{})["Credito"])
				for _, aux := range movimientoContable {

					movimientosContables = append(movimientosContables, aux)
				}

			}

		}
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
		if auxmap, e := infoContrato.(map[string]interface{}); e && auxmap["InfoPersona"] != nil {
			res["infoPersona"], e = auxmap["infoPersona"]
		} else {
			res["infoPersona"] = make(map[string]interface{})
			res["Aprobado"] = false
			res["Code"] = "OPM_E005"
		}
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

		res["Contrato"] = nContrato
		res["VigenciaContrato"] = vigenciaContrato
		res["Aprobado"] = false
		res["Code"] = "OPM_E003"
		if auxmap, e := infoContrato.(map[string]interface{}); e {
			res["infoPersona"], e = auxmap["infoPersona"]
		} else {
			res["infoPersona"] = make(map[string]interface{})
			res["Aprobado"] = false
			res["Code"] = "OPM_E005"
		}
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
		if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/homologacion_concepto?query=ConceptoTitan:"+strconv.Itoa(int(dataConceptoAhomologar["Concepto"].(map[string]interface{})["Id"].(float64)))+",Vigencia:"+strconv.FormatFloat(vigContrato, 'f', -1, 64), &homologacion); err == nil {
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

func formatoPreViewCargueMasivoOpPlanta(liquidacion interface{}, params ...interface{}) (res interface{}) {
	var respuesta []interface{}
	//var rp []map[string]interface{}
	done := make(chan interface{})
	defer close(done)
	_, e := liquidacion.(map[string]interface{})
	if e {
		if liquidacion.(map[string]interface{})["Contratos_por_preliq"] != nil {
			listaLiquidacion := liquidacion.(map[string]interface{})["Contratos_por_preliq"].([]interface{})
			resch := optimize.GenChanInterface(listaLiquidacion...)
			var params2 []interface{}

			params2 = append(params2, liquidacion.(map[string]interface{})["Id_Preliq"].(interface{}))
			//rp, err := GetRpDesdeNecesidadProcesoExternoGeneral(23, "S")
			rp, err := GetRpDesdeNecesidadProcesoExternoGeneral(liquidacion.(map[string]interface{})["Id_Preliq"].(float64), "N")
			beego.Info(err)
			if rp != nil {
				params2 = append(params2, rp)
			} else {
				res = map[string]interface{}{"Code": "E_0458", "Body": err, "Type": "error"}
			}
			f := formatoRegistroOpFunctionDispatcher(liquidacion.(map[string]interface{})["Nombre_tipo_nomina"].(string))
			//beego.Info(liquidacion.(map[string]interface{})["Nombre_tipo_nomina"].(string))
			if f != nil {

				chlistaLiquidacion := optimize.Digest(done, f, resch, params2)
				for dataLiquidacion := range chlistaLiquidacion {
					if dataLiquidacion != nil {
						respuesta = append(respuesta, dataLiquidacion)
					}
				}
				fresumen := resumenOpFunctionDispatcher(liquidacion.(map[string]interface{})["Nombre_tipo_nomina"].(string))
				var resultado interface{}
				if fresumen != nil {
					resultado = fresumen(respuesta)
				}
				OrdenPagoReg := formatoRegistroOpPlanta(respuesta, rp, liquidacion.(map[string]interface{})["Id_Preliq"].(float64))
				res = map[string]interface{}{"DetalleCargueOp": respuesta, "ResumenCargueOp": resultado, "TipoLiquidacion": liquidacion.(map[string]interface{})["Nombre_tipo_nomina"].(string), "OrdenPagoaRegistrar": OrdenPagoReg}
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

	return
}

func formatoRegistroOpPlanta(detalleOP []interface{}, rpForm []map[string]interface{}, idLiquidacion float64) (res interface{}) {
	comp := false
	code := "OPM_S001"
	var valorBase float64
	var conceptoOp []map[string]interface{}
	var MovOp []map[string]interface{}
	acumConceptOp := make(map[float64]map[float64]map[string]interface{})          //acumulacion de estructuras por  idConcepto, IdRegPresDispApr
	acumMovsOp := make(map[float64]map[float64]map[float64]map[string]interface{}) //acumulacion de estructuras por  idConcepto, IdCuentaContable, idCuentaEspecial (0 si no aplica)
	var op interface{}
	if rpForm != nil {
		if rpint, e := rpForm[0]["Rp"].(interface{}); e {
			ordenPago := make(map[string]interface{})
			ordenPago["RegistroPresupuestal"] = rpint
			ordenPago["Liquidacion"] = idLiquidacion
			op = ordenPago
		}
		for _, auxDetalle := range detalleOP {

			if auxMap, e := auxDetalle.(map[string]interface{}); e && auxMap["ConceptoOrdenPago"] != nil && auxMap["MovimientoContable"] != nil && auxMap["Aprobado"].(bool) {
				beego.Info(auxMap["Liquidacion"]) //detalleMap = append(detalleMap, auxMap)
				ConcsOpMap, e := auxMap["ConceptoOrdenPago"].([]map[string]interface{})
				var MovsMap []map[string]interface{}
				err := formatdata.FillStruct(auxMap["MovimientoContable"], &MovsMap)
				if e && err == nil {
					//if e {

					for _, concepOP := range ConcsOpMap {
						//beego.Info(concepOP["Concepto"].(map[string]interface{})["Id"].(float64), concepOP["RegistroPresupuestalDisponibilidadApropiacion"].(map[string]interface{})["Id"].(float64))
						if acumConceptOp[concepOP["Concepto"].(map[string]interface{})["Id"].(float64)][concepOP["RegistroPresupuestalDisponibilidadApropiacion"].(map[string]interface{})["Id"].(float64)] != nil {

							acumConceptOp[concepOP["Concepto"].(map[string]interface{})["Id"].(float64)][concepOP["RegistroPresupuestalDisponibilidadApropiacion"].(map[string]interface{})["Id"].(float64)] = map[string]interface{}{"Concepto": concepOP["Concepto"], "RegistroPresupuestalDisponibilidadApropiacion": concepOP["RegistroPresupuestalDisponibilidadApropiacion"], "Valor": acumConceptOp[concepOP["Concepto"].(map[string]interface{})["Id"].(float64)][concepOP["RegistroPresupuestalDisponibilidadApropiacion"].(map[string]interface{})["Id"].(float64)]["Valor"].(float64) + concepOP["Valor"].(float64)}
						} else {

							if acumConceptOp[concepOP["Concepto"].(map[string]interface{})["Id"].(float64)] == nil {
								acumConceptOp[concepOP["Concepto"].(map[string]interface{})["Id"].(float64)] = make(map[float64]map[string]interface{})
							}
							acumConceptOp[concepOP["Concepto"].(map[string]interface{})["Id"].(float64)][concepOP["RegistroPresupuestalDisponibilidadApropiacion"].(map[string]interface{})["Id"].(float64)] = map[string]interface{}{"Concepto": concepOP["Concepto"], "RegistroPresupuestalDisponibilidadApropiacion": concepOP["RegistroPresupuestalDisponibilidadApropiacion"], "Valor": concepOP["Valor"].(float64)}

						}

					}

					for _, movOp := range MovsMap {
						var idCuentaEspecial float64
						if movOp["CuentaEspecial"] != nil {
							idCuentaEspecial = movOp["CuentaEspecial"].(map[string]interface{})["Id"].(float64)
						}
						if acumMovsOp[movOp["Concepto"].(map[string]interface{})["Id"].(float64)][movOp["CuentaContable"].(map[string]interface{})["Id"].(float64)][idCuentaEspecial] != nil {
							movOp["Debito"] = movOp["Debito"].(float64) + acumMovsOp[movOp["Concepto"].(map[string]interface{})["Id"].(float64)][movOp["CuentaContable"].(map[string]interface{})["Id"].(float64)][idCuentaEspecial]["Debito"].(float64)
							movOp["Credito"] = movOp["Credito"].(float64) + acumMovsOp[movOp["Concepto"].(map[string]interface{})["Id"].(float64)][movOp["CuentaContable"].(map[string]interface{})["Id"].(float64)][idCuentaEspecial]["Credito"].(float64)
							//beego.Info(movOp["Credito"])
						} else {
							if acumMovsOp[movOp["Concepto"].(map[string]interface{})["Id"].(float64)] == nil {
								acumMovsOp[movOp["Concepto"].(map[string]interface{})["Id"].(float64)] = make(map[float64]map[float64]map[string]interface{})
								acumMovsOp[movOp["Concepto"].(map[string]interface{})["Id"].(float64)][movOp["CuentaContable"].(map[string]interface{})["Id"].(float64)] = make(map[float64]map[string]interface{})
								acumMovsOp[movOp["Concepto"].(map[string]interface{})["Id"].(float64)][movOp["CuentaContable"].(map[string]interface{})["Id"].(float64)][idCuentaEspecial] = make(map[string]interface{})

							} else {
								if acumMovsOp[movOp["Concepto"].(map[string]interface{})["Id"].(float64)][movOp["CuentaContable"].(map[string]interface{})["Id"].(float64)] == nil {
									acumMovsOp[movOp["Concepto"].(map[string]interface{})["Id"].(float64)][movOp["CuentaContable"].(map[string]interface{})["Id"].(float64)] = make(map[float64]map[string]interface{})
									acumMovsOp[movOp["Concepto"].(map[string]interface{})["Id"].(float64)][movOp["CuentaContable"].(map[string]interface{})["Id"].(float64)][idCuentaEspecial] = make(map[string]interface{})
								} else {
									acumMovsOp[movOp["Concepto"].(map[string]interface{})["Id"].(float64)][movOp["CuentaContable"].(map[string]interface{})["Id"].(float64)][idCuentaEspecial] = make(map[string]interface{})
								}
							}
						}
						acumMovsOp[movOp["Concepto"].(map[string]interface{})["Id"].(float64)][movOp["CuentaContable"].(map[string]interface{})["Id"].(float64)][idCuentaEspecial] = movOp

					}

				}
			}

		}

		for _, indexIdConc := range acumConceptOp {
			for _, indexIdRpDpAp := range indexIdConc {
				beego.Info(indexIdRpDpAp["Valor"])
				valorBase = valorBase + indexIdRpDpAp["Valor"].(float64)
				conceptoOp = append(conceptoOp, indexIdRpDpAp)
			}
		}

		for _, indexIdConc := range acumMovsOp {
			for _, indexIdCuentaCont := range indexIdConc {
				for _, indexCuentaEsp := range indexIdCuentaCont {
					MovOp = append(MovOp, indexCuentaEsp)
				}
			}
		}
		if acumConceptOp != nil {
			for _, desgrRp := range rpForm {
				if rpdispaprMap, e := desgrRp["RegistroPresupuestalDisponibilidadApropiacion"].(map[string]interface{}); e && rpdispaprMap["Id"] != nil {
					id, _ := rpdispaprMap["Id"].(float64)
					var valorAcum float64
					for _, indexIdConc := range acumConceptOp {
						if indexIdConc[id] != nil {
							valorAcum = valorAcum + indexIdConc[id]["Valor"].(float64)
						}
					}
					if valorAcum <= 0 {
						code = "OPM_E001"
						comp = false

					} else if valorAcum <= desgrRp["Saldo"].(float64) {
						comp = true
						code = "OPM_S001"
					} else {
						code = "OPM_E002"
						comp = false
					}
				}
			}
		} else {
			code = "OPM_E003"
			comp = false
		}
		if acumMovsOp == nil {
			code = "OPM_E004"
			comp = false
		}
	} else {
		code = "OPM_E003"
		comp = false
	}

	return map[string]interface{}{"ConceptoOrdenPago": conceptoOp, "MovimientoContable": MovOp, "OrdenPago": op, "Aprobado": comp, "Code": code, "ValorBase": valorBase}
}

// RegistroOpPlanta ... 
func RegistroOpPlanta(datain map[string]interface{}, params ...interface{}) (res interface{}) {
	//"http://"+beego.AppConfig.String("kronosService")+
	data, _ := datain["OrdenPagoaRegistrar"].(interface{})
	alerts := []models.Alert{}
	alert := models.Alert{}
	if auxmap, e := data.(map[string]interface{}); e {
		if auxbool, e := auxmap["Aprobado"].(bool); e {
			if auxbool {
				valorBase, e2 := auxmap["ValorBase"].(float64)
				if Opmap, e := auxmap["OrdenPago"].(map[string]interface{}); e && e2 {
					Opmap = params[0].([]interface{})[0].(map[string]interface{})
					Opmap["RegistroPresupuestal"] = auxmap["OrdenPago"].(map[string]interface{})["RegistroPresupuestal"]
					Opmap["Liquidacion"] = auxmap["OrdenPago"].(map[string]interface{})["Liquidacion"]
					Opmap["SubTipoOrdenPago"], e = params[0].([]interface{})[0].(map[string]interface{})["SubTipoOrdenPago"]
					Opmap["Documento"], e = params[0].([]interface{})[0].(map[string]interface{})["Documento"]
					Opmap["ValorBase"] = valorBase
					auxmap["OrdenPago"] = Opmap
					beego.Info(Opmap)
					if err := request.SendJson("http://"+beego.AppConfig.String("kronosService")+"orden_pago/RegistrarOpProveedor", "POST", &alert, &auxmap); err == nil {
						alerts = append(alerts, alert)
					} else {
						alerts = append(alerts, models.Alert{Code: "E_0458", Body: data, Type: "error"})
					}
				} else {
					beego.Info(e, e2)
					alerts = append(alerts, models.Alert{Code: "E_0458", Body: "Conversion Problem", Type: "error"})
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

	return alerts
}
