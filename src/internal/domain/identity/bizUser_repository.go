package identity

// UserRepository 用户仓储接口
type UserRepository interface {
	// FindByUk 按唯一键查询：有则直接返回，无则提交插入事务再返回
	FindByUk(appId int64, authType int8, oaid, deviceId string) (*User, error)
	// UpdateByFieldmap 更新用户信息
	UpdateByFieldmap(id BizUserID, fieldmap map[string]interface{}) error
}
