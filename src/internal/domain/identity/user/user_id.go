package user

// UserID 用户标识值对象
type UserID struct {
	value int64
}

// NewUserID 创建 UserID
func NewUserID(value int64) (UserID, error) {
	if value <= 0 {
		return UserID{}, nil
	}
	return UserID{value: value}, nil
}

func (u UserID) Value() int64 {
	return u.value
}
