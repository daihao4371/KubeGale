package nodes

import (
	"KubeGale/model/common/request"

	corev1 "k8s.io/api/core/v1"
)

// NodeListResponse 节点列表响应结构体
// 用于返回节点的列表信息，包含分页信息
type NodeListResponse struct {
	Items *[]corev1.Node `json:"items" form:"items"` // 节点列表
	Total int            `json:"total" form:"total"` // 总数
	request.PageInfo
}

// DescribeNodeInfoResponse 节点详情响应结构体
// 用于返回单个节点的详细信息
type DescribeNodeInfoResponse struct {
	Items *corev1.Node `json:"items" form:"items"` // 节点详情
}
