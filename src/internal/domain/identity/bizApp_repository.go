package identity

type AppRepository interface {
	FindByProjectID(projectId int32) (*BizApp, error)
	Insert(app *BizApp) error
}
