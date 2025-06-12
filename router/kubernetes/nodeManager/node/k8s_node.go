package node

import (
	v1 "KubeGale/api/v1"
	"KubeGale/middleware"
	"github.com/gin-gonic/gin"
)

type K8sNodeRouter struct{}

func (s *K8sNodeRouter) Initk8sNodeRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	k8sNodeRouter := Router.Group("kubernetes").Use(middleware.OperationRecord())
	k8sNodeRouterWithoutRecord := Router.Group("kubernetes")
	var k8sNodeApi = v1.ApiGroupApp.NodeApiGroup.K8sNodeApi
	{
		k8sNodeRouter.PUT("nodes", k8sNodeApi.UpdateNodeInfo)
		k8sNodeRouter.POST("nodes/EvictAllPod", k8sNodeApi.EvictAllNodePod)
	}
	{
		k8sNodeRouterWithoutRecord.GET("nodes", k8sNodeApi.GetNodeList) // 获取node列表
		k8sNodeRouterWithoutRecord.GET("nodes/metrics", k8sNodeApi.GetNodeMetricsList)
		k8sNodeRouterWithoutRecord.GET("nodeDetails", k8sNodeApi.DescribeNodeInfo)
	}
}
