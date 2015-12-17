package models

import (
	"github.com/astaxie/beego/orm"
)

var o orm.Ormer

func InitModels(alias string) {
	registerModel()
	o = orm.NewOrm()
	o.Using(alias)
}

func ClearModels() {
	orm.ResetModelCache()
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
