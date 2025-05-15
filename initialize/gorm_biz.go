package initialize

import (
	"KubeGale/global"
	"KubeGale/model/im"

	"go.uber.org/zap"
)

// RegisterIMTables 注册IM相关数据库表
func RegisterIMTables() error {
	db := global.KUBEGALE_DB
	err := db.AutoMigrate(
		im.NotificationConfig{},
		im.FeiShuConfig{},
		im.CardContentConfig{},
	)
	if err != nil {
		global.KUBEGALE_LOG.Error("register IM tables failed", zap.Error(err))
		return err
	}
	global.KUBEGALE_LOG.Info("register IM tables success")
	return nil
}
