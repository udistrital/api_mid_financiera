package controllers

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"

	"github.com/udistrital/utils_oas/optimize"
	"github.com/udistrital/utils_oas/request"
)

// Orden_pago_hcController operations for Orden_pago_hc
type Orden_pago_hcController struct {
	beego.Controller
}

// URLMapping ...
func (c *Orden_pago_hcController) URLMapping() {

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
		if err := request.GetJson("http://"+beego.AppConfig.String("AdministrativaAmazonService")+"vinculacion_docente?query=NumeroContrato:"+numContrato+",Vigencia:"+strconv.FormatFloat(vigContrato, 'f', -1, 64), &infoVinculacion); err == nil {
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
				if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/homologacion_concepto?query=ConceptoTitan:"+strconv.Itoa(int(dataConceptoAhomologar["Concepto"].(map[string]interface{})["Id"].(float64)))+",ConceptoKronos.ConceptoTesoralFacultadProyecto.Facultad:"+strconv.Itoa(int(idFacultad))+",ConceptoKronos.ConceptoTesoralFacultadProyecto.ProyectoCurricular:"+strconv.Itoa(int(idProyecto))+",Vigencia:"+strconv.FormatFloat(vigContrato, 'f', -1, 64), &homologacion); err == nil {
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

func homologacionDescuentosHC(dataDescuento interface{}, params ...interface{}) (res interface{}) {
	//beego.Info(dataDescuento)
	dataDescuentoAhomologar, e := dataDescuento.(map[string]interface{})
	var vigContrato float64
	var homologacion []interface{}
	out := make(map[string]interface{})
	if params != nil {
		if e {
			if !e {
				fmt.Println("e1")
				return nil
			}
			vigContrato, e = params[2].(float64)
			//fmt.Println(vigContrato)
			if !e {
				fmt.Println("e3")
				return nil
			}
			if err := request.GetJson("http://"+beego.AppConfig.String("kronosService")+"homologacion_descuento?query=ConceptoTitan:"+strconv.Itoa(int(dataDescuentoAhomologar["Concepto"].(map[string]interface{})["Id"].(float64)))+",Vigencia:"+strconv.FormatFloat(vigContrato, 'f', -1, 64), &homologacion); err == nil {
				if homologacion != nil {
					for _, descuentoKronos := range homologacion {
						row, e := descuentoKronos.(map[string]interface{})

						if e && dataDescuentoAhomologar["ValorCalculado"].(float64) > 0 {
							out["Descuento"] = row["CuentaEspecialKronos"]
							out["Valor"] = dataDescuentoAhomologar["ValorCalculado"]
							//beego.Info(out)
						} else {
							fmt.Println("err  concKron")
							return nil
						}

					}
				}
			} else {
				fmt.Println(err.Error())
				return nil
			}
		}
	} else {
		return nil
	}
	return out
}

func formatoRegistroOpHC(dataLiquidacion interface{}, params ...interface{}) (res interface{}) {
	idLiquidacion, e := params[0].(float64)
	if e {
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
		var idPreliquidacion float64
		//consulta del rp asociado al contrato de la persona... strconv.Itoa(int(vigenciaContrato)) ... strconv.Itoa(int(solicitudrp))
		desagregacionrp := formatoInfoRp(nContrato, vigenciaContrato)
		//fin consulta del rp ...

		//consulta de los devengos calculados por persona
		var devengosNomina []interface{}
		if devengosNomina, err = ConsultarDevengosNominaPorContrato(idLiquidacion, nContrato, vigenciaContrato); err != nil {
			//si no se consiguen datos de devengos para el contrato liquidado (nomina planta)
			return nil
		}
		//homlogacion de los devengos de la nomina de planta a conceptos de kronos...
		idPreliquidacion = devengosNomina[0].(map[string]interface{})["Preliquidacion"].(map[string]interface{})["Id"].(float64)
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
		//---------------------------------------------------------------------------
		//consulta de los descuentos calculados por persona
		var descuentosNomina []interface{}
		if descuentosNomina, err = ConsultarDescuentosNominaPorContrato(idLiquidacion, nContrato, vigenciaContrato); err != nil {
			//si no se consiguen datos de devengos para el contrato liquidado (nomina planta)
		}
		//homlogacion de los descuentos de la nomina de planta a conceptos de kronos...
		if descuentosNomina != nil {
			//beego.Info(descuentosNomina)
			idPreliquidacion = descuentosNomina[0].(map[string]interface{})["Preliquidacion"].(map[string]interface{})["Id"].(float64)
			done = make(chan interface{})
			defer close(done)
			resch = optimize.GenChanInterface(descuentosNomina...)
			chDescHomologados := optimize.Digest(done, homologacionDescuentosHC, resch, params)
			for descuentoHomologado := range chDescHomologados {
				//beego.Info(descuentoHomologado)
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
		}

		//-----------------------------------------------------------------------------

		//---------------------

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
		if desagregacionrp != nil {
			if rpint, e := desagregacionrp[0]["Rp"].(interface{}); e {
				ordenPago := make(map[string]interface{})
				ordenPago["RegistroPresupuestal"] = rpint
				ordenPago["Liquidacion"] = idPreliquidacion
				res["OrdenPago"] = ordenPago
			} else {
				ordenPago := make(map[string]interface{})
				ordenPago["RegistroPresupuestal"] = nil
				res["OrdenPago"] = ordenPago
			}
		} else {
			ordenPago := make(map[string]interface{})
			ordenPago["RegistroPresupuestal"] = nil
			res["OrdenPago"] = ordenPago
		}

		if auxmap, e := infoContrato.(map[string]interface{}); e {
			res["infoPersona"], e = auxmap["infoPersona"]
		}
		res["Contrato"] = nContrato
		res["VigenciaContrato"] = vigenciaContrato
		res["MovimientoContable"] = movimientosContables
		res["ConceptoOrdenPago"], res["Aprobado"], res["Code"] = formatoConceptoOrdenPago(desagregacionrp, homologacionConceptos)
		return res

	} else {
		fmt.Println("err")
		return nil
	}
	return
}

func formatoPreViewCargueMasivoOpHc(liquidacion interface{}, params ...interface{}) (res interface{}) {
	var respuesta []interface{}
	done := make(chan interface{})
	defer close(done)
	_, e := liquidacion.(map[string]interface{})
	if e {
		if liquidacion.(map[string]interface{})["Contratos_por_preliq"] != nil {
			listaLiquidacion := liquidacion.(map[string]interface{})["Contratos_por_preliq"].([]interface{})
			resch := optimize.GenChanInterface(listaLiquidacion...)
			var params []interface{}
			params = append(params, liquidacion.(map[string]interface{})["Id_Preliq"].(interface{}))
			//dav01
			f := formatoRegistroOpFunctionDispatcher(liquidacion.(map[string]interface{})["Nombre_tipo_nomina"].(string))
			beego.Info(f)
			if f != nil {

				chlistaLiquidacion := optimize.Digest(done, f, resch, params)
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
	return
}
