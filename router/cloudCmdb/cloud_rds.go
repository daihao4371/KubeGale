package cloudCmdb

import (
	api "KubeGale/api/v1/cloudCmdb"
	"KubeGale/middleware"

	"github.com/gin-gonic/gin"
)

type CloudRDSRouter struct{}

func (r *CloudRDSRouter) InitRDSRouter(Router *gin.RouterGroup) {
	cloudRdsRouter := Router.Group("rds").Use(middleware.OperationRecord())
	cloudRdsRouterWithoutRecord := Router.Group("rds")
	cloudRdsApi := api.ApiGroupApp.CloudRDSApi
	{
		cloudRdsRouter.POST("sync", cloudRdsApi.CloudRDSSync)  // 同步RDS
		cloudRdsRouter.POST("tree", cloudRdsApi.CloudRDSTree)  // 目录树
		cloudRdsRouter.POST("get", cloudRdsApi.GetRDSInstance) // 获取RDS实例详情
	}

	{
		cloudRdsRouterWithoutRecord.POST("list", cloudRdsApi.CloudRDSList) // 分页获取列表
	}
}
