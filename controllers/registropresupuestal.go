package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/udistrital/api_mid_financiera/golog"
	"github.com/udistrital/api_mid_financiera/models"
	"github.com/udistrital/api_mid_financiera/tools"
	"github.com/udistrital/api_mid_financiera/utilidades"
)

// RegistroPresupuestalController operations for RegistroPresupuestal
type RegistroPresupuestalController struct {
	beego.Controller
}

// URLMapping ...
func (c *RegistroPresupuestalController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetSaldoRp", c.GetSaldoRp)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
	c.Mapping("GetSolicitudesRp", c.GetSolicitudesRp)
	c.Mapping("GetSolicitudesRpById", c.GetSolicitudesRpById)
}

// GetSolicitudesRp ...
// @Title GetSolicitudesRp
// @Description get saldo rp by apropiacion
// @Success 200 {object} models.SolicitudRp
// @Failure 403 :id is empty
// @router /GetSolicitudesRp [get]
func (c *RegistroPresupuestalController) GetSolicitudesRp() {
	var solicitudes_rp []models.SolicitudRp
	var respuesta []models.SolicitudRp
	if err := getJson("http://"+beego.AppConfig.String("argoService")+"solicitud_rp?limit=0&query=Expedida:false&sortby=Id&order=desc", &solicitudes_rp); err == nil {
		if solicitudes_rp != nil {
			//encontrar datos del CDP objetivo del RP Solicitado
			for _, solicitud := range solicitudes_rp {
				//recuperar datos del CDP objetivo de la solicitud
				var cdp_objtvo []models.Disponibilidad
				if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/disponibilidad?limit=1&query=Id:"+strconv.Itoa(solicitud.Cdp), &cdp_objtvo); err == nil {
					if cdp_objtvo != nil {
						solicitud.DatosDisponibilidad = &cdp_objtvo[0]
						var necesidad_cdp []models.SolicitudDisponibilidad
						if err := getJson("http://"+beego.AppConfig.String("argoService")+"solicitud_disponibilidad?limit=1&query=Id:"+strconv.Itoa(solicitud.DatosDisponibilidad.Solicitud), &necesidad_cdp); err == nil {
							if necesidad_cdp != nil {
								solicitud.DatosDisponibilidad.DatosNecesidad = necesidad_cdp[0].Necesidad
								var depNes []models.DependenciaNecesidad
								if err := getJson("http://"+beego.AppConfig.String("argoService")+"dependencia_necesidad?limit=0&query=Necesidad.Id:"+strconv.Itoa(solicitud.DatosDisponibilidad.DatosNecesidad.Id), &depNes); err == nil {
									if depNes != nil {
										var depSol []models.Dependencia
										var jefe_dep_sol []models.JefeDependencia
										if err := getJson("http://"+beego.AppConfig.String("coreService")+"jefe_dependencia?limit=1&query=Id:"+strconv.Itoa(depNes[0].JefeDependenciaSolicitante), &jefe_dep_sol); err == nil {
											if jefe_dep_sol != nil {
												if err := getJson("http://"+beego.AppConfig.String("oikosService")+"dependencia?limit=1&query=Id:"+strconv.Itoa(jefe_dep_sol[0].DependenciaId), &depSol); err == nil {
													if depSol != nil {
														solicitud.DatosDisponibilidad.DatosNecesidad.DatosDependenciaSolicitante = &depSol[0]
													} else {
														//si no hay datos de la dependencia
													}
												} else {
													//si hay error al consultar la dependecia solicitante
												}
											} else {
												//no hay datos jefe dep
											}

										} else {
											//jefe_dep
										}
									} else {
										//si no hay datos en la consulta dependencia_necesidad
									}
								} else {
									//si hay error al consultar dependencia_necesidad
								}
							} else {
								//si no hay datos de la necesidad
							}
						} else {
							//si hay error al encontrar datos de la necesidad
						}
					} else {
						//si no hay datos del CDP objetivo
						fmt.Println("error cdp: no hay datos, id : ", solicitud.Cdp)
					}
				} else {
					//si ocurre error al traer datos del CDP objetivo
					fmt.Println("error cdp: ", err)
				}
				//obtener informacion del contrato del rp
				var info_contrato []models.ContratoGeneral
				var contratista []models.InformacionProveedor
				fmt.Println("sol ", solicitud.NumeroContrato)
				if err := getJson("http://"+beego.AppConfig.String("argoService")+"contrato_general?limit=1&query=Id:"+solicitud.NumeroContrato, &info_contrato); err == nil {
					if info_contrato != nil {
						if err := getJson("http://"+beego.AppConfig.String("agoraService")+"informacion_proveedor?limit=1&query=NumDocumento:"+strconv.Itoa(info_contrato[0].Contratista), &contratista); err == nil {
							solicitud.DatosProveedor = &contratista[0]
						} else {
							//error consulta proveedor
							fmt.Println(err.Error())
						}

					} else {
						//si no encuentra datos sobre el contrato
						fmt.Println("error contrato: no hay datos, id : ", solicitud.NumeroContrato)
					}
				} else {
					//si ocurre error al obtener los datos del contrato
					fmt.Println(err.Error())
				}
				//cargar datos del compromiso de la solicitud de rp
				var compromiso_rp []models.Compromiso
				if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/compromiso?limit=1&query=Id:"+strconv.Itoa(solicitud.Compromiso), &compromiso_rp); err == nil {
					if compromiso_rp != nil {
						solicitud.DatosCompromiso = &compromiso_rp[0]
					} else {
						//si no encuentra los datos del compromiso

					}
				} else {
					//si hay error al cargar el compromiso del rp
				}
				respuesta = append(respuesta, solicitud)
			}
			c.Data["json"] = respuesta
		} else {
			//si no hay datos de solicitudes
			c.Data["json"] = "sin datos"
		}
	} else {
		//si ocurre error al traer las solicitudes
		c.Data["json"] = err
	}
	c.ServeJSON()
}

// GetSolicitudesRpById ...
// @Title GetSolicitudesRpById
// @Description get GetSolicitudesRpById by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.SolicitudRp
// @Failure 403 :id is empty
// @router /GetSolicitudesRpById/:id [get]
func (c *RegistroPresupuestalController) GetSolicitudesRpById() {
	var solicitudes_rp []models.SolicitudRp
	var respuesta []models.SolicitudRp
	idStr := c.Ctx.Input.Param(":id")
	if err := getJson("http://"+beego.AppConfig.String("argoService")+"solicitud_rp?limit=1&query=Id:"+idStr, &solicitudes_rp); err == nil {
		if solicitudes_rp != nil {
			//encontrar datos del CDP objetivo del RP Solicitado
			for _, solicitud := range solicitudes_rp {
				//recuperar datos del CDP objetivo de la solicitud
				var cdp_objtvo []models.Disponibilidad
				if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/disponibilidad?limit=1&query=Id:"+strconv.Itoa(solicitud.Cdp), &cdp_objtvo); err == nil {
					if cdp_objtvo != nil {
						solicitud.DatosDisponibilidad = &cdp_objtvo[0]
						var necesidad_cdp []models.SolicitudDisponibilidad
						if err := getJson("http://"+beego.AppConfig.String("argoService")+"solicitud_disponibilidad?limit=1&query=Id:"+strconv.Itoa(solicitud.DatosDisponibilidad.Solicitud), &necesidad_cdp); err == nil {
							if necesidad_cdp != nil {
								solicitud.DatosDisponibilidad.DatosNecesidad = necesidad_cdp[0].Necesidad
								var depNes []models.DependenciaNecesidad
								if err := getJson("http://"+beego.AppConfig.String("argoService")+"dependencia_necesidad?limit=0&query=Necesidad.Id:"+strconv.Itoa(solicitud.DatosDisponibilidad.DatosNecesidad.Id), &depNes); err == nil {
									if depNes != nil {
										var depSol []models.Dependencia
										var jefe_dep_sol []models.JefeDependencia
										if err := getJson("http://"+beego.AppConfig.String("coreService")+"jefe_dependencia?limit=1&query=Id:"+strconv.Itoa(depNes[0].JefeDependenciaSolicitante), &jefe_dep_sol); err == nil {
											if jefe_dep_sol != nil {
												if err := getJson("http://"+beego.AppConfig.String("oikosService")+"dependencia?limit=1&query=Id:"+strconv.Itoa(jefe_dep_sol[0].DependenciaId), &depSol); err == nil {
													if depSol != nil {
														solicitud.DatosDisponibilidad.DatosNecesidad.DatosDependenciaSolicitante = &depSol[0]
													} else {
														//si no hay datos de la dependencia
													}
												} else {
													//si hay error al consultar la dependecia solicitante
												}
											} else {
												//no hay datos jefe dep
											}

										} else {
											//jefe_dep
										}
									} else {
										//si no hay datos en la consulta dependencia_necesidad
									}
								} else {
									//si hay error al consultar dependencia_necesidad
								}
							} else {
								//si no hay datos de la necesidad
							}
						} else {
							//si hay error al encontrar datos de la necesidad
						}
					} else {
						//si no hay datos del CDP objetivo
						fmt.Println("error cdp: no hay datos, id : ", solicitud.Cdp)
					}
				} else {
					//si ocurre error al traer datos del CDP objetivo
					fmt.Println("error cdp: ", err)
				}
				//obtener informacion del contrato del rp
				var info_contrato []models.ContratoGeneral
				var contratista []models.InformacionProveedor
				if err := getJson("http://"+beego.AppConfig.String("argoService")+"contrato_general?limit=1&query=Id:"+solicitud.NumeroContrato, &info_contrato); err == nil {
					if info_contrato != nil {
						if err := getJson("http://"+beego.AppConfig.String("agoraService")+"informacion_proveedor?limit=1&query=NumDocumento:"+strconv.Itoa(info_contrato[0].Contratista), &contratista); err == nil {
							solicitud.DatosProveedor = &contratista[0]
						} else {
							//error consulta proveedor
							fmt.Println(err.Error())
						}
					} else {
						//si no encuentra datos sobre el contrato
						fmt.Println("error contrato: no hay datos, id : ", solicitud.NumeroContrato)
					}
				} else {
					//si ocurre error al obtener los datos del contrato
				}
				//cargar datos del compromiso de la solicitud de rp
				var compromiso_rp []models.Compromiso
				if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/compromiso?limit=1&query=Id:"+strconv.Itoa(solicitud.Compromiso), &compromiso_rp); err == nil {
					if compromiso_rp != nil {
						solicitud.DatosCompromiso = &compromiso_rp[0]
					} else {
						//si no encuentra los datos del compromiso
					}
				} else {
					//si hay error al cargar el compromiso del rp
				}
				respuesta = append(respuesta, solicitud)
			}
			if respuesta != nil {
				c.Data["json"] = respuesta[0]
			} else {
				c.Data["json"] = respuesta
			}
		} else {
			//si no hay datos de solicitudes
			c.Data["json"] = nil
		}
	} else {
		//si ocurre error al traer las solicitudes
		c.Data["json"] = err
	}
	c.ServeJSON()
}

// Post ...
// @Title Create
// @Description create RegistroPresupuestal
// @Param	body		body 	models.DatosRegistroPresupuestal	true		"body for DatosRegistroPresupuestal content"
// @Success 201 {object} models.DatosRegistroPresupuestal
// @Failure 403 body is empty
// @router / [post]
func (c *RegistroPresupuestalController) Post() {
	var rp_a_registrar models.DatosRegistroPresupuestal
	var alertas []string
	alertas = append(alertas, "success")
	var predicados []models.Predicado
	var respuesta interface{}
	var numero_asignado_rp int
	var ultimo_registro []models.RegistroPresupuestal
	var comprobacion models.RegistroPresupuestal
	var saldoCDP map[string]float64
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &rp_a_registrar); err == nil {

		reglasBase := CargarReglasBase("Presupuesto")
		for _, rubros_a_comprobar := range rp_a_registrar.Rubros {

			datos := models.DatosRubroRegistroPresupuestal{Disponibilidad: rubros_a_comprobar.Disponibilidad,
				Apropiacion: rubros_a_comprobar.Apropiacion, FuenteFinanciacion: rubros_a_comprobar.FuenteFinanciacion,
			}
			if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/disponibilidad/SaldoCdp", "POST", &saldoCDP, &datos); err == nil {
				fmt.Println(rubros_a_comprobar.FuenteFinanciacion)
				if rubros_a_comprobar.FuenteFinanciacion == nil {
					predicados = append(predicados, models.Predicado{Nombre: "rubro_cdp(" + strconv.Itoa(rubros_a_comprobar.Disponibilidad.Id) + "," + strconv.Itoa(datos.Apropiacion.Id) + "," + strconv.Itoa(0) +
						"," + strconv.FormatFloat(saldoCDP["saldo"], 'f', -1, 64) + ")."})
				} else {
					predicados = append(predicados, models.Predicado{Nombre: "rubro_cdp(" + strconv.Itoa(rubros_a_comprobar.Disponibilidad.Id) + "," + strconv.Itoa(datos.Apropiacion.Id) + "," + strconv.Itoa(rubros_a_comprobar.FuenteFinanciacion.Id) +
						"," + strconv.FormatFloat(saldoCDP["saldo"], 'f', -1, 64) + ")."})
				}

			} else {
				alertas[0] = "error"
				alertas = append(alertas, "No se pudo cargar el saldo para algunas apropiaciones")
				c.Data["json"] = alertas

				c.ServeJSON()
			}
			if rubros_a_comprobar.FuenteFinanciacion == nil {
				predicados = append(predicados, models.Predicado{Nombre: "valor_rubro_rp(" + strconv.Itoa(rubros_a_comprobar.Disponibilidad.Id) + "," + strconv.Itoa(datos.Apropiacion.Id) + "," + strconv.Itoa(0) + "," + strconv.FormatFloat(rubros_a_comprobar.ValorAsignado, 'f', -1, 64) + ")."})

			} else {
				predicados = append(predicados, models.Predicado{Nombre: "valor_rubro_rp(" + strconv.Itoa(rubros_a_comprobar.Disponibilidad.Id) + "," + strconv.Itoa(datos.Apropiacion.Id) + "," + strconv.Itoa(rubros_a_comprobar.FuenteFinanciacion.Id) + "," + strconv.FormatFloat(rubros_a_comprobar.ValorAsignado, 'f', -1, 64) + ")."})

			}
		}
		reglas := FormatoReglas(predicados) + reglasBase
		fmt.Println("reglas: ", reglas)
		res := golog.GetBoolean(reglas, "aprobacion_rp("+strconv.Itoa(rp_a_registrar.Rubros[0].Disponibilidad.Id)+",Y).", "Y")
		if res {
			if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/registro_presupuestal?limit=1&query=Vigencia:"+strconv.FormatFloat(rp_a_registrar.Rp.Vigencia, 'f', -1, 64)+"&sortby=NumeroRegistroPresupuestal&order=desc", &ultimo_registro); err == nil {
				if ultimo_registro == nil {
					numero_asignado_rp = 1
				} else {
					numero_asignado_rp = ultimo_registro[0].NumeroRegistroPresupuestal + 1
				}
				rp_a_registrar.Rp.NumeroRegistroPresupuestal = numero_asignado_rp
				rp_a_registrar.Rp.FechaMovimiento = time.Now().Local()
			} else {

			}
			if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/registro_presupuestal", "POST", &comprobacion, &rp_a_registrar); err == nil {
				//fmt.Println("entro...", respuesta)
				var acumCDP float64
				acumCDP = 0
				var aux map[string]float64
				var disp_apr []models.DisponibilidadApropiacion
				fmt.Println("dip: ", rp_a_registrar.Rubros[0].Disponibilidad.Id)
				if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/disponibilidad_apropiacion?limit=0&query=Disponibilidad.Id:"+strconv.Itoa(rp_a_registrar.Rubros[0].Disponibilidad.Id), &disp_apr); err == nil {
					for _, rubros_a_comprobar := range disp_apr {
						datos := models.DatosRubroRegistroPresupuestal{Disponibilidad: rubros_a_comprobar.Disponibilidad,
							Apropiacion: rubros_a_comprobar.Apropiacion, FuenteFinanciacion: rubros_a_comprobar.FuenteFinanciamiento,
						}

						if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/disponibilidad/SaldoCdp", "POST", &aux, &datos); err == nil {
							acumCDP = acumCDP + aux["saldo"]
							fmt.Println("saldo: ", aux["saldo"])
						}

					}
				} else {
					fmt.Println("error: ", err)
				}

				if acumCDP == 0 {

					rp_a_registrar.Rubros[0].Disponibilidad.Estado.Id = 3
					alertas = append(alertas, "Estado del CDP Objetivo: Agotado")

				} else {
					rp_a_registrar.Rubros[0].Disponibilidad.Estado.Id = 2
					alertas = append(alertas, "Estado del CDP Objetivo: Parcialmente Comprometido")
				}
				sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/disponibilidad/"+strconv.Itoa(rp_a_registrar.Rubros[0].Disponibilidad.Id), "PUT", &respuesta, &rp_a_registrar.Rubros[0].Disponibilidad)
				rp_a_registrar.Rp.DatosSolicitud.Expedida = true
				var respuesta_mod interface{}
				sendJson("http://"+beego.AppConfig.String("argoService")+"solicitud_rp/"+strconv.Itoa(rp_a_registrar.Rp.Solicitud), "PUT", &respuesta_mod, &rp_a_registrar.Rp.DatosSolicitud)
				fmt.Println("Solicitud: ", rp_a_registrar.Rp.Solicitud)
				fmt.Println("respuesta: ", respuesta_mod)
				alertas = append(alertas, "RP registrado exitosamente con el consecutivo No. "+strconv.Itoa(rp_a_registrar.Rp.NumeroRegistroPresupuestal))
			} else {
				alertas[0] = "error"
				alertas = append(alertas, "No se pudo registrar el RP")
			}

		} else {
			alertas[0] = "error"
			alertas = append(alertas, "Algunos valores asignados superan el saldo del CDP para la(s) apropiaciones asignadas")
		}
		c.Data["json"] = alertas
	} else {
		c.Data["json"] = err
		fmt.Println("error1: ", err)
	}

	c.ServeJSON()

}

// CargueMasivoPr ...
// @Title CargueMasivoPr
// @Description create RegistroPresupuestal
// @Param	body		body 	[]models.DatosRegistroPresupuestal	true		"body for DatosRegistroPresupuestal content"
// @Success 201 {object} models.DatosRegistroPresupuestal
// @Failure 403 body is empty
// @router /CargueMasivoPr [post]
func (c *RegistroPresupuestalController) CargueMasivoPr() {
	//variables a usar globales en el proceso.
	var dataRpRegistro []models.DatosRegistroPresupuestal
	var dataAlertas []models.Alert //array con las alertas generadas en aprobacion masiva de solicitudes
	var saldoCDP map[string]float64
	var comprobacion models.DatosRegistroPresupuestal
	tool := new(tools.EntornoReglas)
	var respuestaServices interface{}
	//------------------------------------------------------
	tool.Agregar_dominio("Presupuesto")
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &dataRpRegistro); err == nil {
		for _, rp_a_registrar := range dataRpRegistro { //recorrer el array de solicitudes cargadas
			for _, rubros_a_comprobar := range rp_a_registrar.Rubros { //recorrer la afectacion de la solicitud para inyeccion de reglas.
				datos := models.DatosRubroRegistroPresupuestal{Disponibilidad: rubros_a_comprobar.Disponibilidad,
					Apropiacion: rubros_a_comprobar.Apropiacion, FuenteFinanciacion: rubros_a_comprobar.FuenteFinanciacion,
				}
				if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/disponibilidad/SaldoCdp", "POST", &saldoCDP, &datos); err == nil {
					fmt.Println(rubros_a_comprobar.FuenteFinanciacion)
					if rubros_a_comprobar.FuenteFinanciacion == nil {
						tool.Agregar_predicado("rubro_cdp(" + strconv.Itoa(rubros_a_comprobar.Disponibilidad.Id) + "," + strconv.Itoa(datos.Apropiacion.Id) + "," + strconv.Itoa(0) +
							"," + strconv.FormatFloat(saldoCDP["saldo"], 'f', -1, 64) + ").")
					} else {
						tool.Agregar_predicado("rubro_cdp(" + strconv.Itoa(rubros_a_comprobar.Disponibilidad.Id) + "," + strconv.Itoa(datos.Apropiacion.Id) + "," + strconv.Itoa(rubros_a_comprobar.FuenteFinanciacion.Id) +
							"," + strconv.FormatFloat(saldoCDP["saldo"], 'f', -1, 64) + ").")
					}

				} else {
					dataAlertas = append(dataAlertas, models.Alert{Code: "E_0458", Body: rp_a_registrar, Type: "error"})
				}
				if rubros_a_comprobar.FuenteFinanciacion == nil {
					tool.Agregar_predicado("valor_rubro_rp(" + strconv.Itoa(rubros_a_comprobar.Disponibilidad.Id) + "," + strconv.Itoa(datos.Apropiacion.Id) + "," + strconv.Itoa(0) + "," + strconv.FormatFloat(rubros_a_comprobar.ValorAsignado, 'f', -1, 64) + ").")

				} else {
					tool.Agregar_predicado("valor_rubro_rp(" + strconv.Itoa(rubros_a_comprobar.Disponibilidad.Id) + "," + strconv.Itoa(datos.Apropiacion.Id) + "," + strconv.Itoa(rubros_a_comprobar.FuenteFinanciacion.Id) + "," + strconv.FormatFloat(rubros_a_comprobar.ValorAsignado, 'f', -1, 64) + ").")

				}
			}
			var res string
			fmt.Println(rp_a_registrar)
			err := utilidades.FillStruct(tool.Ejecutar_result("aprobacion_rp("+strconv.Itoa(rp_a_registrar.Rubros[0].Disponibilidad.Id)+",Y).", "Y"), &res)
			if err == nil { //
				if res == "1" { // si se aprueba la solicitud
					rp_a_registrar.Rp.FechaMovimiento = time.Now().Local()
					if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/registro_presupuestal", "POST", &comprobacion, &rp_a_registrar); err == nil {
						dataAlertas = append(dataAlertas, models.Alert{Code: "S_543", Body: comprobacion, Type: "success"})
						rp_a_registrar.Rp.DatosSolicitud.Expedida = true
						if err := sendJson("http://"+beego.AppConfig.String("argoService")+"solicitud_rp/"+strconv.Itoa(rp_a_registrar.Rp.Solicitud), "PUT", &respuestaServices, &rp_a_registrar.Rp.DatosSolicitud); err == nil {
							//dataAlertas = append(dataAlertas, models.Alert{Code: "S_RP002", Body: respuestaServices, Type: "success"})

						} else {
							dataAlertas = append(dataAlertas, models.Alert{Code: "E_RP002", Body: rp_a_registrar, Type: "success"})
						}
					} else {
						dataAlertas = append(dataAlertas, models.Alert{Code: "E_0458", Body: rp_a_registrar, Type: "error"})
					}
				} else {
					dataAlertas = append(dataAlertas, models.Alert{Code: "E_RP001", Body: rp_a_registrar, Type: "error"})

				}
			} else {
				dataAlertas = append(dataAlertas, models.Alert{Code: "E_0458", Body: rp_a_registrar, Type: "error"})
			}
			//res := golog.GetBoolean(reglas, "aprobacion_rp("+strconv.Itoa(rp_a_registrar.Rubros[0].Disponibilidad.Id)+",Y).", "Y")
		}
	} else {
		fmt.Println("err 2 ", err.Error())
		dataAlertas = append(dataAlertas, models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"})
	}
	c.Data["json"] = dataAlertas //respuesta de las alertas generadas durante el proceso.
	c.ServeJSON()

}

// GetSaldoRp ...
// @Title GetSaldoRp
// @Description get saldo rp by apropiacion
// @Param	body		body 	models.DatosRegistroPresupuestal	true		"body for DatosRegistroPresupuestal content"
// @Success 200 {object} models.RegistroPresupuestal
// @Failure 403 :id is empty
// @router /GetSaldoRp [get]
func (c *RegistroPresupuestalController) GetSaldoRp() {
	var v models.DatosRegistroPresupuestal
	var predicados []models.Predicado
	reglasInyectadas := ""
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		for _, rpRubro := range v.Rubros {

			//cargar ordenes de pago elaboradas para el rp
			var conceptoordenpago []models.ConceptoOrdenPago
			if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/concepto_orden_pago?fields=Valor,OrdenDePago&limit=0&query=Concepto.Rubro.Id:"+strconv.FormatInt(rpRubro.Apropiacion.Rubro.Id, 10)+",OrdenPago.RegistroPresupuestal.Id:"+strconv.Itoa(v.Rp.Id), &conceptoordenpago); err == nil {
				if conceptoordenpago != nil {
					for _, ordenpago := range conceptoordenpago {
						predicados = append(predicados, models.Predicado{Nombre: "ordenPagoRP(" + strconv.Itoa(ordenpago.OrdenDePago.RegistroPresupuestal.Id) + "," + strconv.Itoa(rpRubro.Apropiacion.Id) + "," + strconv.FormatFloat(ordenpago.Valor, 'f', -1, 64) + ")."})
					}

				} else {
					//si no se encuentra nada en orden_pago_conceptos
				}

			} else {
				//error al consumir servicio de ordenes de pago
			}
			//cargar anulaciones del rp en la apropiacion indicada
			var anulacionesrp []models.AnulacionRegistroPresupuestalDisponibilidadApropiacion
			if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/anulacion_registro_presupuestal_disponibilidad_apropiacion?fields=Valor,RegistroPresupuestalDisponibilidadApropiacion&limit=0&query=RegistroPresupuestalDisponibilidadApropiacion.RegistroPresupuestal.Id:"+strconv.Itoa(v.Rp.Id)+",RegistroPresupuestalDisponibilidadApropiacion.DisponibilidadApropiacion.Apropiacion.Id:"+strconv.Itoa(rpRubro.Apropiacion.Id), &anulacionesrp); err == nil {
				if anulacionesrp != nil {
					for _, anulacion := range anulacionesrp {
						predicados = append(predicados, models.Predicado{Nombre: "anulacionRP(" + strconv.Itoa(anulacion.RegistroPresupuestalDisponibilidadApropiacion.RegistroPresupuestal.Id) + "," + strconv.Itoa(anulacion.RegistroPresupuestalDisponibilidadApropiacion.DisponibilidadApropiacion.Apropiacion.Id) + "," + strconv.FormatFloat(anulacion.Valor, 'f', -1, 64) + ")."})
					}
				} else {
					//si no encuentra anulaciones del rp
				}
			} else {
				//error al consumir servicio de anulaciones rp
			}
			//cargar valor original del rp para la apropiacion
			var valoresrp []models.RegistroPresupuestalDisponibilidadApropiacion
			if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/registro_presupuestal_disponibilidad_apropiacion?fields=Valor&limit=0&query=RegistroPresupuestal.Id:"+strconv.Itoa(v.Rp.Id)+",DisponibilidadApropiacion.Apropiacion.Id:"+strconv.Itoa(rpRubro.Apropiacion.Id), &valoresrp); err == nil {
				if valoresrp != nil {
					for _, valorrp := range valoresrp {
						predicados = append(predicados, models.Predicado{Nombre: "valorRP(" + strconv.Itoa(valorrp.RegistroPresupuestal.Id) + "," + strconv.Itoa(valorrp.DisponibilidadApropiacion.Apropiacion.Id) + "," + strconv.FormatFloat(valorrp.Valor, 'f', -1, 64) + ")."})
					}
				} else {
					//si no hay valores del rp
				}
			} else {
				//error al consumir el servicio de valores rp
			}
			reglasInyectadas = FormatoReglas(predicados)
			c.Data["json"] = reglasInyectadas
		}
	} else {

	}
	c.ServeJSON()
}

// GetAll ...
// @Title GetAll
// @Description get RegistroPresupuestal
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.RegistroPresupuestal
// @Failure 403
// @router / [get]
func (c *RegistroPresupuestalController) GetAll() {

}

// Put ...
// @Title Put
// @Description update the RegistroPresupuestal
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.RegistroPresupuestal	true		"body for RegistroPresupuestal content"
// @Success 200 {object} models.RegistroPresupuestal
// @Failure 403 :id is not int
// @router /:id [put]
func (c *RegistroPresupuestalController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the RegistroPresupuestal
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *RegistroPresupuestalController) Delete() {

}
