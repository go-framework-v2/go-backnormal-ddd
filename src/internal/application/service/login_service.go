package service

import (
	"go-backnormal-ddd/src/internal/application/dto"
	"go-backnormal-ddd/src/internal/domain/identity"

	"github.com/gin-gonic/gin"
)

// LoginService 登录服务（应用层：编排领域 + 转 DTO）
type LoginService struct {
	appRepo identity.AppRepository
}

func NewLoginService(appRepo identity.AppRepository) *LoginService {
	return &LoginService{appRepo: appRepo}
}

// ============ 游客登录 业务流程编排（DDD 应用层） ============
// GuestLogin 游客登录
func (s *LoginService) GuestLogin(c *gin.Context, req dto.GuestLoginReq) (*dto.GuestLoginResp, error) {
	// 0. 从请求中获取信息
	// 1. 根据Header projectId获取App配置实例ID
	// 2. 根据唯一标识查询、更新或插入用户
	// 3. 生成token
	// 4. 返回结果
	return nil, nil
}

// ============ 根据 projectId 查 App（DDD Demo） ============
// GetAppByProjectID 根据项目ID查询应用：调领域仓储，将领域对象转为 DTO
func (s *LoginService) GetAppByProjectID(c *gin.Context, req dto.GetAppByProjectIDReq) (*dto.GetAppByProjectIDResp, error) {
	app, err := s.appRepo.FindByProjectID(req.ProjectId)
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, nil
	}
	return &dto.GetAppByProjectIDResp{
		Id:        app.ID().Value(),
		BundleId:  app.BundleId(),
		AppName:   app.AppName(),
		ProjectId: app.ProjectId(),
	}, nil
}
