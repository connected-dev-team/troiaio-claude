package main

import (
	"os"

	"github.com/goccy/go-yaml"
)

type Config struct {
	HOST        string `yaml:"host"`
	PORT        int    `yaml:"port"`
	DBASE_HOST  string `yaml:"dbase_host"`
	DBASE_PORT  int    `yaml:"dbase_port"`
	DBASE_NAME  string `yaml:"dbase_name"`
	DBASE_USER  string `yaml:"dbase_user"`
	DBASE_PASSWD string `yaml:"dbase_passwd"`
	JWT_SECRET  string `yaml:"jwt_secret"`
}

var CONF Config

func init() {
	data, err := os.ReadFile("conf.yaml")
	if err != nil {
		panic("Cannot read conf.yaml: " + err.Error())
	}
	err = yaml.Unmarshal(data, &CONF)
	if err != nil {
		panic("Cannot parse conf.yaml: " + err.Error())
	}
}
