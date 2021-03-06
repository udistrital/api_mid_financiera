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

// RubroController operations for Rubro

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

type rowTotales struct {
	Descripcion interface{}
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
	Ingresos        []*rowPac
	Egresos         []*rowPac
	TotalesIngresos []*rowTotales
	TotalesEgresos  []*rowTotales
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

							if errorEjecutado := formatdata.FillStruct(dataIngresos["ejecutado"], &ejecutado); errorEjecutado != nil {
								beego.Info("error llenando estructura ejecutado")
							}

							if errorProyectado := formatdata.FillStruct(dataIngresos["proyectado"], &proyectado); errorProyectado != nil {
								beego.Info("error llenando estructura proyectado")
							}

							if errorReporteValor := formatdata.FillStruct(ejecutado, &reporteRow.Valores.Valor); errorReporteValor != nil {
								beego.Info("error llenando estructura reporte valor")
							}

							if errorReporteProyeccion := formatdata.FillStruct(proyectado, &reporteRow.Valores.Proyeccion); errorReporteProyeccion != nil {
								beego.Info("error llenando estructura reporte proyeccion")
							}

							if errorVariacion := formatdata.FillStruct(math.Abs(ejecutado-proyectado), &reporteRow.Valores.Variacion); errorVariacion != nil {
								beego.Info("error llenando estructura variacion")
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
							if errorEjecutado := formatdata.FillStruct(dataEgresos["ejecutado"], &ejecutado); errorEjecutado != nil {
								beego.Error("error llenando estructura ejecutado")
							}
							if errorProyectado := formatdata.FillStruct(dataEgresos["proyectado"], &proyectado); errorProyectado != nil {
								beego.Error("error llenando estructura proyectado")
							}
							if errorReporteValor := formatdata.FillStruct(ejecutado, &reporteRow.Valores.Valor); errorReporteValor != nil {
								beego.Error("error llenando estructura reporte valor")
							}
							if errorReporteProyeccion := formatdata.FillStruct(proyectado, &reporteRow.Valores.Proyeccion); errorReporteProyeccion != nil {
								beego.Error("error llenando estructura reporte proyeccion")
							}
							if errorVariacion := formatdata.FillStruct(math.Abs(ejecutado-proyectado), &reporteRow.Valores.Variacion); errorVariacion != nil {
								beego.Error("error llenando estructura variacion")
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

func (c *RubroController) agregarSumaFuenteEgresos(reporteData *cuerpoPac, finicio time.Time, alert *models.Alert) (err error) {
	var idFuente string
	var idFuenteAnt string
	var descripcionF string
	var descripcionAnt string
	var codrubro string
	var i int

	vigencia := finicio.Year()
	if len(reporteData.Egresos) > 0 {
		lastRow := reporteData.Egresos[len(reporteData.Egresos)-1]
		for _, filaIngresos := range reporteData.Egresos {
			err = formatdata.FillStruct(filaIngresos.Idfuente, &idFuente)
			err = formatdata.FillStruct(filaIngresos.Fdescrip, &descripcionF)
			if val := strings.Compare(idFuente, idFuenteAnt); val != 0 && len(idFuenteAnt) > 0 {
				Reporte := getNewRow(filaIngresos.Reporte, idFuenteAnt, codrubro, vigencia)
				nuevaFila := &rowPac{Fdescrip: "Total Rubro " + descripcionAnt,
					Reporte: Reporte}
				nuevoTotal := &rowTotales{Descripcion: "Total Rubro " + descripcionAnt,
					Reporte: Reporte}

				reporteData.TotalesEgresos = append(reporteData.TotalesEgresos, nuevoTotal)

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
		nuevaFila := &rowPac{Fdescrip: "Total Rubro " + descripcionAnt,
			Reporte: Reporte}

		nuevoTotal := &rowTotales{Descripcion: "Total Rubro" + descripcionAnt,
			Reporte: Reporte}

		reporteData.TotalesEgresos = append(reporteData.TotalesEgresos, nuevoTotal)

		reporteData.Egresos = append(reporteData.Egresos, nuevaFila)

		Reporte = getSumTotal(lastRow.Reporte, "3", vigencia)
		nuevaFila = &rowPac{Fdescrip: "Total Egresos ",
			Reporte: Reporte}

		nuevoTotal = &rowTotales{Descripcion: "Total Egresos ",
			Reporte: Reporte}

		reporteData.TotalesEgresos = append(reporteData.TotalesEgresos, nuevoTotal)

		reporteData.Egresos = append(reporteData.Egresos, nuevaFila)
	}
	wg.Done()
	return
}
func (c *RubroController) agregarSumaFuenteIngresos(reporteData *cuerpoPac, finicio time.Time, alert *models.Alert) (err error) {
	var idFuente string
	var idFuenteAnt string
	var descripcionF string
	var descripcionAnt string
	var codrubro string
	var i int

	vigencia := finicio.Year()
	if len(reporteData.Ingresos) > 0 {
		lastRow := reporteData.Ingresos[len(reporteData.Ingresos)-1]
		for _, filaIngresos := range reporteData.Ingresos {
			err = formatdata.FillStruct(filaIngresos.Idfuente, &idFuente)
			err = formatdata.FillStruct(filaIngresos.Fdescrip, &descripcionF)
			if val := strings.Compare(idFuente, idFuenteAnt); val != 0 && len(idFuenteAnt) > 0 {
				Reporte := getNewRow(filaIngresos.Reporte, idFuenteAnt, codrubro, vigencia)
				nuevaFila := &rowPac{Fdescrip: "Total Rubro " + descripcionAnt,
					Reporte: Reporte}
				nuevoTotal := &rowTotales{Descripcion: "Total Rubro " + descripcionAnt,
					Reporte: Reporte}

				reporteData.TotalesIngresos = append(reporteData.TotalesIngresos, nuevoTotal)

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

		nuevoTotal := &rowTotales{Descripcion: "Total Rubro " + descripcionAnt,
			Reporte: Reporte}

		reporteData.TotalesIngresos = append(reporteData.TotalesIngresos, nuevoTotal)

		reporteData.Ingresos = append(reporteData.Ingresos, nuevaFila)

		Reporte = getSumTotal(lastRow.Reporte, "2", vigencia)
		nuevaFila = &rowPac{Fdescrip: "Total Ingresos ",
			Reporte: Reporte}
		nuevoTotal = &rowTotales{Descripcion: "Total Ingresos ",
			Reporte: Reporte}

		reporteData.TotalesIngresos = append(reporteData.TotalesIngresos, nuevoTotal)

		reporteData.Ingresos = append(reporteData.Ingresos, nuevaFila)
	}
	wg.Done()
	return
}

func getNewRow(row []*reportePacData, idFuente string, codrubro string, vigencia int) (Reporte []*reportePacData) {
	var Mes string
	var NMes int
	var valorSumaF interface{}
	var mapValorSumaF map[string]interface{}
	var ejecutado float64
	var proyectado float64

	Reporte = make([]*reportePacData, 0)
	for _, valoresMes := range row {

		err := formatdata.FillStruct(valoresMes.Mes, &Mes)
		err = formatdata.FillStruct(valoresMes.N_mes, &NMes)

		if err = request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/rubro/GetSumbySource?vigencia="+strconv.Itoa(vigencia)+"&mes="+strconv.Itoa(NMes)+"&fuente="+idFuente+"&tipo="+codrubro, &valorSumaF); err == nil {
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
			N_mes:   NMes,
			Valores: &valorSuma}

		Reporte = append(Reporte, valoresN)
	}
	return
}

func getSumTotal(row []*reportePacData, tipo string, vigencia int) (Reporte []*reportePacData) {
	var Mes string
	var NMes int
	var valorSumaF interface{}
	var mapValorSumaF map[string]interface{}
	var ejecutado float64
	var proyectado float64

	tipo = tipo + "%"

	Reporte = make([]*reportePacData, 0)
	for _, valoresMes := range row {

		err := formatdata.FillStruct(valoresMes.Mes, &Mes)
		err = formatdata.FillStruct(valoresMes.N_mes, &NMes)

		if err = request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/rubro/GetSumbySource?vigencia="+strconv.Itoa(vigencia)+"&mes="+strconv.Itoa(NMes)+"&tipo="+tipo, &valorSumaF); err == nil {
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
			N_mes:   NMes,
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
		if err != nil {
			fmt.Println("err1 ", err)
			return
		}

	}
	var ingresos interface{}
	err = formatdata.FillStruct(m, &ingresos)
	if err != nil {
		beego.Error("err 2 ", err.Error())
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
		fmt.Println("err2 ", err.Error())
		return
	}

	err = request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/rubro/GetIngresoCierre?vigencia="+strconv.Itoa(vigencia)+"&codigo=3"+"&mes="+strconv.Itoa(mes), &cierreRowEg)

	err = formatdata.FillStruct(cierreRowEg, &egresos)
	if err != nil {
		fmt.Println("err2 ", err.Error())
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
func proyeccionIngresosCierre(reporte *cuerpoCierre, mes int, vigencia int, nperiodos int, alert *models.Alert) {

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
func proyeccionEgresosCierre(reporte *cuerpoCierre, mes int, vigencia int, nperiodos int, alert *models.Alert) {

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
				if alert.Body != nil {
					beego.Error("alert ", alert)
				}
				go proyeccionEgresosCierre(&cuerpoCierre, m, vigencia, nperiodos, &alert)
				go proyeccionIngresosCierre(&cuerpoCierre, m, vigencia, nperiodos, &alert)
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
				urlcrud := "http://" + beego.AppConfig.String("Urlcrud") + ":" + beego.AppConfig.String("Portcrud") + "/" + beego.AppConfig.String("Nscrud") + "/rubro_rubro"
				if err := request.SendJson(urlcrud, "POST", &res, &rubroData); err == nil {
					//Cuando se registra el rubro, se debe mandar una petición a MongoApi para registrar el nuevo rubro.
					//En este caso se genera un map con la estructura que recibe dicho api.
					//Se debe comprobar si se pudo registrar el rubro y la relacion rubro_rubro en postgres.
					if res["Type"] != nil && res["Type"].(string) == "success" {
						urlmongo := "http://" + beego.AppConfig.String("financieraMongoCurdApiService") + "/arbol_rubro/registrarRubro"
						var data map[string]interface{}
						sendData := res["Body"].(map[string]interface{})
						err := request.SendJson(urlmongo, "POST", &data, &sendData)
						beego.Info("data: ", sendData)
						if data["Type"] != nil && err == nil {
							if data["Type"].(string) == "error" {
								resul := res["Body"].(map[string]interface{})
								ue := resul["RubroHijo"].(map[string]interface{})["UnidadEjecutora"].(float64)
								urlcrud = urlcrud + "/DeleteRubroRelation/" + strconv.Itoa(int(resul["Id"].(float64))) + "/" + strconv.Itoa(int(ue))
								if errorDelete := request.SendJson(urlcrud, "DELETE", &data, nil); errorDelete == nil {
									beego.Info("Data ", data)
									panic("Mongo API Error")
								} else {
									beego.Info("Error delete ", errorDelete)
									panic("Delete API Error")
								}

							} else if data["Type"].(string) == "success" {
								beego.Info("Registrado en Mongo")
							}
						} else {
							resul := res["Body"].(map[string]interface{})
							ue := resul["RubroHijo"].(map[string]interface{})["UnidadEjecutora"].(float64)
							urlcrud = urlcrud + "/DeleteRubroRelation/" + strconv.Itoa(int(resul["Id"].(float64))) + "/" + strconv.Itoa(int(ue))
							if errorDelete := request.SendJson(urlcrud, "DELETE", &data, nil); errorDelete == nil {
								beego.Info("Data ", data)
								panic("Mongo API not Found")
							} else {
								beego.Info("Error delete ", errorDelete)
								panic("Delete API Error")
							}

						}
					} else if res["Type"] == nil {
						panic("Financiera Crud Service Error")
					}

					c.Data["json"] = res
				} else {
					panic("Financiera Crud Service Error")
				}
			} else if rubroData["RubroHijo"] != nil { //Si se registra un rubro Padre
				rubro := rubroData["RubroHijo"]
				urlcrud := "http://" + beego.AppConfig.String("Urlcrud") + ":" + beego.AppConfig.String("Portcrud") + "/" + beego.AppConfig.String("Nscrud") + "/rubro"
				if err := request.SendJson(urlcrud, "POST", &res, &rubro); err == nil {
					//Cuando se registra el rubro, se debe mandar una petición a MongoApi para registrar el nuevo rubro.
					//En este caso se genera un map con la estructura que recibe dicho api.
					//Se debe comprobar si se pudo registrar el rubro en postgres.
					if res["Type"] != nil && res["Type"].(string) == "success" {
						urlmongo := "http://" + beego.AppConfig.String("financieraMongoCurdApiService") + "arbol_rubro/registrarRubro"
						var data map[string]interface{}
						res["Body"] = map[string]interface{}{"RubroHijo": res["Body"].(map[string]interface{}), "RubroPadre": map[string]interface{}{}}
						rubroRow := res["Body"].(map[string]interface{})
						errorPost := request.SendJson(urlmongo, "POST", &data, &rubroRow)
						beego.Info("data: ", urlmongo)
						if data["Type"] != nil && errorPost == nil {
							if data["Type"].(string) == "error" {
								beego.Info("Error en mongo")
								resul := res["Body"].(map[string]interface{})["RubroHijo"].(map[string]interface{})
								beego.Info("Send Data: ", resul)
								urlcrud = urlcrud + "/" + strconv.Itoa(int(resul["Id"].(float64)))
								if errorDelete := request.SendJson(urlcrud, "DELETE", &data, nil); errorDelete == nil {
									beego.Info("Data ", data)
									panic("Mongo API Error")
								} else {
									beego.Info("Error ", errorDelete)
									panic("delete API Error")
								}

							} else if data["Type"].(string) == "success" {
								beego.Info("Registrado en Mongo")
							}
						} else {
							resul := res["Body"].(map[string]interface{})["RubroHijo"].(map[string]interface{})
							urlcrud = urlcrud + "/" + strconv.Itoa(int(resul["Id"].(float64)))
							if errorDelete := request.SendJson(urlcrud, "DELETE", &data, nil); errorDelete == nil {
								beego.Info("Data ", data)
								panic("Mongo API not Found")
							} else {
								beego.Info("Error ", errorDelete)
								panic("delete API Error")
							}

						}

					} else if res["Type"] == nil {
						panic("Financiera Crud Service Error")
					}

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
			c.Data["json"] = map[string]interface{}{"Code": "E_0458", "Body": e, "Type": "error"}
		})
	} else {
		c.Data["json"] = nil
	}
	c.ServeJSON()
}

// EliminarRubro ...
// @Title EliminarRubro
// @Description delete the Rubro
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /EliminarRubro/:id [delete]
func (c *RubroController) EliminarRubro() {
	try.This(func() {
		idStr := c.Ctx.Input.Param(":id")
		urlcrud := "http://" + beego.AppConfig.String("Urlcrud") + ":" + beego.AppConfig.String("Portcrud") + "/" + beego.AppConfig.String("Nscrud") + "/rubro/" + idStr
		var res map[string]interface{}
		if err := request.SendJson(urlcrud, "DELETE", &res, nil); err == nil {
			if res["Type"].(string) == "success" {
				var resMg map[string]interface{}
				urlmongo := "http://" + beego.AppConfig.String("financieraMongoCurdApiService") + "arbol_rubro/eliminarRubro/" + idStr
				if err = request.SendJson(urlmongo, "DELETE", &resMg, nil); err != nil {
					fmt.Println("err ", err)
					panic("Mongo Not Found")
				} else if resMg["Type"].(string) == "error" {
					panic("Mongo CRUD Service Error")
				}
			} else if res["Type"].(string) == "error" {
				c.Data["json"] = res
			} else {
				panic("Financiera CRUD Service Error")
			}
		}
		c.Data["json"] = res
	}).Catch(func(e try.E) {
		fmt.Println("expc ", e)
		c.Data["json"] = map[string]interface{}{"Code": "E_0458", "Body": e, "Type": "error"}
	})
	c.ServeJSON()
}

// ArbolRubros ...
// @Title ArbolRubros
// @Description Get Arbol Rubros By UE
// @Param	unidadEjecutora		path 	int64	true		"unidad ejecutora a consultar"
// @Param	rama		query 	string	false		"rama a consultar"
// @Success 200 {object} models.Rubro
// @Failure 403
// @router /ArbolRubros/:unidadEjecutora [get]
func (c *RubroController) ArbolRubros() {

	try.This(func() {
		ueStr := c.Ctx.Input.Param(":unidadEjecutora")
		rama := c.GetString("rama")
		urlmongo := ""
		var res []map[string]interface{}
		if rama == "" {
			urlmongo = "http://" + beego.AppConfig.String("financieraMongoCurdApiService") + "arbol_rubro/RaicesArbol/" + ueStr
		} else {
			urlmongo = "http://" + beego.AppConfig.String("financieraMongoCurdApiService") + "arbol_rubro/ArbolRubro/" + rama + "/" + ueStr
		}
		beego.Info("Url ", urlmongo)
		if err := request.GetJson(urlmongo, &res); err != nil {
			beego.Info(err.Error())
			panic("Mongo API Service Error")
		}
		c.Data["json"] = res
	}).Catch(func(e try.E) {
		fmt.Println("expc ", e)
		c.Data["json"] = map[string]interface{}{"Code": "E_0458", "Body": e, "Type": "error"}
	})
	c.ServeJSON()
}
