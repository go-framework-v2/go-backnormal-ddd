package res

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
	"go-backnormal-ddd/src/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	MysqlDB *gorm.DB // MysqlDB 全局唯一的数据库实例, 对外提供

	RESOURCE_MYSQLDB_NAME = "res-MysqlDB" // res mysql数据库全局单例资源名称

	mysqlOnce      sync.Once
	initMysqlError error
)

// DBResource 实现全局资源接口的包装器
type DBResource struct {
	DB *gorm.DB
}

func (p *DBResource) Name() string {
	return RESOURCE_MYSQLDB_NAME
}

// GetMysqlDB 获取全局唯一的数据库实例(线程安全)
func GetMysqlDB() (*gorm.DB, error) {
	mysqlOnce.Do(func() {
		instance, err := initDB()
		if err != nil {
			initMysqlError = err
			return
		}
		MysqlDB = instance
	})

	return MysqlDB, initMysqlError
}

// initDB 初始化数据库连接(不导出)
func initDB() (*gorm.DB, error) {
	var err error

	username := config.Cfg.Mysql.Username
	password := config.Cfg.Mysql.Password
	host := config.Cfg.Mysql.Host
	port := config.Cfg.Mysql.Port
	dbname := config.Cfg.Mysql.DbName
	maxconn := config.Cfg.Mysql.Maxconn
	minconn := config.Cfg.Mysql.Minconn
	maxlifetimeSec := config.Cfg.Mysql.Maxlifetime
	maxlifetime := time.Duration(maxlifetimeSec) * time.Second

	// 1. 连接数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, dbname)

	// 2. 配置 GORM
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io.Writer
		logger.Config{
			SlowThreshold:             time.Second, // 慢查询阈值，超过1秒的查询才会打印
			LogLevel:                  logger.Warn, // 设置为 Warn，输出警告和错误日志
			IgnoreRecordNotFoundError: true,        // 不把 record not found 当作 error
			Colorful:                  true,        // 彩色日志
		},
	)

	config := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 全局禁用表名复数
		},
		Logger: newLogger, // 使用我们自定义的 logger
	}

	// 3. 连接数据库
	MysqlDB, err = gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		err = fmt.Errorf("failed to connect to database: %w", err)
		return nil, err
	}

	// 4. 配置数据库连接池
	sqlDB, sqlErr := MysqlDB.DB()
	if sqlErr != nil {
		err = fmt.Errorf("failed to configure database pool: %w", sqlErr)
		return nil, err
	}
	sqlDB.SetMaxOpenConns(maxconn)        // 最大连接数
	sqlDB.SetMaxIdleConns(minconn)        // 最大空闲连接数
	sqlDB.SetConnMaxLifetime(maxlifetime) // 连接的最长存活时间

	// 5. 测试数据库连接
	if pingErr := sqlDB.Ping(); pingErr != nil {
		err = fmt.Errorf("failed to ping database: %w", pingErr)
		return nil, err
	}

	// fmt.Println("Current DSN:", dsn)
	// fmt.Println("Successfully connected to MySQL!")

	return MysqlDB, nil
}

// 新增：确保注册只执行一次
var registerMysqlDBOnce sync.Once

// RegisterMysqlDBToGlobal 将数据库实例注册到全局资源管理器
func RegisterMysqlDBToGlobal() error {
	var registerErr error
	registerMysqlDBOnce.Do(func() {
		db, err := GetMysqlDB()
		if err != nil {
			registerErr = err
			return
		}

		RegisterResource(RESOURCE_MYSQLDB_NAME, func() Resource {
			return &DBResource{DB: db}
		})
	})
	return registerErr
}

// 在 DBResource 结构体中添加 Close 方法
func (p *DBResource) Close() error {
	if p.DB == nil {
		return nil
	}

	sqlDB, err := p.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	return sqlDB.Close()
}
