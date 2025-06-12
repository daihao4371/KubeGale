package cloudtty

import (
	v1 "KubeGale/api/v1"
	"KubeGale/middleware"
	"github.com/gin-gonic/gin"
)

type K8sCloudTTYRouter struct{}

func (s *K8sCloudTTYRouter) Initk8sCloudTTYRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	k8sCloudTTYRouter := Router.Group("kubernetes").Use(middleware.OperationRecord())
	var k8sCloudTTYApi = v1.ApiGroupApp.CloudTTY.K8sCloudTTYApi
	{
		k8sCloudTTYRouter.POST("/cloudtty/get", k8sCloudTTYApi.CloudTTYGet) // CloudTTY
	}
}
