package routers

import (
	"github.com/duguying/blog/routers"
)

func InitRouter() {
	routers.InitApiRouter()
	routers.InitAdminRouter()
	routers.InitDefaultRouter()
}
