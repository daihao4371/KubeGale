package system

import (
	"KubeGale/common"
	sysModel "KubeGale/model/system"
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type InitApiIgnore struct{}

const initOrderApiIgnore = initOrderApi + 1

func (i InitApiIgnore) InitializerName() string {
	return sysModel.SysIgnoreApi{}.TableName()
}

func (i *InitApiIgnore) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, common.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(&sysModel.SysIgnoreApi{})
}

func (i *InitApiIgnore) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&sysModel.SysIgnoreApi{})
}

// 添加InitializeData方法，初始化需要忽略的API
func (i *InitApiIgnore) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, common.ErrMissingDBContext
	}
	entities := []sysModel.SysIgnoreApi{
		{Method: "GET", Path: "/api/swagger/*any"},
		{Method: "GET", Path: "/api/api/freshCasbin"},
		{Method: "GET", Path: "/health"},
		{Method: "POST", Path: "/api/base/login"},
		{Method: "POST", Path: "/api/base/captcha"},
	}
	if err := db.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, sysModel.SysIgnoreApi{}.TableName()+"表数据初始化失败!")
	}
	next := context.WithValue(ctx, i.InitializerName(), entities)
	return next, nil
}
