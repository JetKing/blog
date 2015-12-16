package initial

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/duguying/blog/env"
	"github.com/duguying/blog/env/db"
	"github.com/duguying/blog/models"
	"github.com/duguying/blog/routers/install"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gogather/com/log"
)

func InitSql() {

	orm.Debug = true

	err := InitialSqlFromConfig()

	if err != nil {
		env.EnvSet("install_mode", true)
		log.Pinkf("[install mode]\n")

		install.InitInstallRouter()

	} else {
		env.EnvSet("install_mode", false)
		env.EnvSet("blog_db", "default")
		log.Bluef("[service mode]\n")
		InitRouter()
		models.InitModels()
	}
}

func InitialSqlFromConfig() error {
	user := beego.AppConfig.String("mysqluser")
	passwd := beego.AppConfig.String("mysqlpass")
	host := beego.AppConfig.String("mysqlurls")
	port := beego.AppConfig.String("mysqlport")
	dbname := "kjhk" //beego.AppConfig.String("mysqldb")

	err := db.RegisterDB("blog", "mysql", dbname, user, passwd, host, port)

	return err
}
