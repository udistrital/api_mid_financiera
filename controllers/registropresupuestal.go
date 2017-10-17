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

func formatoSolicitudRP(solicitudintfc interface{}) (res interface{}) {
	//recuperar datos del CDP objetivo de la solicitud
	var rubros []interface{}
	solicitud := models.SolicitudRp{}
	err := utilidades.FillStruct(solicitudintfc, &solicitud)
	fmt.Println(err)
	var afectacion_solicitud []map[string]interface{}
	if err := getJson("http://"+beego.AppConfig.String("argoService")+"disponibilidad_apropiacion_solicitud_rp?limit=-1&query=SolicitudRp:"+strconv.Itoa(solicitud.Id), &afectacion_solicitud); err == nil {
		//consulta de la afectacion presupuestal objetivo.
		fmt.Println(solicitud.Id)
		for _, afect := range afectacion_solicitud {

			var disp_apr_sol []map[string]interface{}
			if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/disponibilidad_apropiacion?limit=1&query=Id:"+fmt.Sprintf("%v", afect["DisponibilidadApropiacion"]), &disp_apr_sol); err == nil {

				for _, disp_apro := range disp_apr_sol {
					disp_apro["ValorAsignado"] = afect["Monto"]
					disp_apro["FuenteFinanciacion"] = disp_apro["FuenteFinanciamiento"]
					rubros = append(rubros, disp_apro)
				}

			} else {
				//si sale mal la consulta de la afectacion del cdp objetivo.

			}
		}
		solicitud.Rubros = rubros
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
	return solicitud
}

//funcion para recopilar datos externos de los rp a listar
func FormatoListaRP(rpintfc interface{}) (res interface{}) {
	rp := rpintfc.(map[string]interface{})
	idSolicitudDisponibilidad := int(rp["RegistroPresupuestalDisponibilidadApropiacion"].([]interface{})[0].(map[string]interface{})["DisponibilidadApropiacion"].(map[string]interface{})["Disponibilidad"].(map[string]interface{})["Solicitud"].(float64))
	solicituddisp, err := DetalleSolicitudDisponibilidadById(strconv.Itoa(idSolicitudDisponibilidad))

	if err == nil {
		rp["InfoSolicitudDisponibilidad"] = solicituddisp
		return rp
	}
	return nil
}

// ListaRp ...
// @Title ListaRp
// @Description get RP by vigencia
// @Param	vigencia	query	string	false	"vigencia de la lista"
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Param	rangoinicio	query	string	false	"rango inicial del periodo a consultar"
// @Param	rangofin	query	string	false	"rango final del periodo a consultar"
// @Success 200 {object} models.RegistroPresupuestal
// @Failure 403
// @router ListaRp/:vigencia [get]
func (c *RegistroPresupuestalController) ListaRp() {
	vigenciaStr := c.Ctx.Input.Param(":vigencia")
	vigencia, err := strconv.Atoi(vigenciaStr)
	var rpresupuestal []interface{}
	var respuesta []map[string]interface{}
	var limit int64 = 10
	var offset int64
	var startrange string
	var endrange string
	var query string
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	if r := c.GetString("rangoinicio"); r != "" {
		startrange = r

	}

	if r := c.GetString("rangofin"); r != "" {
		endrange = r

	}
	if startrange != "" && endrange != "" {
		query = ",FechaRegistro__gte:" + startrange + ",FechaRegistro__lte:" + endrange

	}
	if err = getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/registro_presupuestal?limit="+strconv.FormatInt(limit, 10)+"&offset="+strconv.FormatInt(offset, 10)+"&query=Vigencia:"+strconv.Itoa(vigencia)+query, &rpresupuestal); err == nil {
		if rpresupuestal != nil {
			done := make(chan interface{})
			defer close(done)
			resch := utilidades.GenChanInterface(rpresupuestal...)
			chrpresupuestal := utilidades.Digest(done, FormatoListaRP, resch)
			for rp := range chrpresupuestal {
				respuesta = append(respuesta, rp.(map[string]interface{}))
			}
			c.Data["json"] = respuesta
		} else {
			c.Data["json"] = models.Alert{Code: "E_0458", Body: nil, Type: "error"}
		}
	} else {
		c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
	}
	c.ServeJSON()
}

// GetSolicitudesRp ...
// @Title GetSolicitudesRp
// @Description get saldo rp by apropiacion
// @Param	vigencia		path 	int 	true		"vigencia de las solicitudes a consultar"
// @Success 200 {object} models.SolicitudRp
// @Failure 403 :vigencia is empty
// @router /GetSolicitudesRp/:vigencia [get]
func (c *RegistroPresupuestalController) GetSolicitudesRp() {
	vigenciaStr := c.Ctx.Input.Param(":vigencia")
	vigencia, err := strconv.Atoi(vigenciaStr)
	if err == nil {
		var solicitudes_rp []interface{}
		var respuesta []models.SolicitudRp
		if err := getJson("http://"+beego.AppConfig.String("argoService")+"solicitud_rp?limit=-1&query=Expedida:false,Vigencia:"+strconv.Itoa(vigencia)+"&sortby=Id&order=desc", &solicitudes_rp); err == nil {
			if solicitudes_rp != nil {
				//encontrar datos del CDP objetivo del RP Solicitado

				done := make(chan interface{})
				defer close(done)
				resch := utilidades.GenChanInterface(solicitudes_rp...)
				chsolicitud := utilidades.Digest(done, formatoSolicitudRP, resch)
				for solicitud := range chsolicitud {
					respuesta = append(respuesta, solicitud.(models.SolicitudRp))
				}
				c.Data["json"] = respuesta
			} else {
				//si no hay datos de solicitudes
				c.Data["json"] = "sin datos"
			}
		} else {
			//si ocurre error al traer las solicitudes
			c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
			fmt.Println(err.Error())
		}
	} else {
		c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
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

func ListaNecesidadesByRp(solicitudintfc interface{}) (res interface{}) {
	solicitud, e := solicitudintfc.(map[string]interface{})
	var rp []map[string]interface{}
	if e {
		idSol, e := solicitud["Id"].(float64)
		if e {
			if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/registro_presupuestal?query=RegistroPresupuestalDisponibilidadApropiacion.DisponibilidadApropiacion.Disponibilidad.Solicitud:"+strconv.Itoa(int(idSol)), &rp); err == nil {
				if rp != nil {
					solicitud["InfoRp"] = rp[0]
					return solicitud
				} else {
					return nil
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

// ListaNecesidadesByRp ...
// @Title ListaNecesidadesByRp
// @Description Lsta de las necesidades origen de los rp registrados
// @Param	vigencia	query	string	false	"vigencia de la lista"
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Param	rangoinicio	query	string	false	"rango inicial del periodo a consultar"
// @Param	rangofin	query	string	false	"rango final del periodo a consultar"
// @Param	tipoNecesidad	query	string	false	"tipo de la necesidad origen del rp"
// @Success 200 {object} models.Alert
// @Failure 403 body is empty
// @router /ListaNecesidadesByRp/:vigencia [get]
func (c *RegistroPresupuestalController) ListaNecesidadesByRp() {
	vigenciaStr := c.Ctx.Input.Param(":vigencia")
	vigencia, err := strconv.Atoi(vigenciaStr)
	var solicitudNecesidad []interface{}
	var respuesta []map[string]interface{}
	var limit int64 = 10
	var offset int64
	var startrange string
	var endrange string
	var query string
	var tipoNecesidad string
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	if r := c.GetString("rangoinicio"); r != "" {
		startrange = r

	}

	if r := c.GetString("rangofin"); r != "" {
		endrange = r

	}
	if startrange != "" && endrange != "" {
		query = ",FechaSolicitud__gte:" + startrange + ",FechaSolicitud__lte:" + endrange

	}
	if r := c.GetString("tipoNecesidad"); r != "" {
		tipoNecesidad = r
		//peticion a los rp expedidos
		if err = getJson("http://"+beego.AppConfig.String("argoService")+"solicitud_disponibilidad?limit="+strconv.FormatInt(limit, 10)+"&offset="+strconv.FormatInt(offset, 10)+"&query=Vigencia:"+strconv.Itoa(vigencia)+",Necesidad.TipoNecesidad.CodigoAbreviacion:"+tipoNecesidad+query, &solicitudNecesidad); err == nil {
			if solicitudNecesidad != nil {
				done := make(chan interface{})
				defer close(done)
				resch := utilidades.GenChanInterface(solicitudNecesidad...)
				chsolicitud := utilidades.Digest(done, ListaNecesidadesByRp, resch)
				for solicitud := range chsolicitud {
					if solicitud != nil {
						respuesta = append(respuesta, solicitud.(map[string]interface{}))
					}
				}
				c.Data["json"] = respuesta
			} else {
				//si no encuentra solicitudes...
				c.Data["json"] = models.Alert{Code: "E_0458", Body: nil, Type: "error"}
			}
		} else {
			//si ocurre error al traer las solicitudes
			c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
			fmt.Println(err.Error())
		}

	} else {
		c.Data["json"] = models.Alert{Code: "E_0458", Body: "TipoNecesidad lost", Type: "error"}
	}

	c.ServeJSON()

}
