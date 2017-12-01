package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/astaxie/beego"
)

/*
type Field struct {

}

type MetaData struct {
	TotalProperty string
	Root string
	Id string
	Fields *[]Field
	CacheDate string
}

type DataSet struct {
	MetaData *MetaData
	Results  int
	Rows    *[]Row
}
*/

type ReportesController struct {
	beego.Controller
}

func (c *ReportesController) URLMapping() {
	c.Mapping("GetDataSetFinanciera", c.GetDataSetFinanciera)
}

// GetDataSetFinanciera ...
// @Title GetDataSetFinanciera
// @Description Obtiene el data set de los reportes Financiera de SpagoBi
// @Success 201 {object} string
// @Failure 403 body is empty
// @router GetDataSetFinanciera [get]
func (c *ReportesController) GetDataSetFinanciera() {

	// inicio logueo spagobi
	resp, err := http.Get("https://intelligentia.udistrital.edu.co:8443/SpagoBI/servlet/AdapterHTTP?Page=LoginPage&NEW_SESSION=TRUE&userId=biconsulta&password=biconsulta")

	if err != nil {
		fmt.Printf("Error en resp")
	}

	defer resp.Body.Close()

	//fmt.Println(resp.Header["Set-Cookie"])
	cookie := resp.Header["Set-Cookie"][0]
	//fmt.Println(cookie)

	// inicio solicitud dataset

	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://intelligentia.udistrital.edu.co:8443/SpagoBI/restful-services/1.0/datasets/documentos_financiera_ds/data", nil)
	if err != nil {
		fmt.Printf("Error en req")
	}

	//req.Header.Set("Cookie", cookie)
	req.Header.Add("Cookie", cookie)
	resp2, err := client.Do(req)

	if err != nil {
		fmt.Printf("Error en resp2")
	}

	defer resp2.Body.Close()

	body, err := ioutil.ReadAll(resp2.Body)
	bodyString := string(body[:])
	b := []byte(bodyString)

	var f interface{}
	mi := json.Unmarshal(b, &f)

	if mi != nil {
		fmt.Printf("Hay errores")
	}

	m := f.(map[string]interface{})

	c.Data["json"] = m
	c.ServeJSON()
}
