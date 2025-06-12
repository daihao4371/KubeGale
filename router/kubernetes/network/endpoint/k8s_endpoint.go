package endpoint

import (
	v1 "KubeGale/api/v1"
	"KubeGale/middleware"
	"github.com/gin-gonic/gin"
)

type K8sEndPointRouter struct{}

func (s *K8sEndPointRouter) Initk8sEndPointRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	k8sEndPointRouter := Router.Group("kubernetes").Use(middleware.OperationRecord())
	k8sEndPointRouterWithoutRecord := Router.Group("kubernetes")
	var k8sEndPointApi = v1.ApiGroupApp.Endpoint.K8sEndPointApi
	{
		k8sEndPointRouter.POST("endpoint", k8sEndPointApi.CreateEndPoint)   // 新建k8sCluster表
		k8sEndPointRouter.DELETE("endpoint", k8sEndPointApi.DeleteEndPoint) // 删除k8sCluster表
		k8sEndPointRouter.PUT("endpoint", k8sEndPointApi.UpdateEndPoint)
	}
	{
		k8sEndPointRouterWithoutRecord.GET("endpoint", k8sEndPointApi.GetEndPointList) // 获取EndPoint列表
		k8sEndPointRouterWithoutRecord.GET("endpointDetails", k8sEndPointApi.DescribeEndPointInfo)
	}
}
