package util

import (
	"github.com/Unknwon/goconfig"
)

var (
	cfg        *Config
	BaseConfig map[string]string
)

// Config 配置信息
type Config struct {
	Mysql  MysqlOpts  `json:"mysql,omitempty"`
	Influx InfluxOpts `json:"influx,omitempty"`
}

type InfluxOpts struct {
	Host     string `json:"host,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
	Database string `json:"database,omitempty"`
}

// MysqlOpts mysql配置
type MysqlOpts struct {
	Host     string `json:"host,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
	Port     int    `json:"port,omitempty"`
}

// ProcessConfigFile 从配置文件中读取配置信息
func ProcessConfigFile(configFile string) error {
	if configFile == "" {
		return nil
	}
	config, err := goconfig.LoadConfigFile(configFile)
	if err != nil {
		return err
	}
	cfg = new(Config)
	cfg.Mysql.Host = config.MustValue("mysql", "host", "127.0.0.1")
	cfg.Mysql.User = config.MustValue("mysql", "user", "root")
	cfg.Mysql.Password = config.MustValue("mysql", "password", "")
	cfg.Mysql.Port = config.MustInt("mysql", "port", 3306)

	cfg.Influx.Host = config.MustValue("influxdb", "host", "http://127.0.0.1:33311")
	cfg.Influx.User = config.MustValue("influxdb", "username", "root")
	cfg.Influx.Password = config.MustValue("influxdb", "password", "cloud#123456")
	cfg.Influx.Database = config.MustValue("influxdb", "database", "statdata")

	// BaseConfig, _ = config.GetSection("base")
	return nil
}
