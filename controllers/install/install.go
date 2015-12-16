package install

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/duguying/blog/env"
	"github.com/duguying/blog/env/db"
	"github.com/duguying/blog/models"
	"github.com/duguying/blog/routers"
	"github.com/gogather/com/log"
)

// 系统安装
type InstallController struct {
	beego.Controller
}

func (this *InstallController) prepare() {
	if env.EnvGet("install_mode") == false {
		this.Redirect("/", 302)
	}
}

func (this *InstallController) Index() {
	this.Redirect("/install", 302)
}

func (this *InstallController) Welcome() {
	this.prepare()

	this.TplNames = "admin/install.tpl"
}

func (this *InstallController) StartInstall() {
	this.prepare()

	dbname := this.GetString("dbname", "blog")
	user := this.GetString("user", "root")
	passwd := this.GetString("passwd", "")
	host := this.GetString("host", "127.0.0.1")
	port := this.GetString("port", "3306")

	// name := "install"
	force := false
	verbose := true

	err := db.RegisterDB("default", "mysql", "", user, passwd, host, port)
	models.InitModels()

	err = this.createMysqlDB(dbname)
	if err == nil {
		log.Pinkln("创建数据库")
	} else {
		log.Redln("创建数据库失败")
	}

	err = db.RegisterDB("install", "mysql", dbname, user, passwd, host, port)
	if err == nil {
		log.Pinkln("注册新数据库成功")
	}

	if err != nil {
		fmt.Println(err)
	} else {

		err = orm.RunSyncdb("install", force, verbose)
		if err != nil {
			fmt.Println(err)
		} else {
			env.EnvSet("blog_db", "install")
		}
	}

	if err == nil {
		env.EnvSet("install_mode", false)
		routers.InitRouter()
		this.Ctx.WriteString("安装成功")
	}

}

// create mysql database
func (this *InstallController) createMysqlDB(dbname string) error {
	o := orm.NewOrm()
	o.Using("default")
	p, err := o.Raw("CREATE DATABASE IF NOT EXISTS `" + dbname + "` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci").Prepare()
	_, err = p.Exec()
	return err
}
