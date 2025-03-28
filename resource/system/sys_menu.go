package system

import (
	"KubeGale/global"
	"KubeGale/model/system"
	"gorm.io/gorm"
	"log"
)

type MenuMock struct {
	db *gorm.DB
}

func NewMenuMock(db *gorm.DB) *MenuMock {
	return &MenuMock{db: db}
}

func (m *MenuMock) InitMenu() error {
	// 检查是否已经初始化过菜单
	var count int64
	m.db.Model(&system.Menu{}).Count(&count)
	if count > 0 {
		log.Println("[菜单已经初始化过,跳过Mock]")
		return nil
	}
	log.Println("[菜单 Mock开始]")
	menus := []system.Menu{
		{
			ID:        1,
			Name:      "Dashboard",
			Path:      "/",
			Component: "BasicLayout",
			Meta: system.MetaField{
				Order: -1,
				Title: "page.dashboard.title",
			},
		},
		{
			ID:        2,
			Name:      "Welcome",
			Path:      "/system_welcome",
			Component: "/dashboard/SystemWelcome",
			ParentID:  1,
			Meta: system.MetaField{
				AffixTab: true,
				Icon:     "lucide:area-chart",
				Title:    "欢迎页",
			},
		},
		{
			ID:        3,
			Name:      "用户管理",
			Path:      "/system_user",
			Component: "/dashboard/SystemUser",
			ParentID:  1,
			Meta: system.MetaField{
				Icon:  "lucide:user",
				Title: "用户管理",
			},
		},
		{
			ID:        4,
			Name:      "菜单管理",
			Path:      "/system_menu",
			Component: "/dashboard/SystemMenu",
			ParentID:  1,
			Meta: system.MetaField{
				Icon:  "lucide:menu",
				Title: "菜单管理",
			},
		},
		{
			ID:        5,
			Name:      "接口管理",
			Path:      "/system_api",
			Component: "/dashboard/SystemApi",
			ParentID:  1,
			Meta: system.MetaField{
				Icon:  "lucide:zap",
				Title: "接口管理",
			},
		},
		{
			ID:        6,
			Name:      "角色权限",
			Path:      "/system_role",
			Component: "/dashboard/SystemRole",
			ParentID:  1,
			Meta: system.MetaField{
				Icon:  "lucide:users",
				Title: "角色权限",
			},
		},
	}

	// 批量创建菜单记录
	for _, menu := range menus {
		// 使用FirstOrCreate方法，如果记录存在则跳过，不存在则创建
		result := m.db.Where("id = ?", menu.ID).FirstOrCreate(&menu)
		if result.Error != nil {
			if global.KUBEGALE_LOG != nil {
				global.KUBEGALE_LOG.Error("创建菜单失败: " + result.Error.Error())
			} else {
				log.Printf("创建菜单失败: %v", result.Error)
			}
			return result.Error
		}

		if result.RowsAffected == 1 {
			log.Printf("创建菜单 [%s] 成功", menu.Name)
		} else {
			log.Printf("菜单 [%s] 已存在，跳过创建", menu.Name)
		}
	}

	log.Println("[菜单 Mock结束]")
	return nil
}

// InitMenuData 提供一个可以在main中直接调用的函数
func InitMenuData() {
	// 检查数据库连接是否已初始化
	if global.KUBEGALE_DB == nil {
		log.Println("错误: 数据库连接未初始化，无法初始化菜单数据")
		return
	}
	
	menuMock := NewMenuMock(global.KUBEGALE_DB)
	if err := menuMock.InitMenu(); err != nil {
		if global.KUBEGALE_LOG != nil {
			global.KUBEGALE_LOG.Error("初始化菜单数据失败: " + err.Error())
		} else {
			log.Printf("初始化菜单数据失败: %v", err)
		}
	}
}
