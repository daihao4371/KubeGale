package secret

import (
	v1 "KubeGale/api/v1"
	"KubeGale/middleware"
	"github.com/gin-gonic/gin"
)

type K8sSecretRouter struct{}

func (s *K8sSecretRouter) Initk8sSecretRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	k8sNodeRouter := Router.Group("kubernetes").Use(middleware.OperationRecord())
	k8sSecretRouterWithoutRecord := Router.Group("kubernetes")
	var k8sSecretApi = v1.ApiGroupApp.Secret.K8sSecretApi
	{
		k8sNodeRouter.POST("secret", k8sSecretApi.CreateSecret)   // 新建k8sCluster表
		k8sNodeRouter.DELETE("secret", k8sSecretApi.DeleteSecret) // 新建k8sCluster表
		k8sNodeRouter.PUT("secret", k8sSecretApi.UpdateSecret)    // 新建k8sCluster表
	}
	{
		k8sSecretRouterWithoutRecord.GET("secret", k8sSecretApi.GetSecretList)             // 获取node列表
		k8sSecretRouterWithoutRecord.GET("secretDetails", k8sSecretApi.DescribeSecretInfo) // 获取node列表
	}
}
