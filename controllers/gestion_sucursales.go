package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/api_mid_financiera/models"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/optimize"

)

type GestionSucursalesController struct {
	beego.Controller
}

func (c *GestionSucursalesController) URLMapping() {
	c.Mapping("InsertarSucursales", c.InsertarSucursales)
	c.Mapping("ListarSucursales", c.ListarSucursales)
}

// InsertarSucursales ...
// @Title InsertarSucursales
// @Description InsertarSucursales
// @Param	body		body 	[]models.InformacionSucursales	true		"body for InformacionSucursales  content"
// @Success 201
// @Failure 403 body is empty
// @router insertar_sucursal/ [post]
func (c *GestionSucursalesController) InsertarSucursales() {

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
					c.Data["json"] = "Error al insertar ente"
				}

				c.Data["json"] = respuesta
		}else{

			fmt.Println("error al consultar tipo ente: ", err)
			c.Data["json"] = "Error al insertar ente"
		}

		c.Data["json"] = respuesta
	} else {
		fmt.Println("err: ", err)
		c.Data["json"] = "Error al insertar ente"
	}


	c.ServeJSON()
}

// ListarSucursal ...
// @Title ListarSucursal
// @Description ListarSucursal
// @Param	id_sucursal	query	int	false	"id de la sucursal"
// @Success 201 {object} []models.InformacionSucursal
// @Failure 403 body is empty
// @router listar_sucursal/ [get]
func (c *GestionSucursalesController) ListarSucursal() {

	id_sucursal := c.GetString("id_sucursal")
	var sucursales []models.Organizacion
	beego.Error(beego.AppConfig.String("coreOrganizacionService")+"organizacion?query=Id:"+id_sucursal+",TipoOrganizacion.CodigoAbreviacion:SU")
	if err := request.GetJson(beego.AppConfig.String("coreOrganizacionService")+"organizacion?query=Id:"+id_sucursal+",TipoOrganizacion.CodigoAbreviacion:SU", &sucursales); err == nil {

		var informacion_sucursal  = make([]models.InformacionSucursal, len(sucursales))
		for i, suc := range sucursales{
			informacion_sucursal[i].Nombre = suc.Nombre
			informacion_sucursal[i].Direccion = BuscarDireccion(suc.Ente)
			informacion_sucursal[i].Telefono = BuscarTelefono(suc.Ente)
			ubicaciones := BuscarUbicaciones(suc.Ente)
			informacion_sucursal[i].Pais, informacion_sucursal[i].Departamento, informacion_sucursal[i].Ciudad = BuscarLugar(ubicaciones,suc.Ente)

    }

		c.Data["json"] = informacion_sucursal
	}else{
		beego.Error(err)
		c.Data["json"] = err
	}

		c.ServeJSON()
}

// ListarSucursales ...
// @Title ListarSucursales
// @Description ListarSucursales
// @Success 201 {object} []models.InformacionSucursal
// @Failure 403 body is empty
// @router listar_sucursales/ [get]
func (c *GestionSucursalesController) ListarSucursales() {


	var sucursales []models.Organizacion
	if err := request.GetJson(beego.AppConfig.String("coreOrganizacionService")+"organizacion?query=TipoOrganizacion.CodigoAbreviacion:SU", &sucursales); err == nil {

		var informacion_sucursal  = make([]models.InformacionSucursal, len(sucursales))
		for i, suc := range sucursales{
			informacion_sucursal[i].Nombre = suc.Nombre
			informacion_sucursal[i].Direccion = BuscarDireccion(suc.Ente)
			informacion_sucursal[i].Telefono = BuscarTelefono(suc.Ente)
			ubicaciones := BuscarUbicaciones(suc.Ente)
			informacion_sucursal[i].Pais, informacion_sucursal[i].Departamento, informacion_sucursal[i].Ciudad = BuscarLugar(ubicaciones,suc.Ente)

    }

		c.Data["json"] = informacion_sucursal
	}else{
		c.Data["json"] = err
	}

		c.ServeJSON()
}

// GetOne ...
// @Title ListarSucursalesBanco
// @Description lista sucursales dado id banco
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} []models.InformacionSucursal
// @Failure 403 :idBanco is empty
// @router /ListarSucursalesBanco/:idBanco [get]
func (c *GestionSucursalesController) ListarSucursalesBanco() {
	defer c.ServeJSON()
	idBancoStr := c.Ctx.Input.Param(":idBanco")
	var orgHijas []interface{}
	if err := request.GetJson(beego.AppConfig.String("coreOrganizacionService")+"relacion_organizaciones/?query=OrganizacionPadre:" + idBancoStr, &orgHijas); err == nil {
		if (orgHijas!=nil){
			sucursales := optimize.ProccDigest(orgHijas, getValuesSucursales, nil, 3)
			c.Data["json"] = sucursales
		}
	}else{
		c.Data["json"] = models.Alert{Type: "error", Code: "E_0458", Body: err}
	}
}


func GetBancoSucursal(idSucursalStr string)(res interface{},err error) {
	var orgPadre []interface{}
	if err = request.GetJson(beego.AppConfig.String("coreOrganizacionService")+"relacion_organizaciones/?query=OrganizacionHija:" + idSucursalStr, &orgPadre); err == nil {
		if (orgPadre!=nil){
			res = optimize.ProccDigest(orgPadre, getValuesBancos, nil, 3)
		}
	}
	return
}

func getValuesSucursales(rpintfc interface{}, params ...interface{}) (res interface{}) {
	var resSucursal []map[string]interface{}
	sucursalId := strconv.FormatFloat(rpintfc.(map[string]interface{})["OrganizacionHija"].(float64), 'f', -1, 64)
	if err := request.GetJson(beego.AppConfig.String("coreOrganizacionService")+"organizacion/?query=Id:"+sucursalId, &resSucursal); err == nil {
		if resSucursal[0] != nil {
			rpintfc.(map[string]interface{})["OrganizacionHija"] = resSucursal[0]
		}
	}else{
		beego.Error("Error",err.Error());
	}
	return rpintfc
}

func getValuesBancos(rpintfc interface{}, params ...interface{}) (res interface{}) {
	var resBanco []map[string]interface{}
	sucursalId := strconv.FormatFloat(rpintfc.(map[string]interface{})["OrganizacionPadre"].(float64), 'f', -1, 64)
	if err := request.GetJson(beego.AppConfig.String("coreOrganizacionService")+"organizacion/?query=Id:"+sucursalId, &resBanco); err == nil {
		if resBanco[0] != nil {
			rpintfc = resBanco[0]
		}
	}else{
		beego.Error("Error",err.Error());
	}
	return rpintfc
}



func InsertarSucursal(nombre string, id_ente int)(res interface{}, err error){

	var tipo_organizacion []models.TipoOrganizacion
	var respuesta interface{}
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlorganizacion")+":"+beego.AppConfig.String("Portorganizacion")+"/"+beego.AppConfig.String("Nsorganizacion")+"/tipo_organizacion?query=CodigoAbreviacion:SU", &tipo_organizacion); err == nil {

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

	//FALTA DIRECCION
	/*
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlente")+":"+beego.AppConfig.String("Portente")+"/"+beego.AppConfig.String("Nsente")+"/tipo_relacion_ubicacion_ente?query=CodigoAbreviacion:LR", &tipo_relacion_ubicacion_ente); err == nil {

			objeto_ubicacion_ente := &models.UbicacionEnte {Lugar: ciudad, Ente: &models.Ente {Id: id_ente}, TipoRelacionUbicacionEnte : &models.TipoRelacionUbicacionEnte{Id: tipo_relacion_ubicacion_ente[0].Id}, Activo: true}
			if err := request.SendJson("http://"+beego.AppConfig.String("Urlente")+":"+beego.AppConfig.String("Portente")+"/"+beego.AppConfig.String("Nsente")+"/ubicacion_ente/", "POST", &respuesta, &objeto_ubicacion_ente); err == nil {

			}else{
				fmt.Println("error al insertar ciudad")
			}

	}else{
		fmt.Println("error al consultar tipo_ubicacion")
	}
	*/
	return respuesta, err

}

func BuscarDireccion(id_ente int)(direccion string){

	return "No registrado"
}


func BuscarTelefono(id_ente int)(telefono string){

	var tel string;
	var contacto_ente []models.ContactoEnte
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlente")+":"+beego.AppConfig.String("Portente")+"/"+beego.AppConfig.String("Nsente")+"/contacto_ente?query=Ente:"+strconv.Itoa(id_ente), &contacto_ente); err == nil {
		if(contacto_ente != nil){
			tel = contacto_ente[0].Valor
		}else{
			tel = "No registrado"
		}

	}else{
		tel = "No registrado"
	}

 return tel
}

func BuscarUbicaciones(id_ente int)(ub []models.UbicacionEnte){

	var ubicaciones []models.UbicacionEnte

	if err := request.GetJson("http://"+beego.AppConfig.String("Urlente")+":"+beego.AppConfig.String("Portente")+"/"+beego.AppConfig.String("Nsente")+"/ubicacion_ente?query=Ente:"+strconv.Itoa(id_ente), &ubicaciones); err == nil {
		if(ubicaciones != nil){

		}else{
			ubicaciones = nil;
		}

	}else{
		ubicaciones = nil;
	}

	return ubicaciones
}

func BuscarLugar(ubicaciones []models.UbicacionEnte, id_ente int)(p, c,d string){

	var pais = "No registrado"
	var departamento = "No registrado"
	var ciudad = "No registrado"

  var objeto_lugar []models.Lugar


  if(ubicaciones != nil){
				for _, ubi := range ubicaciones{

				  if err := request.GetJson("http://"+beego.AppConfig.String("Urlubicacion")+":"+beego.AppConfig.String("Portubicacion")+"/"+beego.AppConfig.String("Nsubicacion")+"/lugar?query=Id:"+strconv.Itoa(ubi.Lugar), &objeto_lugar); err == nil {

						if(objeto_lugar != nil && objeto_lugar[0].Id != 0){


							if(objeto_lugar[0].TipoLugar.CodigoAbreviacion == "CIUDAD"){
								ciudad = objeto_lugar[0].Nombre
							}

							if(objeto_lugar[0].TipoLugar.CodigoAbreviacion == "DEPARTAMENTO"){
								departamento = objeto_lugar[0].Nombre
							}

							if(objeto_lugar[0].TipoLugar.CodigoAbreviacion == "PAIS"){
								pais = objeto_lugar[0].Nombre
							}
						}


				  }

			}



	}

	return pais, departamento, ciudad
}
