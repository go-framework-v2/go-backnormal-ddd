package api

import (
	"go-backnormal-ddd/src/config"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SetupRouter() *gin.Engine {
	// 设置运行模式
	if config.Cfg.Gin.Release {
		gin.SetMode(gin.ReleaseMode) // 生产模式
	} else {
		gin.SetMode(gin.DebugMode) // 开发模式
	}

	// 初始化路由引擎
	r := gin.New()
	r.SetTrustedProxies(nil) // 允许所有来源的请求

	// 设置全局中间件
	// zap日志中间件
	// r.Use(ginzap.Ginzap(zap.L(), log.LogTmFmtWithMS, false)) // 可以打印所有的请求日志
	r.Use(ginzap.RecoveryWithZap(zap.L(), true))

	// 注册路由，根据不同的功能模块路由组进行注册
	RouteLogin(r)

	return r
}
