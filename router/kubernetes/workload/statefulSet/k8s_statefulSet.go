package statefulSet

import (
	v1 "KubeGale/api/v1"
	"KubeGale/middleware"
	"github.com/gin-gonic/gin"
)

type K8sStatefulSetRouter struct{}

func (s *K8sStatefulSetRouter) Initk8sStatefulSetRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	k8sStatefulSetRouter := Router.Group("kubernetes").Use(middleware.OperationRecord())
	k8sStatefulSetRouterWithoutRecord := Router.Group("kubernetes")
	var k8sStatefulSetApi = v1.ApiGroupApp.StatefulSet.K8sStatefulSetApi
	{
		k8sStatefulSetRouter.POST("statefulset", k8sStatefulSetApi.CreateStatefulset)   // 新建k8sCluster表
		k8sStatefulSetRouter.DELETE("statefulset", k8sStatefulSetApi.DeleteStatefulSet) // 删除k8sCluster表
		k8sStatefulSetRouter.PUT("statefulset", k8sStatefulSetApi.UpdateStatefulSetInfo)
	}
	{
		k8sStatefulSetRouterWithoutRecord.GET("statefulset", k8sStatefulSetApi.GetStatefulSetList) // 获取StatefulSet列表
		k8sStatefulSetRouterWithoutRecord.GET("statefulsetDetails", k8sStatefulSetApi.DescribeStatefulSetInfo)
	}
}
