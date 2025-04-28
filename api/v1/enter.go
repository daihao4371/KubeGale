package v1

import (
	"KubeGale/api/v1/im"
	"KubeGale/api/v1/system"
)

var ApiGroupApp = new(ApiGroup)

type ApiGroup struct {
	SystemApiGroup system.ApiGroup
	ImApiGroup     im.ApiGroup
}
