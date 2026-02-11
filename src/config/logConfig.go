package config

import "go.uber.org/zap"

type LogConfig struct {
	Root string `yaml:"root"`
}

func (c *LogConfig) AppendFieldMap(fMap map[string]string) {
	fMap["log.root"] = "Log.Root"
}

func (c *LogConfig) Print() {
	zap.L().Info("------------ log ------------")
	zap.L().Info("-- ", zap.String("root", c.Root))
}
