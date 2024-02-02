package mysql

import (
	"fmt"

	"github.com/jinzhu/gorm"

	//  mysql driver
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	dbDriver = "mysql"
	dbURLFmt = "%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local"
)

type Config struct {
	LogMode  string `yaml:"log_mode"` // dev = open debug log
	Driver   string `yaml:"db_driver"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func NewDB(cfg Config) *gorm.DB {
	dbURL := fmt.Sprintf(dbURLFmt, cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	db, err := gorm.Open(dbDriver, dbURL)
	if err != nil {
		panic(fmt.Sprintf("gorm open dbURL=%v, err=%v", dbURL, err))
	}

	db.LogMode(cfg.LogMode == "dev")

	return db
}
