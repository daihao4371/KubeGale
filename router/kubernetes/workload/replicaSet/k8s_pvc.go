package replicaSet

import (
	v1 "KubeGale/api/v1"
	"github.com/gin-gonic/gin"
)

type K8sReplicaSetRouter struct{}

func (s *K8sReplicaSetRouter) Initk8sReplicaSetRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	k8sReplicaSetRouterWithoutRecord := Router.Group("kubernetes")
	var k8sReplicaSetApi = v1.ApiGroupApp.Replicaset.K8sReplicaSetApi

	{
		k8sReplicaSetRouterWithoutRecord.GET("replicaSet", k8sReplicaSetApi.GetReplicaSetList) // 获取node列表
	}
}
