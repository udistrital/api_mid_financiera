package controllers

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
)

// HomologacionController operations for homologation fo liquidation of titan
type HomologacionController struct {
	beego.Controller
}

// URLMapping ...
func (c *HomologacionController) URLMapping() {
	c.Mapping("MidHomologacionLiquidacion", c.MidHomologacionLiquidacion)
}

// MidHomologacionLiquidacion ...
// @Title MidHomologacionLiquidacion
// @Description homologa conceptos de titan con conceptos kronos
// @Param	idPreliquidacion	identificador de la liquidación de titan
// @Success 201 {object} models.Conceptos
// @Failure 403 body is empty
// @router /MidHomologacionLiquidacion/:idPreliquidacion [get]
func (c *HomologacionController) MidHomologacionLiquidacion() {
	idPreliquidacionStr := c.Ctx.Input.Param(":idPreliquidacion")
	idPreliquidacion, _ := strconv.Atoi(idPreliquidacionStr)
	fmt.Println(idPreliquidacion)
	var DetallePreliquidacion []interface{}
	var outputData interface{}

	// get data titan
	if err := getJson("http://"+beego.AppConfig.String("titanService")+"detalle_preliquidacion?query=Preliquidacion:"+idPreliquidacionStr+"&sortby=Concepto&order=desc", &DetallePreliquidacion); err == nil {
	} else {
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}
	// Control si no existe detalle de liquidacion
	if len(DetallePreliquidacion) == 0 {
		c.Data["json"] = "No existe DetallePreliquidacion"
		c.ServeJSON()
		return
	}

	fmt.Println("Para Kronos")
	//Envia data to kronos
	if err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/homologacion_concepto/HomolgacionConceptosTitan/", "POST", &outputData, &DetallePreliquidacion); err == nil {
	} else {
		fmt.Println("AAAAAA")
		fmt.Println(err.Error())
		fmt.Println("AAAAAA")
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	}
	c.Data["json"] = outputData
	c.ServeJSON()
}
