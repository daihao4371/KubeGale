package im

import "KubeGale/global"

// FeiShuConfig 飞书通知配置结构体
type FeiShuConfig struct {
	global.KUBEGALE_MODEL                                                  // Own ID for the im_fei_shu_configs table
	NotificationConfigID uint               `gorm:"not null;uniqueIndex"`     // Foreign key to im_notification_configs.id
	NotificationConfig   NotificationConfig `gorm:"foreignKey:NotificationConfigID;references:ID"` // Defines the relationship
	RobotURL             string             `gorm:"type:varchar(255);not null" json:"robot_url"`   // RobotURL 飞书机器人的接收地址
}

// TableName 设置表名
func (FeiShuConfig) TableName() string {
	return "im_fei_shu_configs"
}
