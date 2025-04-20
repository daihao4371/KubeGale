package system

import (
	"KubeGale/common"
	"context"
	adapter "github.com/casbin/gorm-adapter/v3"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderCasbin = initOrderApiIgnore + 1

type initCasbin struct{}

func NewInitCasbin() *initCasbin {
	return &initCasbin{}
}

func (i *initCasbin) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, common.ErrMissingDBContext
	}

	// 先创建基本表结构
	err := db.AutoMigrate(&adapter.CasbinRule{})
	if err != nil {
		return ctx, err
	}

	// 添加唯一索引和设置字符集为utf8mb4
	err = db.Exec("ALTER TABLE casbin_rule ADD UNIQUE INDEX idx_casbin_rule (ptype, v0, v1, v2, v3, v4, v5)").Error
	if err != nil {
		return ctx, err
	}

	// 修改表字符集为utf8mb4
	err = db.Exec("ALTER TABLE casbin_rule CONVERT TO CHARACTER SET utf8mb4").Error
	if err != nil {
		return ctx, err
	}

	// 设置AUTO_INCREMENT起始值
	err = db.Exec("ALTER TABLE casbin_rule AUTO_INCREMENT = 471").Error
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func (i *initCasbin) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&adapter.CasbinRule{})
}

func (i initCasbin) InitializerName() string {
	var entity adapter.CasbinRule
	return entity.TableName()
}

// 修改InitializeData方法中的刷新Casbin权限API路径
func (i *initCasbin) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, common.ErrMissingDBContext
	}

	// 只保留基本权限，其他权限通过管理界面分配
	entities := []adapter.CasbinRule{
		// 超级管理员权限
		{Ptype: "p", V0: "888", V1: "/api/user/admin_register", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/api/createApi", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/api/getApiList", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/api/getApiById", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/api/deleteApi", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/api/updateApi", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/api/getAllApis", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/api/deleteApisByIds", V2: "DELETE"},
		{Ptype: "p", V0: "888", V1: "/api/api/syncApi", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/api/getApiGroups", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/api/enterSyncApi", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/api/ignoreApi", V2: "POST"},

		// 以下是所有角色都需要的基本权限
		{Ptype: "p", V0: "888", V1: "/api/menu/getMenu", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/jwt/jsonInBlacklist", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/user/getUserInfo", V2: "GET"},

		{Ptype: "p", V0: "9528", V1: "/api/menu/getMenu", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/api/jwt/jsonInBlacklist", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/api/user/getUserInfo", V2: "GET"},

		{Ptype: "p", V0: "8881", V1: "/api/menu/getMenu", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/jwt/jsonInBlacklist", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/user/getUserInfo", V2: "GET"},

		// 添加刷新Casbin的权限，仅超级管理员可用
		{Ptype: "p", V0: "888", V1: "/api/api/freshCasbin", V2: "GET"},
	}

	if err := db.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, "Casbin 表 ("+i.InitializerName()+") 数据初始化失败!")
	}
	next := context.WithValue(ctx, i.InitializerName(), entities)
	return next, nil
}

func (i *initCasbin) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where(adapter.CasbinRule{Ptype: "p", V0: "9528", V1: "/user/getUserInfo", V2: "GET"}).
		First(&adapter.CasbinRule{}).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return true
}
