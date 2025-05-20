package system

import (
	"KubeGale/common"
	sysModel "KubeGale/model/system"
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderApi = common.InitOrderSystem + 1

type InitApi struct{}

func (i InitApi) InitializerName() string {
	return sysModel.SysApi{}.TableName()
}

func (i *InitApi) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, common.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(&sysModel.SysApi{})
}

func (i *InitApi) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&sysModel.SysApi{})
}

func (i *InitApi) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, common.ErrMissingDBContext
	}
	entities := []sysModel.SysApi{
		{ApiGroup: "jwt", Method: "POST", Path: "/jwt/jsonInBlacklist", Description: "jwt加入黑名单(退出，必选)"},

		{ApiGroup: "系统用户", Method: "DELETE", Path: "/user/deleteUser", Description: "删除用户"},
		{ApiGroup: "系统用户", Method: "POST", Path: "/user/admin_register", Description: "用户注册"},
		{ApiGroup: "系统用户", Method: "POST", Path: "/user/getUserList", Description: "获取用户列表"},
		{ApiGroup: "系统用户", Method: "PUT", Path: "/user/setUserInfo", Description: "设置用户信息"},
		{ApiGroup: "系统用户", Method: "PUT", Path: "/user/setSelfInfo", Description: "设置自身信息(必选)"},
		{ApiGroup: "系统用户", Method: "GET", Path: "/user/getUserInfo", Description: "获取自身信息(必选)"},
		{ApiGroup: "系统用户", Method: "POST", Path: "/user/setUserAuthorities", Description: "设置权限组"},
		{ApiGroup: "系统用户", Method: "POST", Path: "/user/changePassword", Description: "修改密码（建议选择)"},
		{ApiGroup: "系统用户", Method: "POST", Path: "/user/setUserAuthority", Description: "修改用户角色(必选)"},
		{ApiGroup: "系统用户", Method: "POST", Path: "/user/resetPassword", Description: "重置用户密码"},
		{ApiGroup: "系统用户", Method: "PUT", Path: "/user/setSelfSetting", Description: "用户界面配置"},

		{ApiGroup: "api", Method: "POST", Path: "/api/createApi", Description: "创建api"},
		{ApiGroup: "api", Method: "POST", Path: "/api/deleteApi", Description: "删除Api"},
		{ApiGroup: "api", Method: "POST", Path: "/api/updateApi", Description: "更新Api"},
		{ApiGroup: "api", Method: "POST", Path: "/api/getApiList", Description: "获取api列表"},
		{ApiGroup: "api", Method: "POST", Path: "/api/getAllApis", Description: "获取所有api"},
		{ApiGroup: "api", Method: "POST", Path: "/api/getApiById", Description: "获取api详细信息"},
		{ApiGroup: "api", Method: "DELETE", Path: "/api/deleteApisByIds", Description: "批量删除api"},
		{ApiGroup: "api", Method: "GET", Path: "/api/syncApi", Description: "获取待同步API"},
		{ApiGroup: "api", Method: "GET", Path: "/api/getApiGroups", Description: "获取路由组"},
		{ApiGroup: "api", Method: "POST", Path: "/api/enterSyncApi", Description: "确认同步API"},
		{ApiGroup: "api", Method: "POST", Path: "/api/ignoreApi", Description: "忽略API"},

		{ApiGroup: "角色", Method: "POST", Path: "/authority/copyAuthority", Description: "拷贝角色"},
		{ApiGroup: "角色", Method: "POST", Path: "/authority/createAuthority", Description: "创建角色"},
		{ApiGroup: "角色", Method: "POST", Path: "/authority/deleteAuthority", Description: "删除角色"},
		{ApiGroup: "角色", Method: "PUT", Path: "/authority/updateAuthority", Description: "更新角色信息"},
		{ApiGroup: "角色", Method: "POST", Path: "/authority/getAuthorityList", Description: "获取角色列表"},
		{ApiGroup: "角色", Method: "POST", Path: "/authority/setDataAuthority", Description: "设置角色资源权限"},

		{ApiGroup: "casbin", Method: "POST", Path: "/casbin/updateCasbin", Description: "更改角色api权限"},
		{ApiGroup: "casbin", Method: "POST", Path: "/casbin/getPolicyPathByAuthorityId", Description: "获取权限列表"},

		{ApiGroup: "菜单", Method: "POST", Path: "/menu/addBaseMenu", Description: "新增菜单"},
		{ApiGroup: "菜单", Method: "POST", Path: "/menu/getMenu", Description: "获取菜单树(必选)"},
		{ApiGroup: "菜单", Method: "POST", Path: "/menu/deleteBaseMenu", Description: "删除菜单"},
		{ApiGroup: "菜单", Method: "POST", Path: "/menu/updateBaseMenu", Description: "更新菜单"},
		{ApiGroup: "菜单", Method: "POST", Path: "/menu/getBaseMenuById", Description: "根据id获取菜单"},
		{ApiGroup: "菜单", Method: "POST", Path: "/menu/getMenuList", Description: "分页获取基础menu列表"},
		{ApiGroup: "菜单", Method: "POST", Path: "/menu/getBaseMenuTree", Description: "获取用户动态路由"},
		{ApiGroup: "菜单", Method: "POST", Path: "/menu/getMenuAuthority", Description: "获取指定角色menu"},
		{ApiGroup: "菜单", Method: "POST", Path: "/menu/addMenuAuthority", Description: "增加menu和角色关联关系"},

		{ApiGroup: "操作记录", Method: "POST", Path: "/sysOperationRecord/createSysOperationRecord", Description: "新增操作记录"},
		{ApiGroup: "操作记录", Method: "GET", Path: "/sysOperationRecord/findSysOperationRecord", Description: "根据ID获取操作记录"},
		{ApiGroup: "操作记录", Method: "GET", Path: "/sysOperationRecord/getSysOperationRecordList", Description: "获取操作记录列表"},
		{ApiGroup: "操作记录", Method: "DELETE", Path: "/sysOperationRecord/deleteSysOperationRecord", Description: "删除操作记录"},
		{ApiGroup: "操作记录", Method: "DELETE", Path: "/sysOperationRecord/deleteSysOperationRecordByIds", Description: "批量删除操作历史"},

		{ApiGroup: "按钮权限", Method: "POST", Path: "/authorityBtn/setAuthorityBtn", Description: "设置按钮权限"},
		{ApiGroup: "按钮权限", Method: "POST", Path: "/authorityBtn/getAuthorityBtn", Description: "获取已有按钮权限"},
		{ApiGroup: "按钮权限", Method: "POST", Path: "/authorityBtn/canRemoveAuthorityBtn", Description: "删除按钮"},

		{ApiGroup: "即时通讯", Method: "POST", Path: "/notification/createFeiShu", Description: "创建飞书通知"},
		{ApiGroup: "即时通讯", Method: "PUT", Path: "/notification/updateFeiShu", Description: "更新飞书通知"},
		{ApiGroup: "即时通讯", Method: "DELETE", Path: "/notification/deleteNotification", Description: "删除通知配置"},
		{ApiGroup: "即时通讯", Method: "POST", Path: "/notification/testNotification", Description: "测试通知发送"},
		{ApiGroup: "即时通讯", Method: "POST", Path: "/notification/createCardContent", Description: "创建卡片内容"},
		{ApiGroup: "即时通讯", Method: "PUT", Path: "/notification/updateCardContent", Description: "更新卡片内容"},
		{ApiGroup: "即时通讯", Method: "POST", Path: "/notification/getNotificationList", Description: "获取通知配置列表"},
		{ApiGroup: "即时通讯", Method: "GET", Path: "/notification/getNotificationById", Description: "根据ID获取通知配置"},
		{ApiGroup: "即时通讯", Method: "GET", Path: "/notification/getCardContent", Description: "根据通知ID获取卡片内容"},

		// CMDB项目管理
		{ApiGroup: "cmdb", Method: "POST", Path: "/cmdb/projects", Description: "新建项目"},
		{ApiGroup: "cmdb", Method: "DELETE", Path: "/cmdb/projects", Description: "删除项目"},
		{ApiGroup: "cmdb", Method: "DELETE", Path: "/cmdb/projectsByIds", Description: "批量删除项目"},
		{ApiGroup: "cmdb", Method: "PUT", Path: "/cmdb/projects", Description: "更新项目"},
		{ApiGroup: "cmdb", Method: "GET", Path: "/cmdb/projectsById", Description: "根据ID获取项目"},
		{ApiGroup: "cmdb", Method: "GET", Path: "/cmdb/projects", Description: "获取项目列表"},

		// CMDB主机管理
		{ApiGroup: "cmdb", Method: "POST", Path: "/cmdb/hosts", Description: "新建主机"},
		{ApiGroup: "cmdb", Method: "DELETE", Path: "/cmdb/hosts", Description: "删除主机"},
		{ApiGroup: "cmdb", Method: "DELETE", Path: "/cmdb/hostsByIds", Description: "批量删除主机"},
		{ApiGroup: "cmdb", Method: "PUT", Path: "/cmdb/hosts", Description: "更新主机"},
		{ApiGroup: "cmdb", Method: "POST", Path: "/cmdb/hosts/authentication", Description: "SSH认证主机"},
		{ApiGroup: "cmdb", Method: "POST", Path: "/cmdb/hosts/import", Description: "导入主机"},
		{ApiGroup: "cmdb", Method: "GET", Path: "/cmdb/hostsById", Description: "根据ID获取主机"},
		{ApiGroup: "cmdb", Method: "GET", Path: "/cmdb/hosts", Description: "获取主机列表"},

		// 批量操作
		{ApiGroup: "cmdb", Method: "POST", Path: "/cmdb/batchOperations/execute", Description: "执行批量操作"},
		{ApiGroup: "cmdb", Method: "GET", Path: "/cmdb/batchOperations/execLogs", Description: "获取执行记录"},

		// 云平台管理
		{ApiGroup: "cloud_platform", Method: "POST", Path: "/cloud_platform/getById", Description: "获取云平台信息"},
		{ApiGroup: "cloud_platform", Method: "POST", Path: "/cloud_platform/create", Description: "创建云平台"},
		{ApiGroup: "cloud_platform", Method: "PUT", Path: "/cloud_platform/update", Description: "更新云平台"},
		{ApiGroup: "cloud_platform", Method: "DELETE", Path: "/cloud_platform/delete", Description: "删除云平台"},
		{ApiGroup: "cloud_platform", Method: "DELETE", Path: "/cloud_platform/deleteByIds", Description: "批量删除云平台"},
		{ApiGroup: "cloud_platform", Method: "POST", Path: "/cloud_platform/list", Description: "获取云平台列表"},

		// 云区域管理
		{ApiGroup: "cloud_region", Method: "POST", Path: "/cloud_region/syncRegion", Description: "同步区域信息"},
		{ApiGroup: "cloud_region", Method: "GET", Path: "/cloud_region/tree", Description: "获取区域树形结构"},

		// 云主机管理
		{ApiGroup: "virtualMachine", Method: "POST", Path: "/virtualMachine/sync", Description: "同步云主机"},
		{ApiGroup: "virtualMachine", Method: "POST", Path: "/virtualMachine/tree", Description: "获取云主机目录树"},
		{ApiGroup: "virtualMachine", Method: "POST", Path: "/virtualMachine/list", Description: "获取云主机列表"},

		// 负载均衡管理
		{ApiGroup: "loadBalancer", Method: "POST", Path: "/loadBalancer/sync", Description: "同步负载均衡"},
		{ApiGroup: "loadBalancer", Method: "POST", Path: "/loadBalancer/tree", Description: "获取负载均衡目录树"},
		{ApiGroup: "loadBalancer", Method: "POST", Path: "/loadBalancer/list", Description: "获取负载均衡列表"},

		// RDS管理
		{ApiGroup: "rds", Method: "POST", Path: "/rds/sync", Description: "同步RDS"},
		{ApiGroup: "rds", Method: "POST", Path: "/rds/tree", Description: "获取RDS目录树"},
		{ApiGroup: "rds", Method: "POST", Path: "/rds/list", Description: "获取RDS列表"},
	}
	if err := db.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, sysModel.SysApi{}.TableName()+"表数据初始化失败!")
	}
	next := context.WithValue(ctx, i.InitializerName(), entities)
	return next, nil
}
