package daemonset

import (
	v1 "KubeGale/api/v1"
	"KubeGale/middleware"
	"github.com/gin-gonic/gin"
)

type K8sDaemonSetRouter struct{}

func (s *K8sDaemonSetRouter) Initk8sDaemonSetRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	k8sDaemonSetRouter := Router.Group("kubernetes").Use(middleware.OperationRecord())
	k8sDaemonSetRouterWithoutRecord := Router.Group("kubernetes")
	var k8sDaemonSetApi = v1.ApiGroupApp.DaemonSet.K8sDaemonSetApi
	{
		k8sDaemonSetRouter.PUT("daemonset", k8sDaemonSetApi.UpdateDaemonsSet)
		k8sDaemonSetRouter.DELETE("daemonset", k8sDaemonSetApi.DeleteDaemonSet)
		k8sDaemonSetRouter.POST("daemonset", k8sDaemonSetApi.CreateDaemonSet)
	}
	{
		k8sDaemonSetRouterWithoutRecord.GET("daemonset", k8sDaemonSetApi.GetDaemonSetList) // 获取DaemonSet列表
		k8sDaemonSetRouterWithoutRecord.GET("daemonsetDetails", k8sDaemonSetApi.DescribeDaemonSetInfo)
	}
}
