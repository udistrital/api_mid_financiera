package models

import (
	"errors"
	"strconv"
	"fmt"
	"reflect"
	"strings"
	"github.com/udistrital/utils_oas/request"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
)

type Homologacion_rubro struct {
	Id               int64  `orm:"auto"`
	CodigoHomologado string `orm:"size(128)"`
	NombreHomologado string `orm:"size(128)"`
	Organizacion     int
	Vigencia         float64
}

func init() {
	orm.RegisterModel(new(Homologacion_rubro))
}

// AddHomologacion_rubro insert a new Homologacion_rubro into database and returns
// last inserted Id on success.
func AddHomologacion_rubro(m *Homologacion_rubro) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetHomologacion_rubroById retrieves Homologacion_rubro by Id. Returns error if
// Id doesn't exist
func GetHomologacion_rubroById(id int64) (v *Homologacion_rubro, err error) {
	o := orm.NewOrm()
	v = &Homologacion_rubro{Id: id}
	if err = o.QueryTable(new(Homologacion_rubro)).Filter("Id", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllHomologacion_rubro retrieves all Homologacion_rubro matches certain condition. Returns empty list if
// no records exist
func GetAllHomologacion_rubro(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Homologacion_rubro))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Homologacion_rubro
	qs = qs.OrderBy(sortFields...).RelatedSel()
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateHomologacion_rubro updates Homologacion_rubro by Id and returns error if
// the record to be updated doesn't exist
func UpdateHomologacion_rubroById(m *Homologacion_rubro) (err error) {
	o := orm.NewOrm()
	v := Homologacion_rubro{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteHomologacion_rubro deletes Homologacion_rubro by Id and returns error if
// the record to be deleted doesn't exist
func DeleteHomologacion_rubro(id int64) (err error) {
	o := orm.NewOrm()
	v := Homologacion_rubro{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Homologacion_rubro{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
//getValue from organizacion for all rows
func GetOrganizacionRubroHomologado(rubroHomol interface{},params ...interface{})(res interface{}){
	var organizacion interface{}
	rubroHomolMap:= rubroHomol.(map[string]interface{})
	idOrganizacion := strconv.FormatFloat(rubroHomolMap["Organizacion"].(float64),'f',-1,64)
		if err := request.GetJson(beego.AppConfig.String("coreOrganizacionService")+"organizacion?limit=-1&query=Id:"+idOrganizacion, &organizacion); err == nil {
			rubroHomolMap["Organizacion"]= organizacion
			res = rubroHomolMap
		}
	return
}
