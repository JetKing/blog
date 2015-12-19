package config

import (
	"github.com/go-ini/ini"
)

var cfg *ini.File

func init() {
	var err error

	cfg, err = ini.Load("custom/app.conf")

	if err != nil {
		cfg = ini.Empty()
	}

}

func iniConfig() {
	Set("Service", "port", "3081")
}

func Set(section, key, value string) error {
	sect, err := cfg.GetSection(section)
	if err != nil {
		return err
	} else {
		_, err = sect.NewKey(key, value)
		return err
	}
}

func Get(section, key string) interface{} {
	return cfg.Section(section).Key(key)
}
