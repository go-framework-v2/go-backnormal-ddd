package repository

import (
	"time"

	"gorm.io/gorm"

	"go-backnormal-ddd/src/internal/infrastructure/persistence/mysql/model"
)

// BizSmsCodeRepositoryImpl biz_sms_code 表仓储实现（防刷、落库供验证）
type BizSmsCodeRepositoryImpl struct {
	db *gorm.DB
}

// NewBizSmsCodeRepository 创建仓储
func NewBizSmsCodeRepository(db *gorm.DB) *BizSmsCodeRepositoryImpl {
	return &BizSmsCodeRepositoryImpl{db: db}
}

// GetTodaySendCount 统计该应用下该手机号今日发送条数（防刷）
func (r *BizSmsCodeRepositoryImpl) GetTodaySendCount(appId int64, mobile string) (int, error) {
	var count int64
	err := r.db.Model(&model.BizSmsCodePO{}).
		Where("app_id = ? AND mobile = ? AND created_at >= CURDATE() AND created_at < CURDATE() + INTERVAL 1 DAY", appId, mobile).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

// GetLastCreatedAt 该应用下该手机号最近一次发送时间（用于间隔限制）
func (r *BizSmsCodeRepositoryImpl) GetLastCreatedAt(appId int64, mobile string) (*time.Time, error) {
	var po model.BizSmsCodePO
	err := r.db.Where("app_id = ? AND mobile = ?", appId, mobile).Order("created_at DESC").Limit(1).First(&po).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &po.CreatedAt, nil
}

// Insert 插入一条验证码记录（scene 如 login/bind/reset）
func (r *BizSmsCodeRepositoryImpl) Insert(appId int64, mobile, code string, expiredAt time.Time, scene string) error {
	po := model.BizSmsCodePO{
		AppId:     appId,
		Mobile:    mobile,
		Code:      code,
		ExpiredAt: expiredAt,
		IsUsed:    0,
		Scene:     scene,
	}
	return r.db.Create(&po).Error
}

// GetLastByAppIdMobile 该应用下该手机号最近一条验证码记录（用于登录校验）
func (r *BizSmsCodeRepositoryImpl) GetLastByAppIdMobile(appId int64, mobile string) (*model.BizSmsCodePO, error) {
	var po model.BizSmsCodePO
	err := r.db.Where("app_id = ? AND mobile = ?", appId, mobile).Order("created_at DESC").Limit(1).First(&po).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &po, nil
}

// MarkUsed 标记验证码已使用
func (r *BizSmsCodeRepositoryImpl) MarkUsed(id int64) error {
	return r.db.Model(&model.BizSmsCodePO{}).Where("id = ?", id).Update("is_used", 1).Error
}
