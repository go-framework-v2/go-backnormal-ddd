package model

import (
	"time"
)

// BizAppWechatConfigPO 持久化对象
type BizAppWechatConfigPO struct {
	Id              int64     `gorm:"column:id;primaryKey;not null"`
	AppId           int64     `gorm:"column:app_id;default:null"`
	WechatAppId     string    `gorm:"column:wechat_app_id;not null"`
	WechatAppSecret string    `gorm:"column:wechat_app_secret;not null"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (BizAppWechatConfigPO) TableName() string {
	return "biz_app_wechat_config"
}
