package initial

import (
	"github.com/astaxie/beego"
)

var env map[string]interface{}

func init() {
	env = make(map[string]interface{})
}

func InitEnv() {
	runmode := beego.AppConfig.String("runmode")
	if runmode == "dev" {
		beego.SetStaticPath("/static", "static")
	}
}

func EnvSet(key string, value bool) {
	env[key] = value
}

func EnvGet(key string) interface{} {
	value, ok := env[key]
	if !ok {
		return nil
	} else {
		return value
	}
}
