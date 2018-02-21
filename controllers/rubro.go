package controllers

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/astaxie/beego"
	"github.com/mitchellh/mapstructure"
	"github.com/udistrital/api_mid_financiera/models"
	"github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/utils_oas/ruler"
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
	Egresos  []*rowPac
}

var (
	wg sync.WaitGroup
)

// GenerarPac ...
// @Title GenerarPac
// @Description Get PAC By Rubro
// @Param	pacData		query 	interface{}	true		"objeto con fechas del rango del PAC"
// @Success 200 {object} interface{}
// @Failure 403
// @router /GenerarPac/ [post]
func (c *RubroController) GenerarPac() {
	defer c.ServeJSON()
	wg.Add(4)
	var pacData map[string]interface{} //definicion de la interface que recibe los datos del reporte y proyecciones
	var finicio time.Time
	var ffin time.Time
	var periodos int
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &pacData); err == nil {
		err := formatdata.FillStruct(pacData["inicio"], &finicio)
		err = formatdata.FillStruct(pacData["fin"], &ffin)
		err = formatdata.FillStruct(pacData["periodosproy"], &periodos)
		if err == nil {
			if reporteData, err := cuerpoReporte(finicio, ffin); err == nil {
				var alert models.Alert
				go c.calcularEjecutadoIngresos(&reporteData, finicio, ffin, &alert)
				go c.calcularEjecutadoEngresos(&reporteData, finicio, ffin, &alert)
				go c.calcularProyeccionIngresos(&reporteData, finicio, ffin, periodos, &alert)
				go c.calcularProyeccionEgresos(&reporteData, finicio, ffin, periodos, &alert)
				wg.Wait()
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

func (c *RubroController) calcularEjecutadoIngresos(reporteData *cuerpoPac, finicio time.Time, ffin time.Time, alert *models.Alert) {

	for _, ingresosRow := range reporteData.Ingresos { //recorrer los datos del reporte de ingresos para el rango actual

		for _, reporteRow := range ingresosRow.Reporte {
			var valor string
			var mes int
			err := formatdata.FillStruct(reporteRow.Valores.Valor, &valor)
			err = formatdata.FillStruct(reporteRow.N_mes, &mes)
			if err == nil {
				fechaInicio := time.Date(finicio.Year(), time.Month(mes), finicio.Day(), 0, 0, 0, 0, time.Local)
				fechaFin := time.Date(finicio.Year(), time.Month(mes+1), finicio.Day(), 0, 0, 0, 0, time.Local)

				if fechaFin.After(ffin) {
					fechaFin = ffin
				}
				var rubro string
				var idFuente string
				err := formatdata.FillStruct(ingresosRow.Idrubro, &rubro)

				err = formatdata.FillStruct(ingresosRow.Idfuente, &idFuente)

				if err == nil {

					/*if rubro == "35488" {
						fmt.Println("rubro: ", rubro)
						fmt.Println("Fuente: ", idFuente)
						fmt.Println("finicio: ", fechaInicio)
						fmt.Println("ffin: ", fechaFin)
						fmt.Println("url ", "http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/rubro/GetRubroIngreso?rubro="+rubro+"&fuente="+idFuente+"&finicio="+fechaInicio.Format("2006-01-02")+"&ffin="+fechaFin.Format("2006-01-02"))

					}*/
					var valorIngresos interface{}
					if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/rubro/GetRubroIngreso?rubro="+rubro+"&fuente="+idFuente+"&finicio="+fechaInicio.Format("2006-01-02")+"&ffin="+fechaFin.Format("2006-01-02"), &valorIngresos); err == nil {
						var dataIngresos []map[string]interface{}
						err := formatdata.FillStruct(valorIngresos, &dataIngresos)
						if err != nil {

						} else {
							for _, valorData := range dataIngresos {
								//fmt.Println("rubroProyData(" + fmt.Sprintf("%v", ingresosRow.Idrubro) + "," + fmt.Sprintf("%v", fechaInicio.Year()) + "," + fmt.Sprintf("%v", int(fechaInicio.Month())) + "," + fmt.Sprintf("%v", valorData["valor"]) + ").")
								formatdata.FillStruct(valorData["valor"], &reporteRow.Valores.Valor)
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
	wg.Done()
	return
}

func (c *RubroController) calcularEjecutadoEngresos(reporteData *cuerpoPac, finicio time.Time, ffin time.Time, alert *models.Alert) {

	for _, ingresosRow := range reporteData.Egresos { //recorrer los datos del reporte de ingresos para el rango actual

		for _, reporteRow := range ingresosRow.Reporte {
			var valor string
			var mes int
			err := formatdata.FillStruct(reporteRow.Valores.Valor, &valor)
			err = formatdata.FillStruct(reporteRow.N_mes, &mes)
			if err == nil {
				fechaInicio := time.Date(finicio.Year(), time.Month(mes), finicio.Day(), 0, 0, 0, 0, time.Local)
				fechaFin := time.Date(finicio.Year(), time.Month(mes+1), finicio.Day(), 0, 0, 0, 0, time.Local)

				if fechaFin.After(ffin) {
					fechaFin = ffin
				}
				var rubro string
				var idFuente string
				err := formatdata.FillStruct(ingresosRow.Idrubro, &rubro)

				err = formatdata.FillStruct(ingresosRow.Idfuente, &idFuente)

				if err == nil {

					/*if rubro == "35585" {
						fmt.Println("rubro: ", rubro)
						fmt.Println("Fuente: ", idFuente)
						fmt.Println("finicio: ", fechaInicio)
						fmt.Println("ffin: ", fechaFin)
						fmt.Println("url ", "http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/rubro/GetRubroIngreso?rubro="+rubro+"&fuente="+idFuente+"&finicio="+fechaInicio.Format("2006-01-02")+"&ffin="+fechaFin.Format("2006-01-02"))

					}*/
					var valorEngresos interface{}
					if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/rubro/GetRubroOrdenPago?rubro="+rubro+"&fuente="+idFuente+"&finicio="+fechaInicio.Format("2006-01-02")+"&ffin="+fechaFin.Format("2006-01-02"), &valorEngresos); err == nil {
						var dataEngresos []map[string]interface{}
						err := formatdata.FillStruct(valorEngresos, &dataEngresos)
						if err != nil {

						} else {
							for _, valorData := range dataEngresos {
								//fmt.Println("rubroProyData(" + fmt.Sprintf("%v", ingresosRow.Idrubro) + "," + fmt.Sprintf("%v", fechaInicio.Year()) + "," + fmt.Sprintf("%v", int(fechaInicio.Month())) + "," + fmt.Sprintf("%v", valorData["valor"]) + ").")
								formatdata.FillStruct(valorData["valor"], &reporteRow.Valores.Valor)
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
	wg.Done()
	return
}

func (c *RubroController) calcularProyeccionIngresos(reporteData *cuerpoPac, finicio time.Time, ffin time.Time, nperiodos int, alert *models.Alert) {
	tool := new(ruler.EntornoReglas)
	tool.Agregar_dominio("Presupuesto")
	for _, ingresosRow := range reporteData.Ingresos { //recorrer los datos del reporte de ingresos para el rango actual

		for _, reporteRow := range ingresosRow.Reporte {
			var valor string
			var mes int
			err := formatdata.FillStruct(reporteRow.Valores.Valor, &valor)
			err = formatdata.FillStruct(reporteRow.N_mes, &mes)
			if err == nil {
				fechaInicio := time.Date(finicio.Year(), time.Month(mes), finicio.Day(), 0, 0, 0, 0, time.Local)
				fechaFin := time.Date(finicio.Year(), time.Month(mes+1), finicio.Day(), 0, 0, 0, 0, time.Local)

				if fechaFin.After(ffin) {
					fechaFin = ffin
				}
				var rubro string
				var idFuente string
				err := formatdata.FillStruct(ingresosRow.Idrubro, &rubro)

				err = formatdata.FillStruct(ingresosRow.Idfuente, &idFuente)

				if err == nil {

					/*if rubro == "35488" {
						fmt.Println("rubro: ", rubro)
						fmt.Println("Fuente: ", idFuente)
						fmt.Println("finicio: ", fechaInicio)
						fmt.Println("ffin: ", fechaFin)
						fmt.Println("url ", "http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/rubro/GetRubroIngreso?rubro="+rubro+"&fuente="+idFuente+"&finicio="+fechaInicio.Format("2006-01-02")+"&ffin="+fechaFin.Format("2006-01-02"))

					}*/
					var valorIngresos interface{}
					for i := 1; i <= nperiodos; i++ {
						Inicio := time.Date(fechaInicio.Year()-i, fechaInicio.Month(), fechaInicio.Day(), 0, 0, 0, 0, time.Local)
						Fin := time.Date(fechaFin.Year()-i, fechaFin.Month(), fechaFin.Day(), 0, 0, 0, 0, time.Local)
						//fmt.Println("url ", "http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/rubro/GetRubroIngreso?rubro="+rubro+"&fuente="+idFuente+"&finicio="+Inicio.Format("2006-01-02")+"&ffin="+Fin.Format("2006-01-02"))
						if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/rubro/GetRubroIngreso?rubro="+rubro+"&fuente="+idFuente+"&finicio="+Inicio.Format("2006-01-02")+"&ffin="+Fin.Format("2006-01-02"), &valorIngresos); err == nil {
							var dataIngresos []map[string]interface{}
							err := formatdata.FillStruct(valorIngresos, &dataIngresos)
							if err != nil {

							} else {
								for _, valorData := range dataIngresos {
									fmt.Println("rubro_proy_data(" + fmt.Sprintf("%v", ingresosRow.Idrubro) + "," + fmt.Sprintf("%v", i) + "," + fmt.Sprintf("%v", fechaInicio.Year()) + "," + fmt.Sprintf("%v", int(fechaInicio.Month())) + "," + fmt.Sprintf("%v", valorData["valor"]) + ").")
									tool.Agregar_predicado("rubro_proy_data(" + fmt.Sprintf("%v", ingresosRow.Idrubro) + "," + fmt.Sprintf("%v", fechaInicio.Year()) + "," + fmt.Sprintf("%v", int(fechaInicio.Month())) + "," + fmt.Sprintf("%v", valorData["valor"]) + ").")
									//formatdata.FillStruct(valorData["valor"], &reporteRow.Valores.Valor)
								}

							}

						} else {
							fmt.Println("err v", err.Error())
							alert = &models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
						}
					}
					tool.Agregar_predicado("rubro_proy_data(" + fmt.Sprintf("%v", ingresosRow.Idrubro) + ",1, 2014, 1, 345667).")
					tool.Agregar_predicado("rubro_proy_data(" + fmt.Sprintf("%v", ingresosRow.Idrubro) + ",2,2015,1,345668).")
					tool.Agregar_predicado("rubro_proy_data(" + fmt.Sprintf("%v", ingresosRow.Idrubro) + ",3,2016,1,345000).")
					reporteRow.Valores.Proyeccion = tool.Ejecutar_result("minimos_cuadrados_rubr("+fmt.Sprintf("%v", ingresosRow.Idrubro)+","+strconv.Itoa(nperiodos)+",R).", "R")

					var ej float64
					var proystr string

					err := formatdata.FillStruct(reporteRow.Valores.Proyeccion, &proystr)
					err1 := formatdata.FillStruct(reporteRow.Valores.Valor, &ej)
					if err == nil && err1 == nil {
						var variacion float64
						var pvar float64

						proy, _ := strconv.ParseFloat(proystr, 64)
						variacion = math.Abs(ej - proy)
						if ej <= 0 {

							pvar = (variacion / variacion)

						} else {

							pvar = (variacion / ej)
						}

						reporteRow.Valores.Pvariacion = pvar
						reporteRow.Valores.Variacion = variacion
					} else {
						fmt.Println("err ", err)
						fmt.Println("err1 ", err1)
					}
					tool.Quitar_predicados()
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
	wg.Done()
	return
}

func (c *RubroController) calcularProyeccionEgresos(reporteData *cuerpoPac, finicio time.Time, ffin time.Time, nperiodos int, alert *models.Alert) {
	tool := new(ruler.EntornoReglas)
	tool.Agregar_dominio("Presupuesto")
	for _, ingresosRow := range reporteData.Egresos { //recorrer los datos del reporte de ingresos para el rango actual

		for _, reporteRow := range ingresosRow.Reporte {
			var valor string
			var mes int
			err := formatdata.FillStruct(reporteRow.Valores.Valor, &valor)
			err = formatdata.FillStruct(reporteRow.N_mes, &mes)
			if err == nil {
				fechaInicio := time.Date(finicio.Year(), time.Month(mes), finicio.Day(), 0, 0, 0, 0, time.Local)
				fechaFin := time.Date(finicio.Year(), time.Month(mes+1), finicio.Day(), 0, 0, 0, 0, time.Local)

				if fechaFin.After(ffin) {
					fechaFin = ffin
				}
				var rubro string
				var idFuente string
				err := formatdata.FillStruct(ingresosRow.Idrubro, &rubro)

				err = formatdata.FillStruct(ingresosRow.Idfuente, &idFuente)

				if err == nil {

					/*if rubro == "35488" {
						fmt.Println("rubro: ", rubro)
						fmt.Println("Fuente: ", idFuente)
						fmt.Println("finicio: ", fechaInicio)
						fmt.Println("ffin: ", fechaFin)
						fmt.Println("url ", "http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/rubro/GetRubroIngreso?rubro="+rubro+"&fuente="+idFuente+"&finicio="+fechaInicio.Format("2006-01-02")+"&ffin="+fechaFin.Format("2006-01-02"))

					}*/
					var valorIngresos interface{}
					for i := 1; i <= nperiodos; i++ {
						Inicio := time.Date(fechaInicio.Year()-i, fechaInicio.Month(), fechaInicio.Day(), 0, 0, 0, 0, time.Local)
						Fin := time.Date(fechaFin.Year()-i, fechaFin.Month(), fechaFin.Day(), 0, 0, 0, 0, time.Local)
						//fmt.Println("url ", "http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/rubro/GetRubroIngreso?rubro="+rubro+"&fuente="+idFuente+"&finicio="+Inicio.Format("2006-01-02")+"&ffin="+Fin.Format("2006-01-02"))
						if err := getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/rubro/GetRubroOrdenPago?rubro="+rubro+"&fuente="+idFuente+"&finicio="+Inicio.Format("2006-01-02")+"&ffin="+Fin.Format("2006-01-02"), &valorIngresos); err == nil {
							var dataIngresos []map[string]interface{}
							err := formatdata.FillStruct(valorIngresos, &dataIngresos)
							if err != nil {

							} else {
								for _, valorData := range dataIngresos {
									fmt.Println("rubro_proy_data(" + fmt.Sprintf("%v", ingresosRow.Idrubro) + "," + fmt.Sprintf("%v", i) + "," + fmt.Sprintf("%v", fechaInicio.Year()) + "," + fmt.Sprintf("%v", int(fechaInicio.Month())) + "," + fmt.Sprintf("%v", valorData["valor"]) + ").")
									tool.Agregar_predicado("rubro_proy_data(" + fmt.Sprintf("%v", ingresosRow.Idrubro) + "," + fmt.Sprintf("%v", fechaInicio.Year()) + "," + fmt.Sprintf("%v", int(fechaInicio.Month())) + "," + fmt.Sprintf("%v", valorData["valor"]) + ").")
									//formatdata.FillStruct(valorData["valor"], &reporteRow.Valores.Valor)
								}

							}

						} else {
							fmt.Println("err v", err.Error())
							alert = &models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
						}
					}
					tool.Agregar_predicado("rubro_proy_data(" + fmt.Sprintf("%v", ingresosRow.Idrubro) + ",1, 2014, 1, 345667).")
					tool.Agregar_predicado("rubro_proy_data(" + fmt.Sprintf("%v", ingresosRow.Idrubro) + ",2,2015,1,345668).")
					tool.Agregar_predicado("rubro_proy_data(" + fmt.Sprintf("%v", ingresosRow.Idrubro) + ",3,2016,1,345000).")
					reporteRow.Valores.Proyeccion = tool.Ejecutar_result("minimos_cuadrados_rubr("+fmt.Sprintf("%v", ingresosRow.Idrubro)+","+strconv.Itoa(nperiodos)+",R).", "R")
					var ej float64
					var proystr string

					err := formatdata.FillStruct(reporteRow.Valores.Proyeccion, &proystr)
					err1 := formatdata.FillStruct(reporteRow.Valores.Valor, &ej)
					if err == nil && err1 == nil {
						var variacion float64
						var pvar float64

						proy, _ := strconv.ParseFloat(proystr, 64)
						variacion = math.Abs(ej - proy)
						if ej <= 0 {

							pvar = (variacion / variacion)

						} else {

							pvar = (variacion / ej)
						}

						reporteRow.Valores.Pvariacion = pvar
						reporteRow.Valores.Variacion = variacion
					} else {
						fmt.Println("err ", err)
						fmt.Println("err1 ", err1)
					}
					tool.Quitar_predicados()
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
	wg.Done()
	return
}

func cuerpoReporte(inicio time.Time, fin time.Time) (res cuerpoPac, err error) {

	mesinicio := int(inicio.Month())
	mesfin := int(fin.Month())
	var m []map[string]interface{}
	cuerpo := make(map[string]interface{})
	err = getJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/apropiacion/GetApropiacionesHijo/"+strconv.Itoa(inicio.Year())+"?tipo=2", &m)
	if err != nil {
		return
	}

	for i := 0; i < len(m); i++ {
		var fechas []map[string]interface{}
		for j := 0; j <= (mesfin - mesinicio); j++ {
			finicio := inicio.AddDate(0, j, 0)
			aux := make(map[string]interface{})

			val := make(map[string]interface{})
			val["valor"] = 0.0
			val["proyeccion"] = 0.0
			val["variacion"] = 0.0
			val["pvariacion"] = 0.0
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
	err = formatdata.FillStruct(m, &ingresos)
	if err != nil {
		fmt.Println("err2 ", err)
		return
	}
	cuerpo["ingresos"] = ingresos

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
			val["valor"] = 0.0
			val["proyeccion"] = 0.0
			val["variacion"] = 0.0
			val["pvariacion"] = 0.0
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

	err = formatdata.FillStruct(m, &ingresos)
	if err != nil {
		fmt.Println("err2 ", err)
		return
	}

	cuerpo["egresos"] = ingresos
	//fmt.Println(cuerpo["egresos"])
	err = mapstructure.Decode(cuerpo, &res)
	if err != nil {
		fmt.Println("err2 ", err)
		return
	}

	return
}
