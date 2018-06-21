package mongoProcess

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/manucorporat/try"
	"github.com/udistrital/api_mid_financiera/controllers"
	"github.com/udistrital/api_mid_financiera/models"
	"github.com/udistrital/utils_oas/optimize"
)

func sendRpInfoToMongo(ctx *context.Context) {
	try.This(func() {
		serviceResponse := ctx.Input.Data()["json"].([]models.Alert)
		var params []interface{}
		for _, data := range serviceResponse {
			if data.Type == "success" {
				info := data.Body
				params = append(params, info)
				work := optimize.WorkRequest{JobParameter: params, Job: (controllers.AddRpMongo)}
				// Push the work onto the queue.
				optimize.WorkQueue <- work
				beego.Info("Job Queued!")
			}

		}
	}).Catch(func(e try.E) {
		beego.Info("Exepc ", e)
	})

}

func RpMongoJobInit() { //inicia los escuchadores de los procesos que deben guardarse simultaneamente en postgres y mongo
	optimize.StartDispatcher(1, 200)

	beego.InsertFilter("/v1/registro_presupuestal/CargueMasivoPr", beego.AfterExec, sendRpInfoToMongo, false)
}
