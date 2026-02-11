package main

import (
	_ "embed"
	"fmt"
	"go-backnormal-ddd/src/api"
	"go-backnormal-ddd/src/config"
	"go-backnormal-ddd/src/res"

	"github.com/go-framework-v2/go-config/viper"
	"github.com/go-framework-v2/go-log/log"

	"go.uber.org/zap"
)

//go:embed config_default.yaml
var configDefault []byte

func main() {
	// 初始化配置，合并生产/开发环境配置文件。通过conf.Cfg对象获取配置信息。
	viper.FDefault = configDefault
	cfg := config.New()
	// cfg.Env = "dev" // TODO: 手动设置环境变量，手动设置>默认环境配置合并default。线上需注释。
	viper.FillConfig(cfg, &cfg.BaseConfig)

	// 初始化日志
	log.InitLogger(config.Cfg.Log.Root, false)
	defer func(l *zap.Logger) {
		_ = l.Sync()
	}(zap.L())

	// 打印配置, 注意需要先初始化日志。
	cfg.Print()

	// 初始化上下文
	// 注册全局资源对象，并在res包中暴露使用
	res.InitResources()
	res.Print()
	// 添加这行：确保程序退出时关闭所有资源
	defer res.CloseAllResources()
	// 初始化定时任务cron
	// ...
	// 初始化全局变量globval
	// ...
	// 初始化go-cache内存缓存cache
	// ...
	// 初始化RPC客户端rpc
	// ...

	// 初始化和启动服务
	r := api.SetupRouter()
	r.Run(":" + fmt.Sprintf("%d", config.Cfg.ProjectConfig.Port))
}
