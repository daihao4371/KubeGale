package cloudCmdb

import (
	api "KubeGale/api/v1/cloudCmdb"
	"KubeGale/middleware"

	"github.com/gin-gonic/gin"
)

type CloudLoadBalancerRouter struct {
}

func (l *CloudLoadBalancerRouter) InitLoadBalancerRouter(Router *gin.RouterGroup) {
	cloudLoadBalancerRouter := Router.Group("loadBalancer").Use(middleware.OperationRecord())
	cloudLoadBalancerRouterWithoutRecord := Router.Group("loadBalancer")
	cloudLoadBalancerApi := api.ApiGroupApp.CloudLoadBalancerApi
	{
		cloudLoadBalancerRouter.POST("sync", cloudLoadBalancerApi.CloudLoadBalancerSync) // 同步负载均衡
		cloudLoadBalancerRouter.POST("tree", cloudLoadBalancerApi.CloudLoadBalancerTree) // 目录树
	}

	{
		cloudLoadBalancerRouterWithoutRecord.POST("list", cloudLoadBalancerApi.CloudLoadBalancerList) // 分页获取列表
	}
}
