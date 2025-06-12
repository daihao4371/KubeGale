package cloudtty

import (
	v1 "KubeGale/api/v1"
	"KubeGale/middleware"
	"github.com/gin-gonic/gin"
)

type K8sNodeTTYRouter struct{}

func (s *K8sNodeTTYRouter) Initk8sNodeTTYRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	k8sNodeTTYRouter := Router.Group("kubernetes").Use(middleware.OperationRecord())
	var k8sNodeTTYApi = v1.ApiGroupApp.CloudTTY.K8SNodeTTYApi
	{
		k8sNodeTTYRouter.POST("/nodetty/get", k8sNodeTTYApi.NodeTTYGet) // NodeTTY
	}

}
