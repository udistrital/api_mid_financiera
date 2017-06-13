package tools

import (
	"fmt"

	"github.com/astaxie/beego"
	//. "github.com/mndrix/golog"
	. "github.com/mndrix/golog"
	"github.com/udistrital/api_mid_financiera/models"
)

type EntornoReglas struct {
	dominio    string
	predicados string
}

func (e EntornoReglas) Agregar_dominio(dominio string) {
	var v []models.Predicado
	if err := getJson("http://"+beego.AppConfig.String("Urlruler")+":"+beego.AppConfig.String("Portruler")+"/"+beego.AppConfig.String("Nsruler")+"/predicado?limit=0&query=Dominio.Nombre:"+dominio, &v); err == nil {
		for i := 0; i < len(v); i++ {
			e.predicados = e.predicados + v[i].Nombre + "\n"
		}
	}
}

func (e EntornoReglas) Agregar_predicado(predicado string) {
	e.predicados = e.predicados + predicado + "\n"
}

func (e EntornoReglas) ejecutar_regla(regla string) {
	var m = NewMachine()
	f := m.Consult(e.predicados)
	solutions := f.ProveAll(regla)
	for _, solution := range solutions {
		fmt.Printf("%s", solution.ByName_("SD"))
		//fmt.Printf("%s", solution.ByName_("SC"))
	}
}
