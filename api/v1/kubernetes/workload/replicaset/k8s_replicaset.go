package replicaset

import (
	"KubeGale/global"
	"KubeGale/model/common/request"
	"KubeGale/model/common/response"
	"KubeGale/model/kubernetes/replicaSet"
	"KubeGale/service"
	"KubeGale/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type K8sReplicaSetApi struct{}

var k8sReplicaSetService = service.ServiceGroupApp.ReplicaSetServiceGroup.K8sReplicaSetService

func (k *K8sReplicaSetApi) GetReplicaSetList(c *gin.Context) {
	req := replicaSet.GetReplicaSetListReq{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if list, total, err := k8sReplicaSetService.GetReplicaSetList(req, utils.GetUserUuid(c)); err != nil {
		global.KUBEGALE_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithDetailed(replicaSet.ReplicaSetListResponse{
			Items: list,
			Total: total,
			PageInfo: request.PageInfo{
				Page:     req.Page,
				PageSize: req.PageSize,
				Keyword:  req.Keyword,
			},
		}, "获取成功", c)
	}
}
