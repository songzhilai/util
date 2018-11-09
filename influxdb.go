package util

import (
	"time"

	client "github.com/influxdata/influxdb/client/v2"
)

//Influxdb 历史数据库连接
var Influxdb client.Client

//InitInflux 连接influxdb
func InitInflux() error {
	var err error
	for i := 0; i < 3; i++ {
		Influxdb, err = client.NewHTTPClient(client.HTTPConfig{
			Addr:     cfg.Influx.Host,
			Username: cfg.Influx.User,
			Password: cfg.Influx.Password,
		})
		if err != nil {
			Logger.Error("influx连接失败" + err.Error())
			time.Sleep(3 * time.Second)
		} else {
			break
		}
	}
	return err
}

//QueryDB  查询influxdb数据
func QueryDB(cmd string) (res []client.Result, err error) {
	database := cfg.Influx.Database
	q := client.Query{
		Command: cmd,
		// Database: "funcInfluxdb",
		Database: database,
	}
	if response, err := Influxdb.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}
