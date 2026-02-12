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
	db := res.MysqlDB // 数据库的连接池入口*gorm.DB.
	loginService = service.NewLoginService(
		db,                                 // 这里是应用层可用的数据库入口. 让应用层可以控制事务边界, 可以用来开事务或只直接查库.
		repository.NewBizAppRepository(db), // 职责清晰，应用服务会用的仓储. 保留db是负责「在不需要事务的场景里直接用、以及方便测试和扩展」.
		repository.NewBizUserRepository(db),
	)
}

func RouteLogin(r *gin.Engine) {
	initLoginService()

	loginGroup := r.Group("/user/login")
	{
		loginGroup.POST("/guest", GuestLogin) // 以此为参考
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
