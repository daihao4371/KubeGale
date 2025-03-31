package middleware

import (
	"KubeGale/global"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func InitCasbin(db *gorm.DB) *casbin.Enforcer {
	// 检查数据库连接是否为空
	if db == nil {
		global.KUBEGALE_LOG.Warn("数据库连接为空，Casbin 初始化跳过")
		return nil
	}

	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		global.KUBEGALE_LOG.Error("Casbin适配器初始化失败", zap.Error(err))
		return nil
	}
	
	e, err := casbin.NewEnforcer("model.conf", adapter)
	if err != nil {
		global.KUBEGALE_LOG.Error("Casbin初始化失败", zap.Error(err))
		return nil
	}
	
	return e
}
