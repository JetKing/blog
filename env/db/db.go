package db

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
)

// register database
// extra args for mysql `user, passwd, host, port`
func RegisterDB(alias string, dbtype string, dbname string, args ...string) error {
	length := len(args)

	if dbtype == "mysql" {
		user := args[0]
		passwd := args[1]
		host := args[2]
		port := args[3]

		if length == 4 {
			return orm.RegisterDataBase(alias, "mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", user, passwd, host, port, dbname))
		} else {
			return errors.New("illeage args")
		}
	} else if dbtype == "sqlite3" {
		if length == 0 {
			return orm.RegisterDataBase(alias, "sqlite3", dbname+".db")
		} else {
			return errors.New("illeage args")
		}
	} else {
		return errors.New("database type: " + dbtype + ", not support")
	}
}
