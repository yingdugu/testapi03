package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type ProjectImage struct {
	Id        string `orm:"column(id);size(255);null"`
	Name      string `orm:"column(name);pk"`
	Imagepath string `orm:"column(imagepath);size(255);null"`
	ImageUrl  string `orm:"column(image_url);size(255);null"`
}

func (t *ProjectImage) TableName() string {
	return "project_image"
}

func init() {
	orm.RegisterModel(new(ProjectImage))
}

// AddProjectImage insert a new ProjectImage into database and returns
// last inserted Id on success.
func AddProjectImage(m *ProjectImage) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetProjectImageById retrieves ProjectImage by Id. Returns error if
// Id doesn't exist
func GetProjectImageByName(name string) (v *ProjectImage, err error) {
	o := orm.NewOrm()
	v = &ProjectImage{Name: name}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllProjectImage retrieves all ProjectImage matches certain condition. Returns empty list if
// no records exist
func GetAllProjectImage(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(ProjectImage))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
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

	var l []ProjectImage
	qs = qs.OrderBy(sortFields...)
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

// UpdateProjectImage updates ProjectImage by Id and returns error if
// the record to be updated doesn't exist
func UpdateProjectImageById(m *ProjectImage) (err error) {
	o := orm.NewOrm()
	v := ProjectImage{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteProjectImage deletes ProjectImage by Id and returns error if
// the record to be deleted doesn't exist
func DeleteProjectImage(id string) (err error) {
	o := orm.NewOrm()
	v := ProjectImage{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&ProjectImage{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
