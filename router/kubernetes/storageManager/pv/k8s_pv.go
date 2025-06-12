package pv

import (
	v1 "KubeGale/api/v1"
	"KubeGale/middleware"
	"github.com/gin-gonic/gin"
)

type K8sPVRouter struct{}

func (s *K8sPVRouter) Initk8sPVRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	k8sPVRouter := Router.Group("kubernetes").Use(middleware.OperationRecord())
	k8sPVRouterWithoutRecord := Router.Group("kubernetes")
	var k8sPVApi = v1.ApiGroupApp.Pv.K8sPVApi
	{
		k8sPVRouter.POST("pv", k8sPVApi.CreatePV)   // 新建k8sCluster表
		k8sPVRouter.DELETE("pv", k8sPVApi.DeletePV) // 删除k8sCluster表
		k8sPVRouter.PUT("pv", k8sPVApi.UpdatePV)
	}
	{
		k8sPVRouterWithoutRecord.GET("pv", k8sPVApi.GetPVList) // 获取PV列表
		k8sPVRouterWithoutRecord.GET("pvDetails", k8sPVApi.DescribePVInfo)
	}
}
