package cmdb

import (
	"KubeGale/middleware"
	"github.com/gin-gonic/gin"
)

type BatchOperationsRouter struct{}

// InitBatchOperationsRouter 初始化 BatchOperations表 路由信息
func (s *BatchOperationsRouter) InitBatchOperationsRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	BatchOperationsRouter := Router.Group("cmdb").Use(middleware.OperationRecord())
	{
		BatchOperationsRouter.POST("batchOperations/execute", batchOperationsApi.ExecuteCommands) // 新建BatchOperations表
		BatchOperationsRouter.GET("batchOperations/execLogs", batchOperationsApi.ExecuteRecords)  // 获取执行记录
	}
}
