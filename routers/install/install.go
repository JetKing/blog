package install

import (
	"github.com/astaxie/beego"
	"github.com/duguying/blog/controllers/index"
	"github.com/duguying/blog/controllers/install"
)

func InitInstallRouter() {
	beego.Router("/install", &install.InstallController{}, "*:Welcome")
	beego.Router("/install/start", &install.InstallController{}, "*:StartInstall")
	beego.Router("/", &index.MainController{})
}
