package service

import (
	"gorm.io/gorm"
	"go-backnormal-ddd/src/internal/domain/identity/app"
	"go-backnormal-ddd/src/internal/domain/identity/user"
)

// LoginService 登录服务（应用层：编排领域 + 转 DTO）
type LoginService struct {
	db       *gorm.DB               // 这里是应用层可用的数据库入口. 让应用层可以控制事务边界, 可以用来开事务或只直接查库.
	appRepo  app.AppRepository // 职责清晰，应用服务会用的仓储. 保留db是负责「在不需要事务的场景里直接用、以及方便测试和扩展」.
	userRepo user.UserRepository
}

func NewLoginService(db *gorm.DB, appRepo app.AppRepository, userRepo user.UserRepository) *LoginService {
	return &LoginService{
		db: db, 
		appRepo: appRepo, 
		userRepo: userRepo,
	}
}
