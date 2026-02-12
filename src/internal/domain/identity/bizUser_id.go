package identity

// BizUserID 用户ID值对象
type BizUserID struct {
	value int64
}

// NewBizUserID 创建 BizUserID
func NewBizUserID(value int64) (BizUserID, error) {
	if value <= 0 {
		return BizUserID{}, nil
	}
	return BizUserID{value: value}, nil
}

func (u BizUserID) Value() int64 {
	return u.value
}
