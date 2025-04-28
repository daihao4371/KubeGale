package im

// DingTalkConfig 钉钉通知配置结构体
// 继承自 NotificationConfig，并增加了钉钉特有的字段，如签名秘钥、机器人地址等。
type DingTalkConfig struct {
	NotificationConfig
	SignatureKey string `gorm:"type:varchar(255);not null" json:"signature_key"` // SignatureKey 钉钉机器人的签名秘钥，用于消息的签名验证
	RobotURL     string `gorm:"type:varchar(255);not null" json:"robot_url"`     // RobotURL 钉钉机器人接收消息的地址
}

// TableName 设置表名
func (DingTalkConfig) TableName() string {
	return "im_ding_talk_configs"
}
