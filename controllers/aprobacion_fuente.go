package controllers

import (
	"strconv"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/udistrital/api_mid_financiera/models"
	. "github.com/mndrix/golog"
)

type AprobacionFuenteController struct {
	beego.Controller
}

// URLMapping ...
func (c *AprobacionFuenteController) URLMapping() {
	c.Mapping("AprobacionFuente", c.ValorMovimientoFuente)
}

//http://localhost:8080/v1/aprobacion_fuente/ValorMovimientoFuenteTraslado?idfuente=37&idapropiacion=256&iddependencia=122&traslado=40000000
//http://localhost:8080/v1/aprobacion_fuente/ValorMovimientoFuente?idfuente=37&idapropiacion=256&iddependencia=122


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

								c.Data["json"] = map[string]interface{}{ "Apropiacion": idapropiacion, "Dependencia": iddependencia, "FuenteFinanciamiento": idfuente, "ValorGastado": valorcdp , "ValorTotal": valor}
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

func (c *AprobacionFuenteController) ValorMovimientoFuenteTraslado() {
	var res []interface{}
	var resfuente []interface{}
	var predicados []models.Predicado
	if idfuente, err := c.GetInt("idfuente"); err == nil {
		fmt.Println(idfuente)
		if iddependencia, err := c.GetInt("iddependencia"); err == nil {
			if idapropiacion, err := c.GetInt("idapropiacion"); err == nil {
				if valortraslado, err := c.GetFloat("traslado"); err == nil {

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
								//reglas
								reglasBase := CargarReglasBase("FuenteFinanciamiento")

								predicados = append(predicados, models.Predicado{Nombre: "movimientofuente(" + strconv.Itoa(idapropiacion) + "," + strconv.Itoa(iddependencia) +  "," + strconv.Itoa(idfuente) + "," + strconv.FormatFloat(valor, 'f', -1, 64) +")."})
								predicados = append(predicados, models.Predicado{Nombre: "saldofuente(" + strconv.Itoa(idapropiacion) + "," + strconv.Itoa(iddependencia) +  "," + strconv.Itoa(idfuente) + "," + strconv.FormatFloat(valorcdp, 'f', -1, 64) +")."})
								//valor que se va a transladar
								predicados = append(predicados, models.Predicado{Nombre: "saldofuente(" + strconv.Itoa(idapropiacion) + "," + strconv.Itoa(iddependencia) +  "," + strconv.Itoa(idfuente) + "," + strconv.FormatFloat(valortraslado, 'f', -1, 64) +")."})
								reglas := FormatoReglas(predicados) + reglasBase

								m := NewMachine().Consult(reglas)
								resultados := m.ProveAll("total_fuente_dependencia_apropiacion_saldo(" + strconv.Itoa(idapropiacion) + "," + strconv.Itoa(iddependencia) +  "," + strconv.Itoa(idfuente) + ",Y).")
								var resp string
								var restante float64
								for _, solution := range resultados {
									resp = fmt.Sprintf("%s", solution.ByName_("Y"))
								}
								f, _ := strconv.ParseFloat(resp, 64)
								restante = f

								var trasladar bool
								if restante >= 0{
									trasladar = true
								}else {
									trasladar = false
								}

								fmt.Println("reglas: ", reglas)
								fmt.Println("RESP: ", resp)

								c.Data["json"] = map[string]interface{}{ "Apropiacion": idapropiacion, "Dependencia": iddependencia, "FuenteFinanciamiento": idfuente, "ValorGastado": valorcdp , "ValorTotal": valor, "ValorRestante": restante, "Trasladar": trasladar }
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
					fmt.Println("err3.5 ", err.Error())
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
