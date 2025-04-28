package initialize

import (
	"KubeGale/global"
	"KubeGale/model/cmdb"
	"KubeGale/model/im"
)

func bizModel() error {
	db := global.KUBEGALE_DB
	err := db.AutoMigrate(
		cmdb.CmdbProjects{},
		cmdb.CmdbHosts{},
		im.NotificationConfig{},
		im.DingTalkConfig{},
		im.FeiShuConfig{},
		im.CardContentConfig{},
	)
	if err != nil {
		return err
	}
	return err
}
