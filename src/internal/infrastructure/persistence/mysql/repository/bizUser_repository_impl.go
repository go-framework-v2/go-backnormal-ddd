package repository

import (
	"fmt"

	"gorm.io/gorm"

	"go-backnormal-ddd/src/cons"
	"go-backnormal-ddd/src/internal/domain/identity"
	"go-backnormal-ddd/src/internal/infrastructure/persistence/mysql/model"
)

// BizUserRepositoryImpl 用户仓储实现
type BizUserRepositoryImpl struct {
	db *gorm.DB
}

// NewBizUserRepository 创建用户仓储
func NewBizUserRepository(db *gorm.DB) *BizUserRepositoryImpl {
	return &BizUserRepositoryImpl{db: db}
}

// FindByUk 根据唯一键 查询用户 存在直接返回，不存在插入再返回
func (r *BizUserRepositoryImpl) FindByUk(appId int64, authType int8, oaid, deviceId string) (*identity.User, error) {
	// 1. 先查询是否存在
	var po model.BizUserPO
	preDb := r.db.Where("app_id =? AND auth_type =? AND is_valid =?", appId, authType, cons.IsValidYes)
	if oaid != "" {
		preDb = preDb.Where("oaid =?", oaid)
	}
	if deviceId != "" {
		preDb = preDb.Where("device_id =?", deviceId)
	}
	err := preDb.First(&po).Error

	if err == nil { // 2. 存在, 直接返回
		return r.toDomain(&po), nil
	} else if err == gorm.ErrRecordNotFound { // 3. 不存在，则插入（使用当前 r.db，可为外层事务 tx）
		po = r.toPO(identity.NewUserForCreate(appId, authType, oaid, deviceId))
		if err = r.db.Create(&po).Error; err != nil {
			return nil, fmt.Errorf("create user error: %v", err)
		}
		// 插入成功, 再次查询
		err = r.db.Where("id = ?", po.Id).First(&po).Error
		if err != nil {
			return nil, fmt.Errorf("find user by id error: %v", err)
		}
		return r.toDomain(&po), nil
	} else {
		return nil, fmt.Errorf("find user by uk error: %v", err)
	}
}

// UpdateByFieldmap 更新用户信息
func (r *BizUserRepositoryImpl) UpdateByFieldmap(id identity.BizUserID, fieldmap map[string]interface{}) error {
	po := model.BizUserPO{
		Id: id.Value(),
	}
	if err := r.db.Model(&po).Updates(fieldmap).Error; err != nil {
		return fmt.Errorf("update user by fieldmap error: %v", err)
	}
	return nil
}

func (r *BizUserRepositoryImpl) toDomain(po *model.BizUserPO) *identity.User {
	id, _ := identity.NewBizUserID(po.Id)
	return identity.RestoreUser(
		id, po.AppId, po.AuthType,
		po.Oaid, po.DeviceId,
		po.WechatUserId, po.MobileUserId,
		po.Nickname, po.AvatarUrl, po.RealName, po.IdCard,
		po.Ip, po.DeviceModel, po.Channel, po.Mobile,
		po.CreatedAt, po.UpdatedAt, po.IsValid,
	)
}

func (r *BizUserRepositoryImpl) toPO(u *identity.User) model.BizUserPO {
	po := model.BizUserPO{
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

	return po
}
