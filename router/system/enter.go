package system

import api "KubeGale/api/v1"

type RouterGroup struct {
	ApiRouter
	JwtRouter
	BaseRouter
	UserRouter
	CasbinRouter
	AuthorityRouter
	OperationRecordRouter
	AuthorityApiRouter
}

var (
	jwtApi             = api.ApiGroupApp.SystemApiGroup.JwtApi
	baseApi            = api.ApiGroupApp.SystemApiGroup.BaseApi
	casbinApi          = api.ApiGroupApp.SystemApiGroup.CasbinApi
	authorityApi       = api.ApiGroupApp.SystemApiGroup.AuthorityApi
	apiRouterApi       = api.ApiGroupApp.SystemApiGroup.SystemApiApi
	operationRecordApi = api.ApiGroupApp.SystemApiGroup.OperationRecordApi
	authorityApiApi    = api.ApiGroupApp.SystemApiGroup.AuthorityApiApi
)
