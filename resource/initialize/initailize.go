package initalize

import (
	"KubeGale/global"
	"KubeGale/resource/system"
	"context"
	"go.uber.org/zap"
)

// InitializeAllSystemData 初始化所有系统数据
func InitializeAllSystemData(ctx context.Context) {
	// 初始化用户表
	userInitalizer := &system.InitUser{}
	ctx, err := userInitalizer.MigrateTable(ctx)
	if err != nil {
		global.KUBEGALE_LOG.Error("初始化用户表失败", zap.Error(err))
	} else {
		global.KUBEGALE_LOG.Info("初始化用户表成功")
		// 设置管理员密码
		ctx = context.WithValue(ctx, "adminPassword", "123456")
		ctx, err = userInitalizer.MigrateTable(ctx)
		if err != nil {
			global.KUBEGALE_LOG.Error("初始化用户数据失败", zap.Error(err))
		} else {
			global.KUBEGALE_LOG.Info("初始化用户数据成功")
		}
	}
}
