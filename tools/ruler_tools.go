package tools

import (
	"errors"
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	//."github.com/mndrix/golog"
	. "github.com/mndrix/golog"
	"github.com/udistrital/api_mid_financiera/models"
	"github.com/udistrital/api_mid_financiera/utilidades"
)

type EntornoReglas struct {
	predicados string
	base       string
}

func (e *EntornoReglas) Agregar_dominio(dominio string) {
	var v []models.Predicado
	if err := getJson("http://"+beego.AppConfig.String("Urlruler")+":"+beego.AppConfig.String("Portruler")+"/"+beego.AppConfig.String("Nsruler")+"/predicado?limit=0&query=Dominio.Nombre:"+dominio, &v); err == nil {
		for i := 0; i < len(v); i++ {
			e.base = e.base + v[i].Nombre + "\n"
		}
	}
}

func (e *EntornoReglas) Agregar_predicado_dinamico(predicados ...string) (err error) {
	result := ""
	//recorrer los predicados que se quieren insertar
	for _, predicadod := range predicados { //se recorren el o los predicados dinamicos
		var rulename string
		for i, rp := range strings.SplitN(predicadod, ":", 3) {
			if len(rp) <= 1 {
				err = errors.New("Error1: invalid query key/value pair")
				return
			}
			if i == 0 {
				rulename = rulename + rp + "("
			} else {
				vr := strings.Split(rp, "|")
				if len(vr) <= 1 {
					err = errors.New("Error2: invalid query key/value pair")
					return
				}
				vs := strings.Split(vr[0], ".")
				if len(vs) < 2 || len(vs) > 4 || len(vs) == 3 {
					err = errors.New("Error2: invalid query key/value pair")
					return
				}
				service := vs[0]
				route := vs[1]
				sort := ""
				fmt.Println(vr)
				if len(vs) == 4 {
					sort = "&query=" + vs[2] + ":" + vs[3]
				}
				fmt.Println("http://" + beego.AppConfig.String(service) + route + "?limit=-1" + sort)
				var serviceresult []map[string]interface{}
				if err = getJson("http://"+beego.AppConfig.String(service)+route+"?limit=-1"+sort, &serviceresult); err == nil {
					//result[]
					fmt.Println("res ", vr)
					for _, res := range serviceresult {
						for j := 1; j < len(vr); j++ {
							if j == 1 {
								if values := strings.Split(vr[j], "."); len(values) > 1 {
									var finalvalue interface{}
									for index, mp := range values {
										if index != 0 {
											var aux map[string]interface{}
											err = utilidades.FillStruct(finalvalue, &aux)
											fmt.Println("finalvalue ", finalvalue)
											if err != nil {
												return
											}
											err = utilidades.FillStruct(aux[mp], &finalvalue)
											fmt.Println("finalvalue ", finalvalue)
											if err != nil {
												return
											}
										} else {
											err = utilidades.FillStruct(res[mp], &finalvalue)
											fmt.Println("finalvalue1 ", finalvalue)
											if err != nil {
												return
											}
										}
									}
									value := fmt.Sprintf("%v", finalvalue) //convertir cualquier interface en string **
									result = result + rulename + value
								} else {
									value := fmt.Sprintf("%v", res[vr[j]]) //convertir cualquier interface en string **
									result = result + rulename + value
								}
							} else {
								value := fmt.Sprintf("%v", res[vr[j]])
								result = result + "," + value
							}
						}
						result = result + ")."
						e.Agregar_predicado(result)
						result = ""
					}

				} else {
					return
				}

			}

		}

	}
	//fmt.Println(result)
	return
}

func (e *EntornoReglas) Agregar_predicado(predicado string) {
	e.predicados = e.predicados + predicado + "\n"
}

func (e *EntornoReglas) Obtener_predicados() (predicados string) {
	return e.predicados + e.base
}

func (e *EntornoReglas) Quitar_predicados() {
	e.predicados = ``
}

func (e *EntornoReglas) Ejecutar_result(regla string, variable string) (res interface{}) {
	f := NewMachine().Consult(e.predicados + e.base)
	solutions := f.ProveAll(regla)
	//fmt.Println(solutions)
	for _, solution := range solutions {
		res = fmt.Sprintf("%v", solution.ByName_(variable))
		//fmt.Printf("%s", solution.ByName_("R"))
	}
	return
}

func (e *EntornoReglas) Ejecutar_all_result(regla string, variable string) (res []interface{}) {
	f := NewMachine().Consult(e.predicados + e.base)
	solutions := f.ProveAll(regla)
	//fmt.Println(solutions)
	for _, solution := range solutions {
		res = append(res, fmt.Sprintf("%v", solution.ByName_(variable)))
		//fmt.Printf("%s", solution.ByName_("R"))
	}
	return
}
