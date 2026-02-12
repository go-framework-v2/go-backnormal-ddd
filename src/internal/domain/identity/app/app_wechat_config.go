package app

import "errors"

// WechatPlatformConfig 微信开放平台配置（值对象，隶属于 App 聚合根，无独立身份）
type WechatPlatformConfig struct {
	wechatAppId     string
	wechatAppSecret string
}

// NewWechatPlatformConfig 构造微信平台配置值对象
func NewWechatPlatformConfig(wechatAppId, wechatAppSecret string) (*WechatPlatformConfig, error) {
	if wechatAppId == "" {
		return nil, errors.New("wechatAppId must not be empty")
	}
	if wechatAppSecret == "" {
		return nil, errors.New("wechatAppSecret must not be empty")
	}
	return &WechatPlatformConfig{
		wechatAppId:     wechatAppId,
		wechatAppSecret: wechatAppSecret,
	}, nil
}

func (c WechatPlatformConfig) WechatAppId() string     { return c.wechatAppId }
func (c WechatPlatformConfig) WechatAppSecret() string { return c.wechatAppSecret }
