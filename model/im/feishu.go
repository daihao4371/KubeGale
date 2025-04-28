package im

// FeiShuConfig 飞书通知配置结构体
// 继承自 NotificationConfig，并增加了飞书特有的字段，如机器人地址等。
type FeiShuConfig struct {
	NotificationConfig
	RobotURL string `gorm:"type:varchar(255);not null" json:"robot_url"` // RobotURL 飞书机器人的接收地址
}

// TableName 设置表名
func (FeiShuConfig) TableName() string {
	return "im_fei_shu_configs"
}
