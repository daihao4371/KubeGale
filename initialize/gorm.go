package initialize

import (
	"KubeGale/global"
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

	//system.SysApi{},
	//system.SysIgnoreApi{},
	//system.SysUser{},
	//system.SysBaseMenu{},
	//system.JwtBlacklist{},
	//system.SysAuthority{},
	//system.SysDictionary{},
	//system.SysOperationRecord{},
	//system.SysDictionaryDetail{},
	//system.SysBaseMenuParameter{},
	//system.SysBaseMenuBtn{},
	//system.SysAuthorityBtn{},
	//system.SysExportTemplate{},
	//system.Condition{},
	//system.JoinTemplate{},
	//
	//example.ExaFile{},
	//example.ExaCustomer{},
	//example.ExaFileChunk{},
	//example.ExaFileUploadAndDownload{},
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
