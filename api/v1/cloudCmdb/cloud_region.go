package cloudCmdb

import (
	"KubeGale/global"
	"KubeGale/model/common/request"
	"KubeGale/model/common/response"
	"KubeGale/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CloudRegionApi struct{}

// CloudPlatformSyncRegion 同步云平台区域
func (r *CloudRegionApi) CloudPlatformSyncRegion(c *gin.Context) {
	var idInfo request.GetById
	_ = c.ShouldBindJSON(&idInfo)
	if err := utils.Verify(idInfo, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err := cloudRegionService.SyncRegion(idInfo.ID)
	if err != nil {
		global.KUBEGALE_LOG.Error("同步操作失败!", zap.Error(err))
		response.FailWithMessage("同步操作失败", c)
	} else {
		response.OkWithMessage("同步操作成功, 数据异步处理中, 请稍后!", c)
	}
}
