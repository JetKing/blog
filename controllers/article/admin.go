package article

import (
	"fmt"
	// "github.com/astaxie/beego"
	// "encoding/json"
	"github.com/duguying/blog/controllers"
	"github.com/duguying/blog/models"
	// "github.com/duguying/blog/utils"
	"github.com/duguying/blog/utils"
	"github.com/gogather/com"
	"github.com/gogather/com/log"
	"strconv"
)

// 管理
type AdminArticleController struct {
	controllers.BaseController
}

func (this *AdminArticleController) ListArticle() {
	s := this.Ctx.Input.Param(":page")
	page, err := strconv.Atoi(s)
	if nil != err || page < 0 {
		page = 1
	}

	maps, nextPage, pages, err := models.TheArticle.ArticleListForAdmin(int(page), 10)
	if nil != err {
		this.Data["json"] = map[string]interface{}{"result": false, "msg": "get list failed", "refer": "/"}
		this.ServeJson()
	} else {
		for _, art := range maps {
			art["time"] = utils.GetDate(art["time"].(string))
		}
		this.Data["json"] = map[string]interface{}{
			"result":   true,
			"msg":      "get list success",
			"refer":    "/",
			"pages":    pages,
			"nextPage": nextPage,
			"data":     maps,
			"page":     page,
		}
		this.ServeJson()
	}

}

func (this *AdminArticleController) GetArticle() {
	s := this.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(s)
	if nil != err || id < 0 {
		id = 1
	}
	art, err := models.TheArticle.Get(int64(id))
	if err != nil {
		this.Data["json"] = map[string]interface{}{"result": false, "msg": "get list failed", "refer": "/"}
	} else {
		this.Data["json"] = map[string]interface{}{
			"result": true,
			"msg":    "get article success",
			"data":   art,
		}
	}
	this.ServeJson()
}

func (this *AdminArticleController) AddArticle() {
	paramsBody := string(this.Ctx.Input.RequestBody)
	var params map[string]interface{}
	p, err := com.JsonDecode(paramsBody)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"result": false, "msg": "parse params failed", "refer": "/"}
		this.ServeJson()
		return
	} else {
		params = p.(map[string]interface{})["params"].(map[string]interface{})
	}

	// log.Pinkln(params)

	title := params["title"].(string)
	content := params["content"].(string)
	keywords := params["keywords"].(string)
	abstract := params["abstract"].(string)

	// if not login, permission deny
	user := this.GetSession("username")
	if user == nil {
		this.Data["json"] = map[string]interface{}{"result": false, "msg": "login first please", "refer": nil}
		this.ServeJson()
		return
	}

	if "" == title {
		this.Data["json"] = map[string]interface{}{"result": false, "msg": "title can not be empty", "refer": "/"}
		this.ServeJson()
		return
	}

	username := user.(string)

	id, err := models.TheArticle.AddArticle(title, content, keywords, abstract, username)
	if nil == err {
		this.Data["json"] = map[string]interface{}{
			"result": true,
			"msg":    "success added, id " + fmt.Sprintf("[%d] ", id),
			"data":   id,
			"refer":  nil,
		}
	} else {
		log.Warnln(err)
		this.Data["json"] = map[string]interface{}{"result": false, "msg": "added failed", "refer": nil}
	}
	this.ServeJson()
}

func (this *AdminArticleController) DelArticle() {
	// if not login, permission deny
	user := this.GetSession("username")
	if user == nil {
		this.Data["json"] = map[string]interface{}{"result": false, "msg": "login first please", "refer": nil}
		this.ServeJson()
		return
	}

	paramsBody := string(this.Ctx.Input.RequestBody)
	var params map[string]interface{}
	p, err := com.JsonDecode(paramsBody)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"result": false, "msg": "parse params failed", "refer": "/"}
		this.ServeJson()
		return
	} else {
		params = p.(map[string]interface{})["params"].(map[string]interface{})
	}

	id := int64(params["id"].(float64))
	// title := this.Ctx.Input.Param(":title")

	if id < 0 {
		id = 0
	}

	num, err := models.TheArticle.DeleteArticle(id, "")

	if nil != err {
		log.Fatal(err)
		this.Data["json"] = map[string]interface{}{"result": false, "msg": "delete faild", "refer": nil}
		this.ServeJson()
	} else if 0 == num {
		this.Data["json"] = map[string]interface{}{"result": false, "msg": "articles dose not exist", "refer": nil}
		this.ServeJson()
	} else {
		this.Data["json"] = map[string]interface{}{"result": true, "msg": fmt.Sprintf("[%d]", num) + " articles deleted", "refer": nil}
		this.ServeJson()
	}
}

func (this *AdminArticleController) UpdateArticle() {
	// if not login, permission deny
	user := this.GetSession("username")
	if user == nil {
		this.Data["json"] = map[string]interface{}{"result": false, "msg": "login first please", "refer": nil}
		this.ServeJson()
		return
	}

	paramsBody := string(this.Ctx.Input.RequestBody)
	var params map[string]interface{}
	p, err := com.JsonDecode(paramsBody)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"result": false, "msg": "parse params failed", "refer": "/"}
		this.ServeJson()
		return
	} else {
		params = p.(map[string]interface{})["params"].(map[string]interface{})
	}

	id := int64(params["id"].(float64))
	newTitle := params["title"].(string)
	newContent := params["content"].(string)
	newKeywords := params["keywords"].(string)

	if "" == newTitle {
		this.Data["json"] = map[string]interface{}{"result": false, "msg": "title can not be empty", "refer": "/"}
		this.ServeJson()
	}

	var art models.Article

	if nil == err {
		art, err = models.TheArticle.Get(id)
	} else {
		this.Ctx.WriteString("not found")
		return
	}

	art.Title = newTitle
	art.Content = newContent
	art.Keywords = newKeywords

	err = models.TheArticle.UpdateArticle(id, "", art)

	if nil != err {
		this.Data["json"] = map[string]interface{}{"result": false, "msg": "update failed", "refer": "/"}
		this.ServeJson()
	} else {
		this.Data["json"] = map[string]interface{}{"result": true, "msg": "update success", "refer": "/"}
		this.ServeJson()
	}

}

// 管理- 项目
type AdminProjectController struct {
	controllers.BaseController
}

func (this *AdminProjectController) GetProject() {
	s := this.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(s)
	if nil != err || id < 0 {
		id = 1
	}

	project, err := models.TheProject.GetProject(id, "")
	if err != nil {
		this.Data["json"] = map[string]interface{}{"result": false, "msg": "get failed", "error": err}
	} else {
		this.Data["json"] = map[string]interface{}{"result": true, "msg": "get success", "data": project}
	}
	this.ServeJson()
}

func (this *AdminProjectController) ListProject() {
	s := this.Ctx.Input.Param(":page")
	page, err := strconv.Atoi(s)
	if nil != err || page < 0 {
		page = 1
	}

	maps, nextPageFlag, totalPages, err := models.TheProject.ListProjects(int(page), 10)

	if err != nil {
		this.Data["json"] = map[string]interface{}{
			"result": false,
			"msg":    "get list failed, " + err.Error(),
		}
	} else {
		this.Data["json"] = map[string]interface{}{
			"has_next":    nextPageFlag,
			"total_pages": totalPages,
			"data":        maps,
			"result":      true,
			"msg":         "get list success",
		}
	}
	this.ServeJson()
}
