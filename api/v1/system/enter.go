package system

import "KubeGale/service"

type ApiGroup struct {
	JwtApi
	BaseApi
	CasbinApi
	SystemApiApi
	AuthorityApi
	OperationRecordApi
}

var (
	apiService             = service.ServiceGroupApp.SystemServiceGroup.ApiService
	jwtService             = service.ServiceGroupApp.SystemServiceGroup.JwtService
	userService            = service.ServiceGroupApp.SystemServiceGroup.UserService
	casbinService          = service.ServiceGroupApp.SystemServiceGroup.CasbinService
	authorityService       = service.ServiceGroupApp.SystemServiceGroup.AuthorityService
	operationRecordService = service.ServiceGroupApp.SystemServiceGroup.OperationRecordService
)
