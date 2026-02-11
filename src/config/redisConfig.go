package config

import "go.uber.org/zap"

type RedisConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Password     string `mapstructure:"password"`
	Db           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
	MaxRetries   int    `mapstructure:"max_retries"`
	MaxConnAge   int    `mapstructure:"max_conn_age"`
}

func (c *RedisConfig) AppendFieldMap(fMap map[string]string) {
	fMap["redis.host"] = "Redis.Host"
	fMap["redis.port"] = "Redis.Port"
	fMap["redis.password"] = "Redis.Password"
	fMap["redis.db"] = "Redis.Db"
	fMap["redis.pool_size"] = "Redis.PoolSize"
	fMap["redis.min_idle_conns"] = "Redis.MinIdleConns"
	fMap["redis.max_retries"] = "Redis.MaxRetries"
	fMap["redis.max_conn_age"] = "Redis.MaxConnAge"
}

func (c *RedisConfig) Print() {
	zap.L().Info("------------ redis ------------")
	zap.L().Info("-- ", zap.String("host", c.Host))
	zap.L().Info("-- ", zap.Int("port", c.Port))
	// zap.L().Info("-- ", zap.String("password", c.Password))
	zap.L().Info("-- ", zap.Int("db", c.Db))
	zap.L().Info("-- ", zap.Int("pool_size", c.PoolSize))
	zap.L().Info("-- ", zap.Int("min_idle_conns", c.MinIdleConns))
	zap.L().Info("-- ", zap.Int("max_retries", c.MaxRetries))
	zap.L().Info("-- ", zap.Int("max_conn_age", c.MaxConnAge))
}
