package dto

import (
	"fmt"

	"go.uber.org/zap"
)

// ============ 根据 projectId 查 App（DDD Demo） ============
type GetAppByProjectIDReq struct {
	ProjectId int32 `json:"projectId"`
}

type GetAppByProjectIDResp struct {
	Id        int64  `json:"id"`
	BundleId  string `json:"bundleId"`
	AppName   string `json:"appName"`
	ProjectId int32  `json:"projectId"`
}

// ============ 游客登录 请求和响应 ============
type GuestLoginReq struct {
	DeviceId    string `json:"deviceId"`
	OAID        string `json:"oaid"`
	Model       string `json:"model"`
	RealChannel string `json:"realChannel"`
}

func (req GuestLoginReq) Validate() error {
	if req.DeviceId == "" && req.OAID == "" {
		err := fmt.Errorf("device_id or oaid is required")
		zap.S().Error(err)
		return err
	}
	if req.Model == "" {
		err := fmt.Errorf("model is required")
		zap.S().Error(err)
		return err
	}
	if req.RealChannel == "" {
		err := fmt.Errorf("real_channel is required")
		zap.S().Error(err)
		return err
	}
	return nil
}

type GuestLoginResp struct {
	UserID   int64  `json:"userId"`
	Nickname string `json:"nickname"`
	Mobile   string `json:"mobile"`
	Token    string `json:"token"`
}
