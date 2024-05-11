package conf

import (
	"fmt"
	"sync/atomic"

	"github.com/spf13/viper"
	"github.com/chuck1024/gd/v2/logger"
)

var (
	v atomic.Value
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(fmt.Sprintf("/etc/%s/", ""))
	viper.AddConfigPath(fmt.Sprintf("$HOME/.%s/", ""))
	viper.AddConfigPath("./conf/")
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Warn("not found config in ['/etc', '$HOME/', './conf/']")
		} else {
			logger.Fatalf("read config failed:%v", err)
		}
	}
	var conf Config
	err = viper.Unmarshal(&conf)
	if err != nil {
		logger.Errorf("unmarshal config to struct failed:%v", err)
	}
	v.Store(&conf)
}
