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
	if idfuente, err := c.GetInt("idfuente"); err == nil {
		if iddependencia, err := c.GetInt("iddependencia"); err == nil {
			if idapropiacion, err := c.GetInt("idapropiacion"); err == nil {

				var Movimiento []map[string]interface{}
				if err := getJson("http://10.20.0.254/financiera_api/v1/movimiento_fuente_financiamiento_apropiacion?query=FuenteFinanciamientoApropiacion.FuenteFinanciamiento.Id:"+strconv.Itoa(idfuente)+"&FuenteFinanciamientoApropiacion.Apropiacion.Id:"+strconv.Itoa(idapropiacion)+"&FuenteFinanciamientoApropiacion.Dependencia:"+strconv.Itoa(iddependencia), &Movimiento); err == nil {
					if Movimiento != nil {
						fmt.Println(Movimiento)

						var valorGastado map[string]interface{}
						if err := getJson("http://10.20.0.254/financiera_mid_api/v1/disponibilidad/ValorDisponibilidadesFuenteRubroDependencia?idfuente="+strconv.Itoa(idfuente)+"&idapropiacion="+strconv.Itoa(idapropiacion)+"iddependencia="+strconv.Itoa(iddependencia), &valorGastado); err == nil {
							fmt.Println(Movimiento)

							for _, valor := range valorGastado {
								res = append(res, valor)
							}

							/*
							for _, dependencia := range dependencias {
								peticion := "solicitud_disponibilidad?"
								peticion = peticion + "limit=-1&query=Necesidad.FuenteReversa.FuenteFinanciamiento:" + strconv.Itoa(idfuente) + ","
								peticion = peticion + "Necesidad.FuenteReversa.Apropiacion:" + strconv.Itoa(idapropiacion) + ","
								peticion = peticion + "Necesidad.DependenciaReversa.JefeDependenciaDestino:" + strconv.Itoa(int(dependencia["Id"].(float64))) + ","
								peticion = peticion + "Expedida:true"
								var solicitud_disponibilidades []map[string]interface{}
								if err := getJson("http://"+beego.AppConfig.String("argoService")+peticion, &solicitud_disponibilidades); err == nil {
									fmt.Println(solicitud_disponibilidades)
									for _, solicitud_disponibilidad := range solicitud_disponibilidades {
										var disponibilidades []map[string]interface{}
										fmt.Println("/disponibilidad_apropiacion?query=Disponibilidad.Solicitud:" + strconv.Itoa(int(solicitud_disponibilidad["Id"].(float64))) + ",Apropiacion.Id:" + strconv.Itoa(idapropiacion))
										if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/disponibilidad_apropiacion?query=Disponibilidad.Solicitud:"+strconv.Itoa(int(solicitud_disponibilidad["Id"].(float64)))+",Apropiacion.Id:"+strconv.Itoa(idapropiacion)+",FuenteFinanciamiento:"+strconv.Itoa(idfuente), &disponibilidades); err == nil {
											if disponibilidades != nil {
												for _, disponibilidad := range disponibilidades {
													res = append(res, disponibilidad)
												}
											}
										} else {
											fmt.Println("err7  ", err.Error())
											c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
										}
									}
								} else {
									fmt.Println("err6 ", err.Error())
									c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
								}
							}
							if res != nil {
								var valorcdp float64
								valorcdp = 0
								for _, row := range res {
									valorcdp = valorcdp + row.(map[string]interface{})["Valor"].(float64)
								}
								c.Data["json"] = map[string]interface{}{"valor": valorcdp}
							}
						*/}else {
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
