package cloudCmdb

import (
	"KubeGale/api/v1/cloudCmdb"
	"KubeGale/middleware"
	"github.com/gin-gonic/gin"
)

type CloudStatsRouter struct{}

func (r *CloudStatsRouter) InitCloudStatsRouter(Router *gin.RouterGroup) {
	cloudStatsApi := cloudCmdb.ApiGroupApp.CloudStatsApi // Assumes ApiGroupApp structure is in api/v1/cloudCmdb/enter.go
	
	// All cloud CMDB stats routes require authentication and casbin checks
	cloudCmdbStatsAuthRouter := Router.Group("stats").Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())
	{
		cloudCmdbStatsAuthRouter.GET("countsByProvider", cloudStatsApi.GetResourceCountsByProviderHandler)
	}
}
