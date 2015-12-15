package models

import (
	"github.com/astaxie/beego/orm"
)

var o orm.Ormer

func init() {
	registerModel()
	o = orm.NewOrm()
	o.Using("blogtest")
}

func registerModel() {
	orm.RegisterModel(new(Article))
	orm.RegisterModel(new(File))
	orm.RegisterModel(new(Project))
	orm.RegisterModel(new(Tags))
	orm.RegisterModel(new(Users))
	orm.RegisterModel(new(UserLog))
	orm.RegisterModel(new(Varify))

	orm.RunCommand()
}
