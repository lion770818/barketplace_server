package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// 讀取 env 檔案
func NewEnvConfig(filePath string) *Config {

	err := godotenv.Load()
	if err != nil {
		log.Fatal(fmt.Printf("Error loading .env file = %s", filePath))
	}

	connectNum, err := strconv.Atoi(os.Getenv("rabbitmq_connectNum"))
	if err != nil {
		log.Fatalf("connectNum rabbitmq_channelNum:%v, err=%v", os.Getenv("rabbitmq_channelNum"), err)
		return nil
	}
	channelNum, err := strconv.Atoi(os.Getenv("rabbitmq_channelNum"))
	if err != nil {
		log.Fatalf("channelNum rabbitmq_channelNum:%v, err=%v", os.Getenv("rabbitmq_channelNum"), err)
		return nil
	}

	max_size, err := strconv.Atoi(os.Getenv("log_max_size"))
	if err != nil {
		log.Fatalf("channelNum rabbitmq_channelNum:%v, err=%v", os.Getenv("log_max_size"), err)
		return nil
	}
	max_age, err := strconv.Atoi(os.Getenv("log_max_age"))
	if err != nil {
		log.Fatalf("max_age log_max_age:%v, err=%v", os.Getenv("log_max_age"), err)
		return nil
	}
	max_backups, err := strconv.Atoi(os.Getenv("log_max_backups"))
	if err != nil {
		log.Fatalf("max_backups max_backups:%v, err=%v", os.Getenv("log_max_backups"), err)
		return nil
	}

	baseConf := &ConfigBase{
		Web: Web{
			Mode: os.Getenv("web_mode"),
			Port: os.Getenv("web_port"),
		},
		Mysql: Mysql{
			LogMode:  os.Getenv("mysql_log_mode"),
			Driver:   os.Getenv("mysql_db_driver"),
			Host:     os.Getenv("mysql_host"),
			Port:     os.Getenv("mysql_port"),
			Database: os.Getenv("mysql_database"),
			User:     os.Getenv("mysql_user"),
			Password: os.Getenv("mysql_password"),
		},
		Auth: Auth{
			Active:     os.Getenv("auth_active"),
			ExpireTime: os.Getenv("auth_expireTime"),
			PrivateKey: os.Getenv("auth_privateKey"),
		},
		Redis: Redis{
			Host:     os.Getenv("redis_host"),
			Port:     os.Getenv("redis_port"),
			Password: os.Getenv("auth_password"),
		},
		RabbitMq: RabbitMq{
			Host:       os.Getenv("rabbitmq_host"),
			Port:       os.Getenv("rabbitmq_port"),
			User:       os.Getenv("rabbitmq_user"),
			Password:   os.Getenv("rabbitmq_password"),
			ConnectNum: connectNum,
			ChannelNum: channelNum,
		},
		Log: Log{
			Env:        os.Getenv("log_env"),
			Path:       os.Getenv("log_path"),
			Encoding:   os.Getenv("log_host"),
			MaxSize:    max_size,
			MaxAge:     max_age,
			MaxBackups: max_backups,
		},
	}

	// AuthExpireTime 解析为 time.Duration
	authExpireTime, err := time.ParseDuration(baseConf.Auth.ExpireTime)
	if err != nil {
		panic(err)
	}

	// 构造 Config
	pConfig := &Config{
		ConfigBase:     baseConf,
		AuthExpireTime: authExpireTime,
	}

	log.Printf("pConfig:%+v", pConfig)
	return pConfig
}
