package service

import (
	v1 "KubeGale/api/v1"
	"KubeGale/middleware"
	"github.com/gin-gonic/gin"
)

type K8sServiceRouter struct{}

func (s *K8sServiceRouter) Initk8sServiceRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	k8sNodeRouter := Router.Group("kubernetes").Use(middleware.OperationRecord())
	k8sServiceRouterWithoutRecord := Router.Group("kubernetes")
	var k8sServiceApi = v1.ApiGroupApp.Service.K8sServiceApi
	{
		k8sNodeRouter.POST("service", k8sServiceApi.CreateService)
		k8sNodeRouter.DELETE("service", k8sServiceApi.DeleteService)
		k8sNodeRouter.PUT("service", k8sServiceApi.UpdateService)

	}
	{
		k8sServiceRouterWithoutRecord.GET("service", k8sServiceApi.GetServiceList)
		k8sServiceRouterWithoutRecord.GET("serviceDetails", k8sServiceApi.DescribeServiceInfo)

	}

}
