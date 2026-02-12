package identity

import (
	"time"
)

// BizAppID 应用ID值对象
type BizAppID struct {
	value int64
}

// NewBizAppID 创建NewBizAppID
func NewBizAppID(value int64) (BizAppID, error) {
	return BizAppID{value: value}, nil
}

// NewBizAppIDFromSeed 根据种子生成BizAppID
func NewBizAppIDFromSeed(seed int64) BizAppID {
	now := time.Now().UnixNano()
	value := int64(seed) + now
	return BizAppID{value: value}
}

func (a BizAppID) Value() int64 {
	return a.value
}
