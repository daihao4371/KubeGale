package request

import (
	"KubeGale/model/common/request"
	"KubeGale/model/im"
)

// CreateDingTalkRequest 创建钉钉通知请求结构
type CreateDingTalkRequest struct {
	Name               string               `json:"name" binding:"required"`                // 通知名称
	NotificationPolicy string               `json:"notification_policy" binding:"required"` // 通知策略
	SignatureKey       string               `json:"signature_key" binding:"required"`       // 签名秘钥
	RobotURL           string               `json:"robot_url" binding:"required"`           // 机器人地址
	SendDailyStats     bool                 `json:"send_daily_stats"`                       // 是否发送每日统计
	CardContent        im.CardContentConfig `json:"card_content,omitempty"`                 // 卡片内容配置
}

// CreateFeiShuRequest 创建飞书通知请求结构
type CreateFeiShuRequest struct {
	Name               string               `json:"name" binding:"required"`                // 通知名称
	NotificationPolicy string               `json:"notification_policy" binding:"required"` // 通知策略
	RobotURL           string               `json:"robot_url" binding:"required"`           // 机器人地址
	SendDailyStats     bool                 `json:"send_daily_stats"`                       // 是否发送每日统计
	CardContent        im.CardContentConfig `json:"card_content,omitempty"`                 // 卡片内容配置
}

// UpdateDingTalkRequest 更新钉钉通知请求结构
type UpdateDingTalkRequest struct {
	ID                 uint                 `json:"id" binding:"required"`                  // 通知ID
	Name               string               `json:"name" binding:"required"`                // 通知名称
	NotificationPolicy string               `json:"notification_policy" binding:"required"` // 通知策略
	SignatureKey       string               `json:"signature_key"`                          // 签名秘钥
	RobotURL           string               `json:"robot_url" binding:"required"`           // 机器人地址
	SendDailyStats     bool                 `json:"send_daily_stats"`                       // 是否发送每日统计
	CardContent        im.CardContentConfig `json:"card_content,omitempty"`                 // 卡片内容配置
}

// UpdateFeiShuRequest 更新飞书通知请求结构
type UpdateFeiShuRequest struct {
	ID                 uint                 `json:"id" binding:"required"`                  // 通知ID
	Name               string               `json:"name" binding:"required"`                // 通知名称
	NotificationPolicy string               `json:"notification_policy" binding:"required"` // 通知策略
	RobotURL           string               `json:"robot_url" binding:"required"`           // 机器人地址
	SendDailyStats     bool                 `json:"send_daily_stats"`                       // 是否发送每日统计
	CardContent        im.CardContentConfig `json:"card_content,omitempty"`                 // 卡片内容配置
}

// TestNotificationRequest 测试通知请求结构
type TestNotificationRequest struct {
	ID      uint   `json:"id" binding:"required"`      // 通知ID
	Type    string `json:"type" binding:"required"`    // 通知类型
	Message string `json:"message" binding:"required"` // 测试消息内容
}

// SearchNotificationParams 通知配置查询参数
type SearchNotificationParams struct {
	Name               string `json:"name" form:"name"`                               // 通知配置名称
	Type               string `json:"type" form:"type"`                               // 通知类型：dingtalk/feishu
	NotificationPolicy string `json:"notification_policy" form:"notification_policy"` // 通知策略
	request.PageInfo
	OrderKey string `json:"orderKey" form:"orderKey"` // 排序字段
	Desc     bool   `json:"desc" form:"desc"`         // 排序方式:升序false(默认)|降序true
}
