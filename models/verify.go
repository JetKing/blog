package models

import (
	// "github.com/astaxie/beego/orm"
	"time"
)

type Varify struct {
	Id       int       `orm:"pk; auto;"`
	Username string    `orm:"varchar(50); not null"`
	Code     string    `orm:"varchar(128)"; not null`
	Overdue  time.Time `orm:"timestamp;"`
}

var TheVarify Varify

// 增加验证码
// insert into varify (`username`, `code`, `overdue`) value ('lijun', 'wasdfgert', '2014-08-23 14:51:14')
func (this *Varify) AddVerify(username string, code string, overdue time.Time) error {
	// 1小时后过期
	overdueTime := overdue.Add(1 * time.Hour).Format("2006-01-02 15:04:05")
	_, err := o.Raw("insert into varify (`username`, `code`, `overdue`) value ('" + username + "', '" + code + "', '" + overdueTime + "')").Exec()
	return err
}

// 检查验证码
func (this *Varify) CheckVarify(code string) (bool, string, error) {
	var varifyItem Varify
	err := o.Raw("select * from varify where code='" + code + "' and overdue > now()").QueryRow(&varifyItem)

	if code == varifyItem.Code {
		o.Raw("delete from varify where code='" + code + "'").Exec()
		return true, varifyItem.Username, err
	} else {
		return false, varifyItem.Username, err
	}
}
