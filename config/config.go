package config

import (
	"io/ioutil"
	"log"
	"vocab8/domain/transport"

	"github.com/BurntSushi/toml"
	"github.com/go-pg/pg/v10"
)

var (
	ConfigPath string
	ConfigDir  string = "config/config.cfg"
	Cfg        *Config
	DB         *pg.DB
	WordPool   []transport.Word
)

type Config struct {
	Production bool    `toml:"production" `
	Port       int     `toml:"port" `
	Host       string  `toml:"host" `
	Context    string  `toml:"context" `
	DB         DBModel `toml:"db" `
	WordPath   string  `toml:"word-path" `
}

type DBModel struct {
	Host          string `toml:"host" `
	Port          int    `toml:"port" `
	Name          string `toml:"dbname" `
	User          string `toml:"user" `
	Password      string `toml:"password" `
	MigrationPath string `toml:"migration-path"`
}

func ParseConfig(path string, dest interface{}) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panic(err)
	}
	if _, err := toml.Decode(string(content), dest); err != nil {
		log.Panic(err)
	}
}
