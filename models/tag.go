package models

import (
	// "errors"
	// "fmt"
	// "github.com/astaxie/beego/orm"
	// "strconv"
	"time"
)

type Tags struct {
	Id   int
	Name string
	Time time.Time
}

var TheTags Tags

func (this *Tags) NewTag(tagName string) (int64, error) {
	tag := new(Tags)
	tag.Name = tagName
	return o.Insert(tag)
}
