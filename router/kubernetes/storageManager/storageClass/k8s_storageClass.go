package storageClass

import (
	v1 "KubeGale/api/v1"
	"KubeGale/middleware"
	"github.com/gin-gonic/gin"
)

type K8sStorageClassRouter struct{}

func (s *K8sStorageClassRouter) Initk8sStorageClassRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	k8sStorageClassRouter := Router.Group("kubernetes").Use(middleware.OperationRecord())
	k8sStorageClassRouterWithoutRecord := Router.Group("kubernetes")
	var k8sStorageClassApi = v1.ApiGroupApp.StorageClass.K8sStorageClassApi
	{
		k8sStorageClassRouter.POST("storageClass", k8sStorageClassApi.CreateStorageClass)   // 新建k8sCluster表
		k8sStorageClassRouter.DELETE("storageClass", k8sStorageClassApi.DeleteStorageClass) // 删除k8sCluster表
		k8sStorageClassRouter.PUT("storageClass", k8sStorageClassApi.UpdateStorageClass)
	}
	{
		k8sStorageClassRouterWithoutRecord.GET("storageClass", k8sStorageClassApi.GetStorageClassList) // 获取StorageClass列表
		k8sStorageClassRouterWithoutRecord.GET("storageClassDetails", k8sStorageClassApi.DescribeStorageClassInfo)
	}
}
