package model

import (
	"time"
)

// BizUserPO 用户表持久化对象
type BizUserPO struct {
	Id           int64     `gorm:"column:id;primaryKey;autoIncrement"`
	AppId        int64     `gorm:"column:app_id;default:null"` // 不可能为0, 默认0值就表示未赋值, 实际应存储为NULL
	AuthType     int8      `gorm:"column:auth_type;not null"`
	Oaid         string    `gorm:"column:oaid;size:255"`
	DeviceId     string    `gorm:"column:device_id;size:255"`
	WechatUserId int64     `gorm:"column:wechat_user_id;default:null"`
	MobileUserId int64     `gorm:"column:mobile_user_id;default:null"`
	Nickname     string    `gorm:"column:nickname;size:255"`
	AvatarUrl    string    `gorm:"column:avatar_url;size:500"`
	RealName     string    `gorm:"column:real_name;size:200"`
	IdCard       string    `gorm:"column:id_card;size:64"`
	Ip           string    `gorm:"column:ip;size:100"`
	DeviceModel  string    `gorm:"column:device_model;size:255"`
	Channel      string    `gorm:"column:channel;size:500"`
	Mobile       string    `gorm:"column:mobile;size:20"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime"`
	IsValid      int8      `gorm:"column:is_valid;not null;default:1"`
}

func (BizUserPO) TableName() string {
	return "biz_user"
}
