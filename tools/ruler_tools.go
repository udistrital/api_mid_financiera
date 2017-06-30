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
}

func (e *EntornoReglas) Agregar_dominio(dominio string) {
	var v []models.Predicado
	if err := getJson("http://"+beego.AppConfig.String("Urlruler")+":"+beego.AppConfig.String("Portruler")+"/"+beego.AppConfig.String("Nsruler")+"/predicado?limit=0&query=Dominio.Nombre:"+dominio, &v); err == nil {
		for i := 0; i < len(v); i++ {
			e.predicados = e.predicados + v[i].Nombre + "\n"
		}
	}
}

func (e *EntornoReglas) Agregar_predicado_dinamico(predicados ...string) (result []map[string]interface{}, err error) {
	var regla string
	regla = ""
	//recorrer los predicados que se quieren insertar
	for _, predicadod := range predicados { //se recorren el o los predicados dinamicos
		for i, rp := range strings.SplitN(predicadod, ":", 3) {
			if len(rp) <= 1 {
				err = errors.New("Error1: invalid query key/value pair")
				return
			}
			if i == 0 {
				regla = regla + rp + "("
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
				if len(vs) == 4 {
					sort = "&query=" + vs[2] + ":" + vs[3]
				}
				fmt.Println("http://" + beego.AppConfig.String(service) + route + "?limit=-1" + sort)
				var serviceresult []interface{}
				if err = getJson("http://"+beego.AppConfig.String(service)+route+"?limit=-1"+sort, &serviceresult); err == nil {
					err = utilidades.FillStruct(serviceresult, &result)
					//fmt.Println("res ", result)
				} else {
					return
				}
				for j := 1; j < len(vr); j++ {
					if j == 1 {

						regla = regla + vr[j]
					} else {
						regla = regla + "," + vr[j]
					}
				}
			}

		}
		regla = regla + "). \n"
	}
	fmt.Println(regla)
	return
}

func (e *EntornoReglas) Agregar_predicado(predicado string) {
	e.predicados = e.predicados + predicado + "\n"
}

func (e *EntornoReglas) Obtener_predicados() (predicados string) {
	return e.predicados
}

func (e *EntornoReglas) ejecutar_regla(regla string, variable string) {
	var m = NewMachine()
	f := m.Consult(e.predicados)
	solutions := f.ProveAll(regla)
	for _, solution := range solutions {
		fmt.Printf("%s", solution.ByName_(variable))
		//fmt.Printf("%s", solution.ByName_("SC"))
	}
}
