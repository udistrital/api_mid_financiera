package controllers

import (
	"strconv"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/udistrital/api_mid_financiera/models"
)

// AprobacionFuenteController operations for homologation fo liquidation of titan
type AprobacionFuenteController struct {
	beego.Controller
}

// URLMapping ...
func (c *AprobacionFuenteController) URLMapping() {
	c.Mapping("AprobacionFuente", c.ValorMovimientoFuente)
}


// @Title MidAprobacionFuenteLiquidacion
// @Description homologa conceptos de titan con conceptos kronos
// @Param	idPreliquidacion	identificador de la liquidación de titan
// @Param	body		body 	models.IdLiquidacion, models.RegistroPresupuestal	"body for AprobacionFuente content"
// @Success 201 {object} models.Conceptos
// @Failure 403 body is empty
// @router MidAprobacionFuenteLiquidacion [post]

//http://10.20.0.254/financiera_api/v1/movimiento_fuente_financiamiento_apropiacion?query=FuenteFinanciamientoApropiacion.FuenteFinanciamiento.Id:38&FuenteFinanciamientoApropiacion.Apropiacion.Id:256&FuenteFinanciamientoApropiacion.Dependencia:95

func (c *AprobacionFuenteController) ValorMovimientoFuente() {
	var res []interface{}
	var resfuente []interface{}
	if idfuente, err := c.GetInt("idfuente"); err == nil {
		fmt.Println(idfuente)
		if iddependencia, err := c.GetInt("iddependencia"); err == nil {
			if idapropiacion, err := c.GetInt("idapropiacion"); err == nil {

				var Movimiento []map[string]interface{}
				if err := getJson("http://10.20.0.254/financiera_api/v1/movimiento_fuente_financiamiento_apropiacion?query=FuenteFinanciamientoApropiacion.FuenteFinanciamiento.Id:"+strconv.Itoa(idfuente)+",FuenteFinanciamientoApropiacion.Apropiacion.Id:"+strconv.Itoa(idapropiacion)+",FuenteFinanciamientoApropiacion.Dependencia:"+strconv.Itoa(iddependencia), &Movimiento); err == nil {
					if Movimiento != nil {

						for _, Movimientos := range Movimiento {
							resfuente = append(resfuente, Movimientos)
						}


						var valorGastado map[string]interface{}
						if err := getJson("http://10.20.0.254/financiera_mid_api/v1/disponibilidad/ValorDisponibilidadesFuenteRubroDependencia?idfuente="+strconv.Itoa(idfuente)+"&idapropiacion="+strconv.Itoa(idapropiacion)+"&iddependencia="+strconv.Itoa(iddependencia), &valorGastado); err == nil {
							fmt.Println(valorGastado)
							if valorGastado != nil {
							for _, valores := range valorGastado {
								res = append(res, valores)
							}

							if resfuente != nil{
								if res != nil{
								var valor float64
								valor = 0
								var valorcdp float64
								valorcdp = 0
								valorcdp = res[0].(float64)

								for _, rowfuente := range resfuente {
									valor = valor + rowfuente.(map[string]interface{})["Valor"].(float64)
								}
								c.Data["json"] = map[string]interface{}{"valorTotal": valor , "valorGastado": valorcdp}
							  }
							}
						}
						}else {
							fmt.Println("err5 ", err.Error())
							c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
						}
					} else {
						fmt.Println("aqui")
						c.Data["json"] = nil
					}
				} else {
					fmt.Println("err4 ", err.Error())
					c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
				}
			} else {
				fmt.Println("err3 ", err.Error())
				c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
			}
		} else {
			fmt.Println("err2 ", err.Error())
			c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
		}
	} else {
		fmt.Println("err1 ", err.Error())
		c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
	}

	c.ServeJSON()
}
