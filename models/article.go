package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/duguying/blog/utils"
	// "github.com/gogather/com"
	"strconv"
	"strings"
	"time"
)

type Article struct {
	Id       int64
	Title    string
	Uri      string
	Keywords string
	Abstract string
	Content  string
	Author   string
	Time     time.Time
	Count    int
	Status   int // 0 草稿 1 发布
}

var TheArticle Article

func (this *Article) Get(id int64) (Article, error) {
	var err error
	var art Article

	err = utils.GetCache("GetArticle.id."+fmt.Sprintf("%d", id), &art)
	if err != nil {
		art = Article{Id: id}
		err = o.Read(&art, "id")
		utils.SetCache("GetArticle.id."+fmt.Sprintf("%d", id), art, 600)
	}

	return art, err
}

// 添加文章
func (this *Article) AddArticle(title string, content string, keywords string, abstract string, author string) (int64, error) {
	sql := "insert into article(title, uri, keywords, abstract, content, author, status) values(?, ?, ?, ?, ?, ?, 1)"
	res, err := o.Raw(sql, title, strings.Replace(title, "/", "-", -1), keywords, abstract, content, author).Exec()
	if nil != err {
		return 0, err
	} else {
		return res.LastInsertId()
	}
}

// 添加草稿
func (this *Article) AddDraft(title string, content string, keywords string, abstract string, author string) (int64, error) {
	sql := "insert into article(title, uri, keywords, abstract, content, author, status) values(?, ?, ?, ?, ?, ?, 0)"
	res, err := o.Raw(sql, title, strings.Replace(title, "/", "-", -1), keywords, abstract, content, author).Exec()
	if nil != err {
		return 0, err
	} else {
		return res.LastInsertId()
	}
}

// 由草稿发布
func (this *Article) DraftPublish(id int64) error {
	article, err := this.Get(id)
	if err != nil {
		return err
	}
	article.Status = 1
	_, err = o.Update(&article, "status")
	return err
}

// 通过uri获取文章-cached
func (this *Article) GetArticleByUri(uri string) (Article, error) {
	var err error
	var art Article

	err = utils.GetCache("GetArticleByUri.uri."+uri, &art)
	if err == nil {
		// get view count
		count, err := this.GetArticleViewCount(art.Id)
		if err == nil {
			art.Count = int(count)
		}

		return art, nil
	} else {
		art = Article{Uri: uri}
		err = o.Read(&art, "uri")
		utils.SetCache("GetArticleByUri.uri."+uri, art, 600)
	}

	return art, err
}

// 通过文章标题获取文章-cached
func (this *Article) GetArticleByTitle(title string) (Article, error) {
	var err error
	var art Article

	err = utils.GetCache("GetArticleByTitle.title."+title, &art)
	if err != nil {
		// get view count
		count, err := this.GetArticleViewCount(art.Id)
		if err == nil {
			art.Count = int(count)
		}

		return art, nil
	} else {
		art = Article{Title: title}
		err = o.Read(&art, "title")
		utils.SetCache("GetArticleByTitle.title."+title, art, 600)
	}

	return art, err
}

// 获取文章浏览量
func (this *Article) GetArticleViewCount(id int64) (int, error) {
	var maps []orm.Params

	sql := `select count from article where id=?`
	num, err := o.Raw(sql, id).Values(&maps)
	if err == nil && num > 0 {
		count := maps[0]["count"].(string)

		return strconv.Atoi(count)
	} else {
		return 0, err
	}
}

// 更新阅览数统计
func (this *Article) UpdateCount(id int64) error {
	art := Article{Id: id}
	err := o.Read(&art)

	o.QueryTable("article").Filter("id", id).Update(orm.Params{
		"count": art.Count + 1,
	})

	return err
}

// 更新文章
func (this *Article) UpdateArticle(id int64, uri string, newArt Article) error {
	var art Article

	if 0 != id {
		art = Article{Id: id}
	} else if "" != uri {
		art = Article{Uri: uri}
	}

	art.Title = newArt.Title
	art.Keywords = newArt.Keywords
	art.Abstract = newArt.Abstract
	art.Content = newArt.Content

	getArt, _ := this.Get(id)
	utils.DelCache("GetArticleByUri.uri." + getArt.Uri)
	utils.DelCache("GetArticle.id." + fmt.Sprintf("%d", art.Id))

	_, err := o.Update(&art, "title", "keywords", "abstract", "content")
	return err
}

// 通过uri删除文章
func (this *Article) DeleteArticle(id int64, uri string) (int64, error) {
	var art Article

	if 0 != id {
		art.Id = id
	} else if "" != uri {
		art.Uri = uri
	}

	getArt, _ := this.Get(id)
	utils.DelCache("GetArticleByUri.uri." + getArt.Uri)
	utils.DelCache("GetArticle.id." + fmt.Sprintf("%d", art.Id))

	return o.Delete(&art)
}

// 按月份统计文章数-cached
func (this *Article) CountByMonth() ([]orm.Params, error) {
	var maps []orm.Params

	err := utils.GetCache("CountByMonth", &maps)
	if nil != err {
		sql := "select DATE_FORMAT(time,'%Y年%m月') as date,count(*) as number ,year(time) as year, month(time) as month from article where status=1 group by date order by year desc, month desc"

		num, err := o.Raw(sql).Values(&maps)
		if err == nil && num > 0 {
			utils.SetCache("CountByMonth", maps, 3600)
			return maps, nil
		} else {
			return nil, err
		}
	} else {
		return maps, err
	}

}

// 获取某月的文章列表-cached
// year 年
// month 月
// page 页码
// numPerPage 每页条数
// 返回值:
// []orm.Params 文章
// bool 是否有下一页
// int 总页数
// error 错误
func (this *Article) ListByMonth(year int, month int, page int, numPerPage int) ([]orm.Params, bool, int, error) {
	if year < 0 {
		year = 1970
	}

	if month < 0 || month > 12 {
		month = 1
	}

	if page < 1 {
		page = 1
	}

	if numPerPage < 1 {
		numPerPage = 10
	}

	var maps, maps2 []orm.Params

	var err error

	// get data - cached
	err = utils.GetCache(fmt.Sprintf("ListByMonth.list.%d.%d.%d", year, month, page), &maps)
	if nil != err {
		sql1 := "select * from article where year(time)=? and month(time)=? and status=1 order by time desc limit ?,?"
		_, err = o.Raw(sql1, year, month, numPerPage*(page-1), numPerPage).Values(&maps)
		utils.SetCache(fmt.Sprintf("ListByMonth.list.%d.%d.%d", year, month, page), maps, 3600)
	}

	err = utils.GetCache(fmt.Sprintf("ListByMonth.count.%d.%d", year, month), &maps2)
	if nil != err {
		sql2 := "select count(*)as number from article where year(time)=? and month(time)=? and status=1"
		_, err = o.Raw(sql2, year, month).Values(&maps2)
		utils.SetCache(fmt.Sprintf("ListByMonth.count.%d.%d", year, month), maps2, 3600)
	}

	// calculate pages
	number, _ := strconv.Atoi(maps2[0]["number"].(string))
	var addFlag int
	if 0 == (number % numPerPage) {
		addFlag = 0
	} else {
		addFlag = 1
	}
	pages := number/numPerPage + addFlag

	var flagNextPage bool
	if pages == page {
		flagNextPage = false
	} else {
		flagNextPage = true
	}

	if err == nil {
		return maps, flagNextPage, pages, nil
	} else {
		return nil, false, pages, err
	}

}

// 文章分页列表
// page 页码
// numPerPage 每页条数
// 返回值:
// []orm.Params 文章
// bool 是否有下一页
// int 总页数
// error 错误
func (this *Article) ListPage(page int, numPerPage int) ([]orm.Params, bool, int, error) {
	// pagePerNum := 6
	sql1 := "select * from article where status=1 order by time desc limit ?," + fmt.Sprintf("%d", numPerPage)
	sql2 := "select count(*) as number from article where status=1"
	var maps, maps2 []orm.Params

	num, err := o.Raw(sql1, numPerPage*(page-1)).Values(&maps)
	if err != nil {
		fmt.Println("execute sql1 error:")
		fmt.Println(err)
		return nil, false, 0, err
	}

	err = utils.GetCache("ArticleNumber", &maps2)
	if nil != err {
		_, err = o.Raw(sql2).Values(&maps2)
		if err != nil {
			fmt.Println("execute sql2 error:")
			fmt.Println(err)
			return nil, false, 0, err
		}
		utils.SetCache("ArticleNumber", maps2, 3600)
	}

	number, err := strconv.Atoi(maps2[0]["number"].(string))

	var addFlag int
	if 0 == (number % numPerPage) {
		addFlag = 0
	} else {
		addFlag = 1
	}

	pages := number/numPerPage + addFlag

	var flagNextPage bool
	if pages == page {
		flagNextPage = false
	} else {
		flagNextPage = true
	}

	if err == nil && num > 0 {
		return maps, flagNextPage, pages, nil
	} else {
		return nil, false, pages, err
	}
}

// 同关键词文章列表
// 返回值:
// []orm.Params 文章
// bool 是否有下一页
// error 错误
func (this *Article) ListByKeyword(keyword string, page int, numPerPage int) ([]orm.Params, bool, int, error) {
	// numPerPage := 6
	sql1 := "select * from article where keywords like ? and status=1 order by time desc limit ?,?"
	sql2 := "select count(*) as number from article where keywords like ? and status=1"
	var maps, maps2 []orm.Params

	num, err := o.Raw(sql1, fmt.Sprintf("%%%s%%", keyword), numPerPage*(page-1), numPerPage).Values(&maps)
	o.Raw(sql2, fmt.Sprintf("%%%s%%", keyword)).Values(&maps2)

	number, _ := strconv.Atoi(maps2[0]["number"].(string))

	var addFlag int
	if 0 == (number % numPerPage) {
		addFlag = 0
	} else {
		addFlag = 1
	}

	pages := number/numPerPage + addFlag

	var flagNextPage bool
	if pages == page {
		flagNextPage = false
	} else {
		flagNextPage = true
	}

	if err == nil && num > 0 {
		return maps, flagNextPage, pages, nil
	} else {
		return nil, false, pages, err
	}
}

// 最热文章列表 - cached
func (this *Article) HottestArticleList() ([]orm.Params, error) {
	var maps []orm.Params

	// get data - cached
	err := utils.GetCache("HottestArticleList", &maps)
	if nil != err {
		sql := "select id,uri,title,count from article where status=1 order by count desc limit 20"
		o := orm.NewOrm()
		_, err = o.Raw(sql).Values(&maps)

		utils.SetCache("HottestArticleList", maps, 3600)
	}

	return maps, err
}

// 列出文章 for admin
func (this *Article) ArticleListForAdmin(page int, numPerPage int) ([]orm.Params, bool, int, error) {
	sql1 := "select id,uri,title,count,time from article where status=1 order by time desc limit ?," + fmt.Sprintf("%d", numPerPage)
	sql2 := "select count(*) as number from article where status=1"
	var maps, maps2 []orm.Params

	num, err := o.Raw(sql1, numPerPage*(page-1)).Values(&maps)
	if err != nil {
		fmt.Println("execute sql1 error:")
		fmt.Println(err)
		return nil, false, 0, err
	}

	err = utils.GetCache("ArticleNumber", &maps2)
	if nil != err {
		_, err = o.Raw(sql2).Values(&maps2)
		if err != nil {
			fmt.Println("execute sql2 error:")
			fmt.Println(err)
			return nil, false, 0, err
		}
		utils.SetCache("ArticleNumber", maps2, 3600)
	}

	number, err := strconv.Atoi(maps2[0]["number"].(string))

	var addFlag int
	if 0 == (number % numPerPage) {
		addFlag = 0
	} else {
		addFlag = 1
	}

	pages := number/numPerPage + addFlag

	var flagNextPage bool
	if pages == page {
		flagNextPage = false
	} else {
		flagNextPage = true
	}

	if err == nil && num > 0 {
		return maps, flagNextPage, pages, nil
	} else {
		return nil, false, pages, err
	}
}
