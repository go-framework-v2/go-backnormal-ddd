package identity

import (
	"errors"
	"time"
)

// BizApp 应用聚合根
type BizApp struct {
	id        BizAppID
	bundleId  string
	appName   string
	projectId int32
	createdAt time.Time
}

// NewBizApp 工厂方法
func NewBizApp(bundleId, appName string, projectId int32) (*BizApp, error) {
	if bundleId == "" || appName == "" || projectId == 0 {
		return nil, errors.New("bundleId, appName, projectId is required")
	}

	id := NewBizAppIDFromSeed(1132431223)

	now := time.Now()

	return &BizApp{
		id:        id,
		bundleId:  bundleId,
		appName:   appName,
		projectId: projectId,
		createdAt: now,
	}, nil
}

// GetBizAppByProjectId 根据项目ID获取应用
func (a *BizApp) GetBizAppByProjectId(projectId int32) (*BizApp, error) {
	return nil, nil
}

// ... 业务方法

func (a BizApp) ID() BizAppID         { return a.id }
func (a BizApp) BundleId() string     { return a.bundleId }
func (a BizApp) AppName() string      { return a.appName }
func (a BizApp) ProjectId() int32     { return a.projectId }
func (a BizApp) CreatedAt() time.Time { return a.createdAt }

// RestoreBizApp 从持久化还原（供Repository使用）
func RestoreBizApp(id BizAppID, bundleId, appName string, projectId int32, createdAt time.Time) *BizApp {
	return &BizApp{
		id:        id,
		bundleId:  bundleId,
		appName:   appName,
		projectId: projectId,
		createdAt: createdAt,
	}
}
