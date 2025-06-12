package limitRange

import (
	v1 "KubeGale/api/v1"
	"KubeGale/middleware"
	"github.com/gin-gonic/gin"
)

type K8sLimitRangeRouter struct{}

func (s *K8sLimitRangeRouter) Initk8sLimitRangeRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	k8sNodeRouter := Router.Group("kubernetes").Use(middleware.OperationRecord())
	k8sLimitRangeRouterWithoutRecord := Router.Group("kubernetes")
	var k8sLimitRangeApi = v1.ApiGroupApp.LimitRange.K8sLimitRangeApi
	{
		k8sNodeRouter.POST("limitRange", k8sLimitRangeApi.CreateLimitRange)
		k8sNodeRouter.DELETE("limitRange", k8sLimitRangeApi.DeleteLimitRange)
		k8sNodeRouter.PUT("limitRange", k8sLimitRangeApi.UpdateLimitRange)
	}
	{
		k8sLimitRangeRouterWithoutRecord.GET("limitRange", k8sLimitRangeApi.GetLimitRangeList)
		k8sLimitRangeRouterWithoutRecord.GET("limitRangeDetails", k8sLimitRangeApi.DescribeLimitRangeInfo)
	}
}
