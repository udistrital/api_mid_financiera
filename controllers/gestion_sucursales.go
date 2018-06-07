package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

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

		ciudad,_ := strconv.Atoi(info_sucursal.Ciudad)
		departamento,_ := 	strconv.Atoi(info_sucursal.Departamento)
		pais,_ := strconv.Atoi(info_sucursal.Pais)
		//Primero, se busca el código del tipo ente correspondiente al código de abreviación
		if err = request.GetJson("http://"+beego.AppConfig.String("Urlente")+":"+beego.AppConfig.String("Portente")+"/"+beego.AppConfig.String("Nsente")+"/tipo_ente?query=CodigoAbreviacion:TE_3", &tipo_ente); err == nil {

			//Se inserta en ente y se devuelve el id registrado
				ente := &models.Ente {TipoEnte: &models.TipoEnte {Id: tipo_ente[0].Id}}
				if err := request.SendJson("http://"+beego.AppConfig.String("Urlente")+":"+beego.AppConfig.String("Portente")+"/"+beego.AppConfig.String("Nsente")+"/ente/", "POST", &respuesta, &ente); err == nil {
					res:= respuesta.(map[string]interface{})
					idTarget := res["Body"].(map[string]interface{})["Id"].(float64)
					id_ente := int(idTarget)
					//SE TOMA ESE ENTE Y SE INSERTA EN SUCURSAL, UBICACION Y CONTACTO
				 	respuesta, err = InsertarSucursal(info_sucursal.Nombre, id_ente)
				  respuesta, err = InsertarContacto(info_sucursal.Telefono,id_ente)
					respuesta, err = InsertarUbicacion(info_sucursal.Direccion,pais, departamento, ciudad, id_ente)

					c.Data["json"] = respuesta

				}else{
					fmt.Println("error al insertar ente: ", err)
					c.Data["json"] = err
				}

				c.Data["json"] = respuesta
		}else{

			fmt.Println("error al consultar tipo ente: ", err)
			c.Data["json"] = err
		}

		c.Data["json"] = respuesta
	} else {
		fmt.Println("err: ", err)
		c.Data["json"] = err
	}


	c.ServeJSON()
}

func InsertarSucursal(nombre string, id_ente int)(res interface{}, err error){

	var tipo_organizacion []models.TipoOrganizacion
	var respuesta interface{}
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlorganizacion")+":"+beego.AppConfig.String("Portorganizacion")+"/"+beego.AppConfig.String("Nsorganizacion")+"/tipo_organizacion?query=CodigoAbreviacion:TO_2", &tipo_organizacion); err == nil {

			objeto_organizacion := &models.Organizacion {Nombre: nombre, Ente: id_ente, TipoOrganizacion : &models.TipoOrganizacion{Id: tipo_organizacion[0].Id}}
			if err := request.SendJson("http://"+beego.AppConfig.String("Urlorganizacion")+":"+beego.AppConfig.String("Portorganizacion")+"/"+beego.AppConfig.String("Nsorganizacion")+"/organizacion/", "POST", &respuesta, &objeto_organizacion); err == nil {

			}else{
				fmt.Println("error al insertar sucursal")
			}

	}else{
		fmt.Println("error al consultar tipo_organizacion")
	}

	return respuesta, err
}


func InsertarContacto(telefono string, id_ente int)(res interface{}, err error){

	var tipo_contacto []models.TipoContacto
	var respuesta interface{}
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlente")+":"+beego.AppConfig.String("Portente")+"/"+beego.AppConfig.String("Nsente")+"/tipo_contacto?query=CodigoAbreviacion:TEL", &tipo_contacto); err == nil {

			objeto_contacto := &models.ContactoEnte {Valor: telefono, Ente: &models.Ente {Id: id_ente}, TipoContacto : &models.TipoContacto{Id: tipo_contacto[0].Id}}
			if err := request.SendJson("http://"+beego.AppConfig.String("Urlente")+":"+beego.AppConfig.String("Portente")+"/"+beego.AppConfig.String("Nsente")+"/contacto_ente/", "POST", &respuesta, &objeto_contacto); err == nil {

			}else{
				fmt.Println("error al insertar contacto")
			}

	}else{
		fmt.Println("error al consultar tipo_contacto")
	}

	return respuesta, err
}


func InsertarUbicacion(direccion string, pais, departamento, ciudad int,id_ente int)(res interface{}, err error){

	var tipo_relacion_ubicacion_ente []models.TipoRelacionUbicacionEnte
	var respuesta interface{}
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlente")+":"+beego.AppConfig.String("Portente")+"/"+beego.AppConfig.String("Nsente")+"/tipo_relacion_ubicacion_ente?query=CodigoAbreviacion:LR", &tipo_relacion_ubicacion_ente); err == nil {

			objeto_ubicacion_ente := &models.UbicacionEnte {Lugar: pais, Ente: &models.Ente {Id: id_ente}, TipoRelacionUbicacionEnte : &models.TipoRelacionUbicacionEnte{Id: tipo_relacion_ubicacion_ente[0].Id}, Activo: true}
			if err := request.SendJson("http://"+beego.AppConfig.String("Urlente")+":"+beego.AppConfig.String("Portente")+"/"+beego.AppConfig.String("Nsente")+"/ubicacion_ente/", "POST", &respuesta, &objeto_ubicacion_ente); err == nil {

			}else{
				fmt.Println("error al insertar pais")
			}

	}else{
		fmt.Println("error al consultar tipo_ubicacion")
	}


	if err := request.GetJson("http://"+beego.AppConfig.String("Urlente")+":"+beego.AppConfig.String("Portente")+"/"+beego.AppConfig.String("Nsente")+"/tipo_relacion_ubicacion_ente?query=CodigoAbreviacion:LR", &tipo_relacion_ubicacion_ente); err == nil {

			objeto_ubicacion_ente := &models.UbicacionEnte {Lugar: departamento, Ente: &models.Ente {Id: id_ente}, TipoRelacionUbicacionEnte : &models.TipoRelacionUbicacionEnte{Id: tipo_relacion_ubicacion_ente[0].Id}, Activo: true}
			if err := request.SendJson("http://"+beego.AppConfig.String("Urlente")+":"+beego.AppConfig.String("Portente")+"/"+beego.AppConfig.String("Nsente")+"/ubicacion_ente/", "POST", &respuesta, &objeto_ubicacion_ente); err == nil {

			}else{
				fmt.Println("error al insertar departamento")
			}

	}else{
		fmt.Println("error al consultar tipo_ubicacion")
	}

	if err := request.GetJson("http://"+beego.AppConfig.String("Urlente")+":"+beego.AppConfig.String("Portente")+"/"+beego.AppConfig.String("Nsente")+"/tipo_relacion_ubicacion_ente?query=CodigoAbreviacion:LR", &tipo_relacion_ubicacion_ente); err == nil {

			objeto_ubicacion_ente := &models.UbicacionEnte {Lugar: ciudad, Ente: &models.Ente {Id: id_ente}, TipoRelacionUbicacionEnte : &models.TipoRelacionUbicacionEnte{Id: tipo_relacion_ubicacion_ente[0].Id}, Activo: true}
			if err := request.SendJson("http://"+beego.AppConfig.String("Urlente")+":"+beego.AppConfig.String("Portente")+"/"+beego.AppConfig.String("Nsente")+"/ubicacion_ente/", "POST", &respuesta, &objeto_ubicacion_ente); err == nil {

			}else{
				fmt.Println("error al insertar ciudad")
			}

	}else{
		fmt.Println("error al consultar tipo_ubicacion")
	}

	return respuesta, err

}
