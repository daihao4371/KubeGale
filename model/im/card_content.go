package im

import (
	"KubeGale/global"
	"time"
)

// CardContentConfig 告警卡片内容配置结构体
// 存储告警卡片的详细内容配置，包括告警等级、名称、内容、通知人等信息。
type CardContentConfig struct {
	global.KUBEGALE_MODEL
	NotificationID     uint      `json:"notification_id" gorm:"index;not null"`                 // 关联的通知配置ID
	AlertLevel         string    `gorm:"type:varchar(50);not null" json:"alert_level"`          // AlertLevel 告警等级，如 "Critical", "Warning", "Info"
	AlertName          string    `gorm:"type:varchar(255);not null" json:"alert_name"`          // AlertName 告警的名称或标题
	NotificationPolicy string    `gorm:"type:varchar(255);not null" json:"notification_policy"` // NotificationPolicy 告警的通知策略
	AlertContent       string    `gorm:"type:text" json:"alert_content"`                        // AlertContent 告警的详细内容
	AlertTime          time.Time `json:"alert_time"`                                            // AlertTime 告警触发的时间
	NotifiedUsers      []string  `gorm:"type:text;serializer:json" json:"notified_users"`       // 通知用户
	LastSimilarAlert   string    `gorm:"type:text" json:"last_similar_alert"`                   // LastSimilarAlert 上次相似的告警信息
	AlertHandler       string    `gorm:"type:varchar(255);not null" json:"alert_handler"`       // AlertHandler 告警处理人
	ClaimAlert         bool      `json:"claim_alert"`                                           // ClaimAlert 是否认领该告警
	ResolveAlert       bool      `json:"resolve_alert"`                                         // ResolveAlert 是否已经解决该告警
	MuteAlert          bool      `json:"mute_alert"`                                            // MuteAlert 是否屏蔽该告警
	UnresolvedAlert    bool      `json:"unresolved_alert"`                                      // UnresolvedAlert 是否为未解决告警
}

// TableName 设置表名
func (CardContentConfig) TableName() string {
	return "im_card_content_configs"
}
