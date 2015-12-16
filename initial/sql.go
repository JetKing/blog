package initial

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gogather/com/log"
)

func InitSql() {

	RegisterDB("test1", "test2", "test3")

	user := beego.AppConfig.String("mysqluser")
	passwd := beego.AppConfig.String("mysqlpass")
	host := beego.AppConfig.String("mysqlurls")
	port := beego.AppConfig.String("mysqlport")
	dbname := beego.AppConfig.String("mysqldb")

	orm.Debug = true

	err := RegisterDB("mysql", dbname, user, passwd, host, port)

	if err != nil {
		EnvSet("install_mode", true)
		log.Pinkf("[install mode]\n")
		err = RegisterDB("mysql", "", user, passwd, host, port)
	}
}

// register database
// extra args for mysql `user, passwd, host, port`
func RegisterDB(dbtype string, dbname string, args ...string) error {
	length := len(args)

	if dbtype == "mysql" {
		user := args[0]
		passwd := args[1]
		host := args[2]
		port := args[3]

		if length == 4 {
			return orm.RegisterDataBase("default", "mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", user, passwd, host, port, dbname))
		} else {
			return errors.New("illeage args")
		}
	} else if dbtype == "sqlite3" {
		if length == 0 {
			return orm.RegisterDataBase("default", "sqlite3", dbname+".db")
		} else {
			return errors.New("illeage args")
		}
	} else {
		return errors.New("database type: " + dbtype + ", not support")
	}
}
