package pvc

import (
	v1 "KubeGale/api/v1"
	"KubeGale/middleware"
	"github.com/gin-gonic/gin"
)

type K8sPvcRouter struct{}

func (s *K8sPvcRouter) Initk8sPvcRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	k8sNodeRouter := Router.Group("kubernetes").Use(middleware.OperationRecord())
	k8sPvcRouterWithoutRecord := Router.Group("kubernetes")
	var k8sPvcApi = v1.ApiGroupApp.PvcGroup.K8sPvcApi
	{
		k8sNodeRouter.POST("pvc", k8sPvcApi.CreatePVC)   // 新建k8sCluster表
		k8sNodeRouter.DELETE("pvc", k8sPvcApi.DeletePVC) // 删除k8sCluster表
		k8sNodeRouter.PUT("pvc", k8sPvcApi.UpdatePVC)    // 删除k8sCluster表
	}
	{
		k8sPvcRouterWithoutRecord.GET("pvc", k8sPvcApi.GetPvcList)             // 获取node列表
		k8sPvcRouterWithoutRecord.GET("pvcDetails", k8sPvcApi.DescribePVCInfo) // 获取node列表
	}
}
