package model

import (
	"time"
)

// BizAppPO 持久化对象
type BizAppPO struct {
	Id        int64     `gorm:"column:id;primaryKey;not null"`
	BundleId  string    `gorm:"column:bundle_id;size:255;not null;index"`
	AppName   string    `gorm:"column:app_name;size:255;not null"`
	ProjectId int32     `gorm:"column:project_id;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (BizAppPO) TableName() string {
	return "biz_app"
}
