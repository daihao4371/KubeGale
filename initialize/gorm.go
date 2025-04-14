package initialize

import (
	"KubeGale/global"
	"KubeGale/model/cmdb"
	"KubeGale/model/system"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
)

func Gorm() *gorm.DB {
	global.KUBEGALE_ACTIVE_DBNAME = &global.KUBEGALE_CONFIG.Mysql.Dbname
	return GormMysql()
}

func RegisterTables() {
	db := global.KUBEGALE_DB
	err := db.AutoMigrate(
		system.SysApi{},
		system.SysIgnoreApi{},
		system.SysUser{},
		system.SysBaseMenu{},
		system.JwtBlacklist{},
		system.SysAuthority{},
		system.SysOperationRecord{},
		system.SysBaseMenuParameter{},
		system.SysBaseMenuBtn{},
		system.SysAuthorityBtn{},

		// 资产管理平台
		cmdb.CmdbProjects{},
		cmdb.CmdbHosts{},
	)
	if err != nil {
		global.KUBEGALE_LOG.Error("register table failed", zap.Error(err))
		os.Exit(0)
	}

	//err = bizModel()

	if err != nil {
		global.KUBEGALE_LOG.Error("register biz_table failed", zap.Error(err))
		os.Exit(0)
	}
	global.KUBEGALE_LOG.Info("register table success")
}
