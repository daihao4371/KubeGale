package system

import (
	"KubeGale/common"
	sysModel "KubeGale/model/system"
	"KubeGale/utils"
	"context"
	"github.com/gofrs/uuid/v5"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type InitUser struct{}

func (i *InitUser) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, common.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(&sysModel.SysUser{})
}

func (i *InitUser) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&sysModel.SysUser{})
}

func (i InitUser) InitializerName() string {
	return sysModel.SysUser{}.TableName()
}

func (i *InitUser) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, common.ErrMissingDBContext
	}

	ap := ctx.Value("adminPassword")
	apStr, ok := ap.(string)
	if !ok {
		apStr = "123456"
	}

	adminPassword := utils.BcryptHash(apStr)

	entities := []sysModel.SysUser{
		{
			UUID:        uuid.Must(uuid.NewV4()),
			Username:    "admin",
			Password:    adminPassword,
			NickName:    "花海",
			HeaderImg:   "https://pic.cnblogs.com/avatar/2399534/20220419203643.png",
			AuthorityId: 888,
			Phone:       "13999999999",
			Email:       "888@qq.com",
		},
	}
	if err = db.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, sysModel.SysUser{}.TableName()+"表数据初始化失败!")
	}
	next = context.WithValue(ctx, i.InitializerName(), entities)
	authorityEntities, ok := ctx.Value(InitAuthority{}.InitializerName()).([]sysModel.SysAuthority)
	if !ok {
		return next, errors.Wrap(common.ErrMissingDependentContext, "创建 [用户-权限] 关联失败, 未找到权限表初始化数据")
	}
	if err = db.Model(&entities[0]).Association("Authorities").Replace(authorityEntities); err != nil {
		return next, err
	}

	// 多个用户注册的时候才需要开启
	/*	if err = db.Model(&entities[1]).Association("Authorities").Replace(authorityEntities[:1]); err != nil {
		return next, err
	}*/
	return next, err
}
