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
			Router: `Solicitudes/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:OrdenPagoNominaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:OrdenPagoNominaController"],
		beego.ControllerComments{
			Method: "CrearOPSeguridadSocial",
			Router: `CrearOPSeguridadSocial`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:OrdenPagoNominaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/api_mid_financiera/controllers:OrdenPagoNominaController"],
		beego.ControllerComments{
			Method: "MidCrearOPNomina",
			Router: `MidCrearOPNomina`,
			AllowHTTPMethods: []string{"post"},
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
			Method: "ListaRp",
			Router: `ListaRp/:vigencia`,
			AllowHTTPMethods: []string{"get"},
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
