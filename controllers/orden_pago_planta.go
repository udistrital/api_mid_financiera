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
	idLiquidacion, e := params[0].(float64)
	if e {
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
		var idPreliquidacion float64
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
		//---------------------------------------------------------------------------
		//consulta de los descuentos calculados por persona
		/*var descuentosNomina []interface{}
		if descuentosNomina, err := ConsultarDescuentosNominaPorContrato(idLiquidacion, nContrato, vigenciaContrato); err != nil {
			//si no se consiguen datos de devengos para el contrato liquidado (nomina planta)
			return nil
		}
		//homlogacion de los descuentos de la nomina de planta a conceptos de kronos...
		idPreliquidacion = listaDetalles[0].(map[string]interface{})["Preliquidacion"].(map[string]interface{})["Id"].(float64)
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
		}
		res["Contrato"] = nContrato
		res["VigenciaContrato"] = vigenciaContrato
		res["MovimientoContable"] = movimientosContables
		res["IdPreliquidacion"] = idPreliquidacion
		//res["ConceptoOrdenPago"], res["Aprobado"], res["Code"] = formatoConceptoOrdenPago(desagregacionrp, homologacionConceptos)
		return res
	} else {
		return nil
	}
	return
}

func homologacionConceptosDocentesPlanta(dataConcepto interface{}, params ...interface{}) (res interface{}) {
	dataConceptoAhomologar, e := dataConcepto.(map[string]interface{})
	//fmt.Println(dataConcepto)
	if !e {
		fmt.Println("e1")
		return nil
	}
	//fmt.Println(dataConceptoAhomologar)
	out := make(map[string]interface{})
	var homologacion []interface{}
	//aqui va la consulta sobre facultad y proyecto para HC (modificar para hacerla de forma genral)

	//fmt.Println("http://" + beego.AppConfig.String("Urlcrud") + ":" + beego.AppConfig.String("Portcrud") + "/" + beego.AppConfig.String("Nscrud") + "/homologacion_concepto?query=ConceptoTitan:" + strconv.Itoa(int(dataConceptoAhomologar["Concepto"].(map[string]interface{})["Id"].(float64))) + ",ConceptoKronos.ConceptoTesoralFacultadProyecto.Facultad:" + strconv.Itoa(int(idFacultad)) + ",ConceptoKronos.ConceptoTesoralFacultadProyecto.ProyectoCurricular:" + strconv.Itoa(int(idProyecto)))
	if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/homologacion_concepto?query=ConceptoTitan:"+strconv.Itoa(int(dataConceptoAhomologar["Concepto"].(map[string]interface{})["Id"].(float64))), &homologacion); err == nil {
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

	return out

}
