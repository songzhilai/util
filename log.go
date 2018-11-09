package util

import (
	seelog "github.com/cihub/seelog"
)

// Logger 日志全局对象
var Logger seelog.LoggerInterface

// SetLogConfig log配置
func SetLogConfig(configFile string) error {
	var err error
	Logger = seelog.Disabled
	Logger, err = seelog.LoggerFromConfigAsFile(configFile)
	if err != nil {
		return err
	}

	return nil
}
