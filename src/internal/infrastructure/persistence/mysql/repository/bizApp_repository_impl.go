package repository

import (
	"fmt"

	"gorm.io/gorm"

	"go-backnormal-ddd/src/internal/domain/identity/app"
	"go-backnormal-ddd/src/internal/infrastructure/persistence/mysql/model"
)

// BizAppRepositoryImpl biz_app 表仓储实现，实现领域 app.AppRepository（聚合根 App，微信配置作为值对象一并加载）
type BizAppRepositoryImpl struct {
	db *gorm.DB
}

// NewBizAppRepository 创建 biz_app 仓储
func NewBizAppRepository(db *gorm.DB) *BizAppRepositoryImpl {
	return &BizAppRepositoryImpl{db: db}
}

// FindByProjectID 按项目ID查询应用，并加载 biz_app_wechat_config 转为值对象
func (r *BizAppRepositoryImpl) FindByProjectID(projectId int32) (*app.App, error) {
	var po model.BizAppPO
	if err := r.db.Where("project_id = ?", projectId).First(&po).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("find biz_app by project_id: %w", err)
	}
	id, _ := app.NewAppID(po.Id)
	var wechatConfig *app.WechatPlatformConfig
	var wechatPO model.BizAppWechatConfigPO
	if err := r.db.Where("app_id = ?", po.Id).First(&wechatPO).Error; err == nil {
		cfg, _ := app.NewWechatPlatformConfig(wechatPO.WechatAppId, wechatPO.WechatAppSecret)
		wechatConfig = cfg
	}
	return app.RestoreApp(id, po.BundleId, po.AppName, po.ProjectId, po.CreatedAt, wechatConfig), nil
}

// Insert 插入 biz_app 表（微信配置需单独维护 biz_app_wechat_config）
func (r *BizAppRepositoryImpl) Insert(a *app.App) error {
	po := model.BizAppPO{
		Id:        a.ID().Value(),
		BundleId:  a.BundleId(),
		AppName:   a.AppName(),
		ProjectId: a.ProjectId(),
		CreatedAt: a.CreatedAt(),
	}
	if err := r.db.Create(&po).Error; err != nil {
		return fmt.Errorf("insert biz_app: %w", err)
	}
	return nil
}
