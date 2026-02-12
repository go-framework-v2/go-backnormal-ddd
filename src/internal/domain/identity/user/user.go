package user

import (
	"time"
	"go-backnormal-ddd/src/cons"
)

// User 用户聚合根
type User struct {
	id           UserID
	appId        int64
	authType     int8
	oaid         string
	deviceId     string
	wechatUserId int64
	mobileUserId int64
	nickname     string
	avatarUrl    string
	realName     string
	idCard       string
	ip           string
	deviceModel  string
	channel      string
	mobile       string
	createdAt    time.Time
	updatedAt    time.Time
	isValid      int8
}

// NewUserForCreate 构造用于首次插入的用户（id=0，由 DB 自增）
func NewUserForCreate(appId int64, authType int8, oaid, deviceId string) *User {
	return &User{
		id:       UserID{},
		appId:    appId,
		authType: authType,
		oaid:     oaid,
		deviceId: deviceId,
		isValid:  cons.IsValidYes,
	}
}

// RestoreUser 从持久化还原（供 Repository 使用）
func RestoreUser(
	id UserID,
	appId int64,
	authType int8,
	oaid, deviceId string,
	wechatUserId, mobileUserId int64,
	nickname, avatarUrl, realName, idCard string,
	ip, deviceModel, channel, mobile string,
	createdAt, updatedAt time.Time,
	isValid int8,
) *User {
	return &User{
		id:           id,
		appId:        appId,
		authType:     authType,
		oaid:         oaid,
		deviceId:     deviceId,
		wechatUserId: wechatUserId,
		mobileUserId: mobileUserId,
		nickname:     nickname,
		avatarUrl:    avatarUrl,
		realName:     realName,
		idCard:       idCard,
		ip:           ip,
		deviceModel:  deviceModel,
		channel:      channel,
		mobile:       mobile,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
		isValid:      isValid,
	}
}

func (u *User) ID() UserID             { return u.id }
func (u *User) AppId() int64           { return u.appId }
func (u *User) AuthType() int8         { return u.authType }
func (u *User) Oaid() string           { return u.oaid }
func (u *User) DeviceId() string       { return u.deviceId }
func (u *User) WechatUserId() int64    { return u.wechatUserId }
func (u *User) MobileUserId() int64    { return u.mobileUserId }
func (u *User) AvatarUrl() string      { return u.avatarUrl }
func (u *User) RealName() string       { return u.realName }
func (u *User) IdCard() string         { return u.idCard }
func (u *User) Nickname() string       { return u.nickname }
func (u *User) Mobile() string         { return u.mobile }
func (u *User) Ip() string             { return u.ip }
func (u *User) DeviceModel() string    { return u.deviceModel }
func (u *User) Channel() string        { return u.channel }
func (u *User) CreatedAt() time.Time   { return u.createdAt }
func (u *User) UpdatedAt() time.Time    { return u.updatedAt }
func (u *User) IsValid() int8          { return u.isValid }

// WithDeviceInfo 返回更新了 ip/deviceModel/channel 的副本，用于按需更新
func (u *User) WithDeviceInfo(ip, deviceModel, channel string) *User {
	return RestoreUser(
		u.id, u.appId, u.authType, u.oaid, u.deviceId,
		u.wechatUserId, u.mobileUserId,
		u.nickname, u.avatarUrl, u.realName, u.idCard,
		ip, deviceModel, channel, u.mobile,
		u.createdAt, u.updatedAt, u.isValid,
	)
}
