package models

import (
	"errors"
	"github.com/astaxie/beego/orm"
	"github.com/gogather/com"
	// "log"
	"regexp"
	"time"
)

type Users struct {
	Id       int
	Username string
	Password string
	Salt     string
	Email    string
}

type Varify struct {
	Id       int
	Username string
	Code     string
	Overdue  time.Time
}

// 添加用户
func AddUser(username string, password string) (int64, error) {
	user := new(Users)
	user.Username = username
	user.Salt = com.RandString(10)
	user.Password = com.Md5(password + user.Salt)
	return o.Insert(user)
}

// 通过用户名查找用户
func FindUser(username string) (Users, error) {
	user := Users{Username: username}
	err := o.Read(&user, "username")

	return user, err
}

// 修改用户名
func ChangeUsername(oldUsername string, newUsername string) error {
	_, err := o.QueryTable("users").Filter("username", oldUsername).Update(orm.Params{
		"username": newUsername,
	})
	return err
}

// 修改邮箱
func ChangeEmail(username string, email string) error {
	reg := regexp.MustCompile(`^(\w)+(\.\w+)*@(\w)+((\.\w+)+)$`)
	result := reg.MatchString(email)
	if !result {
		return errors.New("not a email")
	}

	num, err := o.QueryTable("users").Filter("username", username).Update(orm.Params{"email": email})

	if nil != err {
		return err
	} else if 0 == num {
		return errors.New("not update")
	} else {
		return nil
	}
}

// 增加验证码
// insert into varify (`username`, `code`, `overdue`) value ('lijun', 'wasdfgert', '2014-08-23 14:51:14')
func AddVerify(username string, code string, overdue time.Time) error {
	// 1小时后过期
	overdueTime := overdue.Add(1 * time.Hour).Format("2006-01-02 15:04:05")
	_, err := o.Raw("insert into varify (`username`, `code`, `overdue`) value ('" + username + "', '" + code + "', '" + overdueTime + "')").Exec()
	return err
}

// 检查验证码
func CheckVarify(code string) (bool, string, error) {
	var varifyItem Varify
	err := o.Raw("select * from varify where code='" + code + "' and overdue > now()").QueryRow(&varifyItem)

	if code == varifyItem.Code {
		o.Raw("delete from varify where code='" + code + "'").Exec()
		return true, varifyItem.Username, err
	} else {
		return false, varifyItem.Username, err
	}
}

// 设置密码
func SetPassword(username string, password string) error {
	salt := com.RandString(10)

	num, err := o.QueryTable("users").Filter("username", username).Update(orm.Params{
		"salt":     salt,
		"password": com.Md5(password + salt),
	})
	if 0 == num {
		return errors.New("item not exist")
	}

	return err
}

// 修改密码
func ChangePassword(username string, oldPassword string, newPassword string) error {
	salt := com.RandString(10)

	user := Users{Username: username}
	err := o.Read(&user, "username")
	if nil != err {
		return err
	} else {
		if user.Password == com.Md5(oldPassword+user.Salt) {
			_, err := o.QueryTable("users").Filter("username", username).Update(orm.Params{
				"salt":     salt,
				"password": com.Md5(newPassword + salt),
			})
			return err
		} else {
			return errors.New("verification failed")
		}
	}
}
