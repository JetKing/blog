package models

import (
	"errors"
	"github.com/astaxie/beego/orm"
	"github.com/gogather/com"
	"regexp"
)

type Users struct {
	Id       int    `orm:"pk; auto;"`
	Username string `orm:"varchar(255); not null"`
	Password string `orm:"varchar(255); not null"`
	Salt     string `orm:"varchar(255); not null"`
	Email    string `orm:"varchar(255);"`
}

var TheUsers Users

// 添加用户
func (this *Users) AddUser(username string, password string) (int64, error) {
	user := new(Users)
	user.Username = username
	user.Salt = com.RandString(10)
	user.Password = com.Md5(password + user.Salt)
	return o.Insert(user)
}

// 通过用户名查找用户
func (this *Users) FindUser(username string) (Users, error) {
	user := Users{Username: username}
	err := o.Read(&user, "username")

	return user, err
}

// 修改用户名
func (this *Users) ChangeUsername(oldUsername string, newUsername string) error {
	_, err := o.QueryTable("users").Filter("username", oldUsername).Update(orm.Params{
		"username": newUsername,
	})
	return err
}

// 修改邮箱
func (this *Users) ChangeEmail(username string, email string) error {
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

// 设置密码
func (this *Users) SetPassword(username string, password string) error {
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
func (this *Users) ChangePassword(username string, oldPassword string, newPassword string) error {
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
