package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/mitchellh/mapstructure"
	"github.com/udistrital/api_mid_financiera/models"
	"github.com/udistrital/api_mid_financiera/utilidades"
)

type RubroController struct {
	beego.Controller
}
type rowPac struct {
	Codigo      interface{}
	Descripcion interface{}
	Fdescrip    interface{}
	Id          interface{}
	Idfuente    interface{}
	Idrubro     interface{}
	Reporte     []*reportePacData
}
type reportePacData struct {
	Mes     interface{}
	N_mes   interface{}
	Valores *valoresPac
}
type valoresPac struct {
	Proyeccion interface{}
	Pvariacion interface{}
	Valor      interface{}
	Variacion  interface{}
}
type cuerpoPac struct {
	Ingresos []*rowPac
}

// GenerarPac ...
// @Title GenerarPac
// @Description Get PAC By Rubro
// @Param	pacData		query 	interface{}	true		"objeto con fechas del rango del PAC"
// @Success 200 {object} interface{}
// @Failure 403
// @router /GenerarPac/ [post]
func (c *RubroController) GenerarPac() {
	defer c.ServeJSON()

	var pacData map[string]interface{} //definicion de la interface que recibe los datos del reporte y proyecciones
	var finicio time.Time
	var ffin time.Time
	var periodos int
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &pacData); err == nil {
		err := utilidades.FillStruct(pacData["inicio"], &finicio)
		err = utilidades.FillStruct(pacData["fin"], &ffin)
		err = utilidades.FillStruct(pacData["fin"], &periodos)
		if err != nil {
			if reporteData, err := cuerpoReporte(finicio, ffin); err == nil {

				var alert models.Alert
				var proy []map[string]interface{}
				go c.calcularEjecutado(&reporteData, finicio, ffin, &alert, &proy)

				fmt.Println(proy)
				if alert.Body == nil {
					fmt.Println("no alert")
				} else {
					fmt.Println("alert ", alert)
				}
				c.Data["json"] = reporteData

			} else {
				fmt.Println("err 2")
				c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
			}
		} else {
			fmt.Println("err 1 ", err.Error())
			c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
		}

	} else {
		fmt.Println("err 1")
		c.Data["json"] = models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
	}

}
func (c *RubroController) calcularEjecutado(reporteData *cuerpoPac, finicio time.Time, ffin time.Time, alert *models.Alert, res *[]map[string]interface{}) {

	for _, ingresosRow := range reporteData.Ingresos { //recorrer los datos del reporte de ingresos para el rango actual

		for _, reporteRow := range ingresosRow.Reporte {
			var valor string
			var mes int
			err := utilidades.FillStruct(reporteRow.Valores.Valor, &valor)
			err = utilidades.FillStruct(reporteRow.N_mes, &mes)
			if err == nil {
				fechaInicio := time.Date(finicio.Year(), time.Month(mes), finicio.Day(), 0, 0, 0, 0, time.Local)
				fechaFin := time.Date(finicio.Year(), time.Month(mes+1), finicio.Day(), 0, 0, 0, 0, time.Local)

				if fechaFin.After(ffin) {
					fechaFin = ffin
				}
				var rubro string
				var idFuente string
				err := utilidades.FillStruct(ingresosRow.Idrubro, &rubro)

				err = utilidades.FillStruct(ingresosRow.Idfuente, &idFuente)

				if err == nil {

					if rubro == "35773" {
						fmt.Println("rubro: ", rubro)
						fmt.Println("Fuente: ", idFuente)
						fmt.Println("finicio: ", fechaInicio)
						fmt.Println("ffin: ", fechaFin)
						fmt.Println("url ", "http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/rubro/GetRubroIngreso?rubro="+rubro+"&fuente="+idFuente+"&finicio="+fechaInicio.Format("2006-01-02")+"&ffin="+fechaFin.Format("2006-01-02"))

					}
					var valorIngresos interface{}
					if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/rubro/GetRubroIngreso?rubro="+rubro+"&fuente="+idFuente+"&finicio="+fechaInicio.Format("2006-01-02")+"&ffin="+fechaFin.Format("2006-01-02"), &valorIngresos); err == nil {
						var dataIngresos []map[string]interface{}
						err := utilidades.FillStruct(valorIngresos, &dataIngresos)
						if err != nil {

						} else {
							for _, valorData := range dataIngresos {
								//fmt.Println("rubroProyData(" + fmt.Sprintf("%v", ingresosRow.Idrubro) + "," + fmt.Sprintf("%v", fechaInicio.Year()) + "," + fmt.Sprintf("%v", int(fechaInicio.Month())) + "," + fmt.Sprintf("%v", valorData["valor"]) + ").")
								utilidades.FillStruct(valorData["valor"], &reporteRow.Valores.Valor)
							}

						}

					} else {
						fmt.Println("err v", err.Error())
						alert = &models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
					}
				} else {
					fmt.Println("err ", err.Error())
					alert = &models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
				}

			} else {
				fmt.Println("err 2 ", err.Error())
				alert = &models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
			}
		}

	}

	return
}

func cuerpoReporte(inicio time.Time, fin time.Time) (res cuerpoPac, err error) {

	mesinicio := int(inicio.Month())
	mesfin := int(fin.Month())
	var m []map[string]interface{}
	cuerpo := make(map[string]interface{})
	err = getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/apropiacion/GetApropiacionesHijo/"+strconv.Itoa(inicio.Year())+"?tipo=3", &m)
	if err != nil {
		return
	}

	for i := 0; i < len(m); i++ {
		var fechas []map[string]interface{}
		for j := 0; j <= (mesfin - mesinicio); j++ {
			finicio := inicio.AddDate(0, j, 0)
			aux := make(map[string]interface{})

			val := make(map[string]interface{})
			val["valor"] = "0"
			val["proyeccion"] = "0"
			val["variacion"] = "0"
			val["pvariacion"] = "0"
			aux["valores"] = val

			if aux != nil {
				aux["mes"] = finicio.Format("Jan")
				aux["n_mes"] = int(finicio.Month())
				fechas = append(fechas, aux)
			}

		}
		m[i]["reporte"] = fechas
		//m[i]["egresos"], err = RubroOrdenPago(m[i]["id"])
		if err != nil {
			fmt.Println("err1 ", err)
			return
		}

	}
	var ingresos interface{}
	err = utilidades.FillStruct(m, &ingresos)
	if err != nil {
		fmt.Println("err2 ", err)
		return
	}
	cuerpo["ingresos"] = ingresos
	err = mapstructure.Decode(cuerpo, &res)
	if err != nil {
		fmt.Println("err2 ", err)
		return
	}

	return
}
