package system

import "KubeGale/service"

type ApiGroup struct {
	JwtApi
	BaseApi
	CasbinApi
	SystemApiApi
	AuthorityApi
	AuthorityMenuApi
	OperationRecordApi
	AuthorityBtnApi
}

var (
	apiService             = service.ServiceGroupApp.SystemServiceGroup.ApiService
	jwtService             = service.ServiceGroupApp.SystemServiceGroup.JwtService
	menuService            = service.ServiceGroupApp.SystemServiceGroup.MenuService
	userService            = service.ServiceGroupApp.SystemServiceGroup.UserService
	casbinService          = service.ServiceGroupApp.SystemServiceGroup.CasbinService
	baseMenuService        = service.ServiceGroupApp.SystemServiceGroup.BaseMenuService
	authorityService       = service.ServiceGroupApp.SystemServiceGroup.AuthorityService
	authorityBtnService    = service.ServiceGroupApp.SystemServiceGroup.AuthorityBtnService
	operationRecordService = service.ServiceGroupApp.SystemServiceGroup.OperationRecordService
)
