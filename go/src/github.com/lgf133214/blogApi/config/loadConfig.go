package config

import (
	"gopkg.in/ini.v1"
)

var (
	PerPageNum int64
)

func init() {
	cfg, err := ini.Load("config/config.ini")
	if err != nil {
		panic(err)
	}

	defaultSection, err := cfg.GetSection("default")
	if err != nil {
		panic(err)
	}

	key, err := defaultSection.GetKey("perPageNum")
	if err != nil {
		panic(err)
	}
	PerPageNum, err = key.Int64()
	if err != nil {
		panic(err)
	}


}
