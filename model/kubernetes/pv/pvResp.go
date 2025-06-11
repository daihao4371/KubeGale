package pv

import (
	"KubeGale/model/common/request"

	v1 "k8s.io/api/core/v1"
)

// PVListResponse 持久卷列表响应结构体
// 用于返回持久卷的列表信息，包含分页信息
type PVListResponse struct {
	Items *[]v1.PersistentVolume `json:"items" form:"items"` // 持久卷列表
	Total int                    `json:"total" form:"total"` // 总数
	request.PageInfo
}

// DescribePVResponse 持久卷详情响应结构体
// 用于返回单个持久卷的详细信息
type DescribePVResponse struct {
	Items *v1.PersistentVolume `json:"items" form:"items"` // 持久卷详情
}
