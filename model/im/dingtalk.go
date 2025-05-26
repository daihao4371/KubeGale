package im

import "KubeGale/global"

// DingTalkConfig 钉钉通知配置结构体
type DingTalkConfig struct {
	global.KUBEGALE_MODEL                                                      // Own ID for the im_ding_talk_configs table
	NotificationConfigID uint                   `gorm:"not null;uniqueIndex"` // Foreign key to im_notification_configs.id
	NotificationConfig   NotificationConfig     `gorm:"foreignKey:NotificationConfigID;references:ID"` // Defines the relationship
	WebhookURL           string                 `gorm:"type:varchar(255);not null" json:"webhook_url"` // WebhookURL 钉钉机器人的接收地址
	Secret               string                 `gorm:"type:varchar(255)" json:"secret,omitempty"`     // Secret 钉钉机器人加签密钥 (optional)
}

// TableName 设置表名
func (DingTalkConfig) TableName() string {
	return "im_ding_talk_configs"
}
