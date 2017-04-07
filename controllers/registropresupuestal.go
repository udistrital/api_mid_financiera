package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/udistrital/api_mid_financiera/golog"
	"github.com/udistrital/api_mid_financiera/models"
	"strconv"
	"time"

	"github.com/astaxie/beego"
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
	var saldoCDP float64
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &rp_a_registrar); err == nil {

		reglasBase := CargarReglasBase("Presupuesto")
		for _, rubros_a_comprobar := range rp_a_registrar.Rubros {

			datos := models.DatosRubroRegistroPresupuestal{Disponibilidad: rubros_a_comprobar.Disponibilidad,
				Apropiacion: rubros_a_comprobar.Apropiacion,
			}
			if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/disponibilidad/SaldoCdp", "POST", &saldoCDP, &datos); err == nil {
				predicados = append(predicados, models.Predicado{Nombre: "rubro_cdp(" + strconv.Itoa(rubros_a_comprobar.Disponibilidad.Id) + "," + strconv.Itoa(datos.Apropiacion.Id) + "," + strconv.FormatFloat(saldoCDP, 'f', -1, 64) + ")."})
			} else {
				alertas[0] = "error"
				alertas = append(alertas, "No se pudo cargar el saldo para algunas apropiaciones")
				c.Data["json"] = alertas

				c.ServeJSON()
			}
			predicados = append(predicados, models.Predicado{Nombre: "valor_rubro_rp(" + strconv.Itoa(rubros_a_comprobar.Disponibilidad.Id) + "," + strconv.Itoa(datos.Apropiacion.Id) + "," + strconv.FormatFloat(rubros_a_comprobar.ValorAsignado, 'f', -1, 64) + ")."})
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
				var aux float64
				for _, rubros_a_comprobar := range rp_a_registrar.Rubros {
					datos := models.DatosRubroRegistroPresupuestal{Disponibilidad: rubros_a_comprobar.Disponibilidad,
						Apropiacion: rubros_a_comprobar.Apropiacion,
					}

					if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/disponibilidad/SaldoCdp", "POST", &aux, &datos); err == nil {
						acumCDP = acumCDP + aux
					}

				}
				if acumCDP == 0 {

					rp_a_registrar.Rubros[0].Disponibilidad.Estado.Id = 3
					alertas = append(alertas, "Estado del CDP Objetivo: Agotado")

				} else {
					rp_a_registrar.Rubros[0].Disponibilidad.Estado.Id = 2
					alertas = append(alertas, "Estado del CDP Objetivo: Parcialmente Comprometido")
				}
				sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/disponibilidad/"+strconv.Itoa(rp_a_registrar.Rubros[0].Disponibilidad.Id), "PUT", &respuesta, &rp_a_registrar.Rubros[0].Disponibilidad)

				alertas = append(alertas, "RP registrado exitosamente")
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
