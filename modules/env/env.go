package env

import (
// "github.com/astaxie/beego"
)

var env map[string]interface{}

func init() {
	env = make(map[string]interface{})
}

func EnvSet(key string, value interface{}) {
	env[key] = value
}

func EnvGet(key string) interface{} {
	value, ok := env[key]
	if !ok {
		return nil
	} else {
		return value
	}
}
