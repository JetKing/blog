package models

import (
	// "github.com/astaxie/beego/orm"
	"time"
)

type UserLog struct {
	Id         int64
	User       int64
	Ip         string
	Ua         string
	Location   string
	Action     int
	CreateTime time.Time
}

var TheUserLog UserLog

func (this *UserLog) AddUserlog(user int64, ip string, ua string, location string, action int) (int64, error) {
	userLog := new(UserLog)
	userLog.User = user
	userLog.Ip = ip
	userLog.Ua = ua
	userLog.Location = location
	userLog.Action = action
	return o.Insert(userLog)
}

func (this *UserLog) GetUserLogByIp(ip string) (UserLog, error) {
	userLog := UserLog{Ip: ip}
	err := o.Read(&userLog, "ip")
	return userLog, err
}

func (this *UserLog) IsValidLocation(data map[string]interface{}) bool {
	cityName := data["cityName"].(string)
	countryName := data["countryName"].(string)
	regionName := data["regionName"].(string)
	if len(cityName) == 0 && len(countryName) == 0 && len(regionName) == 0 {
		return false
	} else {
		return true
	}
}
