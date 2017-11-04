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
