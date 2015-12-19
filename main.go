package main

import (
	"github.com/astaxie/beego"
	_ "github.com/duguying/blog/initial"
	_ "github.com/duguying/blog/modules/config"
)

func main() {

	beego.Run()
}
