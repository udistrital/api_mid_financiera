package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/udistrital/api_mid_financiera/golog"
	"github.com/udistrital/api_mid_financiera/models"

	"github.com/astaxie/beego"
)

type DisponibilidadController struct {
	beego.Controller
}

func (c *DisponibilidadController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("Solicitudes", c.InfoSolicitudDisponibilidad)
	c.Mapping("SolicitudById", c.InfoSolicitudDisponibilidadById)
}

// Solicitudes ...
// @Title Solicitudes
// @Description get Solicitudes
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.InfoSolDisp
// @Failure 403
// @router Solicitudes/ [get]
func (this *DisponibilidadController) InfoSolicitudDisponibilidad() {
	var solicitud []models.SolicitudDisponibilidad
	var res []models.InfoSolDisp
	if err := getJson("http://"+beego.AppConfig.String("argoService")+"solicitud_disponibilidad?limit=0&query=Expedida:false&sortby=Id&order=desc", &solicitud); err == nil {

		for _, resultado := range solicitud {

			var depNes []models.DependenciaNecesidad
			var depSol []models.Dependencia
			var depDest []models.Dependencia
			var necesidad []models.Necesidad
			//var temp models.InfoSolDisp
			//temp.SolicitudDisponibilidad = &resultado
			if err := getJson("http://"+beego.AppConfig.String("argoService")+"necesidad?limit=1&query=Id:"+strconv.Itoa(resultado.Necesidad.Id), &necesidad); err == nil {
				necesidadaux := necesidad[0]
				resultado.Necesidad = &necesidadaux

				if err := getJson("http://"+beego.AppConfig.String("argoService")+"dependencia_necesidad?limit=0&query=Necesidad.Id:"+strconv.Itoa(necesidad[0].Id), &depNes); err == nil {
					fmt.Println("http://" + beego.AppConfig.String("oikosService") + "dependencia?limit=0&query=Id:" + strconv.Itoa(depNes[0].DependenciaSolicitante))
					if err := getJson("http://"+beego.AppConfig.String("oikosService")+"dependencia?limit=0&query=Id:"+strconv.Itoa(depNes[0].DependenciaSolicitante), &depSol); err == nil {
						//temp.DependenciaSolicitante = &depSol[0]
						if err := getJson("http://"+beego.AppConfig.String("agoraService")+"informacion_persona_natural/"+strconv.Itoa(depNes[0].JefeDependenciaSolicitante), &depSol[0].InfoJefeDependencia); err == nil {
							//temp.DependenciaSolicitante = &depSol[0]

						} else {
							this.Data["json"] = err.Error()
							fmt.Println("error5: ", err)
						}
						fmt.Println(depNes[0].OrdenadorGasto)
						if err := getJson("http://"+beego.AppConfig.String("agoraService")+"informacion_persona_natural/"+strconv.Itoa(depNes[0].OrdenadorGasto), &depSol[0].InfoOrdenador); err == nil {
							//temp.DependenciaSolicitante = &depSol[0]

						} else {
							this.Data["json"] = err.Error()
							fmt.Println("error5: ", err)
						}
					} else {
						this.Data["json"] = err.Error()
						fmt.Println("error5: ", err)
					}
					if err := getJson("http://"+beego.AppConfig.String("oikosService")+"dependencia?limit=0&query=Id:"+strconv.Itoa(depNes[0].DependenciaDestino), &depDest); err == nil {
						//temp.DependenciaDestino = &depDest[0]
					} else {
						this.Data["json"] = err.Error()
						fmt.Println("error4: ", err)
					}
					if depSol == nil {
						depSol = append(depSol, models.Dependencia{Nombre: "Indefinida"})
					}
					if depDest == nil {
						depDest = append(depDest, models.Dependencia{Nombre: "Indefinida"})
					}
					temp := models.InfoSolDisp{SolicitudDisponibilidad: resultado, DependenciaSolicitante: depSol[0], DependenciaDestino: depDest[0]}
					res = append(res, temp)
				} else {
					this.Data["json"] = err.Error()
					fmt.Println("error3: ", err)
				}
			} else {
				this.Data["json"] = err.Error()
				fmt.Println("error2: ", err)
			}
		}

		this.Data["json"] = res
		fmt.Println("solicitud: ", solicitud)
	} else {
		this.Data["json"] = err.Error()
		fmt.Println("error1: ", err)
	}
	this.ServeJSON()
}

// SolicitudById ...
// @Title SolicitudById
// @Description get SolicitudById
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.InfoSolDisp
// @Failure 403
// @router SolicitudById/:id [get]
func (this *DisponibilidadController) InfoSolicitudDisponibilidadById() {
	var solicitud []models.SolicitudDisponibilidad
	var res []models.InfoSolDisp
	idStr := this.Ctx.Input.Param(":id")
	if err := getJson("http://"+beego.AppConfig.String("argoService")+"solicitud_disponibilidad?limit=0&query=Id:"+idStr, &solicitud); err == nil {

		for _, resultado := range solicitud {
			//var temp models.InfoSolDisp
			//temp.SolicitudDisponibilidad = resultado
			var depNes []models.DependenciaNecesidad
			var depSol []models.Dependencia
			var depDest []models.Dependencia
			var necesidad []models.Necesidad
			if err := getJson("http://"+beego.AppConfig.String("argoService")+"necesidad?limit=1&query=Id:"+strconv.Itoa(resultado.Necesidad.Id), &necesidad); err == nil {
				if necesidad != nil {
					necesidadaux := necesidad[0]
					resultado.Necesidad = &necesidadaux

					if err := getJson("http://"+beego.AppConfig.String("argoService")+"dependencia_necesidad?limit=0&query=Necesidad.Id:"+strconv.Itoa(necesidad[0].Id), &depNes); err == nil {
						fmt.Println("http://" + beego.AppConfig.String("oikosService") + "dependencia?limit=0&query=Id:" + strconv.Itoa(depNes[0].DependenciaSolicitante))
						if err := getJson("http://"+beego.AppConfig.String("oikosService")+"dependencia?limit=0&query=Id:"+strconv.Itoa(depNes[0].DependenciaSolicitante), &depSol); err == nil {
							//temp.DependenciaSolicitante = depSol[0]
							if err := getJson("http://"+beego.AppConfig.String("agoraService")+"informacion_persona_natural/"+strconv.Itoa(depNes[0].JefeDependenciaSolicitante), &depSol[0].InfoJefeDependencia); err == nil {
								//temp.DependenciaSolicitante = &depSol[0]

							} else {
								this.Data["json"] = err.Error()
								fmt.Println("error7: ", err)
							}
							fmt.Println(depNes[0].OrdenadorGasto)
							if err := getJson("http://"+beego.AppConfig.String("agoraService")+"informacion_persona_natural/"+strconv.Itoa(depNes[0].OrdenadorGasto), &depSol[0].InfoOrdenador); err == nil {
								//temp.DependenciaSolicitante = &depSol[0]

							} else {
								this.Data["json"] = err.Error()
								fmt.Println("error6: ", err)
							}
						} else {
							this.Data["json"] = err.Error()
							fmt.Println("error5: ", err)
						}

						if err := getJson("http://"+beego.AppConfig.String("oikosService")+"dependencia?limit=0&query=Id:"+strconv.Itoa(depNes[0].DependenciaDestino), &depDest); err == nil {
							//temp.DependenciaDestino = depDest[0]
						} else {
							this.Data["json"] = err.Error()
							fmt.Println("error4: ", err)
						}
						if depSol == nil {
							depSol = append(depSol, models.Dependencia{Nombre: "Indefinida"})
						}
						if depDest == nil {
							depDest = append(depDest, models.Dependencia{Nombre: "Indefinida"})
						}
						temp := models.InfoSolDisp{SolicitudDisponibilidad: resultado, DependenciaSolicitante: depSol[0], DependenciaDestino: depDest[0]}
						res = append(res, temp)
					} else {
						this.Data["json"] = err.Error()
						fmt.Println("error3: ", err)
					}
				} else {
					this.Data["json"] = "error"
				}
			} else {
				this.Data["json"] = err.Error()
				fmt.Println("error2: ", err)
			}
		}

		this.Data["json"] = res
		fmt.Println("solicitud: ", solicitud)
	} else {
		this.Data["json"] = err.Error()
		fmt.Println("error1: ", err)
	}
	this.ServeJSON()
}

// Post ...
// @Title Post
// @Description create Disponibilidad
// @Param	body		body 	models.InfoSolDisp	true		"body for InfoSolDisp content"
// @Success 201 {int} models.InfoSolDisp
// @Failure 403 body is empty
// @router / [post]
func (this *DisponibilidadController) Post() {
	var predicados []models.Predicado
	var solicitudes_disponibilidad []models.InfoSolDisp
	var ultimo_registro_disponibilidad []models.Disponibilidad //se cargara en esta variable el ultimo registro de disponibilidad
	var numero_asignado_disponibilidad float64
	var disponibilidad models.Disponibilidad
	var respuesta models.Disponibilidad
	var respuesta_mod interface{}
	var respuesta_disponibilidad_rubro interface{}
	var rubros_solicitud []models.FuenteFinanciacionRubroNecesidad
	var alertas []string
	alertas = append(alertas, "success")
	aprobada := true
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &solicitudes_disponibilidad); err == nil {
		for i := 0; i < len(solicitudes_disponibilidad); i++ {
			if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/disponibilidad?limit=1&query=Vigencia:"+strconv.FormatFloat(solicitudes_disponibilidad[i].SolicitudDisponibilidad.Vigencia, 'f', -1, 64)+"&sortby=NumeroDisponibilidad&order=desc", &ultimo_registro_disponibilidad); err == nil {
				if ultimo_registro_disponibilidad == nil {
					numero_asignado_disponibilidad = 1
				} else {
					numero_asignado_disponibilidad = ultimo_registro_disponibilidad[0].NumeroDisponibilidad + 1
				}
				fmt.Println("numero: ", numero_asignado_disponibilidad)
			} else {
				fmt.Println("error1345: ", err.Error())
				alertas[0] = "error"
				alertas = append(alertas, "no se pudo caragar datos del servicio crud")
				this.Data["json"] = alertas
				this.ServeJSON()
			}

			//comprobar apropiacion de los rubros
			if err := getJson("http://"+beego.AppConfig.String("argoService")+"fuente_financiacion_rubro_necesidad?limit=0&query=SolicitudNecesidad.Id:"+strconv.Itoa(solicitudes_disponibilidad[i].SolicitudDisponibilidad.Necesidad.Id), &rubros_solicitud); err == nil {
				fmt.Println("rubros: ", len(rubros_solicitud))
				reglasBase := CargarReglasBase("Presupuesto")
				for j := 0; j < len(rubros_solicitud); j++ {

					var saldo_aprop float64
					fmt.Println("aprop: ", rubros_solicitud[j].Apropiacion)
					if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/apropiacion/SaldoApropiacion/"+strconv.Itoa(rubros_solicitud[j].Apropiacion), &saldo_aprop); err == nil {
						predicados = append(predicados, models.Predicado{Nombre: "rubro_apropiacion(" + strconv.Itoa(rubros_solicitud[j].Apropiacion) + "," + strconv.Itoa(rubros_solicitud[j].Id) + "," + strconv.FormatFloat(saldo_aprop, 'f', -1, 64) + ")."})
					} else {
						alertas[0] = "error"
						alertas = append(alertas, "error al cargar saldo de la apropiacion ")
						aprobada = false
						fmt.Println("error2: ", err)
					}
					predicados = append(predicados, models.Predicado{Nombre: "valor_rubro_cdp(" + strconv.Itoa(rubros_solicitud[j].Apropiacion) + "," + strconv.Itoa(rubros_solicitud[j].Id) + "," + strconv.FormatFloat(rubros_solicitud[j].MontoParcial, 'f', -1, 64) + ")."})

					/*if (apropiacion_rubr != nil){
									this.Data["json"] = "No se puede expedir el cdp de la solicitud No "+strconv.Itoa(solicitudes_disponibilidad[i].Numero)
					 	 		  this.ServeJSON()
								}*/
					fmt.Println("Apropiacion: ", alertas)

				}
				reglas := FormatoReglas(predicados) + reglasBase
				fmt.Println("reglas: ", reglas)
				res := golog.GetBoolean(reglas, "aprobacion_cdp("+strconv.Itoa(rubros_solicitud[0].Apropiacion)+",Y).", "Y")
				if res {
					aprobada = res
				} else {
					aprobada = res
					alertas[0] = "error"
					alertas = append(alertas, "Expedición de CDP no aprobada, algunos valores superan el saldo de las APropiaciones.")
				}
			} else {
				alertas[0] = "error"
				alertas = append(alertas, "No se pudo cargar los rubros de la solicitud "+strconv.Itoa(solicitudes_disponibilidad[i].SolicitudDisponibilidad.Numero))

				aprobada = false
				fmt.Println("error3: ", err)
			}

			//-----------------------------------
			//realizar segunda peticion para los datos de la necesidad.
			if aprobada { // si se aprueba la solicitud, se genera el cdp
				disponibilidad = models.Disponibilidad{
					UnidadEjecutora:      &models.UnidadEjecutora{Id: solicitudes_disponibilidad[i].SolicitudDisponibilidad.Necesidad.UnidadEjecutora},
					Vigencia:             solicitudes_disponibilidad[i].SolicitudDisponibilidad.Necesidad.Vigencia,
					NumeroDisponibilidad: numero_asignado_disponibilidad,
					NumeroOficio:         strconv.Itoa(solicitudes_disponibilidad[i].SolicitudDisponibilidad.Numero),
					FechaRegistro:        time.Now().Local(),
					Estado:               &models.EstadoDisponibilidad{Id: 1},
					Solicitud:            solicitudes_disponibilidad[i].SolicitudDisponibilidad.Id,
				}
				solicitudes_disponibilidad[i].SolicitudDisponibilidad.Expedida = true

				err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/disponibilidad", "POST", &respuesta, &disponibilidad)

				if err == nil {
					sendJson("http://"+beego.AppConfig.String("argoService")+"solicitud_disponibilidad/"+strconv.Itoa(solicitudes_disponibilidad[i].SolicitudDisponibilidad.Id), "PUT", &respuesta_mod, &solicitudes_disponibilidad[i].SolicitudDisponibilidad)
					fmt.Println("err", respuesta_mod)
					for j := 0; j < len(rubros_solicitud); j++ {
						disponibilidad_apropiacion := models.DisponibilidadApropiacion{
							Apropiacion:    &models.Apropiacion{Id: rubros_solicitud[j].Apropiacion},
							Disponibilidad: &respuesta, //&respuesta,
							Valor:          rubros_solicitud[j].MontoParcial}
						sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/disponibilidad_apropiacion", "POST", &respuesta_disponibilidad_rubro, &disponibilidad_apropiacion)
					}
					alertas = append(alertas, "se genero el CDP para la solicitud No "+strconv.Itoa(solicitudes_disponibilidad[i].SolicitudDisponibilidad.Numero))
				} else {
					alertas[0] = "error"
					alertas = append(alertas, "Error al registrar el CDP para la solicitud No "+strconv.Itoa(solicitudes_disponibilidad[i].SolicitudDisponibilidad.Numero))
					fmt.Println("err", err)
				}

			} else {

			}

			aprobada = true
		} //fin for json
		this.Data["json"] = alertas
		this.ServeJSON()
	} else {
		fmt.Println("error: ", err)
		this.Data["json"] = alertas
		this.ServeJSON()
	}
}
