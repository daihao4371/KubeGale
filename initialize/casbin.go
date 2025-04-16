package initialize

import (
	"context"
	"errors"

	adapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

type initCasbin struct{}

func NewInitCasbin() *initCasbin {
	return &initCasbin{}
}

var (
	ErrMissingDBContext        = errors.New("missing db in context")
	ErrMissingDependentContext = errors.New("missing dependent value in context")
	ErrDBTypeMismatch          = errors.New("db type mismatch")
)

func (i *initCasbin) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, ErrMissingDBContext
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
