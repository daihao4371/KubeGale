package initialize

import (
	"KubeGale/global"
	"KubeGale/model/cmdb"
	"KubeGale/model/im"
)

func bizModel() error {
	db := global.KUBEGALE_DB
	err := db.AutoMigrate(
		// CMDB相关表
		cmdb.CmdbProjects{},
		cmdb.CmdbHosts{},

		// IM通知相关表
		im.NotificationConfig{}, // 基础通知配置表
		im.DingTalkConfig{},     // 钉钉通知配置表
		im.FeiShuConfig{},       // 飞书通知配置表
		im.CardContentConfig{},  // 告警卡片内容配置表
	)
	if err != nil {
		return err
	}
	return nil
}
