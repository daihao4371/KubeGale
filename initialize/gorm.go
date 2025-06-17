package initialize

import (
	"KubeGale/global"
	"KubeGale/model/system"
	"os"

	"go.uber.org/zap"
	"gorm.io/gorm"
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
		system.JwtBlacklist{},
		system.SysAuthority{},
		system.SysOperationRecord{},
	)
	if err != nil {
		global.KUBEGALE_LOG.Error("register table failed", zap.Error(err))
		os.Exit(0)
	}

	err = bizModel()

	if err != nil {
		global.KUBEGALE_LOG.Error("register biz_table failed", zap.Error(err))
		os.Exit(0)
	}
	global.KUBEGALE_LOG.Info("register table success")
}
