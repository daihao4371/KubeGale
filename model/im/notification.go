package im

import (
	"KubeGale/global"
	"time"
)

// NotificationType 通知类型枚举
const (
	NotificationTypeDingTalk = "ding_talk" // 钉钉通知类型
	NotificationTypeFeiShu   = "fei_shu"   // 飞书通知类型
)

// NotificationConfig 通用通知配置结构体
// 包含了所有通知类型共有的字段，如名称、地址、类型、通知策略、创建时间等。
type NotificationConfig struct {
	global.KUBEGALE_MODEL
	Name               string    `gorm:"type:varchar(255);not null;unique" json:"name"`         // Name 通知配置名称
	Type               string    `gorm:"type:varchar(50);not null" json:"type"`                 // Type 通知的类型，如 "ding_talk" 或 "fei_shu"
	NotificationPolicy string    `gorm:"type:varchar(255);not null" json:"notification_policy"` // NotificationPolicy 通知策略，比如告警触发时如何通知
	SendDailyStats     bool      `json:"send_daily_stats"`                                      // SendDailyStats 是否发送每日统计
	CreatedAt          time.Time `json:"created_at"`                                            // CreatedAt 创建时间
	UpdatedAt          time.Time `json:"updated_at"`                                            // UpdatedAt 更新时间
}

// TableName 设置表名
func (NotificationConfig) TableName() string {
	return "im_notification_configs"
}
