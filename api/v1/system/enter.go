package system

import "KubeGale/service"

type ApiGroup struct {
	BaseApi
	SystemApiApi
	AuthorityMenuApi
	OperationRecordApi
	AuthorityApi
	RoleApi
}

var (
	apiService             = service.ServiceGroupApp.SystemServiceGroup.ApiService
	userService            = service.ServiceGroupApp.SystemServiceGroup.UserService
	operationRecordService = service.ServiceGroupApp.SystemServiceGroup.OperationRecordService
	menuService            = service.ServiceGroupApp.SystemServiceGroup.MenuService
	authorityService       = service.ServiceGroupApp.SystemServiceGroup.AuthorityService
	roleService            = service.ServiceGroupApp.SystemServiceGroup.RoleService
)
