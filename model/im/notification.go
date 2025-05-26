package im

import (
	"KubeGale/global"
)

// NotificationType 通知类型枚举
const (
	NotificationTypeFeiShu   = "feishu"   // 飞书通知类型
	NotificationTypeDingTalk = "dingtalk" // 钉钉通知类型 (added for future use)
)

// NotificationConfig 通用通知配置结构体
type NotificationConfig struct {
	global.KUBEGALE_MODEL                                                              // ID, CreatedAt, UpdatedAt, DeletedAt
	Name               string `gorm:"type:varchar(255);not null;uniqueIndex:idx_name_type_unique,priority:1" json:"name"` // Name 通知配置名称
	Type               string `gorm:"type:varchar(50);not null;uniqueIndex:idx_name_type_unique,priority:2" json:"type"`  // Type 通知的类型 (feishu, dingtalk)
	NotificationPolicy string `gorm:"type:varchar(255)" json:"notification_policy"`                                      // NotificationPolicy 通知策略，比如告警触发时如何通知, comma-separated
	SendDailyStats     bool   `json:"send_daily_stats"`                                                                  // SendDailyStats 是否发送每日统计
}

// TableName 设置表名
func (NotificationConfig) TableName() string {
	return "im_notification_configs"
}
