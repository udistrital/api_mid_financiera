package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/manucorporat/try"
	"github.com/mitchellh/mapstructure"
	"github.com/udistrital/api_mid_financiera/models"
	"github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/utils_oas/optimize"
	"github.com/udistrital/utils_oas/request"
)

type GestionSucursalesController struct {
	beego.Controller
}

func (c *GestionSucursalesController) URLMapping() {
	c.Mapping("InsertarSucursales", c.InsertarSucursales)
	c.Mapping("ListarSucursales", c.ListarSucursales)
	c.Mapping("Put", c.Put)
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
	var ente models.Ente
	try.This(func() {
		if err := json.Unmarshal(c.Ctx.Input.RequestBody, &info_sucursal); err == nil {
			ciudad := int(info_sucursal.Ciudad.(map[string]interface{})["Id"].(float64))
			departamento := int(info_sucursal.Departamento.(map[string]interface{})["Id"].(float64))
			pais := int(info_sucursal.Pais.(map[string]interface{})["Id"].(float64))
			//Primero, se busca el código del tipo ente correspondiente al código de abreviación
			if err = request.GetJson(beego.AppConfig.String("coreEnteService")+"tipo_ente?query=CodigoAbreviacion:TE_3", &tipo_ente); err == nil {

				respuesta, err = InsertarSucursal(info_sucursal.Organizacion.Nombre)
				ente.Id = int(respuesta.(map[string]interface{})["Ente"].(float64))
				ente.TipoEnte = tipo_ente[0]
				idEnteStr := strconv.Itoa(ente.Id)
				respuesta.(map[string]interface{})["Ente"] = &ente

				if err := request.SendJson(beego.AppConfig.String("coreEnteService")+"ente/"+idEnteStr, "PUT", respuesta.(map[string]interface{})["Ente"], ente); err != nil {
					beego.Error("Error", err.Error())
					c.Data["json"] = map[string]interface{}{"Code": "E_0458", "Body": err.Error(), "Type": "error"}
					return
				}
				if _, err = InsertarContacto(info_sucursal.Telefono.(map[string]interface{})["Valor"].(string), ente.Id); err != nil {
					c.Data["json"] = map[string]interface{}{"Code": "E_0458", "Body": err.Error(), "Type": "error"}
					return
				}

				if _, err = InsertarUbicacion(info_sucursal.Direccion.(map[string]interface{})["Valor"].(string), pais, departamento, ciudad, ente.Id); err != nil {
					c.Data["json"] = map[string]interface{}{"Code": "E_0458", "Body": err.Error(), "Type": "error"}
					return
				}

				c.Data["json"] = map[string]interface{}{"Code": "S_543", "Body": respuesta, "Type": "success"}

			} else {
				c.Data["json"] = map[string]interface{}{"Code": "E_0458", "Body": err.Error(), "Type": "error"}
			}
		} else {
			c.Data["json"] = map[string]interface{}{"Code": "E_0458", "Body": err.Error(), "Type": "error"}
		}
	}).Catch(func(e try.E) {
		beego.Error("expc ", e)
		c.Data["json"] = map[string]interface{}{"Code": "E_0458", "Body": e, "Type": "error"}
	})

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
	if err := request.GetJson(beego.AppConfig.String("coreOrganizacionService")+"organizacion?query=Id:"+id_sucursal+",TipoOrganizacion.CodigoAbreviacion:SU", &sucursales); err == nil {

		var informacion_sucursal = make([]models.InformacionSucursal, len(sucursales))
		for i, suc := range sucursales {
			informacion_sucursal[i].Organizacion.Nombre = suc.Nombre
			informacion_sucursal[i].Telefono = BuscarTelefono(suc.Ente)
			ubicaciones := BuscarUbicaciones(suc.Ente)
			informacion_sucursal[i].Pais, informacion_sucursal[i].Departamento, informacion_sucursal[i].Ciudad, informacion_sucursal[i].Direccion = BuscarLugar(ubicaciones, suc.Ente)
		}
		c.Data["json"] = informacion_sucursal
	} else {
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
	if err := request.GetJson(beego.AppConfig.String("coreOrganizacionService")+"organizacion/?query=TipoOrganizacion.CodigoAbreviacion:SU&limit=-1", &sucursales); err == nil {

		var informacion_sucursal = make([]models.InformacionSucursal, len(sucursales))
		for i, suc := range sucursales {
			informacion_sucursal[i].Telefono = BuscarTelefono(suc.Ente)
			ubicaciones := BuscarUbicaciones(suc.Ente)
			informacion_sucursal[i].Pais, informacion_sucursal[i].Departamento, informacion_sucursal[i].Ciudad, informacion_sucursal[i].Direccion = BuscarLugar(ubicaciones, suc.Ente)
			informacion_sucursal[i].Organizacion = suc

		}

		c.Data["json"] = informacion_sucursal
	} else {
		beego.Error("Error ", err)
		c.Data["json"] = err
	}

	c.ServeJSON()
}

// ListarSoloSucursalesBanco ...
// @Title ListarSoloSucursalesBanco
// @Description lista sucursales dado id banco
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} []models.InformacionSucursal
// @Failure 403 :idBanco is empty
// @router /ListarSoloSucursalesBanco/:idBanco [get]
func (c *GestionSucursalesController) ListarSoloSucursalesBanco() {
	defer c.ServeJSON()
	idBancoStr := c.Ctx.Input.Param(":idBanco")
	var orgHijas []interface{}
	if err := request.GetJson(beego.AppConfig.String("coreOrganizacionService")+"relacion_organizaciones/?query=OrganizacionPadre:"+idBancoStr, &orgHijas); err == nil {
		if orgHijas != nil {
			sucursales := optimize.ProccDigest(orgHijas, getValuesSucursalesOnly, nil, 3)
			c.Data["json"] = sucursales
		}
	} else {
		c.Data["json"] = models.Alert{Type: "error", Code: "E_0458", Body: err}
	}
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
	if err := request.GetJson(beego.AppConfig.String("coreOrganizacionService")+"relacion_organizaciones/?query=OrganizacionPadre:"+idBancoStr, &orgHijas); err == nil {
		if orgHijas != nil {
			sucursales := optimize.ProccDigest(orgHijas, getValuesSucursales, nil, 3)
			c.Data["json"] = sucursales
		}
	} else {
		c.Data["json"] = models.Alert{Type: "error", Code: "E_0458", Body: err}
	}
}

func GetBancoSucursal(idSucursalStr string) (res interface{}, err error) {
	var orgPadre []interface{}
	if err = request.GetJson(beego.AppConfig.String("coreOrganizacionService")+"relacion_organizaciones/?query=OrganizacionHija:"+idSucursalStr, &orgPadre); err == nil {
		if orgPadre != nil {
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
			id_ente := int(resSucursal[0]["Ente"].(float64))
			rpintfc.(map[string]interface{})["OrganizacionHija"] = resSucursal[0]
			rpintfc.(map[string]interface{})["OrganizacionHija"].(map[string]interface{})["Telefono"] = BuscarTelefono(id_ente)
			ubicaciones := BuscarUbicaciones(id_ente)
			rpintfc.(map[string]interface{})["OrganizacionHija"].(map[string]interface{})["Pais"], rpintfc.(map[string]interface{})["OrganizacionHija"].(map[string]interface{})["Departamento"], rpintfc.(map[string]interface{})["OrganizacionHija"].(map[string]interface{})["Ciudad"], rpintfc.(map[string]interface{})["OrganizacionHija"].(map[string]interface{})["Direccion"] = BuscarLugar(ubicaciones, id_ente)
		}
	} else {
		beego.Error("Error", err.Error())
	}
	return rpintfc
}

func getValuesSucursalesOnly(rpintfc interface{}, params ...interface{}) (res interface{}) {
	var resSucursal []map[string]interface{}
	sucursalId := strconv.FormatFloat(rpintfc.(map[string]interface{})["OrganizacionHija"].(float64), 'f', -1, 64)
	if err := request.GetJson(beego.AppConfig.String("coreOrganizacionService")+"organizacion/?query=Id:"+sucursalId, &resSucursal); err == nil {
		if resSucursal[0] != nil {
			id_ente := int(resSucursal[0]["Ente"].(float64))
			rpintfc = resSucursal[0]
			rpintfc.(map[string]interface{})["Telefono"] = BuscarTelefono(id_ente)
			ubicaciones := BuscarUbicaciones(id_ente)
			rpintfc.(map[string]interface{})["Pais"], rpintfc.(map[string]interface{})["Departamento"], rpintfc.(map[string]interface{})["Ciudad"], rpintfc.(map[string]interface{})["Direccion"] = BuscarLugar(ubicaciones, id_ente)
		}
	} else {
		beego.Error("Error", err.Error())
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
	} else {
		beego.Error("Error", err.Error())
	}
	return rpintfc
}

func InsertarSucursal(nombre string) (res interface{}, err error) {

	var tipo_organizacion []models.TipoOrganizacion
	var respuesta interface{}
	if err := request.GetJson(beego.AppConfig.String("coreOrganizacionService")+"tipo_organizacion?query=CodigoAbreviacion:SU", &tipo_organizacion); err == nil {
		objeto_organizacion := &models.Organizacion{Nombre: nombre, TipoOrganizacion: &models.TipoOrganizacion{Id: tipo_organizacion[0].Id}}
		if err := request.SendJson(beego.AppConfig.String("coreOrganizacionService")+"organizacion/", "POST", &respuesta, &objeto_organizacion); err != nil {
			fmt.Println("error al insertar sucursal")
		}
	} else {
		fmt.Println("error al consultar tipo_organizacion")
	}
	return respuesta, err
}

func InsertarContacto(telefono string, id_ente int) (res interface{}, err error) {

	var tipo_contacto []models.TipoContacto
	var respuesta interface{}
	if err := request.GetJson(beego.AppConfig.String("coreEnteService")+"tipo_contacto?query=CodigoAbreviacion:TEL", &tipo_contacto); err == nil {

		objeto_contacto := &models.ContactoEnte{Valor: telefono, Ente: &models.Ente{Id: id_ente}, TipoContacto: &models.TipoContacto{Id: tipo_contacto[0].Id}}
		if err := request.SendJson(beego.AppConfig.String("coreEnteService")+"contacto_ente/", "POST", &respuesta, &objeto_contacto); err != nil {
			beego.Error("error al insertar contacto")
		}
	} else {
		beego.Error("error al consultar tipo_contacto")
	}

	return respuesta, err
}

func InsertarUbicacion(direccion string, pais, departamento, ciudad, id_ente int) (res interface{}, err error) {

	var respuesta interface{}
	var ubicacionEnteCiudad models.UbicacionEnte
	if _, err = InsertarLugar(pais, id_ente); err != nil {
		fmt.Println("error al insertar pais")
		beego.Error(err.Error())
	}

	if _, err = InsertarLugar(departamento, id_ente); err != nil {
		fmt.Println("error al insertar depto")
		beego.Error(err.Error())
	}

	if respuesta, err = InsertarLugar(ciudad, id_ente); err != nil {
		fmt.Println("error al insertar ciudad")
		beego.Error(err.Error())
	}
	beego.Error("respuesta", respuesta)
	if err = formatdata.FillStruct(respuesta.(map[string]interface{})["Body"], &ubicacionEnteCiudad); err != nil {
		beego.Error(err.Error())
	}
	if err = InsertarDireccion(direccion, ubicacionEnteCiudad); err != nil {
		beego.Error(err.Error())
	}

	return respuesta, err

}

func InsertarLugar(lugar, id_ente int) (respuesta interface{}, err error) {
	var tipo_relacion_ubicacion_ente []models.TipoRelacionUbicacionEnte
	if err := request.GetJson(beego.AppConfig.String("coreEnteService")+"tipo_relacion_ubicacion_ente/?query=CodigoAbreviacion:LU", &tipo_relacion_ubicacion_ente); err == nil {
		objeto_ubicacion_ente := &models.UbicacionEnte{Lugar: lugar, Ente: &models.Ente{Id: id_ente}, TipoRelacionUbicacionEnte: &models.TipoRelacionUbicacionEnte{Id: tipo_relacion_ubicacion_ente[0].Id}, Activo: true}
		if err = request.SendJson(beego.AppConfig.String("coreEnteService")+"ubicacion_ente/", "POST", &respuesta, objeto_ubicacion_ente); err != nil {
			beego.Error(err.Error())
		} else {
			beego.Error("respuesta ", respuesta)
		}
	} else {
		beego.Error(err.Error())
		err = errors.New("error al consultar tipo_ubicacion")
	}
	return
}

func InsertarDireccion(direccion string, ubicacionEnte models.UbicacionEnte) (err error) {
	var atributoUbicacion []interface{}
	var respuesta interface{}
	if err := request.GetJson(beego.AppConfig.String("coreEnteService")+"atributo_ubicacion/?query=NumeroOrden:1", &atributoUbicacion); err == nil {
		valAtributoUbicacion := &models.ValorAtributoUbicacion{UbicacionEnte: &ubicacionEnte, AtributoUbicacion: atributoUbicacion[0], Valor: direccion}
		if err = request.SendJson(beego.AppConfig.String("coreEnteService")+"valor_atributo_ubicacion/", "POST", &respuesta, valAtributoUbicacion); err != nil {
			beego.Error(err.Error())
		} else {
			beego.Error("respuesta ", respuesta)
		}
	}
	return
}

func BuscarTelefono(id_ente int) (telefono interface{}) {

	var contacto_ente []models.ContactoEnte
	if err := request.GetJson(beego.AppConfig.String("coreEnteService")+"contacto_ente/?query=Ente:"+strconv.Itoa(id_ente)+",TipoContacto.CodigoAbreviacion:TEL", &contacto_ente); err == nil {
		if contacto_ente != nil {
			telefono = contacto_ente[0]
		}
	} else {
		beego.Error(err.Error())
	}
	return
}

func BuscarUbicaciones(id_ente int) (ub []models.UbicacionEnte) {

	var ubicaciones []models.UbicacionEnte

	if err := request.GetJson(beego.AppConfig.String("coreEnteService")+"ubicacion_ente/?query=Ente:"+strconv.Itoa(id_ente), &ubicaciones); err != nil {
		beego.Error(err.Error())
		ubicaciones = nil
	}

	return ubicaciones
}

func BuscarLugar(ubicaciones []models.UbicacionEnte, id_ente int) (p, c, d, dir interface{}) {

	var pais map[string]interface{}
	var departamento map[string]interface{}
	var ciudad map[string]interface{}
	var direccion interface{}
	var objeto_lugar []models.Lugar
	var valAtribUbic []map[string]interface{}

	if ubicaciones != nil {
		for _, ubi := range ubicaciones {
			if err := request.GetJson(beego.AppConfig.String("coreUbicacionService")+"lugar?query=Id:"+strconv.Itoa(ubi.Lugar), &objeto_lugar); err == nil {
				if objeto_lugar != nil && objeto_lugar[0].Id != 0 {
					if objeto_lugar[0].TipoLugar.NumeroOrden == 3 {
						if err = mapstructure.Decode(objeto_lugar[0], &ciudad); err != nil {
							beego.Error(err.Error)
						}
						ciudad["UbicacionEnte"] = ubi
					}
					if objeto_lugar[0].TipoLugar.NumeroOrden == 2 {
						if err = mapstructure.Decode(objeto_lugar[0], &departamento); err != nil {
							beego.Error(err.Error)
						}
						departamento["UbicacionEnte"] = ubi
					}
					if objeto_lugar[0].TipoLugar.NumeroOrden == 1 {
						if err = mapstructure.Decode(objeto_lugar[0], &pais); err != nil {
							beego.Error(err.Error)
						}
						pais["UbicacionEnte"] = ubi
					}
					if err := request.GetJson(beego.AppConfig.String("coreEnteService")+"valor_atributo_ubicacion/?query=AtributoUbicacion.NumeroOrden:1,UbicacionEnte:"+strconv.Itoa(ubi.Id), &valAtribUbic); err == nil {
						if valAtribUbic != nil {
							direccion = valAtribUbic[0]
						}
					}
				}
			}
		}
	}

	return pais, departamento, ciudad, direccion
}

// EditarSucursal ...
// @Title Editar Sucursal
// @Description update the sucursal Information
// @Param	idEnte		path 	string	true		"The id you want to update"
// @Param	body		body 	interface{} true		"body for Gestion_cheques content"
// @Success 200 {object} interface{}
// @Failure 403 Body is empty
// @router /EditarSucursal/:idEnte [put]
func (c *GestionSucursalesController) EditarSucursal() {
	defer c.ServeJSON()
	var v map[string]interface{}
	var idEnte int
	var telefono interface{}
	idEnteStr := c.Ctx.Input.Param(":idEnte")
	idEnte, _ = strconv.Atoi(idEnteStr)
	var respuesta interface{}
	var ubicacionEnteCiudad models.UbicacionEnte
	try.This(func() {
		if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
			telefono = v["Telefono"]
			if telefono.(map[string]interface{})["Id"] != nil {
				telStr := strconv.FormatFloat(telefono.(map[string]interface{})["Id"].(float64), 'f', -1, 64)
				if err := request.SendJson(beego.AppConfig.String("coreEnteService")+"contacto_ente/"+telStr, "PUT", &respuesta, telefono); err != nil {
					beego.Error(err.Error())
					c.Data["json"] = models.Alert{Type: "error", Code: "E_GS005", Body: err.Error()}
					return
				}
			} else {
				if _, err := InsertarContacto(telefono.(map[string]interface{})["Valor"].(string), idEnte); err != nil {
					beego.Error(err.Error())
					c.Data["json"] = models.Alert{Type: "error", Code: "E_GS005", Body: err.Error()}
					return
				}
			}

			if _, err = updateLugarUbicacion(v["Pais"], idEnte); err != nil {
				beego.Error(err.Error())
				c.Data["json"] = models.Alert{Type: "error", Code: "E_GS002", Body: err.Error()}
				return
			}

			if _, err = updateLugarUbicacion(v["Departamento"], idEnte); err != nil {
				beego.Error(err.Error())
				c.Data["json"] = models.Alert{Type: "error", Code: "E_GS006", Body: err.Error()}
				return
			}

			if respuesta, err = updateLugarUbicacion(v["Ciudad"], idEnte); err != nil {
				c.Data["json"] = models.Alert{Type: "error", Code: "E_GS001", Body: err.Error()}
				return
			} else {
				if err = formatdata.FillStruct(respuesta.(map[string]interface{})["Body"], &ubicacionEnteCiudad); err != nil {
					beego.Error(err.Error())
				} else {
					direccion := v["Direccion"]
					if direccion != nil {
						if direccion.(map[string]interface{})["Id"] != nil {
							dirStr := strconv.FormatFloat(direccion.(map[string]interface{})["Id"].(float64), 'f', -1, 64)
							if err := request.SendJson(beego.AppConfig.String("coreEnteService")+"valor_atributo_ubicacion/"+dirStr, "PUT", &respuesta, direccion); err != nil {
								beego.Error("Error", err.Error())
								c.Data["json"] = models.Alert{Type: "error", Code: "E_GS003", Body: err.Error()}
								return
							}
						} else {
							if err = InsertarDireccion(direccion.(map[string]interface{})["Valor"].(string), ubicacionEnteCiudad); err != nil {
								beego.Error("Error", err.Error())
								c.Data["json"] = models.Alert{Type: "error", Code: "E_GS003", Body: err.Error()}
								return
							}
						}
					}

				}
			}
			organizacion := v["Organizacion"]
			if organizacion != nil {
				idOrgStr := strconv.FormatFloat(organizacion.(map[string]interface{})["Id"].(float64), 'f', -1, 64)
				if err := request.SendJson(beego.AppConfig.String("coreOrganizacionService")+"organizacion/"+idOrgStr, "PUT", &respuesta, organizacion); err != nil {
					beego.Error("Error", err.Error())
					c.Data["json"] = models.Alert{Type: "error", Code: "E_GS004", Body: err.Error()}
				}
			}
			c.Data["json"] = models.Alert{Type: "success", Code: "S_542", Body: respuesta}
		} else {
			c.Data["json"] = models.Alert{Type: "error", Code: "E_GS004", Body: err.Error()}
		}
	}).Catch(func(e try.E) {
		beego.Error("expc ", e)
		c.Data["json"] = map[string]interface{}{"Code": "E_0458", "Body": e, "Type": "error"}
	})
}

func updateLugarUbicacion(lugar interface{}, idEnte int) (respuesta interface{}, err error) {
	if lugar != nil {
		if lugar.(map[string]interface{})["Id"] != nil {
			idLugarStr := strconv.FormatFloat(lugar.(map[string]interface{})["Id"].(float64), 'f', -1, 64)
			if err = request.SendJson(beego.AppConfig.String("coreEnteService")+"ubicacion_ente/"+idLugarStr, "PUT", &respuesta, lugar); err != nil {
				beego.Error("Error", err.Error())
			}
		} else {
			idLugar := int(lugar.(map[string]interface{})["Lugar"].(float64))
			if respuesta, err = InsertarLugar(idLugar, idEnte); err != nil {
				err = errors.New("error al insertar lugar ubicacion ente")
				beego.Error(err.Error())
			}
		}
	}
	return
}

// InsertarSucursales ...
// @Title InsertarSucursales
// @Description InsertarSucursales
// @Param	body		body 	[]models.InformacionSucursales	true		"body for InformacionSucursales  content"
// @Success 201
// @Failure 403 body is empty
// @router DesvincularSucursales/ [post]
func (c *GestionSucursalesController) DesvincularSucursales() {
	defer c.ServeJSON()
	var info_sucursal []interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &info_sucursal); err == nil {
		eliminados := optimize.ProccDigest(info_sucursal, removeSucursal, nil, 3)
		beego.Info(eliminados)
		c.Data["json"] = models.Alert{Type: "success", Code: "S_GS001"}
		for _, respuesta := range eliminados {
			if respuesta.(string) != "OK" {
				c.Data["json"] = models.Alert{Type: "error", Code: "E_GS007"}
				break
			}
		}
	}
}

func removeSucursal(rpintfc interface{}, params ...interface{}) (res interface{}) {
	idSucursal := strconv.Itoa(int(rpintfc.(map[string]interface{})["Id"].(float64)))
	urlEliminar := beego.AppConfig.String("coreOrganizacionService") + "relacion_organizaciones/" + idSucursal
	request.SendJson(urlEliminar, "DELETE", &res, nil)
	return
}

// InsertarSucursales ...
// @Title InsertarSucursales
// @Description InsertarSucursales
// @Param	body		body 	[]models.InformacionSucursales	true		"body for InformacionSucursales  content"
// @Success 201
// @Failure 403 body is empty
// @router VincularSucursales/ [post]
func (c *GestionSucursalesController) VincularSucursales() {
	defer c.ServeJSON()
	var info_sucursal []interface{}
	var params []interface{}
	var tipoRelacion []interface{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &info_sucursal); err == nil {
		if err := request.GetJson(beego.AppConfig.String("coreOrganizacionService")+"tipo_relacion_organizaciones/?query=CodigoAbreviacion:TRO_1&&limit=-1", &tipoRelacion); err == nil {
			params = append(params, tipoRelacion[0])
		} else {
			c.Data["json"] = models.Alert{Type: "error", Code: "E_GS008", Body: err.Error()}
			return
		}
		agregados := optimize.ProccDigest(info_sucursal, addSucursal, params, 3)
		beego.Info(agregados)
		c.Data["json"] = models.Alert{Type: "success", Code: "S_GS002", Body: agregados}
		for _, respuesta := range agregados {
			if respuesta.(string) != "OK" {
				c.Data["json"] = models.Alert{Type: "error", Code: "E_GS008"}
				break
			}
		}
	}
}

func addSucursal(rpintfc interface{}, params ...interface{}) (res interface{}) {
	rpintfc.(map[string]interface{})["TipoRelacionOrganizaciones"] = params[0]
	urlPost := beego.AppConfig.String("coreOrganizacionService") + "relacion_organizaciones"
	request.SendJson(urlPost, "POST", &res, rpintfc)
	return
}

// ListarBancos ...
// @Title ListarSucurListarBancossales
// @Description ListarBancos
// @Success 201 {object} []models.InformacionSucursal
// @Failure 403 body is empty
// @router /ListarBancos/ [get]
func (c *GestionSucursalesController) ListarBancos() {
	defer c.ServeJSON()
	var bancos []interface{}
	if err := request.GetJson(beego.AppConfig.String("coreOrganizacionService")+"organizacion/?query=TipoOrganizacion.CodigoAbreviacion:EB&limit=-1", &bancos); err == nil {
		if bancos != nil {
			informacion_bancos := optimize.ProccDigest(bancos, getInfoAdicionalBanco, nil, 3)
			c.Data["json"] = informacion_bancos
		}
	} else {
		c.Data["json"] = err
	}

}

func getInfoAdicionalBanco(rpintfc interface{}, params ...interface{}) (res interface{}) {
	var infoAdicional []map[string]interface{}
	var infoIden []interface{}
	idBanco := strconv.Itoa(int(rpintfc.(map[string]interface{})["Id"].(float64)))
	idEnteBanco := strconv.Itoa(int(rpintfc.(map[string]interface{})["Ente"].(float64)))
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/informacion_adicional_banco/?limit=-1&query=Banco:"+idBanco, &infoAdicional); err == nil {
		if infoAdicional != nil {
			rpintfc.(map[string]interface{})["CodigoAch"] = infoAdicional[0]["CodigoAch"]
			rpintfc.(map[string]interface{})["CodigoSuperintendencia"] = infoAdicional[0]["CodigoAch"]
			rpintfc.(map[string]interface{})["IdInformacionAdicional"] = infoAdicional[0]["CodigoAch"]
		} else {
			rpintfc.(map[string]interface{})["CodigoAch"] = 0
			rpintfc.(map[string]interface{})["CodigoSuperintendencia"] = 0
			rpintfc.(map[string]interface{})["IdInformacionAdicional"] = 0
		}
	} else {
		rpintfc.(map[string]interface{})["CodigoAch"] = 0
		rpintfc.(map[string]interface{})["CodigoSuperintendencia"] = 0
		rpintfc.(map[string]interface{})["IdInformacionAdicional"] = 0
	}

	if err := request.GetJson(beego.AppConfig.String("coreEnteService")+"identificacion/?query=Ente:"+idEnteBanco, &infoIden); err == nil {
		if infoIden != nil {
			rpintfc.(map[string]interface{})["Identificacion"] = infoIden
		}
	} else {
		beego.Error(err.Error())
	}
	return rpintfc
}
