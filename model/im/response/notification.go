package response

import (
	"KubeGale/model/im"
	"time"
)

// NotificationResponse 通知配置响应结构体
type NotificationResponse struct {
	ID                 uint      `json:"id"`                  // 通知ID
	Name               string    `json:"name"`                // 通知名称
	Type               string    `json:"type"`                // 通知类型
	NotificationPolicy string    `json:"notification_policy"` // 通知策略
	RobotURL           string    `json:"robot_url"`           // 机器人地址
	CreatedAt          time.Time `json:"created_at"`          // 创建时间
}

// NotificationListResponse 通知配置列表响应结构
type NotificationListResponse struct {
	Notifications []interface{} `json:"notifications"` // 通知配置列表
}

// FeiShuResponse 飞书通知响应结构
type FeiShuResponse struct {
	Config      im.FeiShuConfig      `json:"config"`       // 飞书配置
	CardContent im.CardContentConfig `json:"card_content"` // 卡片内容配置
}

// NotificationDetailConfig 通知配置详情
type NotificationDetailConfig struct {
	ID                 uint      `json:"id"`                  // 通知ID
	Name               string    `json:"name"`                // 通知名称
	Type               string    `json:"type"`                // 通知类型
	NotificationPolicy string    `json:"notification_policy"` // 通知策略
	SendDailyStats     bool      `json:"send_daily_stats"`    // 是否发送每日统计
	CreatedAt          time.Time `json:"created_at"`          // 创建时间
	UpdatedAt          time.Time `json:"updated_at"`          // 更新时间
	RobotURL           string    `json:"robot_url"`           // 机器人地址
}

// NotificationDetailResponse 通知配置详情响应
type NotificationDetailResponse struct {
	Config      NotificationDetailConfig `json:"config"`       // 通知配置
	CardContent CardContentDetail        `json:"card_content"` // 卡片内容
}

// CardContentDetail 卡片内容详情
type CardContentDetail struct {
	ID                 uint      `json:"id"`                  // 卡片内容ID
	NotificationID     uint      `json:"notification_id"`     // 关联的通知ID
	AlertLevel         string    `json:"alert_level"`         // 告警等级
	AlertName          string    `json:"alert_name"`          // 告警名称
	NotificationPolicy string    `json:"notification_policy"` // 通知策略
	AlertContent       string    `json:"alert_content"`       // 告警内容
	AlertTime          time.Time `json:"alert_time"`          // 告警时间
	NotifiedUsers      string    `json:"notified_users"`      // 通知用户
	LastSimilarAlert   string    `json:"last_similar_alert"`  // 上次相似告警
	AlertHandler       string    `json:"alert_handler"`       // 告警处理人
	ClaimAlert         bool      `json:"claim_alert"`         // 是否认领告警
	ResolveAlert       bool      `json:"resolve_alert"`       // 是否解决告警
	MuteAlert          bool      `json:"mute_alert"`          // 是否屏蔽告警
	UnresolvedAlert    bool      `json:"unresolved_alert"`    // 是否为未解决告警
}

// TestNotificationResponse 测试通知响应
type TestNotificationResponse struct {
	Success bool   `json:"success"` // 是否成功
	Message string `json:"message"` // 响应消息
}
