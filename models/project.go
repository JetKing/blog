package models

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/gogather/com/log"
	"strconv"
	"time"
)

type Project struct {
	Id          int
	Name        string
	IconUrl     string
	Author      string
	Description string
	Time        time.Time
}

var TheProject Project

// get project by id or name
func (this *Project) GetProject(id int, name string) (*Project, error) {
	var err error

	pro := Project{}

	if id > 0 {
		pro = Project{Id: id}
		err = o.Read(&pro, "Id")
	} else if len(name) > 0 {
		pro = Project{Name: name}
		err = o.Read(&pro, "Name")
	} else {
		err = errors.New("至少有一个条件")
	}

	return &pro, err
}

// 项目分页列表
// select * from project order by time desc limit 0,6
// page 页码
// numPerPage 每页条数
// 返回值:
// []orm.Params 文章
// bool 是否有下一页
// int 总页数
// error 错误
func (this *Project) ListProjects(page int, numPerPage int) ([]orm.Params, bool, int, error) {
	// pagePerNum := 6
	sql1 := "select * from project order by time desc limit ?," + fmt.Sprintf("%d", numPerPage)
	sql2 := "select count(*) as number from project"
	var maps, maps2 []orm.Params

	num, err := o.Raw(sql1, numPerPage*(page-1)).Values(&maps)
	if err != nil {
		fmt.Println("execute sql1 error:")
		fmt.Println(err)
		return nil, false, 0, err
	}

	_, err = o.Raw(sql2).Values(&maps2)
	if err != nil {
		fmt.Println("execute sql2 error:")
		fmt.Println(err)
		return nil, false, 0, err
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

// add project
func (this *Project) AddProject(name string, icon string, author string, description string, createTime time.Time) (int64, error) {
	pro := new(Project)
	pro.Name = name
	pro.IconUrl = icon
	pro.Author = author
	pro.Description = description
	pro.Time = createTime
	return o.Insert(pro)
}

// delete project
func (this *Project) DeleteProject(id int64) error {
	_, err := o.Delete(&Project{Id: int(id)})
	return err
}

// update project
func (this *Project) UpdateProject(id int64, name string, icon string, description string) error {
	var pro *Project
	var err error

	if 0 != id {
		pro, err = this.GetProject(int(id), "")
		if err != nil {
			return err
		}
	} else {
		return errors.New("id should not 0")
	}

	log.Pinkln(pro)

	pro.Name = name
	pro.IconUrl = icon
	pro.Description = description

	_, err = o.Update(pro, "name", "icon_url", "description")
	return err
}
