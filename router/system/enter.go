package system

import (
	api "KubeGale/api/v1"
)

type RouterGroup struct {
	UserRouter
}

var (
	baseApi = api.ApiGroupApp.SystemApiGroup.BaseApi
)
