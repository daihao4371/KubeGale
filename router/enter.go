package router

import (
	"KubeGale/router/cmdb"
	"KubeGale/router/im"
	"KubeGale/router/system"
)

var RouterGroupApp = new(RouterGroup)

type RouterGroup struct {
	System system.RouterGroup
	Im     im.RouterGroup
	Cmdb   cmdb.RouterGroup
}
