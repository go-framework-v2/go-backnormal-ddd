package config

import "go.uber.org/zap"

type MySQLConfig struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	DbName      string `mapstructure:"dbname"`
	Username    string `mapstructure:"username"`
	Password    string `mapstructure:"password"`
	Maxconn     int    `mapstructure:"maxconn"`
	Minconn     int    `mapstructure:"minconn"`
	Maxlifetime int    `mapstructure:"maxlifetime"`
}

func (c *MySQLConfig) AppendFieldMap(fMap map[string]string) {
	fMap["mysql.host"] = "Mysql.Host"
	fMap["mysql.port"] = "Mysql.Port"
	fMap["mysql.dbname"] = "Mysql.DbName"
	fMap["mysql.username"] = "Mysql.Username"
	fMap["mysql.password"] = "Mysql.Password"
	fMap["mysql.maxconn"] = "Mysql.Maxconn"
	fMap["mysql.minconn"] = "Mysql.Minconn"
	fMap["mysql.maxlife"] = "Mysql.Maxlifetime"
}

func (c *MySQLConfig) Print() {
	zap.L().Info("------------ mysql ------------")
	zap.L().Info("-- ", zap.String("host", c.Host))
	zap.L().Info("-- ", zap.Int("port", c.Port))
	zap.L().Info("-- ", zap.String("dbname", c.DbName))
	zap.L().Info("-- ", zap.String("username", c.Username))
	// zap.L().Info("-- ", zap.String("password", c.Password))
	zap.L().Info("-- ", zap.Int("maxconn", c.Maxconn))
	zap.L().Info("-- ", zap.Int("minconn", c.Minconn))
	zap.L().Info("-- ", zap.Int("maxlife", c.Maxlifetime))
}
