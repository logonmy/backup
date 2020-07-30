package config

import (
	"gopkg.in/ini.v1"
	"log"
)

var (
	Address           string
	TemplatePath      string
	StaticPath        string
	Qps               int64
	Release           bool
	SessionExpireTime int64
	Domain            string
	IndexPerPage      int64
	ListPerPage       int64
	From              string
	Password          string
	Host              string
	Port              string
)

// 加载自定义配置
func init() {
	cfg, err := ini.Load("config/config.ini")
	if err != nil {
		log.Fatal(err)
	}

	section, err := cfg.GetSection("")
	if err != nil {
		log.Fatal(err)
	}

	key, err := section.GetKey("staticPath")
	if err != nil {
		log.Fatal(err)
	}
	StaticPath = key.String()

	key, err = section.GetKey("templatePath")
	if err != nil {
		log.Fatal(err)
	}
	TemplatePath = key.String()

	key, err = section.GetKey("qps")
	if err != nil {
		log.Fatal(err)
	}
	Qps, err = key.Int64()
	if err != nil {
		log.Fatal(err)
	}

	key, err = section.GetKey("release")
	if err != nil {
		log.Fatal(err)
	}
	if key.String() == "1" {
		Release = true
	}

	key, err = section.GetKey("address")
	if err != nil {
		log.Fatal(err)
	}
	Address = key.String()

	key, err = section.GetKey("domain")
	if err != nil {
		log.Fatal(err)
	}
	Domain = key.String()

	key, err = section.GetKey("sessionExpireTime")
	if err != nil {
		log.Fatal(err)
	}
	SessionExpireTime, err = key.Int64()
	if err != nil {
		log.Fatal(err)
	}

	key, err = section.GetKey("from")
	if err != nil {
		log.Fatal(err)
	}
	From = key.String()

	key, err = section.GetKey("password")
	if err != nil {
		log.Fatal(err)
	}
	Password = key.String()

	key, err = section.GetKey("host")
	if err != nil {
		log.Fatal(err)
	}
	Host = key.String()

	key, err = section.GetKey("port")
	if err != nil {
		log.Fatal(err)
	}
	Port = key.String()
}
