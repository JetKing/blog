package models

import (
	"time"
)

type Tags struct {
	Id   int       `orm:"pk; auto;"`
	Name string    `orm:"not null; varchar(50);"`
	Time time.Time `orm:"not null; timestamp;"`
}

var TheTags Tags

func (this *Tags) NewTag(tagName string) (int64, error) {
	tag := new(Tags)
	tag.Name = tagName
	return o.Insert(tag)
}
