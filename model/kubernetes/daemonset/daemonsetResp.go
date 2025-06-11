package daemonset

import (
	"KubeGale/model/common/request"

	v1 "k8s.io/api/apps/v1"
)

// DaemonSetListResponse DaemonSet列表响应结构体
// 用于返回DaemonSet的列表信息，包含分页信息
type DaemonSetListResponse struct {
	Items            *[]v1.DaemonSet `json:"items" form:"items"` // DaemonSet列表
	Total            int             `json:"total" form:"total"` // 总数
	request.PageInfo                 // 分页信息
}

// DescribeDaemonSetResponse DaemonSet详情响应结构体
// 用于返回单个DaemonSet的详细信息
type DescribeDaemonSetResponse struct {
	Items *v1.DaemonSet `json:"items" form:"items"` // DaemonSet详情
}
