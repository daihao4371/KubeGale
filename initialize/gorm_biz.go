package initialize

import (
	"KubeGale/global"
	cloudCmdb "KubeGale/model/cloudCmdb"
	cmdb "KubeGale/model/cmdb"
	"KubeGale/model/im"

	"go.uber.org/zap"
)

// RegisterIMTables 注册IM相关数据库表
func bizModel() error {
	db := global.KUBEGALE_DB
	err := db.AutoMigrate(
		im.NotificationConfig{},
		im.FeiShuConfig{},
		im.CardContentConfig{},

		// 资产管理自建机房资源
		cmdb.CmdbHosts{},
		cmdb.CmdbProjects{},

		// 云资源
		cloudCmdb.LoadBalancer{},
		cloudCmdb.CloudPlatform{},
		cloudCmdb.RDS{},
		cloudCmdb.CloudRegions{},
		cloudCmdb.VirtualMachine{},
	)
	if err != nil {
		global.KUBEGALE_LOG.Error("register  tables failed", zap.Error(err))
		return err
	}
	global.KUBEGALE_LOG.Info("register  tables success")
	return nil
}
