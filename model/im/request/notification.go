package request

import (
	"KubeGale/model/common/request"
	"KubeGale/model/im"
	"time"
)

// CreateFeiShuRequest 创建飞书通知请求结构
type CreateFeiShuRequest struct {
	Name           string   `json:"name" binding:"required"`        // 通知名称
	Type           string   `json:"type" binding:"required"`        // 通知类型
	Enabled        bool     `json:"enabled"`                        // 是否启用
	WebhookURL     string   `json:"webhook_url" binding:"required"` // 机器人地址
	Description    string   `json:"description"`                    // 描述
	Tags           []string `json:"tags"`                           // 标签
	NotifyEvents   []string `json:"notify_events"`                  // 通知事件
	Receivers      []string `json:"receivers"`                      // 接收者
	SendDailyStats bool     `json:"send_daily_stats"`               // 是否发送每日统计
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
	Type               string `json:"type" form:"type"`                               // 通知类型：feishu
	NotificationPolicy string `json:"notification_policy" form:"notification_policy"` // 通知策略
	request.PageInfo
	OrderKey string `json:"orderKey" form:"orderKey"` // 排序字段
	Desc     bool   `json:"desc" form:"desc"`         // 排序方式:升序false(默认)|降序true
}

// CreateCardContentRequest 创建卡片内容请求结构
type CreateCardContentRequest struct {
	NotificationID     uint      `json:"notification_id" binding:"required"`     // 关联的通知配置ID
	AlertLevel         string    `json:"alert_level" binding:"required"`         // 告警等级
	AlertName          string    `json:"alert_name" binding:"required"`          // 告警名称
	NotificationPolicy string    `json:"notification_policy" binding:"required"` // 通知策略
	AlertContent       string    `json:"alert_content"`                          // 告警内容
	AlertTime          time.Time `json:"alert_time"`                             // 告警时间
	NotifiedUsers      []string  `json:"notified_users"`                         // 通知用户
	LastSimilarAlert   string    `json:"last_similar_alert"`                     // 上次相似告警
	AlertHandler       string    `json:"alert_handler" binding:"required"`       // 告警处理人
	ClaimAlert         bool      `json:"claim_alert"`                            // 是否认领告警
	ResolveAlert       bool      `json:"resolve_alert"`                          // 是否解决告警
	MuteAlert          bool      `json:"mute_alert"`                             // 是否屏蔽告警
	UnresolvedAlert    bool      `json:"unresolved_alert"`                       // 是否未解决告警
}
