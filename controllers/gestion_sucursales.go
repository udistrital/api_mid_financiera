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

// GestionSucursalesController operations for GestionSucursales
type GestionSucursalesController struct {
	beego.Controller
}
// URLMapping ...
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

	var infoSucursal models.InformacionSucursal
	var tipoEnte []models.TipoEnte
  var respuesta interface{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &infoSucursal); err == nil {

		ciudad,_ := strconv.Atoi(infoSucursal.Ciudad)
		departamento,_ := 	strconv.Atoi(infoSucursal.Departamento)
		pais,_ := strconv.Atoi(infoSucursal.Pais)
		//Primero, se busca el código del tipo ente correspondiente al código de abreviación
		if err = request.GetJson("http://"+beego.AppConfig.String("Urlente")+":"+beego.AppConfig.String("Portente")+"/"+beego.AppConfig.String("Nsente")+"/tipo_ente?query=CodigoAbreviacion:TE_3", &tipoEnte); err == nil {

			//Se inserta en ente y se devuelve el id registrado
				ente := &models.Ente {TipoEnte: &models.TipoEnte {Id: tipoEnte[0].Id}}
				if err := request.SendJson("http://"+beego.AppConfig.String("Urlente")+":"+beego.AppConfig.String("Portente")+"/"+beego.AppConfig.String("Nsente")+"/ente/", "POST", &respuesta, &ente); err == nil {
					res:= respuesta.(map[string]interface{})
					idTarget := res["Body"].(map[string]interface{})["Id"].(float64)
					idEnte := int(idTarget)
					//SE TOMA ESE ENTE Y SE INSERTA EN SUCURSAL, UBICACION Y CONTACTO
				 	respuesta, err = insertarSucursal(infoSucursal.Nombre, idEnte)
				  respuesta, err = insertarContacto(infoSucursal.Telefono,idEnte)
					respuesta, err = insertarUbicacion(infoSucursal.Direccion,pais, departamento, ciudad, idEnte)

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

	idSucur := c.GetString("id_sucursal")
	var sucursales []models.Organizacion
	if err := request.GetJson(beego.AppConfig.String("coreOrganizacionService")+"organizacion?query=Id:"+idSucur+",TipoOrganizacion.CodigoAbreviacion:SU", &sucursales); err == nil {

		var informacionSucursal  = make([]models.InformacionSucursal, len(sucursales))
		for i, suc := range sucursales{
			informacionSucursal[i].Nombre = suc.Nombre
			informacionSucursal[i].Direccion = buscarDireccion(suc.Ente)
			informacionSucursal[i].Telefono = buscarTelefono(suc.Ente)
			ubicaciones := buscarUbicaciones(suc.Ente)
			informacionSucursal[i].Pais, informacionSucursal[i].Departamento, informacionSucursal[i].Ciudad = buscarLugar(ubicaciones,suc.Ente)

    }

		c.Data["json"] = informacionSucursal
	}else{
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

		var informacionSucursal  = make([]models.InformacionSucursal, len(sucursales))
		for i, suc := range sucursales{
			informacionSucursal[i].Nombre = suc.Nombre
			informacionSucursal[i].Direccion = buscarDireccion(suc.Ente)
			informacionSucursal[i].Telefono = buscarTelefono(suc.Ente)
			ubicaciones := buscarUbicaciones(suc.Ente)
			informacionSucursal[i].Pais, informacionSucursal[i].Departamento, informacionSucursal[i].Ciudad = buscarLugar(ubicaciones,suc.Ente)

    }

		c.Data["json"] = informacionSucursal
	}else{
		c.Data["json"] = err
	}

		c.ServeJSON()
}

// ListarSucursalesBanco ...
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

// GetBancoSucursal ....
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
	sucursalID := strconv.FormatFloat(rpintfc.(map[string]interface{})["OrganizacionHija"].(float64), 'f', -1, 64)
	if err := request.GetJson(beego.AppConfig.String("coreOrganizacionService")+"organizacion/?query=Id:"+sucursalID, &resSucursal); err == nil {
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
	sucursalID := strconv.FormatFloat(rpintfc.(map[string]interface{})["OrganizacionPadre"].(float64), 'f', -1, 64)
	if err := request.GetJson(beego.AppConfig.String("coreOrganizacionService")+"organizacion/?query=Id:"+sucursalID, &resBanco); err == nil {
		if resBanco[0] != nil {
			rpintfc = resBanco[0]
		}
	}else{
		beego.Error("Error",err.Error());
	}
	return rpintfc
}



func insertarSucursal(nombre string, idEnte int)(res interface{}, err error){

	var tipoOrganizacion []models.TipoOrganizacion
	var respuesta interface{}
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlorganizacion")+":"+beego.AppConfig.String("Portorganizacion")+"/"+beego.AppConfig.String("Nsorganizacion")+"/tipo_organizacion?query=CodigoAbreviacion:SU", &tipoOrganizacion); err == nil {

			objetoOrganizacion := &models.Organizacion {Nombre: nombre, Ente: idEnte, TipoOrganizacion : &models.TipoOrganizacion{Id: tipoOrganizacion[0].Id}}
			if err := request.SendJson("http://"+beego.AppConfig.String("Urlorganizacion")+":"+beego.AppConfig.String("Portorganizacion")+"/"+beego.AppConfig.String("Nsorganizacion")+"/organizacion/", "POST", &respuesta, &objetoOrganizacion ); err == nil {

			}else{
				fmt.Println("error al insertar sucursal")
			}

	}else{
		fmt.Println("error al consultar tipo_organizacion")
	}

	return respuesta, err
}


func insertarContacto(telefono string, idEnte int)(res interface{}, err error){

	var tipoContacto []models.TipoContacto
	var respuesta interface{}
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlente")+":"+beego.AppConfig.String("Portente")+"/"+beego.AppConfig.String("Nsente")+"/tipo_contacto?query=CodigoAbreviacion:TEL", &tipoContacto); err == nil {

			objetoContacto := &models.ContactoEnte {Valor: telefono, Ente: &models.Ente {Id: idEnte}, TipoContacto : &models.TipoContacto{Id: tipoContacto[0].Id}}
			if err := request.SendJson("http://"+beego.AppConfig.String("Urlente")+":"+beego.AppConfig.String("Portente")+"/"+beego.AppConfig.String("Nsente")+"/contacto_ente/", "POST", &respuesta, &objetoContacto); err == nil {

			}else{
				fmt.Println("error al insertar contacto")
			}

	}else{
		fmt.Println("error al consultar tipo_contacto")
	}

	return respuesta, err
}


func insertarUbicacion(direccion string, pais, departamento, ciudad int,idEnte int)(res interface{}, err error){

	var tipoRelacionUbicacionEnte []models.TipoRelacionUbicacionEnte
	var respuesta interface{}
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlente")+":"+beego.AppConfig.String("Portente")+"/"+beego.AppConfig.String("Nsente")+"/tipo_relacion_ubicacion_ente?query=CodigoAbreviacion:LR", &tipoRelacionUbicacionEnte); err == nil {

			objetoUbicacionEnte := &models.UbicacionEnte {Lugar: pais, Ente: &models.Ente {Id: idEnte}, TipoRelacionUbicacionEnte : &models.TipoRelacionUbicacionEnte{Id: tipoRelacionUbicacionEnte[0].Id}, Activo: true}
			if err := request.SendJson("http://"+beego.AppConfig.String("Urlente")+":"+beego.AppConfig.String("Portente")+"/"+beego.AppConfig.String("Nsente")+"/ubicacion_ente/", "POST", &respuesta, &objetoUbicacionEnte); err == nil {

			}else{
				fmt.Println("error al insertar pais")
			}

	}else{
		fmt.Println("error al consultar tipo_ubicacion")
	}


	if err := request.GetJson("http://"+beego.AppConfig.String("Urlente")+":"+beego.AppConfig.String("Portente")+"/"+beego.AppConfig.String("Nsente")+"/tipo_relacion_ubicacion_ente?query=CodigoAbreviacion:LR", &tipoRelacionUbicacionEnte); err == nil {

			objetoUbicacionEnte := &models.UbicacionEnte {Lugar: departamento, Ente: &models.Ente {Id: idEnte}, TipoRelacionUbicacionEnte : &models.TipoRelacionUbicacionEnte{Id: tipoRelacionUbicacionEnte[0].Id}, Activo: true}
			if err := request.SendJson("http://"+beego.AppConfig.String("Urlente")+":"+beego.AppConfig.String("Portente")+"/"+beego.AppConfig.String("Nsente")+"/ubicacion_ente/", "POST", &respuesta, &objetoUbicacionEnte); err == nil {

			}else{
				fmt.Println("error al insertar departamento")
			}

	}else{
		fmt.Println("error al consultar tipo_ubicacion")
	}

	if err := request.GetJson("http://"+beego.AppConfig.String("Urlente")+":"+beego.AppConfig.String("Portente")+"/"+beego.AppConfig.String("Nsente")+"/tipo_relacion_ubicacion_ente?query=CodigoAbreviacion:LR", &tipoRelacionUbicacionEnte); err == nil {

			objetoUbicacionEnte := &models.UbicacionEnte {Lugar: ciudad, Ente: &models.Ente {Id: idEnte}, TipoRelacionUbicacionEnte : &models.TipoRelacionUbicacionEnte{Id: tipoRelacionUbicacionEnte[0].Id}, Activo: true}
			if err := request.SendJson("http://"+beego.AppConfig.String("Urlente")+":"+beego.AppConfig.String("Portente")+"/"+beego.AppConfig.String("Nsente")+"/ubicacion_ente/", "POST", &respuesta, &objetoUbicacionEnte); err == nil {

			}else{
				fmt.Println("error al insertar ciudad")
			}

	}else{
		fmt.Println("error al consultar tipo_ubicacion")
	}

	//FALTA DIRECCION
	/*
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlente")+":"+beego.AppConfig.String("Portente")+"/"+beego.AppConfig.String("Nsente")+"/tipo_relacion_ubicacion_ente?query=CodigoAbreviacion:LR", &tipoRelacionUbicacionEnte); err == nil {

			objetoUbicacionEnte := &models.UbicacionEnte {Lugar: ciudad, Ente: &models.Ente {Id: idEnte}, TipoRelacionUbicacionEnte : &models.TipoRelacionUbicacionEnte{Id: tipoRelacionUbicacionEnte[0].Id}, Activo: true}
			if err := request.SendJson("http://"+beego.AppConfig.String("Urlente")+":"+beego.AppConfig.String("Portente")+"/"+beego.AppConfig.String("Nsente")+"/ubicacion_ente/", "POST", &respuesta, &objetoUbicacionEnte); err == nil {

			}else{
				fmt.Println("error al insertar ciudad")
			}

	}else{
		fmt.Println("error al consultar tipo_ubicacion")
	}
	*/
	return respuesta, err

}

func buscarDireccion(idEnte int)(direccion string){

	return "No registrado"
}


func buscarTelefono(idEnte int)(telefono string){

	var tel string;
	var contactoEnte []models.ContactoEnte
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlente")+":"+beego.AppConfig.String("Portente")+"/"+beego.AppConfig.String("Nsente")+"/contacto_ente?query=Ente:"+strconv.Itoa(idEnte), &contactoEnte); err == nil {
		if(contactoEnte != nil){
			tel = contactoEnte[0].Valor
		}else{
			tel = "No registrado"
		}

	}else{
		tel = "No registrado"
	}

 return tel
}

func buscarUbicaciones(idEnte int)(ub []models.UbicacionEnte){

	var ubicaciones []models.UbicacionEnte

	if err := request.GetJson("http://"+beego.AppConfig.String("Urlente")+":"+beego.AppConfig.String("Portente")+"/"+beego.AppConfig.String("Nsente")+"/ubicacion_ente?query=Ente:"+strconv.Itoa(idEnte), &ubicaciones); err == nil {
		if(ubicaciones != nil){

		}else{
			ubicaciones = nil;
		}

	}else{
		ubicaciones = nil;
	}

	return ubicaciones
}

func buscarLugar(ubicaciones []models.UbicacionEnte, idEnte int)(p, c,d string){

	var pais = "No registrado"
	var departamento = "No registrado"
	var ciudad = "No registrado"

  var objetoLugar []models.Lugar


  if(ubicaciones != nil){
				for _, ubi := range ubicaciones{

				  if err := request.GetJson("http://"+beego.AppConfig.String("Urlubicacion")+":"+beego.AppConfig.String("Portubicacion")+"/"+beego.AppConfig.String("Nsubicacion")+"/lugar?query=Id:"+strconv.Itoa(ubi.Lugar), &objetoLugar); err == nil {

						if(objetoLugar != nil && objetoLugar[0].Id != 0){


							if(objetoLugar[0].TipoLugar.CodigoAbreviacion == "CIUDAD"){
								ciudad = objetoLugar[0].Nombre
							}

							if(objetoLugar[0].TipoLugar.CodigoAbreviacion == "DEPARTAMENTO"){
								departamento = objetoLugar[0].Nombre
							}

							if(objetoLugar[0].TipoLugar.CodigoAbreviacion == "PAIS"){
								pais = objetoLugar[0].Nombre
							}
						}


				  }

			}



	}

	return pais, departamento, ciudad
}
