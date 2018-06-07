package controllers

import (
	"encoding/json"
	"fmt"
	//"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/api_mid_financiera/models"
	"github.com/udistrital/utils_oas/request"

)

type GestionSucursalesController struct {
	beego.Controller
}

func (c *GestionSucursalesController) URLMapping() {
	c.Mapping("InsertarSucursales", c.InsertarSucursales)
}

// InsertarSucursales ...
// @Title InsertarSucursales
// @Description InsertarSucursales
// @Param	body		body 	[]models.InformacionSucursales	true		"body for InformacionSucursales  content"
// @Success 201
// @Failure 403 body is empty
// @router insertar_sucursal/ [post]
func (c *GestionSucursalesController) InsertarSucursales() {

	fmt.Println("estoy en sucursaleeeeees")
	var info_sucursal models.InformacionSucursal
	var tipo_ente []models.TipoEnte
  var respuesta interface{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &info_sucursal); err == nil {

		//Primero, se busca el código del tipo ente correspondiente al código de abreviación
		if err = request.GetJson("http://"+beego.AppConfig.String("Urlente")+":"+beego.AppConfig.String("Portente")+"/"+beego.AppConfig.String("Nsente")+"/tipo_ente?query=CodigoAbreviacion:TE_3", &tipo_ente); err == nil {

			//Se inserta en ente y se devuelve el id registrado
				ente := &models.Ente {TipoEnte: &models.TipoEnte {Id: tipo_ente[0].Id}}
				if err := request.SendJson("http://"+beego.AppConfig.String("Urlente")+":"+beego.AppConfig.String("Portente")+"/"+beego.AppConfig.String("Nsente")+"/ente/", "POST", &respuesta, &ente); err == nil {

					fmt.Println("id ente", respuesta)
					//SE TOMA ESE ENTE Y SE INSERTA EN SUCURSAL, UBICACION Y CONTACTO

				 //	InsertarSucursal(infoSucursal.Nombre, 85)
				}else{
					fmt.Println("error al insertar ente: ", err)
				}


		}else{

			fmt.Println("error al consultar tipo ente: ", err)

		}
	} else {
		fmt.Println("err: ", err)

	}
	c.ServeJSON()
}
