package system

import (
	api "KubeGale/api/v1"
)

type RouterGroup struct {
	UserRouter
	MenuRouter
	ApiRouter
	OperationRecordRouter
	AuthorityRouter
	RoleRouter
}

var (
	baseApi            = api.ApiGroupApp.SystemApiGroup.BaseApi
	apiRouterApi       = api.ApiGroupApp.SystemApiGroup.SystemApiApi
	authorityApi       = api.ApiGroupApp.SystemApiGroup.AuthorityApi
	roleApi            = api.ApiGroupApp.SystemApiGroup.RoleApi
	authorityMenuApi   = api.ApiGroupApp.SystemApiGroup.AuthorityMenuApi
	operationRecordApi = api.ApiGroupApp.SystemApiGroup.OperationRecordApi
)
