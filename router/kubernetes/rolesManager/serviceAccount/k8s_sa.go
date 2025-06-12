package serviceAccount

import (
	v1 "KubeGale/api/v1"
	"KubeGale/middleware"
	"github.com/gin-gonic/gin"
)

type K8sServiceAccountRouter struct{}

func (s *K8sServiceAccountRouter) Initk8sServiceAccountRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	k8sNodeRouter := Router.Group("kubernetes").Use(middleware.OperationRecord())
	k8sServiceAccountRouterWithoutRecord := Router.Group("kubernetes")

	var k8sServiceAccountApi = v1.ApiGroupApp.ServiceAccount.K8sServiceAccountApi
	{
		k8sNodeRouter.POST("serviceAccount", k8sServiceAccountApi.CreateServiceAccount)
		k8sNodeRouter.DELETE("serviceAccount", k8sServiceAccountApi.DeleteServiceAccount)
		k8sNodeRouter.PUT("serviceAccount", k8sServiceAccountApi.UpdateServiceAccount)
	}
	{
		k8sServiceAccountRouterWithoutRecord.GET("serviceAccount", k8sServiceAccountApi.GetServiceAccount)                 // 获取node列表
		k8sServiceAccountRouterWithoutRecord.GET("serviceAccountDetails", k8sServiceAccountApi.DescribeServiceAccountInfo) // 获取node列表
	}
}
