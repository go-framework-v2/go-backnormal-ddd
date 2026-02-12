package user

// UserRepository 用户仓储接口
type UserRepository interface {
	FindByUk(appId int64, authType int8, oaid, deviceId string) (*User, error)
	UpdateByFieldmap(id UserID, fieldmap map[string]interface{}) error
}
