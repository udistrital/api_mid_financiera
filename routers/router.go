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
		beego.NSNamespace("/Orden_pago_nomina",
			beego.NSInclude(
				&controllers.OrdenPagoNominaController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
