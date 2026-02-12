package app

import "time"

// AppID 应用标识值对象
type AppID struct {
	value int64
}

// NewAppID 根据持久化值创建 AppID
func NewAppID(value int64) (AppID, error) {
	return AppID{value: value}, nil
}

// NewAppIDFromSeed 根据种子生成 AppID（用于新建聚合时）
func NewAppIDFromSeed(seed int64) AppID {
	return AppID{value: int64(seed) + time.Now().UnixNano()}
}

func (a AppID) Value() int64 {
	return a.value
}
