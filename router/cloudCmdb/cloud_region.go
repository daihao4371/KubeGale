package cloudCmdb

import (
	api "KubeGale/api/v1/cloudCmdb"
	"KubeGale/middleware"

	"github.com/gin-gonic/gin"
)

type CloudRegionRouter struct {
}

func (c *CloudRegionRouter) InitCloudRegionRouter(Router *gin.RouterGroup) {
	cloudRegionRouter := Router.Group("cloud_region").Use(middleware.OperationRecord())
	cloudRegionApi := api.ApiGroupApp.CloudRegionApi
	{
		cloudRegionRouter.POST("syncRegion", cloudRegionApi.CloudPlatformSyncRegion) // 同步区域信息
		cloudRegionRouter.GET("tree", cloudRegionApi.GetRegionTree)                  // 获取区域树形结构
	}
}
