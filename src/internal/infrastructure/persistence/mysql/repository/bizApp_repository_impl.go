package repository

import (
	"fmt"

	"gorm.io/gorm"

	"go-backnormal-ddd/src/internal/domain/identity"
	"go-backnormal-ddd/src/internal/infrastructure/persistence/mysql/model"
)

// BizAppRepositoryImpl BizApp仓储实现
type BizAppRepositoryImpl struct {
	db *gorm.DB
}

// NewBizAppRepository 创建BizApp仓储
func NewBizAppRepository(db *gorm.DB) *BizAppRepositoryImpl {
	return &BizAppRepositoryImpl{db: db}
}

// 方法实现
func (r *BizAppRepositoryImpl) FindByProjectID(projectId int32) (*identity.BizApp, error) {
	var po model.BizAppPO
	err := r.db.Where("project_id = ?", projectId).First(&po).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find bizApp by project_id: %w", err)
	}
	return r.toDomain(&po), nil
}

// 数据转换
func (r *BizAppRepositoryImpl) toDomain(po *model.BizAppPO) *identity.BizApp {
	id, _ := identity.NewBizAppID(po.Id)
	return identity.RestoreBizApp(id, po.BundleId, po.AppName, po.ProjectId, po.CreatedAt)
}

func (r *BizAppRepositoryImpl) toPO(b *identity.BizApp) model.BizAppPO {
	return model.BizAppPO{
		Id:        b.ID().Value(),
		BundleId:  b.BundleId(),
		AppName:   b.AppName(),
		ProjectId: b.ProjectId(),
		CreatedAt: b.CreatedAt(),
	}
}
