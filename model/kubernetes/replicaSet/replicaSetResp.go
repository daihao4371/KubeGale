package replicaSet

import (
	"KubeGale/model/common/request"

	v1 "k8s.io/api/apps/v1"
)

// ReplicaSetListResponse 副本集列表响应结构体
// 用于返回副本集的列表信息，包含分页信息
type ReplicaSetListResponse struct {
	Items *[]v1.ReplicaSet `json:"items" form:"items"` // 副本集列表
	Total int              `json:"total" form:"total"` // 总数
	request.PageInfo
}
