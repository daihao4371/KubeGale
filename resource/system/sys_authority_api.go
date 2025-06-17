package system

import (
	"KubeGale/global"
	"KubeGale/model/system"
	"context"

	"gorm.io/gorm/clause"
)

// InitAuthorityApi 角色-API权限初始化
// 1. 自动迁移表结构 2. 为 admin 角色分配所有 API 权限

type InitAuthorityApi struct{}

func (i *InitAuthorityApi) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, global.KUBEGALE_DB.AutoMigrate(&system.SysAuthorityApi{})
}

func (i *InitAuthorityApi) InitializeData(ctx context.Context) (context.Context, error) {
	// 查询 admin 角色（与初始化用户 AuthorityId 保持一致，使用 888）
	var admin system.SysAuthority
	if err := global.KUBEGALE_DB.Where("authority_id = ?", 888).First(&admin).Error; err != nil {
		return ctx, err
	}
	// 查询所有 API
	var apis []system.SysApi
	if err := global.KUBEGALE_DB.Find(&apis).Error; err != nil {
		return ctx, err
	}
	// 分配所有 API 权限给 admin
	var records []system.SysAuthorityApi
	for _, api := range apis {
		records = append(records, system.SysAuthorityApi{
			AuthorityId: admin.AuthorityId,
			ApiId:       api.ID,
		})
	}
	if len(records) > 0 {
		if err := global.KUBEGALE_DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&records).Error; err != nil {
			return ctx, err
		}
	}
	return ctx, nil
}
