package clusterRole

import (
	v1 "KubeGale/api/v1"
	"KubeGale/middleware"
	"github.com/gin-gonic/gin"
)

type K8sClusterRoleRouter struct{}

func (s *K8sClusterRoleRouter) Initk8sClusterRoleRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	k8sClusterRoleRouter := Router.Group("kubernetes").Use(middleware.OperationRecord())
	k8sClusterRoleRouterWithoutRecord := Router.Group("kubernetes")
	var k8sClusterRoleApi = v1.ApiGroupApp.ClusterRole.K8sClusterRoleApi
	{
		k8sClusterRoleRouter.POST("ClusterRole", k8sClusterRoleApi.CreateClusterRole)   // 新建k8sCluster表
		k8sClusterRoleRouter.DELETE("ClusterRole", k8sClusterRoleApi.DeleteClusterRole) // 删除k8sCluster表
		k8sClusterRoleRouter.PUT("ClusterRole", k8sClusterRoleApi.UpdateClusterRole)
	}
	{
		k8sClusterRoleRouterWithoutRecord.GET("ClusterRole", k8sClusterRoleApi.GetClusterRoleList) // 获取ClusterRole列表
		k8sClusterRoleRouterWithoutRecord.GET("ClusterRoleDetails", k8sClusterRoleApi.DescribeClusterRoleInfo)
	}
}
