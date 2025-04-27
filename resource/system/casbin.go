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

	// 查询所有 API
	var apis []struct {
		Path   string
		Method string
	}
	if err := db.Table("sys_apis").Select("path, method").Find(&apis).Error; err != nil {
		return ctx, errors.Wrap(err, "查询所有API失败")
	}

	// 为 admin 生成所有 API 权限
	entities := make([]adapter.CasbinRule, 0, len(apis))
	for _, api := range apis {
		entities = append(entities, adapter.CasbinRule{
			Ptype: "p",
			V0:    "888", // admin 角色
			V1:    api.Path,
			V2:    api.Method,
		})
	}

	// 你可以在这里为其他角色分配权限（可选）

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
	// 修改为检查 admin 角色(888)的权限
	if errors.Is(db.Where(adapter.CasbinRule{Ptype: "p", V0: "888", V1: "/user/getUserInfo", V2: "GET"}).
		First(&adapter.CasbinRule{}).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return true
}
