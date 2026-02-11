package identity

// AuthType 登录方式
type AuthType int8

const (
	AuthTypeGuest  AuthType = 0 // 游客
	AuthTypeMobile AuthType = 1 // 手机验证码
	AuthTypeAli    AuthType = 2 // 阿里一键
	AuthTypeWechat AuthType = 3 // 微信
)

func (a AuthType) String() string {
	switch a {
	case AuthTypeGuest:
		return "guest"
	case AuthTypeMobile:
		return "mobile"
	case AuthTypeAli:
		return "ali"
	case AuthTypeWechat:
		return "wechat"
	default:
		return "unknown"
	}
}

func (a AuthType) IsGuest() bool  { return a == AuthTypeGuest }
func (a AuthType) IsMobile() bool { return a == AuthTypeMobile }
func (a AuthType) IsWechat() bool { return a == AuthTypeWechat }
func (a AuthType) IsAli() bool    { return a == AuthTypeAli }
