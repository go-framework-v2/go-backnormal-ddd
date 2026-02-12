package app

import (
	"errors"
	"time"
)

// App 应用聚合根（含可选微信开放平台配置值对象）
type App struct {
	id           AppID
	bundleId     string
	appName      string
	projectId    int32
	createdAt    time.Time
	wechatConfig *WechatPlatformConfig // 值对象，nil 表示未配置微信
}

// NewApp 工厂方法（新建时通常无微信配置，后续通过仓储加载或单独更新）
func NewApp(bundleId, appName string, projectId int32) (*App, error) {
	if bundleId == "" || appName == "" || projectId == 0 {
		return nil, errors.New("bundleId, appName, projectId is required")
	}
	return &App{
		id:        NewAppIDFromSeed(1132431223),
		bundleId:  bundleId,
		appName:   appName,
		projectId: projectId,
		createdAt: time.Now(),
	}, nil
}

func (a *App) ID() AppID                              { return a.id }
func (a *App) BundleId() string                       { return a.bundleId }
func (a *App) AppName() string                        { return a.appName }
func (a *App) ProjectId() int32                       { return a.projectId }
func (a *App) CreatedAt() time.Time                   { return a.createdAt }
func (a *App) WechatConfig() *WechatPlatformConfig    { return a.wechatConfig }

// RestoreApp 从持久化还原（供 Repository 使用，wechatConfig 可为 nil）
func RestoreApp(id AppID, bundleId, appName string, projectId int32, createdAt time.Time, wechatConfig *WechatPlatformConfig) *App {
	return &App{
		id:           id,
		bundleId:     bundleId,
		appName:      appName,
		projectId:    projectId,
		createdAt:    createdAt,
		wechatConfig: wechatConfig,
	}
}
