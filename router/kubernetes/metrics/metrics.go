package metrics

import (
	api "KubeGale/api/v1"
	"KubeGale/middleware"
	"github.com/gin-gonic/gin"
)

type MetricRouter struct{}

func (u *MetricRouter) InitMetricRouter(Router *gin.RouterGroup) {
	metricRouter := Router.Group("metrics").Use(middleware.OperationRecord())
	//metricsApi := api.ApiGroupApp.MetricsApi
	metricsApi := api.ApiGroupApp.Metrics
	{
		metricRouter.POST("get", metricsApi.MetricsGet) // 监控数据获取
	}
}
