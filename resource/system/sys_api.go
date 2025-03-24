package system

import (
	"KubeGale/global"
	"KubeGale/model/system"
	"gorm.io/gorm"
	"log"
)

type ApiMock struct {
	db *gorm.DB
}

func NewApiMock(db *gorm.DB) *ApiMock {
	return &ApiMock{db: db}
}

func (m *ApiMock) InitApi() error {
	// 检查是否已经初始化过API
	var count int64
	m.db.Model(&system.Api{}).Count(&count)
	if count > 0 {
		log.Println("[API已经初始化过,跳过Mock]")
		return nil
	}

	log.Println("[API Mock开始]")

	apis := []system.Api{
		// 基础权限
		{ID: 1, Path: "/*", Method: 1, Name: "所有接口GET权限", Description: "所有接口GET权限", Version: "v1", Category: 1, IsPublic: 1},
		{ID: 2, Path: "/*", Method: 2, Name: "所有接口POST权限", Description: "所有接口POST权限", Version: "v1", Category: 1, IsPublic: 1},
		{ID: 3, Path: "/*", Method: 3, Name: "所有接口PUT权限", Description: "所有接口PUT权限", Version: "v1", Category: 1, IsPublic: 1},
		{ID: 4, Path: "/*", Method: 4, Name: "所有接口DELETE权限", Description: "所有接口DELETE权限", Version: "v1", Category: 1, IsPublic: 1},

		// 用户相关接口
		{ID: 5, Path: "/api/user/logout", Method: 2, Name: "用户登出", Description: "用户退出登录", Version: "v1", Category: 1, IsPublic: 1},
		{ID: 6, Path: "/api/user/codes", Method: 1, Name: "获取权限码", Description: "获取用户权限码", Version: "v1", Category: 1, IsPublic: 0},
		{ID: 7, Path: "/api/user/list", Method: 1, Name: "用户列表", Description: "获取用户列表", Version: "v1", Category: 1, IsPublic: 0},
		{ID: 8, Path: "/api/user/profile", Method: 1, Name: "获取用户信息", Description: "获取用户信息", Version: "v1", Category: 1, IsPublic: 0},
		{ID: 9, Path: "/api/user/signup", Method: 2, Name: "用户注册", Description: "用户注册", Version: "v1", Category: 1, IsPublic: 1},
		{ID: 10, Path: "/api/user/login", Method: 2, Name: "用户登录", Description: "用户登录", Version: "v1", Category: 1, IsPublic: 1},
		{ID: 11, Path: "/api/user/refresh_token", Method: 2, Name: "刷新令牌", Description: "刷新用户令牌", Version: "v1", Category: 1, IsPublic: 1},
		{ID: 12, Path: "/api/user/change_password", Method: 2, Name: "修改密码", Description: "修改密码", Version: "v1", Category: 1, IsPublic: 0},
		{ID: 13, Path: "/api/user/write_off", Method: 2, Name: "注销账号", Description: "注销用户账号", Version: "v1", Category: 1, IsPublic: 0},
		{ID: 14, Path: "/api/user/profile/update", Method: 2, Name: "更新用户信息", Description: "更新用户信息", Version: "v1", Category: 1, IsPublic: 0},
		{ID: 15, Path: "/api/user/:id", Method: 4, Name: "删除用户", Description: "删除用户", Version: "v1", Category: 1, IsPublic: 0},

		// 菜单相关接口
		{ID: 16, Path: "/api/menus/list", Method: 2, Name: "获取菜单列表", Description: "获取菜单列表", Version: "v1", Category: 2, IsPublic: 0},
		{ID: 17, Path: "/api/menus/create", Method: 2, Name: "创建菜单", Description: "创建菜单", Version: "v1", Category: 2, IsPublic: 0},
		{ID: 18, Path: "/api/menus/update", Method: 2, Name: "更新菜单", Description: "更新菜单", Version: "v1", Category: 2, IsPublic: 0},
		{ID: 19, Path: "/api/menus/:id", Method: 4, Name: "删除菜单", Description: "删除菜单", Version: "v1", Category: 2, IsPublic: 0},
		{ID: 20, Path: "/api/menus/update_related", Method: 2, Name: "更新用户菜单关联", Description: "更新用户菜单关联", Version: "v1", Category: 2, IsPublic: 0},

		// API管理相关接口
		{ID: 21, Path: "/api/apis/list", Method: 1, Name: "获取API列表", Description: "获取API列表", Version: "v1", Category: 3, IsPublic: 0},
		{ID: 22, Path: "/api/apis/create", Method: 2, Name: "创建API", Description: "创建API", Version: "v1", Category: 3, IsPublic: 0},
		{ID: 23, Path: "/api/apis/update", Method: 2, Name: "更新API", Description: "更新API", Version: "v1", Category: 3, IsPublic: 0},
		{ID: 24, Path: "/api/apis/:id", Method: 4, Name: "删除API", Description: "删除API", Version: "v1", Category: 3, IsPublic: 0},

		// 角色相关接口
		{ID: 25, Path: "/api/roles/list", Method: 2, Name: "获取角色列表", Description: "获取角色列表", Version: "v1", Category: 4, IsPublic: 0},
		{ID: 26, Path: "/api/roles/create", Method: 2, Name: "创建角色", Description: "创建角色", Version: "v1", Category: 4, IsPublic: 0},
		{ID: 27, Path: "/api/roles/update", Method: 2, Name: "更新角色", Description: "更新角色", Version: "v1", Category: 4, IsPublic: 0},
		{ID: 28, Path: "/api/roles/:id", Method: 4, Name: "删除角色", Description: "删除角色", Version: "v1", Category: 4, IsPublic: 0},
		{ID: 29, Path: "/api/roles/:id", Method: 1, Name: "获取角色详情", Description: "获取角色详情", Version: "v1", Category: 4, IsPublic: 0},
		{ID: 30, Path: "/api/roles/user/:id", Method: 1, Name: "获取用户角色", Description: "获取用户角色", Version: "v1", Category: 4, IsPublic: 0},

		// 权限相关接口
		{ID: 31, Path: "/api/permissions/user/assign", Method: 2, Name: "分配用户角色", Description: "分配用户角色", Version: "v1", Category: 5, IsPublic: 0},
		{ID: 32, Path: "/api/permissions/users/assign", Method: 2, Name: "批量分配用户角色", Description: "批量分配用户角色", Version: "v1", Category: 5, IsPublic: 0},

		// 操作记录相关接口
		{ID: 33, Path: "/api/sysOperationRecord/createSysOperationRecord", Method: 2, Name: "创建操作记录", Description: "创建系统操作记录", Version: "v1", Category: 6, IsPublic: 0},
		{ID: 34, Path: "/api/sysOperationRecord/deleteSysOperationRecord", Method: 4, Name: "删除操作记录", Description: "删除系统操作记录", Version: "v1", Category: 6, IsPublic: 0},
		{ID: 35, Path: "/api/sysOperationRecord/deleteSysOperationRecordByIds", Method: 4, Name: "批量删除操作记录", Description: "批量删除系统操作记录", Version: "v1", Category: 6, IsPublic: 0},
		{ID: 36, Path: "/api/sysOperationRecord/findSysOperationRecord", Method: 1, Name: "查找操作记录", Description: "查找系统操作记录", Version: "v1", Category: 6, IsPublic: 0},
		{ID: 37, Path: "/api/sysOperationRecord/getSysOperationRecordList", Method: 1, Name: "获取操作记录列表", Description: "获取系统操作记录列表", Version: "v1", Category: 6, IsPublic: 0},

		// 系统健康检查
		{ID: 38, Path: "/health", Method: 1, Name: "健康检查", Description: "系统健康状态检查", Version: "v1", Category: 7, IsPublic: 1},
	}

	// 批量创建API记录
	for _, api := range apis {
		// 使用FirstOrCreate方法，如果记录存在则跳过，不存在则创建
		result := m.db.Where("id = ?", api.ID).FirstOrCreate(&api)
		if result.Error != nil {
			global.KUBEGALE_LOG.Error("创建API失败: " + result.Error.Error())
			return result.Error
		}

		if result.RowsAffected == 1 {
			log.Printf("创建API [%s] 成功", api.Name)
		} else {
			log.Printf("API [%s] 已存在，跳过创建", api.Name)
		}
	}

	log.Println("[API Mock结束]")
	return nil
}

// 添加初始化函数，方便在main函数中调用
func InitApiData() {
	// 检查数据库连接是否已初始化
	if global.KUBEGALE_DB == nil {
		log.Println("错误: 数据库连接未初始化，无法初始化API数据")
		return
	}
	
	apiMock := NewApiMock(global.KUBEGALE_DB)
	if err := apiMock.InitApi(); err != nil {
		if global.KUBEGALE_LOG != nil {
			global.KUBEGALE_LOG.Error("初始化API数据失败: " + err.Error())
		} else {
			log.Printf("初始化API数据失败: %v", err)
		}
	}
}
