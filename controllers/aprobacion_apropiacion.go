package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/udistrital/api_mid_financiera/golog"
	"github.com/udistrital/api_mid_financiera/models"
	"strconv"

	"github.com/astaxie/beego"
)

type AprobacionController struct {
	beego.Controller
}

func (c *AprobacionController) URLMapping() {
	c.Mapping("Aprobar", c.Aprobar)
}

// Aprobar ...
// @Title Aprobar
// @Description Aprobar Apropiacion
// @Param	body		body 	[]models.Apropiacion	true		"body for Apropiacion  content"
// @Success 201 {int} models.InfoSolDisp
// @Failure 403 body is empty
// @router Aprobar/ [post]
func (this *AprobacionController) Aprobar() {

	var predicados []models.Predicado
	//var predicados_apropiacion []models.Predicado
	var alertas []string
	alertas = append(alertas, "success")
	var postdominio string = ""

	if tdominio := this.GetString("tdominio"); tdominio != "" {
		postdominio = postdominio + "&query=Dominio.Id:" + tdominio
	} else {
		this.Data["json"] = "no se especifico el domino del ruler"
		this.ServeJSON()
	}
	var apropiacion []models.Apropiacion
	//var respuesta interface{}
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &apropiacion); err == nil {
		if err := getJson("http://"+beego.AppConfig.String("Urlruler")+":"+beego.AppConfig.String("Portruler")+"/"+beego.AppConfig.String("Nsruler")+"/predicado?limit=0"+postdominio, &predicados); err == nil {
			//var reglas string = ""
			var reglasbase string = ""
			//var reglasinyectadas string = ""
			var arregloReglas = make([]string, len(predicados))
			var comprobacion string
			var estado_ap int
			//var respuesta []models.FormatoPreliqu
			for i := 0; i < len(predicados); i++ {
				arregloReglas[i] = predicados[i].Nombre
			}

			for i := 0; i < len(arregloReglas); i++ {
				reglasbase = reglasbase + arregloReglas[i]
			}
			for i := len(apropiacion) - 1; i >= 0; i-- {

				comprobacion = comprobar_apropiacion(apropiacion[i])
				if comprobacion == "" {
					alertas = append(alertas, "Apropiacion del rubro "+apropiacion[i].Rubro.Codigo+" No aprobada, algunas apropiaciones hijo no se encuentran aprobadas")
				} else {
					estado_ap, _ = strconv.Atoi(golog.Comprobar_aprobacion(reglasbase, comprobacion))
					estado := models.EstadoApropiacion{Id: estado_ap}
					apropiacion[i].Estado = &estado
					var respuesta interface{}
					if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/apropiacion/"+strconv.Itoa(apropiacion[i].Id), "PUT", &respuesta, &apropiacion[i]); err == nil {
						if estado_ap == 1 {
							alertas[0] = "error"
							alertas = append(alertas, "Apropiacion del rubro "+apropiacion[i].Rubro.Codigo+" No aprobada, la suma de las apropiaciones hijo no corresponde al valor de esta apropiacion")
						} else {
							alertas = append(alertas, "Apropiacion del rubro "+apropiacion[i].Rubro.Codigo+" Aprobada")
						}
					} else {
						alertas[0] = "error"
						alertas = append(alertas, "no se pudo cambiar el estado de la apropiacion")
						this.Data["json"] = alertas
						this.ServeJSON()
					}
				}
			}
			this.Data["json"] = alertas
			this.ServeJSON()
		} else {
			fmt.Println("err: ", err)
			alertas[0] = "error"
			alertas = append(alertas, "no se pudo cargar la informacion de la base de datos")
			this.Data["json"] = alertas
			this.ServeJSON()
		}
	} else {
		fmt.Println("err: ", err)
		alertas[0] = "error"
		alertas = append(alertas, "no se resivieron los datos correctamente")
		this.Data["json"] = alertas
		this.ServeJSON()
	}

}

func comprobar_apropiacion(padre models.Apropiacion) string {
	var rubro_hijo []models.RubroRubro
	var lista_valores []string
	var regla string
	var apropiacion_hijo []models.Apropiacion
	var hoja int
	getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/rubro_rubro?limit=0&query=RubroPadre.Id:"+strconv.FormatInt(padre.Rubro.Id, 10)+",RubroPadre.Vigencia:"+strconv.FormatFloat(padre.Vigencia, 'f', -1, 64), &rubro_hijo)
	if rubro_hijo != nil {
		for i := 0; i < len(rubro_hijo); i++ {

			getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/apropiacion?limit=0&query=Vigencia:"+strconv.FormatFloat(padre.Vigencia, 'f', -1, 64)+",Rubro.Id:"+strconv.FormatInt(rubro_hijo[i].RubroHijo.Id, 10)+"", &apropiacion_hijo)
			if apropiacion_hijo != nil {
				hoja = 0
				for i := 0; i < len(apropiacion_hijo); i++ {
					if apropiacion_hijo[i].Estado.Id == 2 {
						lista_valores = append(lista_valores, strconv.FormatFloat(apropiacion_hijo[i].Valor, 'f', -1, 64))
						fmt.Println("apro: ", lista_valores)
					}

				}

			}
		}
	} else {
		lista_valores = append(lista_valores, strconv.FormatFloat(padre.Valor, 'f', -1, 64))
		hoja = 1
	}
	fmt.Println("hijo: ", len(apropiacion_hijo))
	fmt.Println("valor: ", len(lista_valores))
	if lista_valores != nil && len(apropiacion_hijo) > 0 {
		for i := 0; i < len(lista_valores); i++ {
			if len(lista_valores) == 1 {
				regla = "verifica_hijos([" + lista_valores[i] + "],[" + strconv.FormatFloat(padre.Valor, 'f', -1, 64) + ", 3],Y)."
			} else if len(lista_valores) > 1 {
				if i == 0 {
					regla = "verifica_hijos([" + lista_valores[i]
				} else if i == (len(lista_valores) - 1) {
					regla = regla + "," + lista_valores[i] + "],[" + strconv.FormatFloat(padre.Valor, 'f', -1, 64) + ", 3],Y)."
				} else {
					regla = regla + "," + lista_valores[i]
				}
			}
		}
	}
	if hoja == 1 {
		for i := 0; i < len(lista_valores); i++ {
			if len(lista_valores) == 1 {
				regla = "verifica_hijos([" + lista_valores[i] + "],[" + strconv.FormatFloat(padre.Valor, 'f', -1, 64) + ", 3],Y)."
			} else if len(lista_valores) > 1 {
				if i == 0 {
					regla = "verifica_hijos([" + lista_valores[i]
				} else if i == (len(lista_valores) - 1) {
					regla = regla + "," + lista_valores[i] + "],[" + strconv.FormatFloat(padre.Valor, 'f', -1, 64) + ", 3],Y)."
				} else {
					regla = regla + "," + lista_valores[i]
				}
			}
		}
	}

	return regla
}
