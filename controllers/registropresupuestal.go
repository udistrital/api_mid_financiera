package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego"
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

				var afectacion_solicitud []map[string]interface{}
				if err := getJson("http://"+beego.AppConfig.String("argoService")+"disponibilidad_apropiacion_solicitud_rp?limit=0&query=SolicitudRp:"+strconv.Itoa(solicitud.Id), &afectacion_solicitud); err == nil {
					//consulta de la afectacion presupuestal objetivo.
					for _, afect := range afectacion_solicitud {
						var disp_apr_sol []map[string]interface{}
						if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/disponibilidad_apropiacion?limit=1&query=Id:"+fmt.Sprintf("%v", afect["DisponibilidadApropiacion"]), &disp_apr_sol); err == nil {
							var rubros []interface{}
							for _, disp_apro := range disp_apr_sol {
								disp_apro["ValorAsignado"] = afect["Monto"]
								disp_apro["FuenteFinanciacion"] = disp_apro["FuenteFinanciamiento"]
								rubros = append(rubros, disp_apro)
							}
							solicitud.Rubros = rubros
						} else {
							//si sale mal la consulta de la afectacion del cdp objetivo.

						}
					}
				} else {
					//si sale mal la consulta de la afectacion de la solicitud.
				}

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

				var afectacion_solicitud []map[string]interface{}
				if err := getJson("http://"+beego.AppConfig.String("argoService")+"disponibilidad_apropiacion_solicitud_rp?limit=0&query=SolicitudRp:"+strconv.Itoa(solicitud.Id), &afectacion_solicitud); err == nil {
					//consulta de la afectacion presupuestal objetivo.
					for _, afect := range afectacion_solicitud {
						var disp_apr_sol []map[string]interface{}
						if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/disponibilidad_apropiacion?limit=1&query=Id:"+fmt.Sprintf("%v", afect["DisponibilidadApropiacion"]), &disp_apr_sol); err == nil {
							for _, disp_apro := range disp_apr_sol {
								disp_apro["ValorAsignado"] = afect["Monto"]
								disp_apro["FuenteFinanciacion"] = disp_apro["FuenteFinanciamiento"]
							}
						} else {
							//si sale mal la consulta de la afectacion del cdp objetivo.

						}
					}
				} else {
					//si sale mal la consulta de la afectacion de la solicitud.
				}

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

// CargueMasivoPr ...
// @Title CargueMasivoPr
// @Description create RegistroPresupuestal
// @Param	body		body 	models.DatosRegistroPresupuestal	true		"body for DatosRegistroPresupuestal content"
// @Success 200 {object} models.Alert
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
			if rp_a_registrar.Rubros != nil {
				err := utilidades.FillStruct(tool.Ejecutar_result("aprobacion_rp("+strconv.Itoa(rp_a_registrar.Rubros[0].Disponibilidad.Id)+",Y).", "Y"), &res)
				if err == nil { //
					if res == "1" { // si se aprueba la solicitud
						rp_a_registrar.Rp.FechaRegistro = time.Now().Local()
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
