package app

// AppRepository 应用仓储（聚合根为 App，微信配置作为值对象由仓储一并加载）
type AppRepository interface {
	FindByProjectID(projectId int32) (*App, error)
	Insert(a *App) error
}
