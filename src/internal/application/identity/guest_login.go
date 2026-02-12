package service

import (
	"fmt"

	"go-backnormal-ddd/src/cons"
	"go-backnormal-ddd/src/internal/application/identity/dto"
	"go-backnormal-ddd/src/internal/domain/identity/user"
	"go-backnormal-ddd/src/internal/infrastructure/persistence/mysql/repository"
	"go-backnormal-ddd/src/middleware/jwt"
	"go-backnormal-ddd/src/tool"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ============ 游客登录 业务流程编排（DDD 应用层） ============// ============ 游客登录 业务流程编排（DDD 应用层） ============
func (s *LoginService) GuestLogin(c *gin.Context, req dto.GuestLoginReq) (*dto.GuestLoginResp, error) {
	// 0. 从请求中获取信息
	ip := tool.GetIp(c)                    // 获取客户端IP地址
	projectId, err := tool.GetProjectId(c) // 获取header里面的参数projectId
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}
	oaid := req.OAID // 获取req请求参数deviceId/oaid,model,realChannel
	deviceId := req.DeviceId
	model := req.Model
	realChannel := req.RealChannel
	authType := cons.AuthTypeGuest // 游客身份登录方式

	// 1. 根据 projectId 获取 App，得到 appId
	app, err := s.appRepo.FindByProjectID(int32(projectId))
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}
	if app == nil {
		err = fmt.Errorf("app not found by projectId: %d", projectId)
		zap.S().Error(err)
		return nil, err
	}
	appId := app.ID().Value()

	// 显式事务：Begin → 业务操作 → 成功则 Commit，失败或 panic 则 Rollback
	tx := s.db.Begin()
	if tx.Error != nil {
		zap.S().Error(tx.Error)
		return nil, fmt.Errorf("tx begin failed: %w", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var user *user.User
	userRepoTx := repository.NewBizUserRepository(tx)
	// 2. 根据唯一键查询用户，有则返回，无则插入再返回
	u, err := userRepoTx.FindByUk(appId, authType, oaid, deviceId)
	if err != nil {
		zap.S().Error(err)
		tx.Rollback()
		return nil, err
	}
	if u == nil {
		err = fmt.Errorf("user nil after FindByUk (appId=%d, authType=%d, oaid=%s, deviceId=%s)", appId, authType, oaid, deviceId)
		zap.S().Error(err)
		tx.Rollback()
		return nil, err
	}
	// 按需更新用户信息
	userUpateParam := make(map[string]interface{})
	if u.Ip() != ip {
		userUpateParam["ip"] = ip
	}
	if u.DeviceModel() != model {
		userUpateParam["device_model"] = model
	}
	if u.Channel() != realChannel {
		userUpateParam["channel"] = realChannel
	}
	if len(userUpateParam) > 0 {
		if err = userRepoTx.UpdateByFieldmap(u.ID(), userUpateParam); err != nil {
			zap.S().Error(err)
			tx.Rollback()
			return nil, fmt.Errorf("update user failed: %w", err)
		}
	}
	user = u

	if err = tx.Commit().Error; err != nil {
		zap.S().Error(err)
		tx.Rollback()
		return nil, fmt.Errorf("tx commit failed: %w", err)
	}

	// 3. 生成 token
	token, err := jwt.GenerateToken(user.ID().Value())
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}

	// 4. 返回结果
	return &dto.GuestLoginResp{
		UserID:   user.ID().Value(),
		Nickname: user.Nickname(),
		Mobile:   user.Mobile(),
		Token:    token,
	}, nil
}
