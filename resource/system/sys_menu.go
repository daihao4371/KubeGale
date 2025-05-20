package system

import (
	"KubeGale/common"
	. "KubeGale/model/system"
	"context"
	"fmt"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderMenu = initOrderAuthority + 1

type InitMenu struct{}

func (i InitMenu) InitializerName() string {
	return SysBaseMenu{}.TableName()
}

func (i *InitMenu) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, common.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(
		&SysBaseMenu{},
		&SysBaseMenuParameter{},
		&SysBaseMenuBtn{},
	)
}

func (i *InitMenu) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	m := db.Migrator()
	return m.HasTable(&SysBaseMenu{}) &&
		m.HasTable(&SysBaseMenuParameter{}) &&
		m.HasTable(&SysBaseMenuBtn{})
}

func (i *InitMenu) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, common.ErrMissingDBContext
	}

	// 获取现有的菜单列表
	var existingMenus []SysBaseMenu
	if err := db.Find(&existingMenus).Error; err != nil {
		return ctx, errors.Wrap(err, "获取现有菜单列表失败")
	}

	// 创建菜单映射，用于快速查找
	menuMap := make(map[string]SysBaseMenu)
	for _, menu := range existingMenus {
		key := fmt.Sprintf("%s:%s", menu.Path, menu.Name)
		menuMap[key] = menu
	}

	// 定义需要初始化的菜单
	menus := []SysBaseMenu{
		{MenuLevel: 0, Hidden: false, ParentId: 0, Path: "dashboard", Name: "dashboard", Component: "dashboardIndex.vue", Sort: 1, Meta: Meta{Title: "仪表盘", Icon: "odometer"}},
		{MenuLevel: 0, Hidden: false, ParentId: 0, Path: "admin", Name: "superAdmin", Component: "superAdminIndex.vue", Sort: 3, Meta: Meta{Title: "系统管理", Icon: "user"}},
		{MenuLevel: 0, Hidden: false, ParentId: 3, Path: "authority", Name: "authority", Component: "superAdmin/authority/authority.vue", Sort: 1, Meta: Meta{Title: "角色管理", Icon: "avatar"}},
		{MenuLevel: 0, Hidden: false, ParentId: 3, Path: "menu", Name: "menu", Component: "superAdmin/menu/menu.vue", Sort: 2, Meta: Meta{Title: "菜单管理", Icon: "tickets", KeepAlive: true}},
		{MenuLevel: 0, Hidden: false, ParentId: 3, Path: "api", Name: "api", Component: "superAdmin/api/api.vue", Sort: 3, Meta: Meta{Title: "api管理", Icon: "platform", KeepAlive: true}},
		{MenuLevel: 0, Hidden: false, ParentId: 3, Path: "user", Name: "user", Component: "superAdmin/user/user.vue", Sort: 4, Meta: Meta{Title: "用户管理", Icon: "coordinate"}},
		{MenuLevel: 0, Hidden: false, ParentId: 3, Path: "operation", Name: "operation", Component: "superAdmin/operation/sysOperationRecord.vue", Sort: 6, Meta: Meta{Title: "操作历史", Icon: "pie-chart"}},
	}

	// 用于存储需要新增的菜单
	var newMenus []SysBaseMenu

	// 检查每个菜单是否需要新增
	for _, menu := range menus {
		key := fmt.Sprintf("%s:%s", menu.Path, menu.Name)
		if _, exists := menuMap[key]; !exists {
			newMenus = append(newMenus, menu)
		}
	}

	// 如果有新的菜单需要添加
	if len(newMenus) > 0 {
		if err := db.Create(&newMenus).Error; err != nil {
			return ctx, errors.Wrap(err, SysBaseMenu{}.TableName()+"表数据初始化失败!")
		}
	}

	// 合并现有菜单和新菜单
	allMenus := append(existingMenus, newMenus...)
	next = context.WithValue(ctx, i.InitializerName(), allMenus)
	return next, nil
}

func (i *InitMenu) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where("path = ?", "autoPkg").First(&SysBaseMenu{}).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return true
}
