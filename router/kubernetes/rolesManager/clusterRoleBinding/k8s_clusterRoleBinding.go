package clusterRoleBinding

import (
	v1 "KubeGale/api/v1"
	"KubeGale/middleware"
	"github.com/gin-gonic/gin"
)

type K8sClusterRoleBindingRouter struct{}

func (s *K8sClusterRoleBindingRouter) Initk8sClusterRoleBindingRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	k8sClusterRoleBindingRouter := Router.Group("kubernetes").Use(middleware.OperationRecord())
	k8sClusterRoleBindingRouterWithoutRecord := Router.Group("kubernetes")
	var k8sClusterRoleBindingApi = v1.ApiGroupApp.ClusterRoleBinding.K8sClusterRoleBindingApi
	{
		k8sClusterRoleBindingRouter.POST("ClusterRoleBinding", k8sClusterRoleBindingApi.CreateClusterRoleBinding)   // 新建k8sCluster表
		k8sClusterRoleBindingRouter.DELETE("ClusterRoleBinding", k8sClusterRoleBindingApi.DeleteClusterRoleBinding) // 删除k8sCluster表
		k8sClusterRoleBindingRouter.PUT("ClusterRoleBinding", k8sClusterRoleBindingApi.UpdateClusterRoleBinding)
	}
	{
		k8sClusterRoleBindingRouterWithoutRecord.GET("ClusterRoleBinding", k8sClusterRoleBindingApi.GetClusterRoleBindingList) // 获取ClusterRoleBinding列表
		k8sClusterRoleBindingRouterWithoutRecord.GET("ClusterRoleBindingDetails", k8sClusterRoleBindingApi.DescribeClusterRoleBindingInfo)
	}
}
