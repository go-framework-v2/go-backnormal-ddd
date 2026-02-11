package config

import (
	"go.uber.org/zap"

	"github.com/go-framework-v2/go-config/viper"
)

// 业务全局配置
var Cfg *Config

type Config struct {
	viper.BaseConfig

	Gin GinConfig `mapstructure:"gin"`
	Log LogConfig `mapstructure:"log"`

	Mysql MySQLConfig `mapstructure:"mysql"`
	Redis RedisConfig `mapstructure:"redis"`
}

func New() *Config {
	if Cfg == nil {
		Cfg = &Config{}
	}

	return Cfg
}

func (c *Config) AppendFieldMap(fMap map[string]string) {
	c.Gin.AppendFieldMap(fMap)
	c.Log.AppendFieldMap(fMap)

	c.Mysql.AppendFieldMap(fMap)
	c.Redis.AppendFieldMap(fMap)
}

func (c *Config) Print() {
	zap.L().Info("++++++++++++++ config info begin: ++++++++++++++")
	c.BaseConfig.Print()

	c.Gin.Print()
	c.Log.Print()

	c.Mysql.Print()
	c.Redis.Print()

	zap.L().Info("++++++++++++++ config info end: ++++++++++++++")
}
