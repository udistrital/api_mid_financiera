package controllers

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/astaxie/beego"
	"github.com/manucorporat/try"
	"github.com/mitchellh/mapstructure"
	"github.com/udistrital/api_mid_financiera/models"
	"github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/utils_oas/request"
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
type rowCierre struct {
	Proyeccion interface{}
	Pvariacion interface{}
	Valor      interface{}
	Variacion  interface{}
	IdAprop    interface{}
	Idrubro    interface{}
	CodigoRub  interface{}
	Descrubro  interface{}
	Idfuente   interface{}
	Fdescrip   interface{}
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
type cuerpoCierre struct {
	Ingresos []*rowCierre
	Egresos  []*rowCierre
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
				go c.calcularIngresos(&reporteData, finicio, &alert)
				go c.calcularEgresos(&reporteData, finicio, &alert)
				go c.agregarSumaFuenteEgresos(&reporteData, finicio, &alert)
				go c.agregarSumaFuenteIngresos(&reporteData, finicio, &alert)
				//go c.calcularProyeccionIngresos(&reporteData, finicio, ffin, periodos, &alert)
				//go c.calcularProyeccionEgresos(&reporteData, finicio, ffin, periodos, &alert)
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

func (c *RubroController) calcularIngresos(reporteData *cuerpoPac, finicio time.Time, alert *models.Alert) {

	for _, ingresosRow := range reporteData.Ingresos { //recorrer los datos del reporte de ingresos para el rango actual

		for _, reporteRow := range ingresosRow.Reporte {
			var valor string
			var mes int
			var proyectado float64
			var ejecutado float64
			err := formatdata.FillStruct(reporteRow.Valores.Valor, &valor)
			err = formatdata.FillStruct(reporteRow.N_mes, &mes)
			if err == nil {
				vigencia := finicio.Year()

				var rubro string
				var idFuente string
				err := formatdata.FillStruct(ingresosRow.Idrubro, &rubro)

				err = formatdata.FillStruct(ingresosRow.Idfuente, &idFuente)

				if err == nil {
					var valorIngresos interface{}
					if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/rubro/GetRubroPac?idRubro="+rubro+"&idFuente="+idFuente+"&vigencia="+strconv.Itoa(vigencia)+"&mes="+strconv.Itoa(mes), &valorIngresos); err == nil {
						var dataIngresos map[string]interface{}
						err := formatdata.FillStruct(valorIngresos, &dataIngresos)
						if err == nil {

							formatdata.FillStruct(dataIngresos["ejecutado"], &ejecutado)
							formatdata.FillStruct(dataIngresos["proyectado"], &proyectado)
							formatdata.FillStruct(ejecutado, &reporteRow.Valores.Valor)
							formatdata.FillStruct(proyectado, &reporteRow.Valores.Proyeccion)
							formatdata.FillStruct(math.Abs(ejecutado-proyectado), &reporteRow.Valores.Variacion)
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

func (c *RubroController) calcularEgresos(reporteData *cuerpoPac, finicio time.Time, alert *models.Alert) {

	for _, egresosRow := range reporteData.Egresos { //recorrer los datos del reporte de ingresos para el rango actual

		for _, reporteRow := range egresosRow.Reporte {
			var valor string
			var mes int
			var proyectado float64
			var ejecutado float64
			err := formatdata.FillStruct(reporteRow.Valores.Valor, &valor)
			err = formatdata.FillStruct(reporteRow.N_mes, &mes)
			if err == nil {
				vigencia := finicio.Year()

				var rubro string
				var idFuente string
				err := formatdata.FillStruct(egresosRow.Idrubro, &rubro)
				err = formatdata.FillStruct(egresosRow.Idfuente, &idFuente)

				if err == nil {
					var valorEgresos interface{}
					if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/rubro/GetRubroPac?idRubro="+rubro+"&idFuente="+idFuente+"&vigencia="+strconv.Itoa(vigencia)+"&mes="+strconv.Itoa(mes), &valorEgresos); err == nil {
						var dataEgresos map[string]interface{}
						err := formatdata.FillStruct(valorEgresos, &dataEgresos)
						if err == nil {
							formatdata.FillStruct(dataEgresos["ejecutado"], &ejecutado)
							formatdata.FillStruct(dataEgresos["proyectado"], &proyectado)
							formatdata.FillStruct(proyectado, &reporteRow.Valores.Proyeccion)
							formatdata.FillStruct(ejecutado, &reporteRow.Valores.Valor)
							formatdata.FillStruct(math.Abs(ejecutado-proyectado), &reporteRow.Valores.Variacion)
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
					var valorEngresos interface{}
					if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/rubro/GetRubroOrdenPago?rubro="+rubro+"&fuente="+idFuente+"&finicio="+fechaInicio.Format("2006-01-02")+"&ffin="+fechaFin.Format("2006-01-02"), &valorEngresos); err == nil {
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
						if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/rubro/GetRubroIngreso?rubro="+rubro+"&fuente="+idFuente+"&finicio="+Inicio.Format("2006-01-02")+"&ffin="+Fin.Format("2006-01-02"), &valorIngresos); err == nil {
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
						if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/rubro/GetRubroOrdenPago?rubro="+rubro+"&fuente="+idFuente+"&finicio="+Inicio.Format("2006-01-02")+"&ffin="+Fin.Format("2006-01-02"), &valorIngresos); err == nil {
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
func (c *RubroController) agregarSumaFuenteEgresos(reporteData *cuerpoPac, finicio time.Time, alert *models.Alert) (err error) {
	var idFuente string
	var idFuenteAnt string
	var descripcionF string
	var descripcionAnt string
	var codrubro string
	var i int

	vigencia := finicio.Year()

	lastRow := reporteData.Egresos[len(reporteData.Egresos)-1]
	for _, filaIngresos := range reporteData.Egresos {
		err = formatdata.FillStruct(filaIngresos.Idfuente, &idFuente)
		err = formatdata.FillStruct(filaIngresos.Fdescrip, &descripcionF)
		if val := strings.Compare(idFuente, idFuenteAnt); val != 0 && len(idFuenteAnt) > 0 {
			Reporte := getNewRow(filaIngresos.Reporte, idFuenteAnt, codrubro, vigencia)
			nuevaFila := &rowPac{Fdescrip: "Total Rubro" + descripcionAnt,
				Reporte: Reporte}

			reporteData.Egresos = append(reporteData.Egresos, nil)
			copy(reporteData.Egresos[i+1:], reporteData.Egresos[i:])
			reporteData.Egresos[i] = nuevaFila
			i++
		}
		idFuenteAnt = idFuente
		descripcionAnt = descripcionF
		err = formatdata.FillStruct(filaIngresos.Codigo, &codrubro)
		codrubro = codrubro[:1]
		i++
	}
	Reporte := getNewRow(lastRow.Reporte, idFuenteAnt, codrubro, vigencia)
	nuevaFila := &rowPac{Fdescrip: "Total Rubro" + descripcionAnt,
		Reporte: Reporte}
	reporteData.Egresos = append(reporteData.Egresos, nuevaFila)

	Reporte = getSumTotal(lastRow.Reporte, "2", vigencia)
	nuevaFila = &rowPac{Fdescrip: "Total Egresos ",
		Reporte: Reporte}
	reporteData.Egresos = append(reporteData.Egresos, nuevaFila)

	wg.Done()
	return
}
func (c *RubroController) agregarSumaFuenteIngresos(reporteData *cuerpoPac, finicio time.Time, alert *models.Alert) (err error) {
	//var valores valoresPac
	//var valorSuma int64
	var idFuente string
	var idFuenteAnt string
	var descripcionF string
	var descripcionAnt string
	var codrubro string
	var i int

	vigencia := finicio.Year()

	lastRow := reporteData.Ingresos[len(reporteData.Ingresos)-1]
	for _, filaIngresos := range reporteData.Ingresos {
		err = formatdata.FillStruct(filaIngresos.Idfuente, &idFuente)
		err = formatdata.FillStruct(filaIngresos.Fdescrip, &descripcionF)
		if val := strings.Compare(idFuente, idFuenteAnt); val != 0 && len(idFuenteAnt) > 0 {
			fmt.Println("valores  cambian de ", "fuenteAnt = "+idFuenteAnt, "fuente "+idFuente+" valor i "+strconv.Itoa(i))
			Reporte := getNewRow(filaIngresos.Reporte, idFuenteAnt, codrubro, vigencia)
			nuevaFila := &rowPac{Fdescrip: "Total Rubro" + descripcionAnt,
				Reporte: Reporte}

			reporteData.Ingresos = append(reporteData.Ingresos, nil)
			copy(reporteData.Ingresos[i+1:], reporteData.Ingresos[i:])
			reporteData.Ingresos[i] = nuevaFila
			i++
		}
		idFuenteAnt = idFuente
		descripcionAnt = descripcionF
		err = formatdata.FillStruct(filaIngresos.Codigo, &codrubro)
		codrubro = codrubro[:1]
		i++
	}
	Reporte := getNewRow(lastRow.Reporte, idFuenteAnt, codrubro, vigencia)
	nuevaFila := &rowPac{Fdescrip: "Total Rubro " + descripcionAnt,
		Reporte: Reporte}
	reporteData.Ingresos = append(reporteData.Ingresos, nuevaFila)

	Reporte = getSumTotal(lastRow.Reporte, "2", vigencia)
	nuevaFila = &rowPac{Fdescrip: "Total Ingresos ",
		Reporte: Reporte}
	reporteData.Ingresos = append(reporteData.Ingresos, nuevaFila)
	wg.Done()
	return
}

func getNewRow(row []*reportePacData, idFuente string, codrubro string, vigencia int) (Reporte []*reportePacData) {
	var Mes string
	var N_mes int
	var valorSumaF interface{}
	var mapValorSumaF map[string]interface{}
	var ejecutado float64
	var proyectado float64

	Reporte = make([]*reportePacData, 0)
	for _, valoresMes := range row {

		err := formatdata.FillStruct(valoresMes.Mes, &Mes)
		err = formatdata.FillStruct(valoresMes.N_mes, &N_mes)

		if err = request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/rubro/GetSumbySource?vigencia="+strconv.Itoa(vigencia)+"&mes="+strconv.Itoa(N_mes)+"&fuente="+idFuente+"&tipo="+codrubro, &valorSumaF); err == nil {
			err = formatdata.FillStruct(valorSumaF, &mapValorSumaF)
			err = formatdata.FillStruct(mapValorSumaF["ejecutado"], &ejecutado)
			err = formatdata.FillStruct(mapValorSumaF["proyectado"], &proyectado)
		}
		valorSuma := valoresPac{Proyeccion: proyectado,
			Valor:     ejecutado,
			Variacion: math.Abs(ejecutado - proyectado)}

		if err != nil {
			fmt.Println("Error", err.Error())
		}
		valoresN := &reportePacData{Mes: Mes,
			N_mes:   N_mes,
			Valores: &valorSuma}

		Reporte = append(Reporte, valoresN)
	}
	return
}

func getSumTotal(row []*reportePacData, tipo string, vigencia int) (Reporte []*reportePacData) {
	var Mes string
	var N_mes int
	var valorSumaF interface{}
	var mapValorSumaF map[string]interface{}
	var ejecutado float64
	var proyectado float64

	tipo = tipo + "%"

	Reporte = make([]*reportePacData, 0)
	for _, valoresMes := range row {

		err := formatdata.FillStruct(valoresMes.Mes, &Mes)
		err = formatdata.FillStruct(valoresMes.N_mes, &N_mes)

		if err = request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/rubro/GetSumbySource?vigencia="+strconv.Itoa(vigencia)+"&mes="+strconv.Itoa(N_mes)+"&tipo="+tipo, &valorSumaF); err == nil {
			err = formatdata.FillStruct(valorSumaF, &mapValorSumaF)
			err = formatdata.FillStruct(mapValorSumaF["ejecutado"], &ejecutado)
			err = formatdata.FillStruct(mapValorSumaF["proyectado"], &proyectado)
		}
		valorSuma := valoresPac{Proyeccion: proyectado,
			Valor:     ejecutado,
			Variacion: math.Abs(ejecutado - proyectado)}

		if err != nil {
			fmt.Println("Error", err.Error())
		}
		valoresN := &reportePacData{Mes: Mes,
			N_mes:   N_mes,
			Valores: &valorSuma}

		Reporte = append(Reporte, valoresN)
	}
	return
}

func cuerpoReporte(inicio time.Time, fin time.Time) (res cuerpoPac, err error) {

	mesinicio := int(inicio.Month())
	mesfin := int(fin.Month())
	var m []map[string]interface{}
	cuerpo := make(map[string]interface{})
	err = request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/apropiacion/GetApropiacionesHijo/"+strconv.Itoa(inicio.Year())+"?tipo=2", &m)
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

	err = request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/apropiacion/GetApropiacionesHijo/"+strconv.Itoa(inicio.Year())+"?tipo=3", &m)
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

	err = mapstructure.Decode(cuerpo, &res)
	if err != nil {
		fmt.Println("err2 ", err)
		return
	}

	return
}
func cierreIngresosEgresos(vigencia int, mes int, alert *models.Alert) (res cuerpoCierre, err error) {
	var cierreRow []map[string]interface{}
	var cierreRowEg []map[string]interface{}
	var ingresos interface{}
	var egresos interface{}
	mapCierre := make(map[string]interface{})
	err = request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/rubro/GetIngresoCierre?vigencia="+strconv.Itoa(vigencia)+"&codigo=2"+"&mes="+strconv.Itoa(mes), &cierreRow)
	if err != nil {
		fmt.Println("err ", err)
		alert = &models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
		return
	}

	err = formatdata.FillStruct(cierreRow, &ingresos)
	if err != nil {
		fmt.Println("err2 ", err)
		return
	}

	err = request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/rubro/GetIngresoCierre?vigencia="+strconv.Itoa(vigencia)+"&codigo=3"+"&mes="+strconv.Itoa(mes), &cierreRowEg)

	err = formatdata.FillStruct(cierreRowEg, &egresos)
	if err != nil {
		fmt.Println("err2 ", err)
		return
	}
	mapCierre["ingresos"] = ingresos
	mapCierre["egresos"] = egresos
	err = mapstructure.Decode(mapCierre, &res)

	if err != nil {
		fmt.Println("error decode ", err)
		alert = &models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
		return
	}
	return
}
func ProyeccionIngresosCierre(reporte *cuerpoCierre, mes int, vigencia int, nperiodos int, alert *models.Alert) {

	var rubro string
	var fuente string
	var cantPredicados int64
	var valProyectado float64
	var valEjecutado float64
	var variacion float64
	var pVariacion float64

	tool := new(ruler.EntornoReglas)
	tool.Agregar_dominio("Presupuesto")

	for _, ingresosRow := range reporte.Ingresos {
		err := formatdata.FillStruct(ingresosRow.Idrubro, &rubro)
		err = formatdata.FillStruct(ingresosRow.Idfuente, &fuente)
		if err == nil {
			var valorIngresos interface{}
			for i := 1; i <= nperiodos; i++ {
				if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/rubro/GetPacValue?vigencia="+strconv.Itoa(vigencia-i)+"&mes="+strconv.Itoa(mes)+"&rubro="+rubro+"&fuente="+fuente, &valorIngresos); err == nil {
					var dataIngresos []map[string]interface{}
					err := formatdata.FillStruct(valorIngresos, &dataIngresos)

					if err == nil {
						for _, valorData := range dataIngresos {
							fmt.Println("rubro_proy_data(" + fmt.Sprintf("%v", ingresosRow.Idrubro) + "," + fmt.Sprintf("%v", vigencia-i) + "," + fmt.Sprintf("%v", mes) + "," + fmt.Sprintf("%v", valorData["valor"]) + ").")
							tool.Agregar_predicado("rubro_proy_data(" + fmt.Sprintf("%v", ingresosRow.Idrubro) + "," + fmt.Sprintf("%v", vigencia-i) + "," + fmt.Sprintf("%v", mes) + "," + fmt.Sprintf("%v", valorData["valor"]) + ").")
							cantPredicados++
						}
					} else {
						fmt.Println("err v", err.Error())
						alert = &models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
					}
				}
			}
		} else {
			fmt.Println("err ", err.Error())
			alert = &models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
		}

		err = formatdata.FillStruct(ingresosRow.Valor, &valEjecutado)
		if cantPredicados > 0 {
			ingresosRow.Proyeccion = tool.Ejecutar_result("minimos_cuadrados_rubr("+fmt.Sprintf("%v", ingresosRow.Idrubro)+","+strconv.Itoa(nperiodos)+",R).", "R")
		}
		err = formatdata.FillStruct(ingresosRow.Proyeccion, &valProyectado)

		ingresosRow.Variacion = math.Abs(valEjecutado - valProyectado)

		err = formatdata.FillStruct(ingresosRow.Variacion, &variacion)

		if valEjecutado != 0 {
			pVariacion = (variacion / valEjecutado)
		}

		ingresosRow.Pvariacion = pVariacion
		tool.Quitar_predicados()
		cantPredicados = 0
	}
	wg.Done()
	return
}
func ProyeccionEgresosCierre(reporte *cuerpoCierre, mes int, vigencia int, nperiodos int, alert *models.Alert) {

	var rubro string
	var fuente string
	var valProyectado float64
	var valEjecutado float64
	var variacion float64
	var pVariacion float64
	var cantPredicados int64

	tool := new(ruler.EntornoReglas)
	tool.Agregar_dominio("Presupuesto")

	for _, egresosRow := range reporte.Egresos {
		err := formatdata.FillStruct(egresosRow.Idrubro, &rubro)
		err = formatdata.FillStruct(egresosRow.Idfuente, &fuente)
		if err == nil {
			var valorEgresos interface{}
			for i := 1; i <= nperiodos; i++ {
				if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/rubro/GetPacValue?vigencia="+strconv.Itoa(vigencia-i)+"&mes="+strconv.Itoa(mes)+"&rubro="+rubro+"&fuente="+fuente, &valorEgresos); err == nil {
					var dataEgresos []map[string]interface{}
					err := formatdata.FillStruct(valorEgresos, &dataEgresos)
					if err != nil {
						fmt.Println("error v", err.Error())
					} else {

						for _, valorData := range dataEgresos {
							fmt.Println("rubro_proy_data(" + fmt.Sprintf("%v", egresosRow.Idrubro) + "," + fmt.Sprintf("%v", vigencia-i) + "," + fmt.Sprintf("%v", mes) + "," + fmt.Sprintf("%v", valorData["valor"]) + ").")
							tool.Agregar_predicado("rubro_proy_data(" + fmt.Sprintf("%v", egresosRow.Idrubro) + "," + fmt.Sprintf("%v", vigencia-i) + "," + fmt.Sprintf("%v", mes) + "," + fmt.Sprintf("%v", valorData["valor"]) + ").")
							cantPredicados++
						}
					}
				} else {
					fmt.Println("err ", err.Error())
					alert = &models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
				}
			}
		} else {
			fmt.Println("err ", err.Error())
			alert = &models.Alert{Code: "E_0458", Body: err.Error(), Type: "error"}
		}

		err = formatdata.FillStruct(egresosRow.Valor, &valEjecutado)
		if cantPredicados > 0 {
			beego.Error("proyeccion egresos")
			fmt.Println("minimos_cuadrados_rubr("+fmt.Sprintf("%v", egresosRow.Idrubro)+","+strconv.Itoa(nperiodos)+",R).", "R")
			egresosRow.Proyeccion = tool.Ejecutar_result("minimos_cuadrados_rubr("+fmt.Sprintf("%v", egresosRow.Idrubro)+","+strconv.Itoa(nperiodos)+",R).", "R")
		}
		err = formatdata.FillStruct(egresosRow.Proyeccion, &valProyectado)

		egresosRow.Variacion = math.Abs(valEjecutado - valProyectado)

		err = formatdata.FillStruct(egresosRow.Variacion, &variacion)

		if valEjecutado != 0 {
			pVariacion = (variacion / valEjecutado)
		}

		egresosRow.Pvariacion = pVariacion
		tool.Quitar_predicados()
		cantPredicados = 0
	}
	wg.Done()
	return
}

// GenerarCierre ...
// @Title GenerarCierre
// @Description Get all information to close PAC
// @Param    request        query     interface{}    true        "objeto con parametros cierre"
// @Success 200 {object} interface{}
// @Failure 403
// @router /GenerarCierre/ [post]
func (c *RubroController) GenerarCierre() {
	defer c.ServeJSON()
	wg.Add(2)
	fmt.Println("GenerarCierre")
	var request map[string]interface{} //definicion de la interface que recibe los datos del reporte y proyecciones
	var vigencia int
	var mes string
	var nperiodos int
	//var periodospr string
	var m int
	var alert models.Alert
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &request); err == nil {
		err = formatdata.FillStruct(request["vigencia"], &vigencia)
		err = formatdata.FillStruct(request["mes"], &mes)
		err = formatdata.FillStruct(request["periodosproy"], &nperiodos)

		if err == nil {
			m, err = strconv.Atoi(mes)
			//nperiodos, err = strconv.Atoi(periodospr)
			if err != nil {
				fmt.Println("error val", err.Error())
			}
			if cuerpoCierre, err := cierreIngresosEgresos(vigencia, m, &alert); err == nil {
				if alert.Body == nil {
					fmt.Println("no alert")
				} else {
					fmt.Println("alert ", alert)
				}
				go ProyeccionEgresosCierre(&cuerpoCierre, m, vigencia, nperiodos, &alert)
				go ProyeccionIngresosCierre(&cuerpoCierre, m, vigencia, nperiodos, &alert)
				wg.Wait()
				c.Data["json"] = cuerpoCierre
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

// RegistrarRubro ...
// @Title RegistrarRubro
// @Description Registra Rubro en postgres y mongo
// @Param       body            body    models.Rubro    true            "body for Rubro content"
// @Success 200 {object} models.Alert
// @Failure 403 body is empty
// @router /RegistrarRubro/ [post]
func (c *RubroController) RegistrarRubro() {
	var v interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		try.This(func() {
			res := make(map[string]interface{})
			rubroData := v.(map[string]interface{})
			if rubroData["RubroPadre"] != nil { //Si se registra Un rubro con padre
				url := "http://" + beego.AppConfig.String("Urlcrud") + ":" + beego.AppConfig.String("Portcrud") + "/" + beego.AppConfig.String("Nscrud") + "/rubro_rubro"
				if err := request.SendJson(url, "POST", &res, &rubroData); err == nil {
					//Cuando se registra el rubro, se debe mandar una petición a MongoApi para registrar el nuevo rubro.
					//En este caso se genera un map con la estructura que recibe dicho api.

					c.Data["json"] = res
				} else {
					panic("Service Error")
				}
			} else if rubroData["RubroHijo"] != nil { //Si se registra un rubro Padre
				rubro := rubroData["RubroHijo"]
				url := "http://" + beego.AppConfig.String("Urlcrud") + ":" + beego.AppConfig.String("Portcrud") + "/" + beego.AppConfig.String("Nscrud") + "/rubro"
				if err := request.SendJson(url, "POST", &res, &rubro); err == nil {
					//Cuando se registra el rubro, se debe mandar una petición a MongoApi para registrar el nuevo rubro.
					//En este caso se genera un map con la estructura que recibe dicho api.

					res["Body"] = map[string]interface{}{"RubroHijo": rubro, "RubroPadre": nil}
					c.Data["json"] = res
				} else {
					panic("Service Error")
				}
			} else {
				panic("Data Undefined")
			}

		}).Catch(func(e try.E) {
			// Print crash
			fmt.Println("expc ", e)
			c.Data["json"] = e
		})
	} else {
		c.Data["json"] = nil
	}
	c.ServeJSON()
}
