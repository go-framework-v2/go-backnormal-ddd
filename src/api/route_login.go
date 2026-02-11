package api

import (
	"go-backnormal-ddd/src/internal/application/dto"
	"go-backnormal-ddd/src/internal/application/service"
	"go-backnormal-ddd/src/internal/infrastructure/persistence/mysql/repository"
	"go-backnormal-ddd/src/res"
	"go-backnormal-ddd/src/tool"

	"github.com/gin-gonic/gin"
)

var loginService *service.LoginService

func initLoginService() {
	db := res.MysqlDB

	loginService = service.NewLoginService(
		repository.NewBizAppRepository(db),
	)
}

func RouteLogin(r *gin.Engine) {
	initLoginService()

	loginGroup := r.Group("/user/login")
	{
		loginGroup.POST("/app", GetAppByProjectID) // DDD Demo: 根据 projectId 查 app

		loginGroup.POST("/guest", GuestLogin)
		loginGroup.POST("/wechat", nil)
		loginGroup.POST("/aliMobile", nil)
		loginGroup.POST("/sms/send", nil)
		loginGroup.POST("/sms/verify", nil)
	}
}

// ============ 游客登录 请求转换 ============
func GuestLogin(c *gin.Context) {
	tool.HandleWithBindWithC(c, loginService.GuestLogin, dto.GuestLoginResp{})
}

// ============ 根据 projectId 查 App（DDD Demo） ============
// GetAppByProjectID 从 query 读取 projectId，调应用层，返回 JSON
func GetAppByProjectID(c *gin.Context) {
	tool.HandleWithBindWithC(c, loginService.GetAppByProjectID, dto.GetAppByProjectIDResp{})
}
