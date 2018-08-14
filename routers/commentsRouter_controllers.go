package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:AprobacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:AprobacionController"],
		beego.ControllerComments{
			Method: "AprobacionAsignacionInicial",
			Router: `/AprobacionAsignacionInicial/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:AprobacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:AprobacionController"],
		beego.ControllerComments{
			Method: "InformacionAsignacionInicial",
			Router: `/InformacionAsignacionInicial/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:AprobacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:AprobacionController"],
		beego.ControllerComments{
			Method: "Aprobar",
			Router: `Aprobar/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:AprobacionFuenteController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:AprobacionFuenteController"],
		beego.ControllerComments{
			Method: "ValorMovimientoFuente",
			Router: `/ValorMovimientoFuente`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:AprobacionFuenteController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:AprobacionFuenteController"],
		beego.ControllerComments{
			Method: "ValorMovimientoFuenteLista",
			Router: `/ValorMovimientoFuenteLista`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:AprobacionFuenteController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:AprobacionFuenteController"],
		beego.ControllerComments{
			Method: "ValorMovimientoFuenteListaFunc",
			Router: `/ValorMovimientoFuenteListaFunc`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:AprobacionFuenteController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:AprobacionFuenteController"],
		beego.ControllerComments{
			Method: "ValorMovimientoFuenteTraslado",
			Router: `/ValorMovimientoFuenteTraslado`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:DevolucionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:DevolucionesController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:DevolucionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:DevolucionesController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:DevolucionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:DevolucionesController"],
		beego.ControllerComments{
			Method: "GetOne",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:DevolucionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:DevolucionesController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:DevolucionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:DevolucionesController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:DevolucionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:DevolucionesController"],
		beego.ControllerComments{
			Method: "GetTransformRequest",
			Router: `/GetTransformRequest/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:DisponibilidadController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:DisponibilidadController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:DisponibilidadController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:DisponibilidadController"],
		beego.ControllerComments{
			Method: "AprobarAnulacion",
			Router: `/AprobarAnulacion`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:DisponibilidadController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:DisponibilidadController"],
		beego.ControllerComments{
			Method: "ExpedirDisponibilidad",
			Router: `/ExpedirDisponibilidad`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:DisponibilidadController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:DisponibilidadController"],
		beego.ControllerComments{
			Method: "ValorDisponibilidadesFuenteRubroDependencia",
			Router: `/ValorDisponibilidadesFuenteRubroDependencia`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:DisponibilidadController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:DisponibilidadController"],
		beego.ControllerComments{
			Method: "DisponibilidadByNecesidad",
			Router: `DisponibilidadByNecesidad/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:DisponibilidadController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:DisponibilidadController"],
		beego.ControllerComments{
			Method: "ListaDisponibilidades",
			Router: `ListaDisponibilidades/:vigencia`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:DisponibilidadController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:DisponibilidadController"],
		beego.ControllerComments{
			Method: "InfoSolicitudDisponibilidadById",
			Router: `SolicitudById/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:DisponibilidadController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:DisponibilidadController"],
		beego.ControllerComments{
			Method: "InfoSolicitudDisponibilidad",
			Router: `Solicitudes/:vigencia`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:GestionSucursalesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:GestionSucursalesController"],
		beego.ControllerComments{
			Method: "ListarSucursalesBanco",
			Router: `/ListarSucursalesBanco/:idBanco`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:GestionSucursalesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:GestionSucursalesController"],
		beego.ControllerComments{
			Method: "InsertarSucursales",
			Router: `insertar_sucursal/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:GestionSucursalesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:GestionSucursalesController"],
		beego.ControllerComments{
			Method: "ListarSucursal",
			Router: `listar_sucursal/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:GestionSucursalesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:GestionSucursalesController"],
		beego.ControllerComments{
			Method: "ListarSucursales",
			Router: `listar_sucursales/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:HomologacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:HomologacionController"],
		beego.ControllerComments{
			Method: "MidHomologacionLiquidacion",
			Router: `MidHomologacionLiquidacion`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:HomologacionRubroController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:HomologacionRubroController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:HomologacionRubroController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:HomologacionRubroController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:HomologacionRubroController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:HomologacionRubroController"],
		beego.ControllerComments{
			Method: "GetOne",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:HomologacionRubroController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:HomologacionRubroController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:HomologacionRubroController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:HomologacionRubroController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:HomologacionRubroController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:HomologacionRubroController"],
		beego.ControllerComments{
			Method: "CreateHomologacion",
			Router: `/CreateHomologacion`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:HomologacionRubroController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:HomologacionRubroController"],
		beego.ControllerComments{
			Method: "CreateRubroHomologado",
			Router: `/CreateRubroHomologado`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:HomologacionRubroController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:HomologacionRubroController"],
		beego.ControllerComments{
			Method: "GetAllRubrosHomologado",
			Router: `/GetAllRubrosHomologado/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:HomologacionRubroController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:HomologacionRubroController"],
		beego.ControllerComments{
			Method: "GetArbolRubrosHomologado",
			Router: `/GetArbolRubrosHomologado`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:HomologacionRubroController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:HomologacionRubroController"],
		beego.ControllerComments{
			Method: "GetHomologationNumberEntity",
			Router: `/GetHomologationNumberEntity`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:HomologacionRubroController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:HomologacionRubroController"],
		beego.ControllerComments{
			Method: "GetHomologationNumberRubro",
			Router: `/GetHomologationNumberRubro/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:IngresoSinSituacionFondosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:IngresoSinSituacionFondosController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:IngresoSinSituacionFondosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:IngresoSinSituacionFondosController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:IngresoSinSituacionFondosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:IngresoSinSituacionFondosController"],
		beego.ControllerComments{
			Method: "GetOne",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:IngresoSinSituacionFondosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:IngresoSinSituacionFondosController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:IngresoSinSituacionFondosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:IngresoSinSituacionFondosController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:IngresoSinSituacionFondosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:IngresoSinSituacionFondosController"],
		beego.ControllerComments{
			Method: "ChangeState",
			Router: `ChangeState/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:InversionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:InversionController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:InversionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:InversionController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:InversionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:InversionController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:InversionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:InversionController"],
		beego.ControllerComments{
			Method: "GetOne",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:InversionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:InversionController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:InversionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:InversionController"],
		beego.ControllerComments{
			Method: "CreateInversion",
			Router: `CreateInversion/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:InversionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:InversionController"],
		beego.ControllerComments{
			Method: "GetAllCancelaciones",
			Router: `GetAllCancelaciones/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:InversionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:InversionController"],
		beego.ControllerComments{
			Method: "GetCancelationQuantity",
			Router: `GetCancelationQuantity/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:OrdenPagoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:OrdenPagoController"],
		beego.ControllerComments{
			Method: "GetOrdenPagoByFuenteFinanciamiento",
			Router: `/GetOrdenPagoByFuenteFinanciamiento`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:OrdenPagoNominaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:OrdenPagoNominaController"],
		beego.ControllerComments{
			Method: "ListaConceptosNominaHomologados",
			Router: `/ListaConceptosNominaHomologados`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:OrdenPagoNominaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:OrdenPagoNominaController"],
		beego.ControllerComments{
			Method: "ListaLiquidacionNominaHomologada",
			Router: `/ListaLiquidacionNominaHomologada`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:OrdenPagoNominaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:OrdenPagoNominaController"],
		beego.ControllerComments{
			Method: "PreviewCargueMasivoOp",
			Router: `/PreviewCargueMasivoOp`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:OrdenPagoNominaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:OrdenPagoNominaController"],
		beego.ControllerComments{
			Method: "RegistroCargueMasivoOp",
			Router: `/RegistroCargueMasivoOp`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:OrdenPagoSsController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:OrdenPagoSsController"],
		beego.ControllerComments{
			Method: "GetConceptosMovimeintosContablesSs",
			Router: `/GetConceptosMovimeintosContablesSs`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:OrdenPagoSsController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:OrdenPagoSsController"],
		beego.ControllerComments{
			Method: "TestJota01",
			Router: `/TestJota01`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:OrganizacionesCoreNewController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:OrganizacionesCoreNewController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:OrganizacionesCoreNewController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:OrganizacionesCoreNewController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:OrganizacionesCoreNewController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:OrganizacionesCoreNewController"],
		beego.ControllerComments{
			Method: "GetOne",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:OrganizacionesCoreNewController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:OrganizacionesCoreNewController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:OrganizacionesCoreNewController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:OrganizacionesCoreNewController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:OrganizacionesCoreNewController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:OrganizacionesCoreNewController"],
		beego.ControllerComments{
			Method: "GetOrganizacion",
			Router: `getOrganizacion/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:PartidadobleController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:PartidadobleController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:PartidadobleController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:PartidadobleController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:PartidadobleController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:PartidadobleController"],
		beego.ControllerComments{
			Method: "GetOne",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:PartidadobleController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:PartidadobleController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:PartidadobleController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:PartidadobleController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:RegistroPresupuestalController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:RegistroPresupuestalController"],
		beego.ControllerComments{
			Method: "CargueMasivoPr",
			Router: `/CargueMasivoPr`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:RegistroPresupuestalController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:RegistroPresupuestalController"],
		beego.ControllerComments{
			Method: "GetSolicitudesRp",
			Router: `/GetSolicitudesRp/:vigencia`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:RegistroPresupuestalController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:RegistroPresupuestalController"],
		beego.ControllerComments{
			Method: "GetSolicitudesRpById",
			Router: `/GetSolicitudesRpById/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:RegistroPresupuestalController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:RegistroPresupuestalController"],
		beego.ControllerComments{
			Method: "ListaNecesidadesByRp",
			Router: `/ListaNecesidadesByRp/:vigencia`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:RegistroPresupuestalController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:RegistroPresupuestalController"],
		beego.ControllerComments{
			Method: "SolicitudesRpByDependencia",
			Router: `/SolicitudesRpByDependencia/:vigencia`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:RegistroPresupuestalController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:RegistroPresupuestalController"],
		beego.ControllerComments{
			Method: "ListaRp",
			Router: `ListaRp/:vigencia`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:ReintegroController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:ReintegroController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:ReintegroController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:ReintegroController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:ReintegroController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:ReintegroController"],
		beego.ControllerComments{
			Method: "GetOne",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:ReintegroController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:ReintegroController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:id`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:ReintegroController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:ReintegroController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:id`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:ReintegroController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:ReintegroController"],
		beego.ControllerComments{
			Method: "Create",
			Router: `/Create`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:ReintegroController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:ReintegroController"],
		beego.ControllerComments{
			Method: "GetReintegroDisponible",
			Router: `/GetReintegroDisponible`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:RubroController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:RubroController"],
		beego.ControllerComments{
			Method: "GenerarCierre",
			Router: `/GenerarCierre/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:RubroController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:RubroController"],
		beego.ControllerComments{
			Method: "GenerarPac",
			Router: `/GenerarPac/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

}
