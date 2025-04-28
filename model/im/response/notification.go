package response

import (
	"KubeGale/model/im"
)

// NotificationResponse 通知配置响应结构
type NotificationResponse struct {
	Notification interface{} `json:"notification"` // 通知配置，可能是钉钉或飞书
}

// NotificationListResponse 通知配置列表响应结构
type NotificationListResponse struct {
	Notifications []interface{} `json:"notifications"` // 通知配置列表
}

// DingTalkResponse 钉钉通知响应结构
type DingTalkResponse struct {
	Config      im.DingTalkConfig     `json:"config"`       // 钉钉配置
	CardContent im.CardContentConfig  `json:"card_content"` // 卡片内容配置
}

// FeiShuResponse 飞书通知响应结构
type FeiShuResponse struct {
	Config      im.FeiShuConfig       `json:"config"`       // 飞书配置
	CardContent im.CardContentConfig  `json:"card_content"` // 卡片内容配置
}

// TestNotificationResponse 测试通知响应结构
type TestNotificationResponse struct {
	Success bool   `json:"success"` // 是否成功
	Message string `json:"message"` // 消息
}