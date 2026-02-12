package repository

import (
	"fmt"

	"gorm.io/gorm"

	"go-backnormal-ddd/src/cons"
	"go-backnormal-ddd/src/internal/domain/identity/user"
	"go-backnormal-ddd/src/internal/infrastructure/persistence/mysql/model"
)

// BizUserRepositoryImpl biz_user 表仓储实现，实现领域 user.UserRepository
type BizUserRepositoryImpl struct {
	db *gorm.DB
}

// NewBizUserRepository 创建 biz_user 仓储
func NewBizUserRepository(db *gorm.DB) *BizUserRepositoryImpl {
	return &BizUserRepositoryImpl{db: db}
}

// FindByUk 根据唯一键查询用户：存在则返回，不存在则插入再返回
func (r *BizUserRepositoryImpl) FindByUk(appId int64, authType int8, oaid, deviceId string) (*user.User, error) {
	if appId <= 0 {
		return nil, fmt.Errorf("app_id cannot be less than or equal to 0")
	}
	if oaid == "" && deviceId == "" {
		return nil, fmt.Errorf("oaid and deviceId cannot both be empty")
	}
	var po model.BizUserPO
	preDb := r.db.Where("app_id = ? AND auth_type = ? AND is_valid = ?", appId, authType, cons.IsValidYes)
	if oaid != "" {
		preDb = preDb.Where("oaid = ?", oaid)
	}
	if deviceId != "" {
		preDb = preDb.Where("device_id = ?", deviceId)
	}
	err := preDb.First(&po).Error
	if err == nil {
		return r.toDomain(&po), nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("find biz_user by uk: %w", err)
	}
	po = r.toPO(user.NewUserForCreate(appId, authType, oaid, deviceId))
	if err = r.db.Create(&po).Error; err != nil {
		return nil, fmt.Errorf("create biz_user: %w", err)
	}
	if err = r.db.Where("id = ?", po.Id).First(&po).Error; err != nil {
		return nil, fmt.Errorf("find biz_user by id: %w", err)
	}
	return r.toDomain(&po), nil
}

// UpdateByFieldmap 按字段映射更新 biz_user 表
func (r *BizUserRepositoryImpl) UpdateByFieldmap(id user.UserID, fieldmap map[string]interface{}) error {
	po := model.BizUserPO{Id: id.Value()}
	if err := r.db.Model(&po).Updates(fieldmap).Error; err != nil {
		return fmt.Errorf("update biz_user by fieldmap: %w", err)
	}
	return nil
}

func (r *BizUserRepositoryImpl) toDomain(po *model.BizUserPO) *user.User {
	id, _ := user.NewUserID(po.Id)
	return user.RestoreUser(
		id, po.AppId, po.AuthType,
		po.Oaid, po.DeviceId,
		po.WechatUserId, po.MobileUserId,
		po.Nickname, po.AvatarUrl, po.RealName, po.IdCard,
		po.Ip, po.DeviceModel, po.Channel, po.Mobile,
		po.CreatedAt, po.UpdatedAt, po.IsValid,
	)
}

func (r *BizUserRepositoryImpl) toPO(u *user.User) model.BizUserPO {
	return model.BizUserPO{
		Id:           u.ID().Value(),
		AppId:        u.AppId(),
		AuthType:     u.AuthType(),
		Oaid:         u.Oaid(),
		DeviceId:     u.DeviceId(),
		WechatUserId: u.WechatUserId(),
		MobileUserId: u.MobileUserId(),
		Nickname:     u.Nickname(),
		AvatarUrl:    u.AvatarUrl(),
		RealName:     u.RealName(),
		IdCard:       u.IdCard(),
		Ip:           u.Ip(),
		DeviceModel:  u.DeviceModel(),
		Channel:      u.Channel(),
		Mobile:       u.Mobile(),
		CreatedAt:    u.CreatedAt(),
		UpdatedAt:    u.UpdatedAt(),
		IsValid:      u.IsValid(),
	}
}
