package response

import "KubeGale/model/system"

type SysUserResponse struct {
	User system.SysUser `json:"user"`
}

// LoginResponse 用户登录返回结构
type LoginResponse struct {
	User           system.SysUser `json:"user"`
	Token          string         `json:"token"`
	ExpiresAt      int64          `json:"expiresAt"`
	BtnPermissions interface{}    `json:"btnPermissions"` // 按钮权限数据
}
