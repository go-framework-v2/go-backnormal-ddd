package config

import "go.uber.org/zap"

type GinConfig struct {
	Release bool `yaml:"release"`
}

func (c *GinConfig) AppendFieldMap(fMap map[string]string) {
	fMap["gin.release"] = "Gin.Release"
}

func (c *GinConfig) Print() {
	zap.L().Info("------------ gin ------------")
	zap.L().Info("-- ", zap.Bool("release", c.Release))
}
