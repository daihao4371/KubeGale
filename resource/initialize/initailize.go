package initialize

import (
	"KubeGale/global"
	"KubeGale/resource/system"
	"context"

	"go.uber.org/zap"
)

// InitializeAllSystemData 初始化所有系统数据
func InitializeAllSystemData(ctx context.Context) {
	// 初始化 API 表
	apiInitializer := &system.InitApi{}
	ctx, err := apiInitializer.MigrateTable(ctx)
	if err != nil {
		global.KUBEGALE_LOG.Error("初始化 API 表失败", zap.Error(err))
	} else {
		global.KUBEGALE_LOG.Info("初始化 API 表成功")
		ctx, err = apiInitializer.InitializeData(ctx)
		if err != nil {
			global.KUBEGALE_LOG.Error("初始化 API 数据失败", zap.Error(err))
		} else {
			global.KUBEGALE_LOG.Info("初始化 API 数据成功")
		}
	}

	// 初始化 API 忽略表
	apiIgnoreInitializer := &system.InitApiIgnore{}
	ctx, err = apiIgnoreInitializer.MigrateTable(ctx)
	if err != nil {
		global.KUBEGALE_LOG.Error("初始化 API 忽略表失败", zap.Error(err))
	} else {
		global.KUBEGALE_LOG.Info("初始化 API 忽略表成功")
	}

	// 初始化 Casbin 表
	casbinInitializer := system.NewInitCasbin()
	ctx, err = casbinInitializer.MigrateTable(ctx)
	if err != nil {
		global.KUBEGALE_LOG.Error("初始化 Casbin 表失败", zap.Error(err))
	} else {
		global.KUBEGALE_LOG.Info("初始化 Casbin 表成功")
		// 添加这部分代码，调用 InitializeData 方法为 admin 初始化权限
		ctx, err = casbinInitializer.InitializeData(ctx)
		if err != nil {
			global.KUBEGALE_LOG.Error("初始化 Casbin 权限数据失败", zap.Error(err))
		} else {
			global.KUBEGALE_LOG.Info("初始化 Casbin 权限数据成功")
		}
	}

	// 初始化权限表
	authorityInitializer := &system.InitAuthority{}
	ctx, err = authorityInitializer.MigrateTable(ctx)
	if err != nil {
		global.KUBEGALE_LOG.Error("初始化权限表失败", zap.Error(err))
	} else {
		global.KUBEGALE_LOG.Info("初始化权限表成功")
		ctx, err = authorityInitializer.InitializeData(ctx)
		if err != nil {
			global.KUBEGALE_LOG.Error("初始化权限数据失败", zap.Error(err))
		} else {
			global.KUBEGALE_LOG.Info("初始化权限数据成功")
		}
	}

	// 初始化菜单表
	menuInitializer := &system.InitMenu{}
	ctx, err = menuInitializer.MigrateTable(ctx)
	if err != nil {
		global.KUBEGALE_LOG.Error("初始化菜单表失败", zap.Error(err))
	} else {
		global.KUBEGALE_LOG.Info("初始化菜单表成功")
		ctx, err = menuInitializer.InitializeData(ctx)
		if err != nil {
			global.KUBEGALE_LOG.Error("初始化菜单数据失败", zap.Error(err))
		} else {
			global.KUBEGALE_LOG.Info("初始化菜单数据成功")
		}
	}

	// 初始化用户表
	userInitializer := &system.InitUser{}
	ctx, err = userInitializer.MigrateTable(ctx)
	if err != nil {
		global.KUBEGALE_LOG.Error("初始化用户表失败", zap.Error(err))
	} else {
		global.KUBEGALE_LOG.Info("初始化用户表成功")
		// 设置管理员密码
		ctx = context.WithValue(ctx, "adminPassword", "123456")
		ctx, err = userInitializer.InitializeData(ctx)
		if err != nil {
			global.KUBEGALE_LOG.Error("初始化用户数据失败", zap.Error(err))
		} else {
			global.KUBEGALE_LOG.Info("初始化用户数据成功")
		}
	}
}
