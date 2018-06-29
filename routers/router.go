// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/udistrital/api_mid_financiera/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/aprobacion_apropiacion",
			beego.NSInclude(
				&controllers.AprobacionController{},
			),
		),
		beego.NSNamespace("/disponibilidad",
			beego.NSInclude(
				&controllers.DisponibilidadController{},
			),
		),
		beego.NSNamespace("/registro_presupuestal",
			beego.NSInclude(
				&controllers.RegistroPresupuestalController{},
			),
		),
		beego.NSNamespace("/partida_doble",
			beego.NSInclude(
				&controllers.PartidadobleController{},
			),
		),
		beego.NSNamespace("/orden_pago_nomina",
			beego.NSInclude(
				&controllers.OrdenPagoNominaController{},
			),
		),
		beego.NSNamespace("/orden_pago_ss",
			beego.NSInclude(
				&controllers.OrdenPagoSsController{},
			),
		),
		beego.NSNamespace("/orden_pago",
			beego.NSInclude(
				&controllers.OrdenPagoController{},
			),
		),
		beego.NSNamespace("/rubro",
			beego.NSInclude(
				&controllers.RubroController{},
			),
		),
		beego.NSNamespace("/homologacion",
			beego.NSInclude(
				&controllers.HomologacionController{},
			),
		),
		beego.NSNamespace("/aprobacion_fuente",
			beego.NSInclude(
				&controllers.AprobacionFuenteController{},
			),
		),
		beego.NSNamespace("/devoluciones",
			beego.NSInclude(
				&controllers.DevolucionesController{},
			),
		),
		beego.NSNamespace("/ingreso_sin_situacion_fondos",
			beego.NSInclude(
				&controllers.IngresoSinSituacionFondosController{},
			),
		),
<<<<<<< HEAD
		beego.NSNamespace("/organizaciones_core_new",
			beego.NSInclude(
				&controllers.OrganizacionesCoreNewController{},
			),
		),
		beego.NSNamespace("/rubro_homologado",
			beego.NSInclude(
				&controllers.HomologacionRubroController{},
=======

		beego.NSNamespace("/gestion_sucursales",
			beego.NSInclude(
				&controllers.GestionSucursalesController{},
>>>>>>> 5f5cb4796dceb12d29b6fe7ced46c8bc6535e799
			),
		),
	)
	beego.AddNamespace(ns)
}
