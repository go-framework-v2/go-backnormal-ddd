package model

import "time"

// BizSmsCodePO biz_sms_code 表持久化对象
type BizSmsCodePO struct {
	Id        int64     `gorm:"column:id;primaryKey;autoIncrement"`
	AppId     int64     `gorm:"column:app_id;not null"`
	Mobile    string    `gorm:"column:mobile;not null;size:20"`
	Code      string    `gorm:"column:code;not null;size:10"`
	ExpiredAt time.Time `gorm:"column:expired_at;not null"`
	IsUsed    int8      `gorm:"column:is_used;not null;default:0"`
	Scene     string    `gorm:"column:scene;size:20"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (BizSmsCodePO) TableName() string {
	return "biz_sms_code"
}
