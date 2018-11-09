package util

import (
	"database/sql"
	"fmt"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

var Mysqldb *sql.DB

const (
	// DefaultDataBase 默认数据库
	DefaultDataBase = "dc_mind"

	connectCount = 3
)

// InitMysql 初始化mysql连接池
func InitMysql() error {

	datasourceName := mysqlDataSourceName(
		cfg.Mysql.Host,
		cfg.Mysql.User,
		cfg.Mysql.Password,
		cfg.Mysql.Port,
		DefaultDataBase)

	var err error
	for i := 0; i < connectCount; i++ {
		Mysqldb, err = sql.Open("mysql", datasourceName)
		if err != nil {
			Logger.Error("mysql connect fail:" + datasourceName)
		}
		Mysqldb.SetMaxOpenConns(100)
		Mysqldb.SetMaxIdleConns(100)
		Mysqldb.Ping()
	}

	return err
}

// MysqlDataSourceName 返回mysql的连接字符串
func mysqlDataSourceName(host string, user string, password string, port int, dbName string) string {

	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8",
		user,
		password,
		host,
		port,
		dbName)
}
