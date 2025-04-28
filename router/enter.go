package router

import (
	"KubeGale/router/im"
	"KubeGale/router/system"
)

var RouterGroupApp = new(RouterGroup)

type RouterGroup struct {
	System system.RouterGroup
	Im     im.RouterGroup
}
