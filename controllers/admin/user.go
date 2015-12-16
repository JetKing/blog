package admin

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/duguying/blog/controllers"
	"github.com/duguying/blog/models"
	"github.com/duguying/blog/utils"
	"github.com/gogather/com"
	"time"
)

// 注册
type RegistorController struct {
	controllers.BaseController
}

func (this *RegistorController) Get() {
	registorable, err := beego.AppConfig.Bool("registorable")
	if registorable || nil != err {
		this.TplNames = "registor.tpl"
	} else {
		this.Ctx.WriteString("registor closed")
	}
}

func (this *RegistorController) Post() {
	registorable, err := beego.AppConfig.Bool("registorable")
	if nil != err {
		// default registorable is true, do nothing
	} else if !registorable {
		this.Data["json"] = map[string]interface{}{"result": false, "msg": "registorable is false", "refer": "/"}
		this.ServeJson()
		return
	}

	username := this.GetString("username")
	password := this.GetString("password")

	if !utils.CheckUsername(username) {
		this.Data["json"] = map[string]interface{}{"result": false, "msg": "illegal username", "refer": "/"}
		this.ServeJson()
		return
	}

	id, err := models.TheUsers.AddUser(username, password)
	if nil != err {
		this.Data["json"] = map[string]interface{}{"result": false, "msg": "registor failed", "refer": "/"}
	} else {
		this.Data["json"] = map[string]interface{}{"result": true, "msg": fmt.Sprintf("[%d] ", id) + "registor success", "refer": "/"}
	}
	this.ServeJson()
}

// 登录
type LoginController struct {
	controllers.BaseController
}

func (this *LoginController) Get() {
	// if not login, permission deny
	user := this.GetSession("username")
	if user != nil {
		this.Redirect("/admin", 302)
	} else {
		this.TplNames = "login.tpl"
	}
}

func (this *LoginController) Post() {
	username := this.GetString("username")
	password := this.GetString("password")

	if username == "" || password == "" {
		this.Data["json"] = map[string]interface{}{"result": false, "msg": "invalid request", "refer": "/"}
	}

	user, err := models.TheUsers.FindUser(username)

	if err != nil {
		this.Data["json"] = map[string]interface{}{"result": false, "msg": "user does not exist", "refer": "/"}
	} else {
		passwd := com.Md5(password + user.Salt)
		if passwd == user.Password {
			this.SetSession("username", username)
			this.Redirect("/admin", 302)
			return
		} else {
			this.Data["json"] = map[string]interface{}{"result": false, "msg": "login failed ", "refer": "/"}
		}
	}
	this.ServeJson()
}

// 登出
type LogoutController struct {
	controllers.BaseController
}

func (this *LogoutController) Get() {
	this.DelSession("username")
	this.Ctx.WriteString("you have logout")
}

func (this *LogoutController) Post() {
	this.Data["json"] = map[string]interface{}{"result": false, "msg": "invalid request ", "refer": "/"}
	this.ServeJson()
}

// 测试暂用页
type TestController struct {
	controllers.BaseController
}

func (this *TestController) Get() {
	this.Data["username"] = this.GetSession("username")
	this.TplNames = "test.tpl"
}

func (this *TestController) Post() {
	this.Data["username"] = this.GetSession("username")
	this.TplNames = "test.tpl"
}

// 修改用户名
type ChangeUsernameController struct {
	controllers.BaseController
}

func (this *ChangeUsernameController) Get() {
	this.Data["json"] = map[string]interface{}{"result": false, "msg": "invalid request ", "refer": "/"}
	this.ServeJson()
}

func (this *ChangeUsernameController) Post() {
	// if not login, permission deny
	user := this.GetSession("username")
	if user == nil {
		this.Data["json"] = map[string]interface{}{"result": false, "msg": "login first please", "refer": nil}
		this.ServeJson()
		return
	}

	oldUsername := user.(string)
	newUsername := this.GetString("username")

	err := models.TheUsers.ChangeUsername(oldUsername, newUsername)

	if nil != err {
		// log.Println(err)
		this.Data["json"] = map[string]interface{}{"result": false, "msg": "change username failed", "refer": "/"}
		this.ServeJson()
	} else {
		this.SetSession("username", newUsername)
		this.Data["json"] = map[string]interface{}{"result": true, "msg": "change username success", "refer": "/"}
		this.ServeJson()
	}
}

// 修改Email
type SetEmailController struct {
	controllers.BaseController
}

func (this *SetEmailController) Get() {
	this.Data["json"] = map[string]interface{}{"result": false, "msg": "invalid request ", "refer": "/"}
	this.ServeJson()
}

func (this *SetEmailController) Post() {
	user := this.GetSession("username")
	if user == nil {
		this.Data["json"] = map[string]interface{}{"result": false, "msg": "login first please", "refer": nil}
		this.ServeJson()
		return
	}
	username := user.(string)

	email := this.GetString("email")
	if "" == email {
		this.Data["json"] = map[string]interface{}{"result": false, "msg": "email is needed", "refer": nil}
		this.ServeJson()
		return
	}

	err := models.TheUsers.ChangeEmail(username, email)

	if nil != err {
		this.Data["json"] = map[string]interface{}{"result": false, "msg": "set email failed", "refer": nil}
		this.ServeJson()
		return
	} else {
		this.Data["json"] = map[string]interface{}{"result": true, "msg": "set email success", "refer": "/"}
		this.ServeJson()
	}
}

// 找回密码
type GetBackPasswordController struct {
	controllers.BaseController
}

func (this *GetBackPasswordController) Get() {
	this.TplNames = "getbackpasswd.tpl"
}

func (this *GetBackPasswordController) Post() {
	this.Data["json"] = map[string]interface{}{"result": false, "msg": "invalid request", "refer": "/"}
	this.ServeJson()
}

// 发送找回密码验证邮件
type SendEmailToGetBackPasswordController struct {
	controllers.BaseController
}

func (this *SendEmailToGetBackPasswordController) Get() {
	username := this.GetString("username")
	if "" == username {
		this.Data["json"] = map[string]interface{}{"result": false, "msg": "username could not be empty", "refer": "/"}
		this.ServeJson()
		return
	}

	time := time.Now()
	code := com.Md5(com.RandString(20) + time.String())

	err := models.TheVarify.AddVerify(username, code, time)

	if nil != err {
		this.Data["json"] = map[string]interface{}{"result": false, "msg": "create varify failed", "refer": "/"}
		this.ServeJson()
	} else {
		host := beego.AppConfig.String("host")
		subject := "blog system get your password back"
		body := `click the following link to get your password back <font color="red"><a href="` + host + `/password/reset/` + code + `">` + host + `/password/reset/` + code + `</a></font>`
		currentUser, _ := models.TheUsers.FindUser(username)
		email := currentUser.Email

		err := utils.SendMail(email, subject, body)
		if nil != err {
			this.Data["json"] = map[string]interface{}{"result": false, "msg": "send mail failed", "refer": "/"}
			this.ServeJson()
		} else {
			this.Data["json"] = map[string]interface{}{"result": true, "msg": "create varify success", "refer": "/"}
			this.ServeJson()
		}
	}

}

func (this *SendEmailToGetBackPasswordController) Post() {
	this.Data["json"] = map[string]interface{}{"result": false, "msg": "invalid request ", "refer": "/"}
	this.ServeJson()
}

// 设置密码
type SetPasswordController struct {
	controllers.BaseController
}

func (this *SetPasswordController) Get() {
	varify := this.Ctx.Input.Param(":varify")

	if "" == varify {
		this.Data["json"] = map[string]interface{}{"result": false, "msg": "invalid request ", "refer": "/"}
		this.ServeJson()
	}

	result, username, err := models.TheVarify.CheckVarify(varify)

	if nil != err {
		this.Ctx.WriteString("找回密码已过期")
	} else if !result {
		this.Ctx.WriteString("验证错误")
	} else {
		this.Data["username"] = username
		this.SetSession("username", username)
		this.SetSession("reset", true)
		this.TplNames = "resetpasswd.tpl"
	}
}

func (this *SetPasswordController) Post() {
	user := this.GetSession("username")
	reset := this.GetSession("reset")
	resetable := reset.(bool)
	if !resetable {
		this.Data["json"] = map[string]interface{}{"result": false, "msg": "resetable is false", "refer": nil}
		this.ServeJson()
		return
	}
	if user == nil {
		this.Data["json"] = map[string]interface{}{"result": false, "msg": "session failed", "refer": nil}
		this.ServeJson()
		return
	}
	username := user.(string)
	newPassword := this.GetString("password")
	// fmt.Println(username)
	if "" == newPassword {
		this.Data["json"] = map[string]interface{}{"result": false, "msg": "password is needed", "refer": nil}
		this.ServeJson()
		return
	}

	err := models.TheUsers.SetPassword(username, newPassword)
	if nil != err {
		this.Data["json"] = map[string]interface{}{"result": false, "msg": "set password failed", "refer": nil}
		this.ServeJson()
		return
	} else {
		this.DelSession("reset")
		this.Data["json"] = map[string]interface{}{"result": true, "msg": "set password success", "refer": "/"}
		this.ServeJson()
	}

}

// 修改密码
type ChangePasswordController struct {
	controllers.BaseController
}

func (this *ChangePasswordController) Get() {
	this.Data["json"] = map[string]interface{}{"result": false, "msg": "invalid request ", "refer": "/"}
	this.ServeJson()
}

func (this *ChangePasswordController) Post() {
	// if not login, permission deny
	user := this.GetSession("username")
	if user == nil {
		this.Data["json"] = map[string]interface{}{"result": false, "msg": "login first please", "refer": nil}
		this.ServeJson()
		return
	}

	username := user.(string)
	oldPassword := this.GetString("old_password")
	newPassword := this.GetString("password")

	err := models.TheUsers.ChangePassword(username, oldPassword, newPassword)

	if nil != err {
		this.Data["json"] = map[string]interface{}{"result": false, "msg": "change password faild", "refer": nil}
		this.ServeJson()
	} else {
		this.Data["json"] = map[string]interface{}{"result": true, "msg": "change password success", "refer": nil}
		this.ServeJson()
	}
}
