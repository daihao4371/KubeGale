package system

import (
	api "KubeGale/api/v1"
)

type RouterGroup struct {
	ApiRouter
	JwtRouter
	BaseRouter
	MenuRouter
	UserRouter
	CasbinRouter
	AuthorityRouter
	OperationRecordRouter
	AuthorityBtnRouter
}

var (
	jwtApi             = api.ApiGroupApp.SystemApiGroup.JwtApi
	baseApi            = api.ApiGroupApp.SystemApiGroup.BaseApi
	casbinApi          = api.ApiGroupApp.SystemApiGroup.CasbinApi
	authorityApi       = api.ApiGroupApp.SystemApiGroup.AuthorityApi
	apiRouterApi       = api.ApiGroupApp.SystemApiGroup.SystemApiApi
	authorityBtnApi    = api.ApiGroupApp.SystemApiGroup.AuthorityBtnApi
	authorityMenuApi   = api.ApiGroupApp.SystemApiGroup.AuthorityMenuApi
	operationRecordApi = api.ApiGroupApp.SystemApiGroup.OperationRecordApi
)
