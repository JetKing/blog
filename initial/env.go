package initial

import (
	"github.com/astaxie/beego"
)

func InitEnv() {
	beego.AppName = "goblog"
	beego.HttpPort = 8081
	beego.CopyRequestBody = true
	beego.TemplateLeft = "{{{"
	beego.TemplateRight = "}}}"
	beego.SessionOn = true
	beego.SessionProvider = "file"
	beego.SessionSavePath = "./tmp"
	beego.SessionGCMaxLifetime = 31536000
	beego.SessionCookieLifeTime = 31536000
	beego.RunMode = "dev"

	runmode := beego.RunMode
	if runmode == "dev" {
		beego.SetStaticPath("/static", "static")
	}
}
