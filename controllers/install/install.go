package install

import (
	"github.com/astaxie/beego"
)

// 系统安装
type InstallController struct {
	beego.Controller
}

func (this *InstallController) Welcome() {
	this.TplNames = "admin/install.tpl"
}
