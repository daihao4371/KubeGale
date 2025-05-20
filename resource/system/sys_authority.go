package system

import (
	"KubeGale/common"
	sysModel "KubeGale/model/system"
	"KubeGale/utils"
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderAuthority = initOrderCasbin + 1

type InitAuthority struct{}

func (i *InitAuthority) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, common.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(&sysModel.SysAuthority{})
}

func (i *InitAuthority) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&sysModel.SysAuthority{})
}

func (i InitAuthority) InitializerName() string {
	return sysModel.SysAuthority{}.TableName()
}

func (i *InitAuthority) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, common.ErrMissingDBContext
	}

	// 获取现有的角色列表
	var existingAuthorities []sysModel.SysAuthority
	if err := db.Find(&existingAuthorities).Error; err != nil {
		return ctx, errors.Wrap(err, "获取现有角色列表失败")
	}

	// 创建角色映射，用于快速查找
	authorityMap := make(map[uint]sysModel.SysAuthority)
	for _, authority := range existingAuthorities {
		authorityMap[authority.AuthorityId] = authority
	}

	// 定义需要初始化的角色
	authorities := []sysModel.SysAuthority{
		{AuthorityId: 888, AuthorityName: "超管", ParentId: utils.Pointer[uint](0), DefaultRouter: "dashboard"},
		{AuthorityId: 9528, AuthorityName: "开发负责人", ParentId: utils.Pointer[uint](0), DefaultRouter: "dashboard"},
		{AuthorityId: 8881, AuthorityName: "运维", ParentId: utils.Pointer[uint](888), DefaultRouter: "dashboard"},
	}

	// 用于存储需要新增的角色
	var newAuthorities []sysModel.SysAuthority

	// 检查每个角色是否需要新增
	for _, authority := range authorities {
		if _, exists := authorityMap[authority.AuthorityId]; !exists {
			newAuthorities = append(newAuthorities, authority)
		}
	}

	// 如果有新的角色需要添加
	if len(newAuthorities) > 0 {
		if err := db.Create(&newAuthorities).Error; err != nil {
			return ctx, errors.Wrapf(err, "%s表数据初始化失败!", sysModel.SysAuthority{}.TableName())
		}

		// 为新角色设置数据权限
		for _, authority := range newAuthorities {
			var dataAuthorities []*sysModel.SysAuthority
			switch authority.AuthorityId {
			case 888: // admin角色可以访问所有数据
				dataAuthorities = []*sysModel.SysAuthority{
					{AuthorityId: 888},
					{AuthorityId: 9528},
					{AuthorityId: 8881},
				}
			case 9528: // 开发负责人角色可以访问自己和运维的数据
				dataAuthorities = []*sysModel.SysAuthority{
					{AuthorityId: 9528},
					{AuthorityId: 8881},
				}
			case 8881: // 运维角色只能访问自己的数据
				dataAuthorities = []*sysModel.SysAuthority{
					{AuthorityId: 8881},
				}
			}

			if err := db.Model(&authority).Association("DataAuthorityId").Replace(dataAuthorities); err != nil {
				return ctx, errors.Wrapf(err, "%s表数据初始化失败!",
					db.Model(&authority).Association("DataAuthorityId").Relationship.JoinTable.Name)
			}
		}
	}

	// 合并现有角色和新角色
	allAuthorities := append(existingAuthorities, newAuthorities...)
	next := context.WithValue(ctx, i.InitializerName(), allAuthorities)
	return next, nil
}
