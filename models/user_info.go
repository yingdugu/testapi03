package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type UserInfo struct {
	Id              int     `orm:"column(id);size(100);"`
	Name            string  `orm:"column(name);size(100);null"`
	Gender          string  `orm:"column(gender);size(10);null"`
	Phone           string  `orm:"column(phone);size(100);null"`
	Place           string  `orm:"column(place);size(100);null"`
	Idimage         string  `orm:"column(idimage);size(255);null"`
	Currfunds       float32 `orm:"column(currfunds);null"`
	Ishunter        int     `orm:"column(ishunter);null"`
	Expectposition  string  `orm:"column(expectposition);size(100);null"`
	IdentifyNum     string  `orm:"column(identifyNum);pk"`
	ReaderReward    float32 `orm:"column(reader_reward);null"`
	TrainsmitReward float32 `orm:"column(trainsmit_reward);null"`
	InterviewReward float32 `orm:"column(interview_reward);null"`
	EntryReward1    float32 `orm:"column(entry_reward1);null"`
	EntryReward2    float32 `orm:"column(entry_reward2);null"`
}

func (t *UserInfo) TableName() string {
	return "user_info"
}

func init() {
	orm.RegisterModel(new(UserInfo))
}

// AddUserInfo insert a new UserInfo into database and returns
// last inserted Id on success.
func AddUserInfo(m *UserInfo) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetUserInfoById retrieves UserInfo by Id. Returns error if
// Id doesn't exist
func GetUserInfoById(identifyNum string) (v *UserInfo, err error) {
	o := orm.NewOrm()
	v = &UserInfo{IdentifyNum: identifyNum}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllUserInfo retrieves all UserInfo matches certain condition. Returns empty list if
// no records exist
func GetAllUserInfo(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(UserInfo))
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

	var l []UserInfo
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

// UpdateUserInfo updates UserInfo by Id and returns error if
// the record to be updated doesn't exist
func UpdateUserInfoByNum(m *UserInfo) (err error) {
	o := orm.NewOrm()
	v := UserInfo{IdentifyNum: m.IdentifyNum}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteUserInfo deletes UserInfo by Id and returns error if
// the record to be deleted doesn't exist
func DeleteUserInfo(identifyNum string) (err error) {
	o := orm.NewOrm()
	v := UserInfo{IdentifyNum: identifyNum}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&UserInfo{IdentifyNum: identifyNum}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
